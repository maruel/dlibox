// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// package pins exposes generic pins.
package pins

import (
	"sync"

	"github.com/maruel/dlibox/go/pio/host"
	"github.com/maruel/dlibox/go/pio/host/internal"
	"github.com/maruel/dlibox/go/pio/host/internal/allwinner"
	"github.com/maruel/dlibox/go/pio/host/internal/bcm283x"
	"github.com/maruel/dlibox/go/pio/host/internal/sysfs"
)

// pin implements host.Pin.
type pin struct {
	name string
}

var (
	GROUND      host.Pin = &pin{"GROUND"}
	V3_3        host.Pin = &pin{"V3_3"}
	V5          host.Pin = &pin{"V5"}
	DC_IN       host.Pin = &pin{"DC_IN"}
	TEMP_SENSOR host.Pin = &pin{"TEMP_SENSOR"}
	BAT_PLUS    host.Pin = &pin{"BAT_PLUS"}
	IR_RX       host.Pin = &pin{"IR_RX"}
	EAROUTP     host.Pin = &pin{"EAROUTP"}
	EAROUT_N    host.Pin = &pin{"EAROUT_N"}
	CHARGER_LED host.Pin = &pin{"CHARGER_LED"}
	RESET       host.Pin = &pin{"RESET"}
	PWR_SWITCH  host.Pin = &pin{"PWR_SWITCH"}
	KEY_ADC     host.Pin = &pin{"KEY_ADC"}
	X32KFOUT    host.Pin = &pin{"X32KFOUT"}
	VCC         host.Pin = &pin{"VCC"}
	IOVCC       host.Pin = &pin{"IOVCC"}
)

func (p *pin) Number() int {
	return -1
}

func (p *pin) String() string {
	return p.name
}

func (p *pin) Function() string {
	return ""
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

// ByNumber returns a GPIO pin from its number.
//
// This excludes non-GPIO pins like GROUND, V3_3, etc.
//
// Returns nil in case of failure.
func ByNumber(number int) host.PinIO {
	if Init(true) != nil {
		return nil
	}
	p, _ := All[number]
	return p
}

func Init(fallback bool) error {
	lock.Lock()
	defer lock.Unlock()
	if All != nil {
		return nil
	}
	if internal.IsBCM283x() {
		if err := bcm283x.Init(); err != nil {
			if !fallback {
				return err
			}
		} else {
			All = make(map[int]host.PinIO, len(bcm283x.Pins))
			for i := range bcm283x.Pins {
				All[i] = &bcm283x.Pins[i]
			}
			Functional = bcm283x.Functional
			return nil
		}
	}
	if internal.IsAllWinner() {
		if err := allwinner.Init(); err != nil {
			if !fallback {
				return err
			}
		} else {
			All = make(map[int]host.PinIO, len(allwinner.Pins))
			for i := range allwinner.Pins {
				All[i] = &allwinner.Pins[i]
			}
			Functional = allwinner.Functional
			return nil
		}
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

var lock sync.Mutex
