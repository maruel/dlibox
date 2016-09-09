// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// package pins exposes generic pins.
package pins

import (
	"errors"

	"github.com/maruel/dlibox/go/pio/host"
)

// InvalidPinErr is returned when trying to use INVALID.
var InvalidPinErr = errors.New("invalid pin")

// pin implements host.Pin.
type pin struct {
	name string
}

var (
	INVALID     invalidPin
	GROUND      *pin
	V3_3        *pin
	V5          *pin
	DC_IN       *pin
	TEMP_SENSOR *pin
	BAT_PLUS    *pin
	IR_RX       *pin
	EAROUTP     *pin
	EAROUT_N    *pin
	CHARGER_LED *pin
	RESET       *pin
	PWR_SWITCH  *pin
	KEY_ADC     *pin
	X32KFOUT    *pin
	VCC         *pin
	IOVCC       *pin
)

func (p *pin) Number() int {
	return -1
}

func (p *pin) String() string {
	return p.name
}

// invalidPin implements host.PinIO for compability but fails on all access.
type invalidPin struct {
}

func (invalidPin) Number() int {
	return -1
}

func (invalidPin) String() string {
	return "INVALID"
}

func (invalidPin) In(host.Pull) error {
	return InvalidPinErr
}

func (invalidPin) Read() host.Level {
	return host.Low
}

func (invalidPin) Edges() (chan host.Level, error) {
	return nil, InvalidPinErr
}

func (invalidPin) Out() error {
	return InvalidPinErr
}

func (invalidPin) Set(host.Level) {
}

func init() {
	GROUND = &pin{"GROUND"}
	V3_3 = &pin{"V3_3"}
	V5 = &pin{"V5"}
	DC_IN = &pin{"DC_IN"}
	TEMP_SENSOR = &pin{"TEMP_SENSOR"}
	BAT_PLUS = &pin{"BAT_PLUS"}
	IR_RX = &pin{"IR_RX"}
	EAROUTP = &pin{"EAROUTP"}
	EAROUT_N = &pin{"EAROUT_N"}
	CHARGER_LED = &pin{"CHARGER_LED"}
	RESET = &pin{"RESET"}
	PWR_SWITCH = &pin{"PWR_SWITCH"}
	KEY_ADC = &pin{"KEY_ADC"}
	X32KFOUT = &pin{"X32KFOUT"}
	VCC = &pin{"VCC"}
	IOVCC = &pin{"IOVCC"}
}
