// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Package pins declare well known pins.
//
// Pins is about physical pins, not about their logical function.
//
// While not a protocol strictly speaking, these are "well known constants".
package pins

import "fmt"

var (
	INVALID Pin = &BasicPin{Name: "INVALID"}
	GROUND  Pin = &BasicPin{Name: "GROUND"}
	V3_3    Pin = &BasicPin{Name: "V3_3"}
	V5      Pin = &BasicPin{Name: "V5"}
)

// Pin is the minimal common interface shared between gpio.PinIO and
// analog.PinIO.
type Pin interface {
	fmt.Stringer
	// Number returns the logical pin number or a negative number if the pin is
	// not a GPIO, e.g. GROUND, V3_3, etc.
	Number() int
	// Function returns a user readable string representation of what the pin is
	// configured to do. Common case is In and Out but it can be bus specific pin
	// name.
	Function() string
}

//

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
