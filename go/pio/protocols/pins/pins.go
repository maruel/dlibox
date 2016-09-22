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

	"github.com/maruel/dlibox/go/pio/protocols/analog"
	"github.com/maruel/dlibox/go/pio/protocols/gpio"
)

var (
	// INVALID implements gpio.PinIO and fails on all access.
	INVALID     invalidPin
	GROUND      gpio.PinIO   = &pin{name: "GROUND"}
	V3_3        gpio.PinIO   = &pin{name: "V3_3"}
	V5          gpio.PinIO   = &pin{name: "V5"}
	VCC         gpio.PinIO   = &pin{name: "VCC"}         //
	DC_IN       gpio.PinIO   = &pin{name: "DC_IN"}       //
	TEMP_SENSOR gpio.PinIO   = &pin{name: "TEMP_SENSOR"} //
	BAT_PLUS    gpio.PinIO   = &pin{name: "BAT_PLUS"}    //
	IR_RX       gpio.PinIO   = &pin{name: "IR_RX"}       // IR Data Receive
	CHARGER_LED gpio.PinIO   = &pin{name: "CHARGER_LED"} //
	RESET       gpio.PinIO   = &pin{name: "RESET"}       //
	PWR_SWITCH  gpio.PinIO   = &pin{name: "PWR_SWITCH"}  //
	X32KFOUT    gpio.PinIO   = &pin{name: "X32KFOUT"}    // Clock output of 32Khz crystal
	IOVCC       gpio.PinIO   = &pin{name: "IOVCC"}       // Power supply for port A
	KEY_ADC     analog.PinIO = &pin{name: "KEY_ADC"}     // 6 bits resolution ADC for key application; can work up to 250Hz conversion rate; reference voltage is 2.0V
	EAROUTP     analog.PinIO = &pin{name: "EAROUTP"}     // Earpiece amplifier negative differential output
	EAROUTN     analog.PinIO = &pin{name: "EAROUTN"}     // Earpiece amplifier positive differential output
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

func (invalidPin) In(gpio.Pull) error {
	return invalidPinErr
}

func (invalidPin) Read() gpio.Level {
	return gpio.Low
}

func (invalidPin) Edges() (<-chan gpio.Level, error) {
	return nil, invalidPinErr
}

func (invalidPin) DisableEdges() {
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

// pin implements gpio.PinIO.
type pin struct {
	invalidPin
	name string
}

func (p *pin) String() string {
	return p.name
}

func (p *pin) In(gpio.Pull) error {
	return fmt.Errorf("%s cannot be used as input", p.name)
}

func (p *pin) Edges() (<-chan gpio.Level, error) {
	return nil, fmt.Errorf("%s cannot be used as input", p.name)
}

func (p *pin) Out(gpio.Level) error {
	return fmt.Errorf("%s cannot be used as output", p.name)
}

func (p *pin) ADC() error {
	return fmt.Errorf("%s cannot be used as analog input", p.name)
}

func (p *pin) Range() (int32, int32) {
	return 0, 0
}

func (p *pin) Measure() int32 {
	return 0
}

func (p *pin) DAC(v int32) error {
	return fmt.Errorf("%s cannot be used as analog output", p.name)
}
