// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Package analog defines analog pins, both DAC and ADC.
package analog

import "fmt"

// ADC is an analog-to-digital-conversion input.
type ADC interface {
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
	fmt.Stringer
	ADC
	DAC

	// Number returns the logical pin number or a negative number if the pin is
	// not a GPIO, e.g. GROUND, V3_3, etc.
	Number() int
	// Function returns a user readable string representation of what the pin is
	// configured to do. Common case is ADC and DAC but it can be bus specific pin
	// name.
	Function() string
}
