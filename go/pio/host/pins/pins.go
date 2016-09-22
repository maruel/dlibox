// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Package pins exposes generic pins.
package pins

import (
	"fmt"
	"sort"
	"sync"

	"github.com/maruel/dlibox/go/pio/devices/ir/lirc"
	"github.com/maruel/dlibox/go/pio/host/drivers/allwinner"
	"github.com/maruel/dlibox/go/pio/host/drivers/bcm283x"
	"github.com/maruel/dlibox/go/pio/host/drivers/sysfs"
	"github.com/maruel/dlibox/go/pio/host/internal"
	"github.com/maruel/dlibox/go/pio/protocols/gpio"
)

var (
	GROUND      gpio.PinIO = &pin{"GROUND"}
	V3_3        gpio.PinIO = &pin{"V3_3"}
	V5          gpio.PinIO = &pin{"V5"}
	DC_IN       gpio.PinIO = &pin{"DC_IN"}
	TEMP_SENSOR gpio.PinIO = &pin{"TEMP_SENSOR"}
	BAT_PLUS    gpio.PinIO = &pin{"BAT_PLUS"}
	IR_RX       gpio.PinIO = &pin{"IR_RX"}
	EAROUTP     gpio.PinIO = &pin{"EAROUTP"}
	EAROUT_N    gpio.PinIO = &pin{"EAROUT_N"}
	CHARGER_LED gpio.PinIO = &pin{"CHARGER_LED"}
	RESET       gpio.PinIO = &pin{"RESET"}
	PWR_SWITCH  gpio.PinIO = &pin{"PWR_SWITCH"}
	KEY_ADC     gpio.PinIO = &pin{"KEY_ADC"}
	X32KFOUT    gpio.PinIO = &pin{"X32KFOUT"}
	VCC         gpio.PinIO = &pin{"VCC"}
	IOVCC       gpio.PinIO = &pin{"IOVCC"}
	IR_IN       gpio.PinIO = gpio.INVALID // (any GPIO)
	IR_OUT      gpio.PinIO = gpio.INVALID // (any GPIO)
)

// All returns all the GPIO pins available on this host.
//
// The list is guaranteed to be in order of number.
//
// This list excludes non-GPIO pins like GROUND, V3_3, etc.
func All() []gpio.PinIO {
	Init(true)
	return all
}

// Functional returns a map of all pins implementing hardware provided special
// functionality, like I²C, SPI, ADC.
func Functional() map[string]gpio.PinIO {
	Init(true)
	lock.Lock()
	defer lock.Unlock()
	return byFunction
}

// ByNumber returns a GPIO pin from its number.
//
// Returns nil in case of failure.
func ByNumber(number int) gpio.PinIO {
	Init(true)
	pin, _ := byNumber[number]
	return pin
}

// ByName returns a GPIO pin from its name.
//
// This can be strings like GPIO2, PB8, etc.
//
// Returns nil in case of failure.
func ByName(name string) gpio.PinIO {
	Init(true)
	lock.Lock()
	defer lock.Unlock()
	if byName == nil {
		byName = make(map[string]gpio.PinIO, len(all))
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
func ByFunction(fn string) gpio.PinIO {
	Init(true)
	pin, _ := byFunction[fn]
	return pin
}

// Init initializes the pins and returns the name of the subsystem used.
func Init(fallback bool) (string, error) {
	lock.Lock()
	defer lock.Unlock()
	if all != nil {
		return subsystem, nil
	}

	all = []gpio.PinIO{}
	byFunction = map[string]gpio.PinIO{}
	if internal.IsBCM283x() {
		if err := bcm283x.Init(); err != nil {
			if !fallback {
				return "", err
			}
		} else {
			all = make([]gpio.PinIO, len(bcm283x.Pins))
			for i := range bcm283x.Pins {
				all[i] = &bcm283x.Pins[i]
			}
			sort.Sort(all)
			byFunction = bcm283x.Functional
			subsystem = "bcm283x"
			setIR()
			return subsystem, nil
		}
	}
	if internal.IsAllwinner() {
		if err := allwinner.Init(); err != nil {
			if !fallback {
				return "", err
			}
		} else {
			all = make([]gpio.PinIO, len(allwinner.Pins))
			for i := range allwinner.Pins {
				all[i] = &allwinner.Pins[i]
			}
			sort.Sort(all)
			byFunction = allwinner.Functional
			subsystem = "allwinner"
			setIR()
			return subsystem, nil
		}
	}

	// Fallback to sysfs gpio.
	if err := sysfs.Init(); err != nil {
		return "", err
	}
	all = make([]gpio.PinIO, 0, len(sysfs.Pins))
	for _, p := range sysfs.Pins {
		all = append(all, p)
	}
	sort.Sort(all)
	// sysfs doesn't expose enough information to fill byFunction.
	subsystem = "sysfs"
	setIR()
	return subsystem, nil
}

//

var (
	lock       sync.Mutex
	subsystem  string
	all        pins
	byNumber   map[int]gpio.PinIO
	byName     map[string]gpio.PinIO
	byFunction map[string]gpio.PinIO
)

type pins []gpio.PinIO

func (p pins) Len() int           { return len(p) }
func (p pins) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p pins) Less(i, j int) bool { return p[i].Number() < p[j].Number() }

// pin implements gpio.PinIO.
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

func (p *pin) In(gpio.Pull) error {
	return fmt.Errorf("%s cannot be used as input", p.name)
}

func (p *pin) Read() gpio.Level {
	return gpio.Low
}

func (p *pin) Edges() (<-chan gpio.Level, error) {
	return nil, fmt.Errorf("%s cannot be used as input", p.name)
}

func (p *pin) DisableEdges() {
}

func (p *pin) Pull() gpio.Pull {
	return gpio.PullNoChange
}

func (p *pin) Out(gpio.Level) error {
	return fmt.Errorf("%s cannot be used as output", p.name)
}

func setIR() {
	byNumber = make(map[int]gpio.PinIO, len(all))
	for _, pin := range all {
		byNumber[pin.Number()] = pin
	}
	in, out := lirc.Pins()
	if in != -1 {
		if pin, ok := byNumber[in]; ok {
			IR_IN = pin
		}
	}
	if out != -1 {
		if pin, ok := byNumber[out]; ok {
			IR_OUT = pin
		}
	}
	byFunction["IR_IN"] = IR_IN
	byFunction["IR_OUT"] = IR_OUT
}
