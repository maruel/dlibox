// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package host

import (
	"sort"

	"github.com/maruel/dlibox/go/pio/devices/ir/lirc"
	"github.com/maruel/dlibox/go/pio/host/drivers/allwinner"
	"github.com/maruel/dlibox/go/pio/host/drivers/bcm283x"
	"github.com/maruel/dlibox/go/pio/host/drivers/sysfs"
	"github.com/maruel/dlibox/go/pio/host/internal"
	"github.com/maruel/dlibox/go/pio/host/internal/pins2"
	"github.com/maruel/dlibox/go/pio/protocols/gpio"
	"github.com/maruel/dlibox/go/pio/protocols/pins"
)

var (
	IR_IN  gpio.PinIO = pins.INVALID // (any GPIO)
	IR_OUT gpio.PinIO = pins.INVALID // (any GPIO)
)

// Pins returns all the GPIO pins available on this host.
//
// The list is guaranteed to be in order of number.
//
// This list excludes non-GPIO pins like GROUND, V3_3, etc.
func Pins() []gpio.PinIO {
	Init()
	return pins2.All
}

// PinsFunctional returns a map of all pins implementing hardware provided
// special functionality, like I²C, SPI, ADC.
func PinsFunctional() map[string]gpio.PinIO {
	Init()
	lock.Lock()
	defer lock.Unlock()
	return pins2.ByFunction
}

// PinByNumber returns a GPIO pin from its number.
//
// Returns nil in case of failure.
func PinByNumber(number int) gpio.PinIO {
	Init()
	pin, _ := pins2.ByNumber[number]
	return pin
}

// PinByName returns a GPIO pin from its name.
//
// This can be strings like GPIO2, PB8, etc.
//
// Returns nil in case of failure.
func PinByName(name string) gpio.PinIO {
	Init()
	pins2.Lock.Lock()
	defer pins2.Lock.Unlock()
	if pins2.ByName == nil {
		pins2.ByName = make(map[string]gpio.PinIO, len(pins2.All))
		for _, pin := range pins2.All {
			// This assumes there is not 2 pins with the same name and that String()
			// returns the pin name.
			pins2.ByName[pin.String()] = pin
		}
	}
	pin, _ := pins2.ByName[name]
	return pin
}

// ByFunction returns a GPIO pin from its function.
//
// This can be strings like I2C1_SDA, SPI0_MOSI, etc.
//
// Returns nil in case of failure.
func PinByFunction(fn string) gpio.PinIO {
	Init()
	pin, _ := pins2.ByFunction[fn]
	return pin
}

// Init initializes the pins and returns the name of the subsystem used.
func Init() (string, error) {
	pins2.Lock.Lock()
	defer pins2.Lock.Unlock()
	if pins2.All != nil {
		return subsystem, nil
	}

	pins2.All = []gpio.PinIO{}
	pins2.ByFunction = map[string]gpio.PinIO{}
	if internal.IsBCM283x() {
		if err := bcm283x.Init(); err != nil {
		} else {
			pins2.All = make([]gpio.PinIO, len(bcm283x.Pins))
			for i := range bcm283x.Pins {
				pins2.All[i] = &bcm283x.Pins[i]
			}
			sort.Sort(pins2.All)
			pins2.ByFunction = bcm283x.Functional
			subsystem = "bcm283x"
			setIR()
			return subsystem, nil
		}
	}
	if internal.IsAllwinner() {
		if err := allwinner.Init(); err != nil {
		} else {
			pins2.All = make([]gpio.PinIO, len(allwinner.Pins))
			for i := range allwinner.Pins {
				pins2.All[i] = &allwinner.Pins[i]
			}
			sort.Sort(pins2.All)
			pins2.ByFunction = allwinner.Functional
			subsystem = "allwinner"
			setIR()
			return subsystem, nil
		}
	}

	// Fallback to sysfs gpio.
	if err := sysfs.Init(); err != nil {
		return "", err
	}
	pins2.All = make([]gpio.PinIO, 0, len(sysfs.Pins))
	for _, p := range sysfs.Pins {
		pins2.All = append(pins2.All, p)
	}
	sort.Sort(pins2.All)
	// sysfs doesn't expose enough information to fill pins2.ByFunction.
	subsystem = "sysfs"
	setIR()
	return subsystem, nil
}

//

var (
	subsystem string
)

func setIR() {
	pins2.ByNumber = make(map[int]gpio.PinIO, len(pins2.All))
	for _, pin := range pins2.All {
		pins2.ByNumber[pin.Number()] = pin
	}
	in, out := lirc.Pins()
	if in != -1 {
		if pin, ok := pins2.ByNumber[in]; ok {
			IR_IN = pin
		}
	}
	if out != -1 {
		if pin, ok := pins2.ByNumber[out]; ok {
			IR_OUT = pin
		}
	}
	pins2.ByFunction["IR_IN"] = IR_IN
	pins2.ByFunction["IR_OUT"] = IR_OUT
}
