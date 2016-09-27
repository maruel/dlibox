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
func Init() ([]Driver, []error) {
	lockActual.Lock()
	defer lockActual.Unlock()
	if actualDrivers != nil {
		return actualDrivers, nil
	}
	actualDrivers = []Driver{}
	var errs []error
	cD := make(chan Driver)
	cE := make(chan error)
	var wg1 sync.WaitGroup
	wg1.Add(1)
	go func() {
		defer wg1.Done()
		for d := range cD {
			actualDrivers = append(actualDrivers, d)
		}
	}()
	wg1.Add(1)
	go func() {
		defer wg1.Done()
		for err := range cE {
			errs = append(errs, err)
		}
	}()

	stages, err := getStages()
	if err != nil {
		return nil, []error{err}
	}
	for _, drivers := range stages {
		var wg2 sync.WaitGroup
		for _, driver := range drivers {
			wg2.Add(1)
			go func(d Driver) {
				defer wg2.Done()
				if ok, err := d.Init(); ok {
					if err != nil {
						cE <- fmt.Errorf("drivers: %s.Init() failed: %v", d, err)
					} else {
						cD <- d
					}
				}
			}(driver)
		}
		wg2.Wait()
	}
	close(cD)
	close(cE)
	wg1.Wait()
	return actualDrivers, errs
}

// Register registers a driver to be initialized automatically on Init().
//
// The d.String() value must be unique across all registered drivers.
//
// Can't call Register() while Init() is running.
func Register(d Driver) error {
	lockActual.Lock()
	loaded := actualDrivers != nil
	lockActual.Unlock()
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
	lockDrivers   sync.Mutex
	allDrivers    [nbPriorities][]Driver
	byName        = map[string]Driver{}
	lockActual    sync.Mutex
	actualDrivers []Driver
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

// Create multiple intermediate stages if needed.
func explodeStages(drivers []Driver) ([][]Driver, error) {
	// Search if there's any driver than has dependency on a driver from this
	// stage. This will create multiple intermediate stages.
	dependencies := map[string][]string{}
	for _, d := range drivers {
		dependencies[d.String()] = []string{}
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
			dependencies[name] = append(dependencies[name], depName)
		}
	}

	// Create a reverse dependency map.
	reverse := map[string]map[string]struct{}{}
	for name, deps := range dependencies {
		if reverse[name] == nil {
			reverse[name] = map[string]struct{}{}
		}
		for _, dep := range deps {
			if reverse[dep] == nil {
				reverse[dep] = map[string]struct{}{}
			}
			reverse[dep][name] = struct{}{}
		}
	}

	var stages [][]Driver
	for len(reverse) != 0 {
		// Create a stage.
		var stage []string
		for name, deps := range reverse {
			if len(deps) == 0 {
				stage = append(stage, name)
				delete(reverse, name)
			}
		}
		if len(stage) == 0 {
			return nil, fmt.Errorf("drivers: found cycle(s) in drivers dependencies; %v", dependencies)
		}
		l := make([]Driver, 0, len(stage))
		for _, n := range stage {
			l = append(l, byName[n])
		}
		stages = append(stages, l)

		// Trim off.
		for _, passed := range stage {
			for name := range reverse {
				delete(reverse[name], passed)
			}
		}
	}
	return stages, nil
}
