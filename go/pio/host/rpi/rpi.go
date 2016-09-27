// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Raspberry Pi pin out.

package rpi

import (
	"fmt"
	"strconv"

	"github.com/maruel/dlibox/go/pio/drivers"
	"github.com/maruel/dlibox/go/pio/host/bcm283x"
	"github.com/maruel/dlibox/go/pio/host/headers"
	"github.com/maruel/dlibox/go/pio/host/internal"
	"github.com/maruel/dlibox/go/pio/protocols/gpio"
	"github.com/maruel/dlibox/go/pio/protocols/pins"
)

// Version is the Raspberry Pi version 1, 2 or 3.
//
// Is set to 0 when detection (currently primitive) failed.
var Version int

// Pin as connect on the 40 pins extention header.
//
// Schematics are useful to know what is connected to what:
// https://www.raspberrypi.org/documentation/hardware/raspberrypi/schematics/README.md
//
// The actual pin mapping depends on the board revision! The default values are
// set as the 40 pins header on Raspberry Pi 2 and Raspberry Pi 3.
//
// P1 is also known as J8.
var (
	P1_1  gpio.PinIO = pins.V3_3      // 3.3 volt; max 30mA
	P1_2  gpio.PinIO = pins.V5        // 5 volt (after filtering)
	P1_3  gpio.PinIO = bcm283x.GPIO2  // High, I2C1_SDA
	P1_4  gpio.PinIO = pins.V5        //
	P1_5  gpio.PinIO = bcm283x.GPIO3  // High, I2C1_SCL
	P1_6  gpio.PinIO = pins.GROUND    //
	P1_7  gpio.PinIO = bcm283x.GPIO4  // High, GPCLK0
	P1_8  gpio.PinIO = bcm283x.GPIO14 // Low,  UART0_TXD, UART1_TXD
	P1_9  gpio.PinIO = pins.GROUND    //
	P1_10 gpio.PinIO = bcm283x.GPIO15 // Low,  UART0_RXD, UART1_RXD
	P1_11 gpio.PinIO = bcm283x.GPIO17 // Low,  UART0_RTS, SPI1_CE1, UART1_RTS
	P1_12 gpio.PinIO = bcm283x.GPIO18 // Low,  PCM_CLK, SPI1_CE0, PWM0_OUT
	P1_13 gpio.PinIO = bcm283x.GPIO27 // Low,
	P1_14 gpio.PinIO = pins.GROUND    //
	P1_15 gpio.PinIO = bcm283x.GPIO22 // Low,
	P1_16 gpio.PinIO = bcm283x.GPIO23 // Low,
	P1_17 gpio.PinIO = pins.V3_3      //
	P1_18 gpio.PinIO = bcm283x.GPIO24 // Low,
	P1_19 gpio.PinIO = bcm283x.GPIO10 // Low, SPI0_MOSI
	P1_20 gpio.PinIO = pins.GROUND    //
	P1_21 gpio.PinIO = bcm283x.GPIO9  // Low, SPI0_MISO
	P1_22 gpio.PinIO = bcm283x.GPIO25 // Low,
	P1_23 gpio.PinIO = bcm283x.GPIO11 // Low, SPI0_CLK
	P1_24 gpio.PinIO = bcm283x.GPIO8  // High, SPI0_CE0
	P1_25 gpio.PinIO = pins.GROUND    //
	P1_26 gpio.PinIO = bcm283x.GPIO7  // High, SPI0_CE1

	// Raspberry Pi 2 and later:
	P1_27 gpio.PinIO = bcm283x.GPIO0  // High, I2C0_SDA used to probe for HAT EEPROM, see https://github.com/raspberrypi/hats
	P1_28 gpio.PinIO = bcm283x.GPIO1  // High, I2C0_SCL
	P1_29 gpio.PinIO = bcm283x.GPIO5  // High, GPCLK1
	P1_30 gpio.PinIO = pins.GROUND    //
	P1_31 gpio.PinIO = bcm283x.GPIO6  // High, GPCLK2
	P1_32 gpio.PinIO = bcm283x.GPIO12 // Low,  PWM0_OUT
	P1_33 gpio.PinIO = bcm283x.GPIO13 // Low,  PWM1_OUT
	P1_34 gpio.PinIO = pins.GROUND    //
	P1_35 gpio.PinIO = bcm283x.GPIO19 // Low,  PCM_FS, SPI1_MISO, PWM1_OUT
	P1_36 gpio.PinIO = bcm283x.GPIO16 // Low,  UART0_CTS, SPI1_CE2, UART1_CTS
	P1_37 gpio.PinIO = bcm283x.GPIO26 //
	P1_38 gpio.PinIO = bcm283x.GPIO20 // Low,  PCM_DIN, SPI1_MOSI, GPCLK0
	P1_39 gpio.PinIO = pins.GROUND    //
	P1_40 gpio.PinIO = bcm283x.GPIO21 // Low,  PCM_DOUT, SPI1_CLK, GPCLK1

	// Raspberry Pi 1 header:
	P5_1 gpio.PinIO = pins.V5
	P5_2 gpio.PinIO = pins.V3_3
	P5_3 gpio.PinIO = bcm283x.GPIO28 // Float, I2C0_SDA, PCM_CLK
	P5_4 gpio.PinIO = bcm283x.GPIO29 // Float, I2C0_SCL, PCM_FS
	P5_5 gpio.PinIO = bcm283x.GPIO30 // Low,   PCM_DIN, UART0_CTS, UART1_CTS
	P5_6 gpio.PinIO = bcm283x.GPIO31 // Low,   PCM_DOUT, UART0_RTS, UART1_RTS
	P5_7 gpio.PinIO = pins.GROUND
	P5_8 gpio.PinIO = pins.GROUND

	AUDIO_LEFT          gpio.PinIO = bcm283x.GPIO41 // Low,   PWM1_OUT, SPI2_MOSI, UART1_RXD
	AUDIO_RIGHT         gpio.PinIO = bcm283x.GPIO40 // Low,   PWM0_OUT, SPI2_MISO, UART1_TXD
	HDMI_HOTPLUG_DETECT gpio.PinIO = bcm283x.GPIO46 // High,
)

