// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package internal

import (
	"os"
	"strings"
)

// OS

// IsArmbian returns true if running on a Armbian distribution.
//
// http://www.armbian.com/
func IsArmbian() bool {
	// This is iffy at best.
	// Armbian presents itself as debian in /etc/os-release.
	_, err := os.Stat("/etc/armbian.txt")
	return err == nil
}

// IsRaspbian returns true if running on a Raspbian distribution.
//
// https://raspbian.org/
func IsRaspbian() bool {
	id, _ := OSRelease()["ID"]
	return id == "raspbian"
}

// CPU

// IsBCM283x returns true if running on a Broadcom bcm283x based CPU.
func IsBCM283x() bool {
	//_, err := os.Stat("/sys/bus/platform/drivers/bcm2835_thermal")
	//return err == nil
	hardware, ok := CPUInfo()["Hardware"]
	return ok && strings.HasPrefix(hardware, "BCM")
}

// IsAllwinner returns true if running on an Allwinner based CPU.
//
// https://en.wikipedia.org/wiki/Allwinner_Technology
func IsAllwinner() bool {
	// TODO(maruel): This is too vague.
	hardware, ok := CPUInfo()["Hardware"]
	return ok && strings.HasPrefix(hardware, "sun")
	// /sys/class/sunxi_info/sys_info
}

// Board

// IsRaspberryPi returns true if running on a raspberry pi board.
//
// https://www.raspberrypi.org/
func IsRaspberryPi() bool {
	// This is iffy at best.
	_, err := os.Stat("/sys/bus/platform/drivers/raspberrypi-firmware")
	return err == nil
}

// IsPine64 returns true if running on a pine64 board.
//
// https://www.pine64.org/
func IsPine64() bool {
	// This is iffy at best.
	_, err := os.Stat("/boot/pine64.dtb")
	return err == nil
}
