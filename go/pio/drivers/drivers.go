// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

//go:generate stringer -type Type

package drivers

import (
	"errors"
	"fmt"
	"sync"
)

// Type represent the type of driver.
//
// Lower is more important.
type Type int

const (
	// Processor is the first driver to be loaded.
	Processor Type = iota
	// Pins is basic pin functionality driver, additional to Processor.
	//
	// This includes all headers description.
	Pins
	// Functional is for functionality pin driver, additional to Pins.
	Functional
	// Bus is higher level protocol drivers.
	Bus
	// Device is drivers connecting to buses.
	Device
	nbPriorities
)

// Driver is an implementation for a protocol.
type Driver interface {
	// String returns the name of the driver, as to be presented to the user. It
	// should be unique.
	String() string
	// Type is the type of driver.
	//
	// This is used to load the drivers in order.
	//
	// If a driver implements multiple levels of functionality, it should return
	// the most important one, the one with the lowest value.
	Type() Type
	// Prerequisites returns a list of drivers that must be successfully loaded
	// first before attempting to load this driver.
	Prerequisites() []string
	// Init initializes the driver. It should return false, nil when the driver
	// is irrelevant on the platform.
	Init() (bool, error)
}

// DriverFailure is a driver that failed loaded.
type DriverFailure struct {
	D   Driver
	Err error
}

// State is the state of loaded device drivers.
type State struct {
	Loaded  []Driver
	Skipped []Driver
	Failed  []DriverFailure
}

// Init initially all the relevant drivers.
//
// Drivers are started concurrently for Type.
//
// It returns the list of all drivers loaded and errors on the first call, if
// any. They are ordered by Type but unordered within each type.
//
// Second call is ignored and errors are discarded.
//
// Users will want to use host.Init(), which guarantees a baseline of included
// drivers.
func Init() (*State, error) {
	lockState.Lock()
	defer lockState.Unlock()
	if state != nil {
		return state, nil
	}
	state = &State{}
	cD := make(chan Driver)
	cS := make(chan Driver)
	cE := make(chan DriverFailure)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for d := range cD {
			state.Loaded = append(state.Loaded, d)
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for d := range cS {
			state.Skipped = append(state.Skipped, d)
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for f := range cE {
			state.Failed = append(state.Failed, f)
		}
	}()

	stages, err := getStages()
	if err != nil {
		return state, err
	}
	loaded := map[string]struct{}{}
	for _, drivers := range stages {
		loadStage(drivers, loaded, cD, cS, cE)
	}
	close(cD)
	close(cS)
	close(cE)
	wg.Wait()
	return state, nil
}

// Register registers a driver to be initialized automatically on Init().
//
// The d.String() value must be unique across all registered drivers.
//
// Can't call Register() while Init() is running.
func Register(d Driver) error {
	lockState.Lock()
	loaded := state != nil
	lockState.Unlock()
	if loaded {
		return errors.New("drivers: can't call Register() after Init()")
	}

	lockDrivers.Lock()
	defer lockDrivers.Unlock()
	n := d.String()
	if _, ok := byName[n]; ok {
		return fmt.Errorf("drivers.Register(%s): driver with same name was already registered", d)
	}
	byName[n] = d
	t := d.Type()
	allDrivers[t] = append(allDrivers[t], d)
	return nil
}

// MustRegister calls Register and panics if registration fails.
func MustRegister(d Driver) {
	if err := Register(d); err != nil {
		panic(err)
	}
}

//

var (
	lockDrivers sync.Mutex
	allDrivers  [nbPriorities][]Driver
	byName      = map[string]Driver{}
	lockState   sync.Mutex
	state       *State
)

// getStages returns a set of stages to load the drivers.
//
// Loading is done using two blocking mechanism:
// - By type
// - By prerequisites
// So create a DAG but reduce it as a list of stages.
//
// This cannot be done in Register() since the drivers are not registered in
// order.
func getStages() ([][]Driver, error) {
	lockDrivers.Lock()
	defer lockDrivers.Unlock()
	var stages [][]Driver
	for _, drivers := range allDrivers {
		if len(drivers) == 0 {
			// No driver registered for this type.
			continue
		}
		inner, err := explodeStages(drivers)
		if err != nil {
			return nil, err
		}
		if len(inner) != 0 {
			stages = append(stages, inner...)
		}
	}
	return stages, nil
}

// explodeStages creates multiple intermediate stages if needed.
//
// It searches if there's any driver than has dependency on another driver from
// this stage and creates intermediate stage if so.
func explodeStages(drivers []Driver) ([][]Driver, error) {
	dependencies := map[string]map[string]struct{}{}
	for _, d := range drivers {
		dependencies[d.String()] = map[string]struct{}{}
	}
	for _, d := range drivers {
		name := d.String()
		t := d.Type()
		for _, depName := range d.Prerequisites() {
			dep, ok := byName[depName]
			if !ok {
				return nil, fmt.Errorf("drivers: unsatified dependency %#v->%#v; it is missing; skipping", name, depName)
			}
			dt := dep.Type()
			if dt > t {
				return nil, fmt.Errorf("drivers: inversed dependency %#v(%s)->%#v(%s); skipping", name, t, depName, dt)
			}
			if dt < t {
				// Staging already takes care of this.
				continue
			}
			// Dependency between two drivers of the same type. This can happen
			// when there's a process class driver and a processor specialization
			// driver. As an example, allwinner->R8, allwinner->A64, etc.
			dependencies[name][depName] = struct{}{}
		}
	}

	var stages [][]Driver
	for len(dependencies) != 0 {
		// Create a stage.
		var stage []string
		var l []Driver
		for name, deps := range dependencies {
			if len(deps) == 0 {
				stage = append(stage, name)
				l = append(l, byName[name])
				delete(dependencies, name)
			}
		}
		if len(stage) == 0 {
			return nil, fmt.Errorf("drivers: found cycle(s) in drivers dependencies; %v", dependencies)
		}
		stages = append(stages, l)

		// Trim off.
		for _, passed := range stage {
			for name := range dependencies {
				delete(dependencies[name], passed)
			}
		}
	}
	return stages, nil
}

// loadStage loads all the drivers in this stage concurrently.
func loadStage(drivers []Driver, loaded map[string]struct{}, cD chan<- Driver, cS chan<- Driver, cE chan<- DriverFailure) {
	var wg sync.WaitGroup
	// Use int for concurrent access.
	skip := make([]int, len(drivers))
	for i, driver := range drivers {
		// Load only the driver if prerequisites were loaded. They are
		// guaranteed to be in a previous stage by getStages().
		for _, dep := range driver.Prerequisites() {
			if _, ok := loaded[dep]; !ok {
				skip[i] = 1
				//log.Printf("drivers: skipping %s because missing %s", driver, dep)
				break
			}
		}
	}

	for i, driver := range drivers {
		if skip[i] != 0 {
			cS <- driver
			continue
		}
		wg.Add(1)
		go func(d Driver, j int) {
			defer wg.Done()
			if ok, err := d.Init(); ok {
				if err == nil {
					cD <- d
					return
				}
				cE <- DriverFailure{d, fmt.Errorf("drivers: %s.Init() failed: %v", d, err)}
			} else {
				cS <- d
				//log.Printf("drivers: %s.Init() skipped initialization", d)
			}
			skip[j] = 1
		}(driver, i)
	}
	wg.Wait()

	for i, driver := range drivers {
		if skip[i] != 0 {
			continue
		}
		loaded[driver.String()] = struct{}{}
	}
}
