// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// package pins exposes generic pins.
package pins

import (
	"fmt"
	"sync"

	"github.com/maruel/dlibox/go/pio/host"
	"github.com/maruel/dlibox/go/pio/host/internal"
	"github.com/maruel/dlibox/go/pio/host/internal/allwinner"
	"github.com/maruel/dlibox/go/pio/host/internal/bcm283x"
	"github.com/maruel/dlibox/go/pio/host/internal/sysfs"
)

var (
	GROUND      host.PinIO = &pin{"GROUND"}
	V3_3        host.PinIO = &pin{"V3_3"}
	V5          host.PinIO = &pin{"V5"}
	DC_IN       host.PinIO = &pin{"DC_IN"}
	TEMP_SENSOR host.PinIO = &pin{"TEMP_SENSOR"}
	BAT_PLUS    host.PinIO = &pin{"BAT_PLUS"}
	IR_RX       host.PinIO = &pin{"IR_RX"}
	EAROUTP     host.PinIO = &pin{"EAROUTP"}
	EAROUT_N    host.PinIO = &pin{"EAROUT_N"}
	CHARGER_LED host.PinIO = &pin{"CHARGER_LED"}
	RESET       host.PinIO = &pin{"RESET"}
	PWR_SWITCH  host.PinIO = &pin{"PWR_SWITCH"}
	KEY_ADC     host.PinIO = &pin{"KEY_ADC"}
	X32KFOUT    host.PinIO = &pin{"X32KFOUT"}
	VCC         host.PinIO = &pin{"VCC"}
	IOVCC       host.PinIO = &pin{"IOVCC"}
)

// All refers to all the GPIO pins available on this host.
//
// This gets populated by Init().
//
// This list excludes non-GPIO pins like GROUND, V3_3, etc.
var All map[int]host.PinIO

// Functional lists all pins implementing hardware provided special
// functionality, like I²C, SPI, ADC.
var Functional map[string]host.PinIO

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
	if internal.IsAllwinner() {
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

// pin implements host.PinIO.
type pin struct {
	name string
}

func (p *pin) Number() int {
	return -1
}

func (p *pin) String() string {
	return p.name
}

func (p *pin) Function() string {
	return ""
}

func (p *pin) In(host.Pull) error {
	return fmt.Errorf("%s cannot be used as input", p.name)
}

func (p *pin) Read() host.Level {
	return host.Low
}

func (p *pin) Edges() (<-chan host.Level, error) {
	return nil, fmt.Errorf("%s cannot be used as input", p.name)
}

func (p *pin) DisableEdges() {
}

func (p *pin) Pull() host.Pull {
	return host.PullNoChange
}

func (p *pin) Out() error {
	return fmt.Errorf("%s cannot be used as output", p.name)
}

func (p *pin) Set(host.Level) {
}
