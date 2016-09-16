// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Raspberry Pi pin out.

package rpi

import (
	"github.com/maruel/dlibox/go/pio/host"
	"github.com/maruel/dlibox/go/pio/host/internal/bcm283x"
	"github.com/maruel/dlibox/go/pio/host/pins"
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
	P1_1  host.Pin   = pins.V3_3      // 3.3 volt; max 30mA
	P1_2  host.Pin   = pins.V5        // 5 volt (after filtering)
	P1_3  host.PinIO = bcm283x.GPIO2  // High, I2C_SDA1
	P1_4  host.Pin   = pins.V5        //
	P1_5  host.PinIO = bcm283x.GPIO3  // High, I2C_SCL1
	P1_6  host.Pin   = pins.GROUND    //
	P1_7  host.PinIO = bcm283x.GPIO4  // High, GPCLK0
	P1_8  host.PinIO = bcm283x.GPIO14 // Low,  UART_TXD0, UART_TXD1
	P1_9  host.Pin   = pins.GROUND    //
	P1_10 host.PinIO = bcm283x.GPIO15 // Low,  UART_RXD0, UART_RXD1
	P1_11 host.PinIO = bcm283x.GPIO17 // Low,  UART_RTS0, SPI1_CE1, UART_RTS1
	P1_12 host.PinIO = bcm283x.GPIO18 // Low,  PCM_CLK, SPI1_CE0, PWM0_OUT
	P1_13 host.PinIO = bcm283x.GPIO27 // Low,
	P1_14 host.Pin   = pins.GROUND    //
	P1_15 host.PinIO = bcm283x.GPIO22 // Low,
	P1_16 host.PinIO = bcm283x.GPIO23 // Low,
	P1_17 host.Pin   = pins.V3_3      //
	P1_18 host.PinIO = bcm283x.GPIO24 // Low,
	P1_19 host.PinIO = bcm283x.GPIO10 // Low, SPI0_MOSI
	P1_20 host.Pin   = pins.GROUND    //
	P1_21 host.PinIO = bcm283x.GPIO9  // Low, SPI0_MISO
	P1_22 host.PinIO = bcm283x.GPIO25 // Low,
	P1_23 host.PinIO = bcm283x.GPIO11 // Low, SPI0_CLK
	P1_24 host.PinIO = bcm283x.GPIO8  // High, SPI0_CE0
	P1_25 host.Pin   = pins.GROUND    //
	P1_26 host.PinIO = bcm283x.GPIO7  // High, SPI0_CE1

	// Raspberry Pi 2 and later:
	P1_27 host.PinIO = bcm283x.GPIO0  // High, I2C_SDA0 used to probe for HAT EEPROM, see https://github.com/raspberrypi/hats
	P1_28 host.PinIO = bcm283x.GPIO1  // High, I2C_SCL0
	P1_29 host.PinIO = bcm283x.GPIO5  // High, GPCLK1
	P1_30 host.Pin   = pins.GROUND    //
	P1_31 host.PinIO = bcm283x.GPIO6  // High, GPCLK2
	P1_32 host.PinIO = bcm283x.GPIO12 // Low,  PWM0_OUT
	P1_33 host.PinIO = bcm283x.GPIO13 // Low,  PWM1_OUT
	P1_34 host.Pin   = pins.GROUND    //
	P1_35 host.PinIO = bcm283x.GPIO19 // Low,  PCM_FS, SPI1_MISO, PWM1_OUT
	P1_36 host.PinIO = bcm283x.GPIO16 // Low,  UART_CTS0, SPI1_CE2, UART_CTS1
	P1_37 host.PinIO = bcm283x.GPIO26 //
	P1_38 host.PinIO = bcm283x.GPIO20 // Low,  PCM_DIN, SPI1_MOSI, GPCLK0
	P1_39 host.Pin   = pins.GROUND    //
	P1_40 host.PinIO = bcm283x.GPIO21 // Low,  PCM_DOUT, SPI1_CLK, GPCLK1

	// Raspberry Pi 1 header:
	P5_1 host.Pin   = pins.V5
	P5_2 host.Pin   = pins.V3_3
	P5_3 host.PinIO = bcm283x.GPIO28 // Float, I2C_SDA0, PCM_CLK
	P5_4 host.PinIO = bcm283x.GPIO29 // Float, I2C_SCL0, PCM_FS
	P5_5 host.PinIO = bcm283x.GPIO30 // Low,   PCM_DIN, UART_CTS0, UARTS_CTS1
	P5_6 host.PinIO = bcm283x.GPIO31 // Low,   PCM_DOUT, UART_RTS0, UARTS_RTS1
	P5_7 host.Pin   = pins.GROUND
	P5_8 host.Pin   = pins.GROUND

	AUDIO_LEFT          host.Pin = bcm283x.GPIO41 // Low,   PWM1_OUT, SPI2_MOSI, UART_RXD1
	AUDIO_RIGHT         host.Pin = bcm283x.GPIO40 // Low,   PWM0_OUT, SPI2_MISO, UART_TXD1
	HDMI_HOTPLUG_DETECT host.Pin = bcm283x.GPIO46 // High,
)

// IsConnected returns true if the pin is phyisically connected.
func IsConnected(p host.Pin) bool {
	// TODO(maruel): A bit slow, create a lookup table.
	bcm, ok := p.(host.Pin)
	if !ok {
		return false
	}
	switch bcm {
	case P1_1, P1_2, P1_3, P1_4, P1_5, P1_6, P1_7, P1_8, P1_9, P1_10,
		P1_11, P1_12, P1_13, P1_14, P1_15, P1_16, P1_17, P1_18, P1_19, P1_20,
		P1_21, P1_22, P1_23, P1_24, P1_25, P1_26, P1_27, P1_28, P1_29, P1_30,
		P1_31, P1_32, P1_33, P1_34, P1_35, P1_36, P1_37, P1_38, P1_39, P1_40,
		P5_1, P5_2, P5_3, P5_4, P5_5, P5_6, P5_7, P5_8, AUDIO_LEFT, AUDIO_RIGHT, HDMI_HOTPLUG_DETECT:
		return true
	default:
		return false
	}
}
