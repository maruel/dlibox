// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// TODO(maruel): This could work on FreeBSD too. Make it build on ARM only
// instead.

package rpi

import (
	"log"
	"strconv"

	"github.com/maruel/dlibox/go/pio/host"
	"github.com/maruel/dlibox/go/pio/host/internal"
	"github.com/maruel/dlibox/go/pio/host/internal/bcm283x"
)

func init() {
	// TODO(maruel): Do not run this code when running from non BCM CPU.

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
		Headers["P1"] = [][]host.Pin{
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
		}
		Headers["P5"] = [][]host.Pin{
			{P5_1, P5_2},
			{P5_3, P5_4},
			{P5_5, P5_6},
			{P5_7, P5_8},
		}
		// TODO(maruel): Models from 2012 and earlier have P1_3=GPIO0, P1_5=GPIO1 and P1_13=GPIO21.
		// P2 and P3 are not useful.
		// P6 has a RUN pin for reset but it's not available after Pi version 1.

		P1_27 = host.INVALID
		P1_28 = host.INVALID
		P1_29 = host.INVALID
		P1_30 = host.INVALID
		P1_31 = host.INVALID
		P1_32 = host.INVALID
		P1_33 = host.INVALID
		P1_34 = host.INVALID
		P1_35 = host.INVALID
		P1_36 = host.INVALID
		P1_37 = host.INVALID
		P1_38 = host.INVALID
		P1_39 = host.INVALID
		P1_40 = host.INVALID
	} else {
		Headers["P1"] = [][]host.Pin{
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
		}
		P5_1 = host.INVALID
		P5_2 = host.INVALID
		P5_3 = host.INVALID
		P5_4 = host.INVALID
		P5_5 = host.INVALID
		P5_6 = host.INVALID
		P5_7 = host.INVALID
		P5_8 = host.INVALID
	}
	if Version < 3 {
		AUDIO_LEFT = bcm283x.GPIO45
	}
	Headers["AUDIO"] = [][]host.Pin{
		{AUDIO_LEFT},
		{AUDIO_RIGHT},
	}
}
