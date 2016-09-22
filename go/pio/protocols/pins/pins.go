// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Package pins contains well known pins.
//
// While not a protocol strictly speaking, these are "well known constants".
package pins

import (
	"errors"
	"fmt"

	"github.com/maruel/dlibox/go/pio/protocols/gpio"
)

var (
	// INVALID implements gpio.PinIO and fails on all access.
	INVALID     invalidPin
	GROUND      gpio.PinIO = &pin{name: "GROUND"}
	V3_3        gpio.PinIO = &pin{name: "V3_3"}
	V5          gpio.PinIO = &pin{name: "V5"}
	DC_IN       gpio.PinIO = &pin{name: "DC_IN"}
	TEMP_SENSOR gpio.PinIO = &pin{name: "TEMP_SENSOR"}
	BAT_PLUS    gpio.PinIO = &pin{name: "BAT_PLUS"}
	IR_RX       gpio.PinIO = &pin{name: "IR_RX"}
	EAROUTP     gpio.PinIO = &pin{name: "EAROUTP"} // Earpiece amplifier negative differential output
	EAROUTN     gpio.PinIO = &pin{name: "EAROUTN"} // Earpiece amplifier positive differential output
	CHARGER_LED gpio.PinIO = &pin{name: "CHARGER_LED"}
	RESET       gpio.PinIO = &pin{name: "RESET"}
	PWR_SWITCH  gpio.PinIO = &pin{name: "PWR_SWITCH"}
	KEY_ADC     gpio.PinIO = &pin{name: "KEY_ADC"}
	X32KFOUT    gpio.PinIO = &pin{name: "X32KFOUT"}
	VCC         gpio.PinIO = &pin{name: "VCC"}
	IOVCC       gpio.PinIO = &pin{name: "IOVCC"}
)

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
