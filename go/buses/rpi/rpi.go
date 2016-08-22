// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Raspberry Pi pin out.

package rpi

import (
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"
)

// MaxSpeed is the maximum speed for CPU0 in Hertz. The value is expected to be
// in the range of 700Mhz to 1.2Ghz.
var MaxSpeed int64

// Version is the Raspberry Pi version 1, 2 or 3.
//
// Is set to 0 when detection (currently primitive) failed.
var Version int

// Pin as connect on the 40 pins extention header.
//
// Schematics are useful to know what is connected to what:
// https://www.raspberrypi.org/documentation/hardware/raspberrypi/schematics/README.md
//
// The actual pin mapping depends on the board revision! The values are set as
// the default for the 40 pins header on Raspberry Pi 2 and Raspberry Pi 3.
//
// P1 is also known as J8.
var (
	P1_1  Pin = V3_3   // 3.3 volt; max 30mA
	P1_2  Pin = V5     // 5 volt (after filtering)
	P1_3  Pin = GPIO2  // I2C_SDA1
	P1_4  Pin = V5     //
	P1_5  Pin = GPIO3  // I2C_SCL1
	P1_6  Pin = GROUND //
	P1_7  Pin = GPIO4  // GPCLK0
	P1_8  Pin = GPIO14 // UART_TXD1
	P1_9  Pin = GROUND //
	P1_10 Pin = GPIO15 // UART_RXD1
	P1_11 Pin = GPIO17 //
	P1_12 Pin = GPIO18 //
	P1_13 Pin = GPIO27 //
	P1_14 Pin = GROUND //
	P1_15 Pin = GPIO22 //
	P1_16 Pin = GPIO23 //
	P1_17 Pin = V3_3   //
	P1_18 Pin = GPIO24 //
	P1_19 Pin = GPIO10 // SPI0_MOSI
	P1_20 Pin = GROUND //
	P1_21 Pin = GPIO9  // SPI0_MISO
	P1_22 Pin = GPIO25 //
	P1_23 Pin = GPIO11 // SPI0_CLK
	P1_24 Pin = GPIO8  // SPI0_CE0
	P1_25 Pin = GROUND //
	P1_26 Pin = GPIO7  // SPI0_CE1

	// Raspberry Pi 2 and later:

	P1_27 Pin = GPIO0  // I2C_SDA0 used to probe for HAT EEPROM, see https://github.com/raspberrypi/hats
	P1_28 Pin = GPIO1  // I2C_SCL0
	P1_29 Pin = GPIO5  // GPCLK1
	P1_30 Pin = GROUND //
	P1_31 Pin = GPIO6  // GPCLK2
	P1_32 Pin = GPIO12 // PWM0_OUT
	P1_33 Pin = GPIO13 // PWM1_OUT
	P1_34 Pin = GROUND //
	P1_35 Pin = GPIO19 // SPI1_MISO
	P1_36 Pin = GPIO16 // SPI1_CE2
	P1_37 Pin = GPIO26 //
	P1_38 Pin = GPIO20 // SPI1_MOSI
	P1_39 Pin = GROUND //
	P1_40 Pin = GPIO21 // SPI1_CLK

	P5_1 = V5
	P5_2 = V3_3
	P5_3 = GPIO28 // I2C0_SDA
	P5_4 = GPIO29 // I2C0_SCL
	P5_5 = GPIO30 // PCM_DIN, UART_CTS0, UART_CST1
	P5_6 = GPIO31 // PCM_DOUT, UART_RTS0, UART_RTS1
	P5_7 = GROUND
	P5_8 = GROUND

	AUDIO_LEFT          = GPIO41
	AUDIO_RIGHT         = GPIO40
	HDMI_HOTPLUG_DETECT = GPIO46
)

// IsConnected returns true if the pin is phyisically connected.
func (p Pin) IsConnected() bool {
	// TODO(maruel): A bit slow, create a lookup table.
	switch p {
	case INVALID:
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
	if i, err := strconv.ParseInt(loadCPUInfo()["Revision"], 16, 32); err == nil {
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

	if bytes, err := ioutil.ReadFile("/sys/devices/system/cpu/cpu0/cpufreq/scaling_max_freq"); err == nil {
		s := strings.TrimSpace(string(bytes))
		if i, err := strconv.ParseInt(s, 10, 64); err == nil {
			// Weirdly, the speed is listed as khz. :(
			MaxSpeed = i * 1000
			sleep160cycles = time.Second * 160 / time.Duration(MaxSpeed)
		} else {
			log.Printf("Failed to parse scaling_max_freq: %s", s)
		}
	} else {
		log.Printf("Failed to read scaling_max_freq: %v", err)
	}

	if Version == 1 {
		// TODO(maruel): Models from 2012 and earlier have P1_3=GPIO0, P1_5=GPIO1 and P1_13=GPIO21.
		// P2 and P3 are not useful.
		// P6 has a RUN pin for reset but it's not available afterward.

		P1_27 = INVALID
		P1_28 = INVALID
		P1_29 = INVALID
		P1_30 = INVALID
		P1_31 = INVALID
		P1_32 = INVALID
		P1_33 = INVALID
		P1_34 = INVALID
		P1_35 = INVALID
		P1_36 = INVALID
		P1_37 = INVALID
		P1_38 = INVALID
		P1_39 = INVALID
		P1_40 = INVALID
	} else {
		P5_1 = INVALID
		P5_2 = INVALID
		P5_3 = INVALID
		P5_4 = INVALID
		P5_5 = INVALID
		P5_6 = INVALID
		P5_7 = INVALID
		P5_8 = INVALID
	}
	if Version < 3 {
		AUDIO_LEFT = GPIO45
	}
}

func loadCPUInfo() map[string]string {
	values := map[string]string{}
	bytes, err := ioutil.ReadFile("/proc/cpuinfo")
	if err != nil {
		return values
	}
	for _, line := range strings.Split(string(bytes), "\n") {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		// Ignore information for other processors than the #0.
		if len(values[key]) == 0 {
			values[key] = strings.TrimSpace(parts[1])
		}
	}
	return values
}
