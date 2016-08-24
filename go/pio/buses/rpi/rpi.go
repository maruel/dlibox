// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Raspberry Pi pin out.

package rpi

import (
	"log"
	"strconv"

	"github.com/maruel/dlibox/go/pio/buses/bcm283x"
	"github.com/maruel/dlibox/go/pio/buses/internal"
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
	P1_1  = bcm283x.V3_3   // 3.3 volt; max 30mA
	P1_2  = bcm283x.V5     // 5 volt (after filtering)
	P1_3  = bcm283x.GPIO2  // High, I2C_SDA1
	P1_4  = bcm283x.V5     //
	P1_5  = bcm283x.GPIO3  // High, I2C_SCL1
	P1_6  = bcm283x.GROUND //
	P1_7  = bcm283x.GPIO4  // High, GPCLK0
	P1_8  = bcm283x.GPIO14 // Low,  UART_TXD0, UART_TXD1
	P1_9  = bcm283x.GROUND //
	P1_10 = bcm283x.GPIO15 // Low,  UART_RXD0, UART_RXD1
	P1_11 = bcm283x.GPIO17 // Low,  UART_RTS0, SPI1_CE1, UART_RTS1
	P1_12 = bcm283x.GPIO18 // Low,  PCM_CLK, SPI1_CE0, PWM0_OUT
	P1_13 = bcm283x.GPIO27 // Low,
	P1_14 = bcm283x.GROUND //
	P1_15 = bcm283x.GPIO22 // Low,
	P1_16 = bcm283x.GPIO23 // Low,
	P1_17 = bcm283x.V3_3   //
	P1_18 = bcm283x.GPIO24 // Low,
	P1_19 = bcm283x.GPIO10 // Low, SPI0_MOSI
	P1_20 = bcm283x.GROUND //
	P1_21 = bcm283x.GPIO9  // Low, SPI0_MISO
	P1_22 = bcm283x.GPIO25 // Low,
	P1_23 = bcm283x.GPIO11 // Low, SPI0_CLK
	P1_24 = bcm283x.GPIO8  // High, SPI0_CE0
	P1_25 = bcm283x.GROUND //
	P1_26 = bcm283x.GPIO7  // High, SPI0_CE1

	// Raspberry Pi 2 and later:
	P1_27 = bcm283x.GPIO0  // High, I2C_SDA0 used to probe for HAT EEPROM, see https://github.com/raspberrypi/hats
	P1_28 = bcm283x.GPIO1  // High, I2C_SCL0
	P1_29 = bcm283x.GPIO5  // High, GPCLK1
	P1_30 = bcm283x.GROUND //
	P1_31 = bcm283x.GPIO6  // High, GPCLK2
	P1_32 = bcm283x.GPIO12 // Low,  PWM0_OUT
	P1_33 = bcm283x.GPIO13 // Low,  PWM1_OUT
	P1_34 = bcm283x.GROUND //
	P1_35 = bcm283x.GPIO19 // Low,  PCM_FS, SPI1_MISO, PWM1_OUT
	P1_36 = bcm283x.GPIO16 // Low,  UART_CTS0, SPI1_CE2, UART_CTS1
	P1_37 = bcm283x.GPIO26 //
	P1_38 = bcm283x.GPIO20 // Low,  PCM_DIN, SPI1_MOSI, GPCLK0
	P1_39 = bcm283x.GROUND //
	P1_40 = bcm283x.GPIO21 // Low,  PCM_DOUT, SPI1_CLK, GPCLK1

	// Raspberry Pi 1 header:
	P5_1 = bcm283x.V5
	P5_2 = bcm283x.V3_3
	P5_3 = bcm283x.GPIO28 // Float, I2C_SDA0, PCM_CLK
	P5_4 = bcm283x.GPIO29 // Float, I2C_SCL0, PCM_FS
	P5_5 = bcm283x.GPIO30 // Low,   PCM_DIN, UART_CTS0, UARTS_CTS1
	P5_6 = bcm283x.GPIO31 // Low,   PCM_DOUT, UART_RTS0, UARTS_RTS1
	P5_7 = bcm283x.GROUND
	P5_8 = bcm283x.GROUND

	AUDIO_LEFT          = bcm283x.GPIO41 // Low,   PWM1_OUT, SPI2_MOSI, UART_RXD1
	AUDIO_RIGHT         = bcm283x.GPIO40 // Low,   PWM0_OUT, SPI2_MISO, UART_TXD1
	HDMI_HOTPLUG_DETECT = bcm283x.GPIO46 // High,
)

// IsConnected returns true if the pin is phyisically connected.
func IsConnected(p bcm283x.Pin) bool {
	// TODO(maruel): A bit slow, create a lookup table.
	switch p {
	case bcm283x.INVALID:
		return false
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

//

func init() {
	// Initialize Version. This function is not futureproof, it will return 0 on
	// a Raspberry Pi 4 whenever it comes out.
	if i, err := strconv.ParseInt(internal.CPUInfo["Revision"], 16, 32); err == nil {
		// Ignore the overclock bit.
		i &= 0xFFFFFF
		if i < 0x20 {
			Version = 1
		} else if i == 0xa01041 || i == 0xa21041 {
			Version = 2
		} else if i == 0xa02082 || i == 0xa22082 {
			Version = 3
		} else {
			log.Printf("Unknown hardware version: 0x%x", i)
		}
	} else {
		log.Printf("Failed to read cpu_info: %v", err)
	}

	if Version == 1 {
		// TODO(maruel): Models from 2012 and earlier have P1_3=GPIO0, P1_5=GPIO1 and P1_13=GPIO21.
		// P2 and P3 are not useful.
		// P6 has a RUN pin for reset but it's not available after Pi version 1.

		P1_27 = bcm283x.INVALID
		P1_28 = bcm283x.INVALID
		P1_29 = bcm283x.INVALID
		P1_30 = bcm283x.INVALID
		P1_31 = bcm283x.INVALID
		P1_32 = bcm283x.INVALID
		P1_33 = bcm283x.INVALID
		P1_34 = bcm283x.INVALID
		P1_35 = bcm283x.INVALID
		P1_36 = bcm283x.INVALID
		P1_37 = bcm283x.INVALID
		P1_38 = bcm283x.INVALID
		P1_39 = bcm283x.INVALID
		P1_40 = bcm283x.INVALID
	} else {
		P5_1 = bcm283x.INVALID
		P5_2 = bcm283x.INVALID
		P5_3 = bcm283x.INVALID
		P5_4 = bcm283x.INVALID
		P5_5 = bcm283x.INVALID
		P5_6 = bcm283x.INVALID
		P5_7 = bcm283x.INVALID
		P5_8 = bcm283x.INVALID
	}
	if Version < 3 {
		AUDIO_LEFT = bcm283x.GPIO45
	}
}
