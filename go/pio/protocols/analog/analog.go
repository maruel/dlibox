// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Package analog defines analog pins, both DAC and ADC.
package analog

import (
	"errors"
	"fmt"

	"github.com/maruel/dlibox/go/pio/protocols/pins"
)

// ADC is an analog-to-digital-conversion input.
type ADC interface {
	pins.Pin
	// In setups a pin as an input.
	ADC() error
	// Range returns the maximum supported range [min, max] of the values.
	Range() (int32, int32)
	// Measure return the current pin level.
	//
	// Behavior is undefined if In() wasn't used before.
	//
	// In some rare case, it is possible that Read() fails silently. This happens
	// if another process on the host messes up with the pin after In() was
	// called. In this case, call In() again.
	Measure() int32
}

// DAC is an digital-to-analog-conversion output.
type DAC interface {
	pins.Pin
	// DAC sets a pin as output if it wasn't already and sets the value.
	//
	// After the initial call to ensure that the pin has been set as output, it
	// is generally safe to ignore the error returned.
	DAC(v int32) error
}

// PinIO is a pin that supports both input and output.
//
// It may fail at either input and or output, for example ground, vcc and other
// similar pins.
type PinIO interface {
	pins.Pin
	ADC() error
	Range() (int32, int32)
	Measure() int32
	DAC(v int32) error
}

// INVALID implements PinIO and fails on all access.
var INVALID PinIO = invalidPin{}

// BasicPin implements Pin as a non-functional pin.
type BasicPin struct {
	Name string
}

func (b *BasicPin) Number() int {
	return -1
}

func (b *BasicPin) String() string {
	return b.Name
}

func (b *BasicPin) Function() string {
	return ""
}

func (b *BasicPin) ADC() error {
	return fmt.Errorf("%s cannot be used as ADC", b.Name)
}

func (b *BasicPin) Range() (int32, int32) {
	return 0, 0
}

func (b *BasicPin) Measure() int32 {
	return 0
}

func (b *BasicPin) DAC(v int32) error {
	return fmt.Errorf("%s cannot be used as DAC", b.Name)
}

//

// invalidPinErr is returned when trying to use INVALID.
var invalidPinErr = errors.New("invalid pin")

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

func (invalidPin) ADC() error {
	return invalidPinErr
}

func (invalidPin) Range() (int32, int32) {
	return 0, 0
}

func (invalidPin) Measure() int32 {
	return 0
}

func (invalidPin) DAC(v int32) error {
	return invalidPinErr
}
