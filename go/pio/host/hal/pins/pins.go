// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// package pins exposes generic pins.
package pins

import (
	"fmt"
	"sort"
	"sync"

	"github.com/maruel/dlibox/go/pio/host"
	"github.com/maruel/dlibox/go/pio/host/drivers/allwinner"
	"github.com/maruel/dlibox/go/pio/host/drivers/bcm283x"
	"github.com/maruel/dlibox/go/pio/host/drivers/sysfs"
	"github.com/maruel/dlibox/go/pio/host/internal"
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

// All returns all the GPIO pins available on this host.
//
// The list is guaranteed to be in order of number.
//
// This list excludes non-GPIO pins like GROUND, V3_3, etc.
func All() []host.PinIO {
	Init(true)
	return all
}

// Functional returns a map of all pins implementing hardware provided special
// functionality, like I²C, SPI, ADC.
func Functional() map[string]host.PinIO {
	Init(true)
	lock.Lock()
	defer lock.Unlock()
	return byFunction
}

// ByNumber returns a GPIO pin from its number.
//
// Returns nil in case of failure.
func ByNumber(number int) host.PinIO {
	Init(true)
	lock.Lock()
	defer lock.Unlock()
	if byNumber == nil {
		byNumber = make(map[int]host.PinIO, len(all))
		for _, pin := range all {
			byNumber[pin.Number()] = pin
		}
	}
	pin, _ := byNumber[number]
	return pin
}

// ByName returns a GPIO pin from its name.
//
// This can be strings like GPIO2, PB8, etc.
//
// Returns nil in case of failure.
func ByName(name string) host.PinIO {
	Init(true)
	lock.Lock()
	defer lock.Unlock()
	if byName == nil {
		byName = make(map[string]host.PinIO, len(all))
		for _, pin := range all {
			// This assumes there is not 2 pins with the same name and that String()
			// returns the pin name.
			byName[pin.String()] = pin
		}
	}
	pin, _ := byName[name]
	return pin
}

// ByFunction returns a GPIO pin from its function.
//
// This can be strings like I2C1_SDA, SPI0_MOSI, etc.
//
// Returns nil in case of failure.
func ByFunction(fn string) host.PinIO {
	Init(true)
	pin, _ := byFunction[fn]
	return pin
}

func Init(fallback bool) error {
	lock.Lock()
	defer lock.Unlock()
	if all != nil {
		return nil
	}
	all = []host.PinIO{}
	byFunction = map[string]host.PinIO{}
	if internal.IsBCM283x() {
		if err := bcm283x.Init(); err != nil {
			if !fallback {
				return err
			}
		} else {
			all = make([]host.PinIO, len(bcm283x.Pins))
			for i := range bcm283x.Pins {
				all[i] = &bcm283x.Pins[i]
			}
			sort.Sort(all)
			byFunction = bcm283x.Functional
			return nil
		}
	}
	if internal.IsAllwinner() {
		if err := allwinner.Init(); err != nil {
			if !fallback {
				return err
			}
		} else {
			all = make([]host.PinIO, len(allwinner.Pins))
			for i := range allwinner.Pins {
				all[i] = &allwinner.Pins[i]
			}
			sort.Sort(all)
			byFunction = allwinner.Functional
			return nil
		}
	}

	// Fallback to sysfs gpio.
	if err := sysfs.Init(); err != nil {
		return err
	}
	all = make([]host.PinIO, 0, len(sysfs.Pins))
	for _, p := range sysfs.Pins {
		all = append(all, p)
	}
	sort.Sort(all)
	// sysfs doesn't expose enough information to fill byFunction.
	return nil
}

//

var (
	lock       sync.Mutex
	all        pins
	byNumber   map[int]host.PinIO
	byName     map[string]host.PinIO
	byFunction map[string]host.PinIO
)

type pins []host.PinIO

func (p pins) Len() int           { return len(p) }
func (p pins) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p pins) Less(i, j int) bool { return p[i].Number() < p[j].Number() }

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
