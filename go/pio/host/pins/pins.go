// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// package pins exposes generic pins.
package pins

import (
	"strings"

	"github.com/maruel/dlibox/go/pio/host"
	"github.com/maruel/dlibox/go/pio/host/internal"
	"github.com/maruel/dlibox/go/pio/host/internal/bcm283x"
	"github.com/maruel/dlibox/go/pio/host/internal/sysfs"
)

// pin implements host.Pin.
type pin struct {
	name string
}

var (
	GROUND      host.Pin
	V3_3        host.Pin
	V5          host.Pin
	DC_IN       host.Pin
	TEMP_SENSOR host.Pin
	BAT_PLUS    host.Pin
	IR_RX       host.Pin
	EAROUTP     host.Pin
	EAROUT_N    host.Pin
	CHARGER_LED host.Pin
	RESET       host.Pin
	PWR_SWITCH  host.Pin
	KEY_ADC     host.Pin
	X32KFOUT    host.Pin
	VCC         host.Pin
	IOVCC       host.Pin
)

func (p *pin) Number() int {
	return -1
}

func (p *pin) String() string {
	return p.name
}

func (p *pin) Function() string {
	return p.name
}

// All refers to all the GPIO pins available on this host.
//
// This gets populated by Init().
//
// This list excludes non-GPIO pins like GROUND, V3_3, etc.
var All map[int]host.PinIO

// Functional lists all pins implementing hardware provided special
// functionality, like I²C, SPI, ADC.
var Functional map[string]host.Pin

// ByName returns a GPIO pin from its name.
//
// This excludes non-GPIO pins like GROUND, V3_3, etc.
//
// Returns nil in case of failure.
//
// TODO(maruel): Remove?
func ByName(name string) host.PinIO {
	// TODO(maruel): Create a map on first use?
	if Init() != nil {
		return nil
	}
	for _, p := range All {
		if p.String() == name {
			return p
		}
	}
	return nil
}

// ByNumber returns a GPIO pin from its number.
//
// This excludes non-GPIO pins like GROUND, V3_3, etc.
//
// Returns nil in case of failure.
func ByNumber(number int) host.PinIO {
	if Init() != nil {
		return nil
	}
	p, _ := All[number]
	return p
}

func Init() error {
	// TODO(maruel): concurrency safety.
	if All != nil {
		return nil
	}
	// Try to detect the CPU and use the right runtime module.
	hardware, ok := internal.CPUInfo["Hardware"]
	if ok && strings.HasPrefix(hardware, "BCM") {
		if err := bcm283x.Init(); err != nil {
			// TODO(maruel): Fallback to sysfs.
			return err
		}
		All = make(map[int]host.PinIO, len(bcm283x.Pins))
		for i := range bcm283x.Pins {
			All[i] = &bcm283x.Pins[i]
		}
		Functional = bcm283x.Functional
		return nil
	}
	if ok && strings.HasPrefix(hardware, "sun") {
		/* TODO(maruel): Implement
		if err := a64.Init(); err != nil {
			// TODO(maruel): Fallback to sysfs.
			return err
		}
		All = make([]host.PinIO, len(a64.Pins))
		for i := range a64.Pins {
			All[i] = &a64.Pins[i]
		}
		Functional = bcm283x.Functional
		return nil
		*/
	}

	// Fallback to sysfs gpio.
	if err := sysfs.Init(); err != nil {
		return err
	}
	All = make(map[int]host.PinIO, len(sysfs.Pins))
	for id, p := range sysfs.Pins {
		All[id] = p
	}
	// Functional cannot be populated.
	return nil
}

//

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
