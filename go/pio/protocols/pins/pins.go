// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Package pins declare well known pins.
//
// Pins is about physical pins, not about their logical function.
//
// While not a protocol strictly speaking, these are "well known constants".
package pins

import (
	"errors"
	"fmt"
	"time"

	"github.com/maruel/dlibox/go/pio/protocols/gpio"
)

var (
	// INVALID implements gpio.PinIO and fails on all access.
	INVALID invalidPin
	GROUND  gpio.PinIO = &BasicPin{Name: "GROUND"}
	V3_3    gpio.PinIO = &BasicPin{Name: "V3_3"}
	V5      gpio.PinIO = &BasicPin{Name: "V5"}
)

// Pin is the minimal common interface shared between gpio.PinIO and
// analog.PinIO.
type Pin interface {
	fmt.Stringer
	Number() int
	Function() string
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

func (invalidPin) In(gpio.Pull, gpio.Edge) error {
	return invalidPinErr
}

func (invalidPin) Read() gpio.Level {
	return gpio.Low
}

func (invalidPin) WaitForEdge(timeout time.Duration) bool {
	return false
}

func (invalidPin) Pull() gpio.Pull {
	return gpio.PullNoChange
}

func (invalidPin) Out(gpio.Level) error {
	return invalidPinErr
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

// BasicPin implements gpio.PinIO as a non-functional pin.
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

func (b *BasicPin) In(gpio.Pull, gpio.Edge) error {
	return fmt.Errorf("%s cannot be used as input", b.Name)
}

func (b *BasicPin) Read() gpio.Level {
	return gpio.Low
}

func (b *BasicPin) WaitForEdge(timeout time.Duration) bool {
	return false
}

func (b *BasicPin) Pull() gpio.Pull {
	return gpio.PullNoChange
}

func (b *BasicPin) Out(gpio.Level) error {
	return fmt.Errorf("%s cannot be used as output", b.Name)
}

func (b *BasicPin) ADC() error {
	return fmt.Errorf("%s cannot be used as analog input", b.Name)
}

func (b *BasicPin) Range() (int32, int32) {
	return 0, 0
}

func (b *BasicPin) Measure() int32 {
	return 0
}

func (b *BasicPin) DAC(v int32) error {
	return fmt.Errorf("%s cannot be used as analog output", b.Name)
}
