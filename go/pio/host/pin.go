// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package host

import (
	"errors"
	"fmt"
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
	Edges() (<-chan Level, error)
	// DisableEdges() closes a previous Edges() channel and stops polling.
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

// invalidPinErr is returned when trying to use INVALID.
var invalidPinErr = errors.New("invalid pin")

// INVALID implements PinIO and fails on all access.
var INVALID invalidPin

// invalidPin implements PinIO for compability but fails on all access.
type invalidPin struct {
}

func (invalidPin) Number() int {
	return -1
}

func (invalidPin) String() string {
	return "INVALID"
}

func (invalidPin) Function() string {
	return ""
}

func (invalidPin) In(Pull) error {
	return invalidPinErr
}

func (invalidPin) Read() Level {
	return Low
}

func (invalidPin) Edges() (<-chan Level, error) {
	return nil, invalidPinErr
}

func (invalidPin) DisableEdges() {
}

func (invalidPin) Pull() Pull {
	return PullNoChange
}

func (invalidPin) Out(Level) error {
	return invalidPinErr
}
