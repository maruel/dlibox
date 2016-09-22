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
	return allPins
}

// PinsFunctional returns a map of all pins implementing hardware provided
// special functionality, like I²C, SPI, ADC.
func PinsFunctional() map[string]gpio.PinIO {
	Init()
	lock.Lock()
	defer lock.Unlock()
	return pinByFunction
}

// PinByNumber returns a GPIO pin from its number.
//
// Returns nil in case of failure.
func PinByNumber(number int) gpio.PinIO {
	Init()
	pin, _ := pinByNumber[number]
	return pin
}

// PinByName returns a GPIO pin from its name.
//
// This can be strings like GPIO2, PB8, etc.
//
// Returns nil in case of failure.
func PinByName(name string) gpio.PinIO {
	Init()
	lock.Lock()
	defer lock.Unlock()
	if pinByName == nil {
		pinByName = make(map[string]gpio.PinIO, len(allPins))
		for _, pin := range allPins {
			// This assumes there is not 2 pins with the same name and that String()
			// returns the pin name.
			pinByName[pin.String()] = pin
		}
	}
	pin, _ := pinByName[name]
	return pin
}

// ByFunction returns a GPIO pin from its function.
//
// This can be strings like I2C1_SDA, SPI0_MOSI, etc.
//
// Returns nil in case of failure.
func PinByFunction(fn string) gpio.PinIO {
	Init()
	pin, _ := pinByFunction[fn]
	return pin
}

// Init initializes the pins and returns the name of the subsystem used.
func Init() (string, error) {
	lock.Lock()
	defer lock.Unlock()
	if allPins != nil {
		return subsystem, nil
	}

	allPins = []gpio.PinIO{}
	pinByFunction = map[string]gpio.PinIO{}
	if internal.IsBCM283x() {
		if err := bcm283x.Init(); err != nil {
		} else {
			allPins = make([]gpio.PinIO, len(bcm283x.Pins))
			for i := range bcm283x.Pins {
				allPins[i] = &bcm283x.Pins[i]
			}
			sort.Sort(allPins)
			pinByFunction = bcm283x.Functional
			subsystem = "bcm283x"
			setIR()
			return subsystem, nil
		}
	}
	if internal.IsAllwinner() {
		if err := allwinner.Init(); err != nil {
		} else {
			allPins = make([]gpio.PinIO, len(allwinner.Pins))
			for i := range allwinner.Pins {
				allPins[i] = &allwinner.Pins[i]
			}
			sort.Sort(allPins)
			pinByFunction = allwinner.Functional
			subsystem = "allwinner"
			setIR()
			return subsystem, nil
		}
	}

	// Fallback to sysfs gpio.
	if err := sysfs.Init(); err != nil {
		return "", err
	}
	allPins = make([]gpio.PinIO, 0, len(sysfs.Pins))
	for _, p := range sysfs.Pins {
		allPins = append(allPins, p)
	}
	sort.Sort(allPins)
	// sysfs doesn't expose enough information to fill pinByFunction.
	subsystem = "sysfs"
	setIR()
	return subsystem, nil
}

//

var (
	subsystem     string
	allPins       pinList
	pinByNumber   map[int]gpio.PinIO
	pinByName     map[string]gpio.PinIO
	pinByFunction map[string]gpio.PinIO
)

type pinList []gpio.PinIO

func (p pinList) Len() int           { return len(p) }
func (p pinList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p pinList) Less(i, j int) bool { return p[i].Number() < p[j].Number() }

func setIR() {
	pinByNumber = make(map[int]gpio.PinIO, len(allPins))
	for _, pin := range allPins {
		pinByNumber[pin.Number()] = pin
	}
	in, out := lirc.Pins()
	if in != -1 {
		if pin, ok := pinByNumber[in]; ok {
			IR_IN = pin
		}
	}
	if out != -1 {
		if pin, ok := pinByNumber[out]; ok {
			IR_OUT = pin
		}
	}
	pinByFunction["IR_IN"] = IR_IN
	pinByFunction["IR_OUT"] = IR_OUT
}