//

func zapPins() {
	P1_1 = pins.INVALID
	P1_2 = pins.INVALID
	P1_3 = pins.INVALID
	P1_4 = pins.INVALID
	P1_5 = pins.INVALID
	P1_6 = pins.INVALID
	P1_7 = pins.INVALID
	P1_8 = pins.INVALID
	P1_9 = pins.INVALID
	P1_10 = pins.INVALID
	P1_11 = pins.INVALID
	P1_12 = pins.INVALID
	P1_13 = pins.INVALID
	P1_14 = pins.INVALID
	P1_15 = pins.INVALID
	P1_16 = pins.INVALID
	P1_17 = pins.INVALID
	P1_18 = pins.INVALID
	P1_19 = pins.INVALID
	P1_20 = pins.INVALID
	P1_21 = pins.INVALID
	P1_22 = pins.INVALID
	P1_23 = pins.INVALID
	P1_24 = pins.INVALID
	P1_25 = pins.INVALID
	P1_26 = pins.INVALID
	P1_27 = pins.INVALID
	P1_28 = pins.INVALID
	P1_29 = pins.INVALID
	P1_30 = pins.INVALID
	P1_31 = pins.INVALID
	P1_32 = pins.INVALID
	P1_33 = pins.INVALID
	P1_34 = pins.INVALID
	P1_35 = pins.INVALID
	P1_36 = pins.INVALID
	P1_37 = pins.INVALID
	P1_38 = pins.INVALID
	P1_39 = pins.INVALID
	P1_40 = pins.INVALID
	P5_1 = pins.INVALID
	P5_2 = pins.INVALID
	P5_3 = pins.INVALID
	P5_4 = pins.INVALID
	P5_5 = pins.INVALID
	P5_6 = pins.INVALID
	P5_7 = pins.INVALID
	P5_8 = pins.INVALID
	AUDIO_LEFT = pins.INVALID
	AUDIO_RIGHT = pins.INVALID
	HDMI_HOTPLUG_DETECT = pins.INVALID
}

