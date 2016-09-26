// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package drivers

import (
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
	// Init initializes the driver. It should return false, nil when the driver
	// is irrelevant on the platform.
	Init() (bool, error)
}

// Init initially all the relevant drivers.
//
// Drivers are started concurrently for each of their group.
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
	if actual != nil {
		return actual, nil
	}
	actual = []Driver{}
	var errs []error
	cD := make(chan Driver)
	cE := make(chan error)
	var wg1 sync.WaitGroup
	wg1.Add(1)
	go func() {
		defer wg1.Done()
		for d := range cD {
			actual = append(actual, d)
		}
	}()
	wg1.Add(1)
	go func() {
		defer wg1.Done()
		for err := range cE {
			errs = append(errs, err)
		}
	}()
	for _, driversType := range drivers {
		if len(driversType) != 0 {
			var wg2 sync.WaitGroup
			for _, driver := range driversType {
				wg2.Add(1)
				go func(d Driver) {
					defer wg2.Done()
					ok, err := d.Init()
					if !ok {
						// Driver is not applicable to this host.
						return
					}
					if err != nil {
						cE <- fmt.Errorf("drivers: %s.Init() failed: %v", d, err)
						return
					}
					cD <- d
				}(driver)
			}
			wg2.Wait()
		}
	}
	close(cD)
	close(cE)
	wg1.Wait()
	return actual, errs
}

// Register registers a driver to be initialized automatically on Init().
//
// Calls to Register() after Init() are effectively ignored.
func Register(d Driver) {
	lockDrivers.Lock()
	defer lockDrivers.Unlock()
	t := d.Type()
	drivers[t] = append(drivers[t], d)
}

var (
	lockDrivers sync.Mutex
	drivers     [nbPriorities][]Driver
	lockActual  sync.Mutex
	actual      []Driver
)
