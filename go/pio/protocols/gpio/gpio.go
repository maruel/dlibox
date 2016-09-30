// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Package gpio defines digital pins.
//
// The GPIO pins are described in their logical functionality, not in their
// physical position.
package gpio

import (
	"fmt"
	"sort"
	"sync"
)

// Level is the level of the pin: Low or High.
type Level bool

const (
	// Low represents 0v.
	Low Level = false
	// High represents Vin, generally 3.3v or 5v.
	High Level = true
)

func (l Level) String() string {
	if l == Low {
		return "Low"
	}
	return "High"
}

// Pull specifies the internal pull-up or pull-down for a pin set as input.
type Pull uint8

const (
	Float        Pull = 0 // Let the input float
	Down         Pull = 1 // Apply pull-down
	Up           Pull = 2 // Apply pull-up
	PullNoChange Pull = 3 // Do not change the previous pull resistor setting or an unknown value
)

const pullName = "FloatDownUpPullNoChange"

var pullIndex = [...]uint8{0, 5, 9, 11, 23}

func (i Pull) String() string {
	if i >= Pull(len(pullIndex)-1) {
		return fmt.Sprintf("Pull(%d)", i)
	}
	return pullName[pullIndex[i]:pullIndex[i+1]]
}

// PinIn is an input GPIO pin.
//
// It may optionally support internal pull resistor and edge based triggering.
type PinIn interface {
	// In setups a pin as an input.
	In(pull Pull) error
	// Read return the current pin level.
	//
	// Behavior is undefined if In() wasn't used before.
	//
	// In some rare case, it is possible that Read() fails silently. This happens
	// if another process on the host messes up with the pin after In() was
	// called. In this case, call In() again.
	Read() Level
	// Edges returns a channel that sends level changes.
	//
	// Behavior is undefined if In() wasn't used before.
	//
	// BUG(maruel): Still undecided on the form of edge detection and this
	// function may be refactored.
	Edges() (<-chan Level, error)
	// DisableEdges() closes a previous Edges() channel and stops polling.
	//
	// BUG(maruel): This function is the equivalent of In(PullNoChange). As such,
	// it is redundant.
	DisableEdges()
	// Pull returns the internal pull resistor if the pin is set as input pin.
	// Returns PullNoChange if the value cannot be read.
	Pull() Pull
}

// PinOut is an output GPIO pin.
type PinOut interface {
	// Out sets a pin as output if it wasn't already and sets the initial value.
	//
	// After the initial call to ensure that the pin has been set as output, it
	// is generally safe to ignore the error returned.
	Out(l Level) error
}

// PinIO is a GPIO pin that supports both input and output.
//
// It may fail at either input and or output, for example ground, vcc and other
// similar pins.
type PinIO interface {
	fmt.Stringer
	PinIn
	PinOut

	// Number returns the logical pin number or a negative number if the pin is
	// not a GPIO, e.g. GROUND, V3_3, etc.
	Number() int
	// Function returns a user readable string representation of what the pin is
	// configured to do. Common case is In and Out but it can be bus specific pin
	// name.
	Function() string
}

// ByNumber returns a GPIO pin from its number.
//
// Returns nil in case the pin is not present.
func ByNumber(number int) PinIO {
	pin, _ := byNumber[number]
	return pin
}

// PinByName returns a GPIO pin from its name.
//
// This can be strings like GPIO2, PB8, etc.
//
// Returns nil in case the pin is not present.
func ByName(name string) PinIO {
	pin, _ := byName[name]
	return pin
}

// ByFunction returns a GPIO pin from its function.
//
// This can be strings like I2C1_SDA, SPI0_MOSI, etc.
//
// Returns nil in case there is no pin setup with this function.
func ByFunction(fn string) PinIO {
	pin, _ := byFunction[fn]
	return pin
}

// All returns all the GPIO pins available on this host.
//
// The list is guaranteed to be in order of number.
//
// This list excludes non-GPIO pins like GROUND, V3_3, etc.
func All() []PinIO {
	// TODO(maruel): Return a copy?
	return all
}

// Functional returns a map of all pins implementing hardware provided
// special functionality, like I²C, SPI, ADC.
func Functional() map[string]PinIO {
	// TODO(maruel): Return a copy?
	return byFunction
}

// Register registers a GPIO pin.
//
// Registering the same pin number or name twice is an error.
func Register(pin PinIO) error {
	lock.Lock()
	defer lock.Unlock()
	number := pin.Number()
	if _, ok := byNumber[number]; ok {
		return fmt.Errorf("registering the same pin %d twice", number)
	}
	name := pin.String()
	if _, ok := byName[name]; ok {
		return fmt.Errorf("registering the same pin %s twice", name)
	}
	all = append(all, pin)
	sort.Sort(all)
	byNumber[number] = pin
	byName[name] = pin
	return nil
}

// MapFunction registers a GPIO pin for a specific function.
func MapFunction(function string, pin PinIO) {
	lock.Lock()
	defer lock.Unlock()
	byFunction[function] = pin
}

//

var (
	lock       sync.Mutex
	all        = pinList{}
	byNumber   = map[int]PinIO{}
	byName     = map[string]PinIO{}
	byFunction = map[string]PinIO{}
)

type pinList []PinIO

func (p pinList) Len() int           { return len(p) }
func (p pinList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p pinList) Less(i, j int) bool { return p[i].Number() < p[j].Number() }