// driver implements drivers.Driver.
type driver struct {
}

func (d *driver) String() string {
	return "rpi"
}

func (d *driver) Type() drivers.Type {
	return drivers.Pins
}

func (d *driver) Prerequisites() []string {
	return []string{"bcm283x"}
}

func (d *driver) Init() (bool, error) {
	if !internal.IsRaspberryPi() {
		zapPins()
		return false, nil
	}

	// Initialize Version.
	//
	// This code is not futureproof, it will error out on a Raspberry Pi 4
	// whenever it comes out.
	rev, _ := internal.CPUInfo()["Revision"]
	if i, err := strconv.ParseInt(rev, 16, 32); err == nil {
		// Ignore the overclock bit.
		i &= 0xFFFFFF
		if i < 0x20 {
			Version = 1
		} else if i == 0xa01041 || i == 0xa21041 {
			Version = 2
		} else if i == 0xa02082 || i == 0xa22082 {
			Version = 3
		} else {
			return true, fmt.Errorf("rpi: unknown hardware version: 0x%x", i)
		}
	} else {
		return true, fmt.Errorf("rpi: failed to read cpu_info: %v", err)
	}

	if Version == 1 {
		headers.Register("P1", [][]pins.Pin{
			{P1_1, P1_2},
			{P1_3, P1_4},
			{P1_5, P1_6},
			{P1_7, P1_8},
			{P1_9, P1_10},
			{P1_11, P1_12},
			{P1_13, P1_14},
			{P1_15, P1_16},
			{P1_17, P1_18},
			{P1_19, P1_20},
			{P1_21, P1_22},
			{P1_23, P1_24},
			{P1_25, P1_26},
		})
		headers.Register("P5", [][]pins.Pin{
			{P5_1, P5_2},
			{P5_3, P5_4},
			{P5_5, P5_6},
			{P5_7, P5_8},
		})

		// TODO(maruel): Models from 2012 and earlier have P1_3=GPIO0, P1_5=GPIO1 and P1_13=GPIO21.
		// P2 and P3 are not useful.
		// P6 has a RUN pin for reset but it's not available after Pi version 1.
		P1_27 = pins.INVALID
		P1_28 = pins.INVALID
		P1_29 = pins.INVALID
		P1_30 = pins.INVALID
		P1_31 = pins.INVALID
		P1_32 = pins.INVALID
		P1_33 = pins.INVALID
		P1_34 = pins.INVALID
		P1_35 = pins.INVALID
		P1_36 = pins.INVALID
		P1_37 = pins.INVALID
		P1_38 = pins.INVALID
		P1_39 = pins.INVALID
		P1_40 = pins.INVALID
	} else {
		headers.Register("P1", [][]pins.Pin{
			{P1_1, P1_2},
			{P1_3, P1_4},
			{P1_5, P1_6},
			{P1_7, P1_8},
			{P1_9, P1_10},
			{P1_11, P1_12},
			{P1_13, P1_14},
			{P1_15, P1_16},
			{P1_17, P1_18},
			{P1_19, P1_20},
			{P1_21, P1_22},
			{P1_23, P1_24},
			{P1_25, P1_26},
			{P1_27, P1_28},
			{P1_29, P1_30},
			{P1_31, P1_32},
			{P1_33, P1_34},
			{P1_35, P1_36},
			{P1_37, P1_38},
			{P1_39, P1_40},
		})
		P5_1 = pins.INVALID
		P5_2 = pins.INVALID
		P5_3 = pins.INVALID
		P5_4 = pins.INVALID
		P5_5 = pins.INVALID
		P5_6 = pins.INVALID
		P5_7 = pins.INVALID
		P5_8 = pins.INVALID
	}
	if Version < 3 {
		AUDIO_LEFT = bcm283x.GPIO45
	}
	headers.Register("AUDIO", [][]pins.Pin{
		{AUDIO_LEFT},
		{AUDIO_RIGHT},
	})
	headers.Register("HDMI", [][]pins.Pin{{HDMI_HOTPLUG_DETECT}})
	return true, nil
}

var _ drivers.Driver = &driver{}
