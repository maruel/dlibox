// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package rpi

import (
	"log"
	"strconv"

	"github.com/maruel/dlibox/go/pio/host/bcm283x"
	"github.com/maruel/dlibox/go/pio/host/internal"
	"github.com/maruel/dlibox/go/pio/host/pins"
)

func init() {
	// Initialize Version. This function is not futureproof, it will return 0 on
	// a Raspberry Pi 4 whenever it comes out.
	rev, _ := internal.CPUInfo["Revision"]
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
			log.Printf("Unknown hardware version: 0x%x", i)
		}
	} else {
		log.Printf("Failed to read cpu_info: %v", err)
	}

	if Version == 1 {
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
}
