// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Package internal implements common functionality to auto-detect the host.
package internal

import (
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"
	"unicode"
)

// IsArmbian returns true if running on a Armbian distribution.
//
// http://www.armbian.com/
func IsArmbian() bool {
	if isArm && isLinux {
		// This is iffy at best.
		// Armbian presents itself as debian in /etc/os-release.
		_, err := os.Stat("/etc/armbian.txt")
		return err == nil
	}
	return false
}

// IsRaspbian returns true if running on a Raspbian distribution.
//
// https://raspbian.org/
func IsRaspbian() bool {
	if isArm && isLinux {
		id, _ := OSRelease()["ID"]
		return id == "raspbian"
	}
	return false
}

// CPU

// IsBCM283x returns true if running on a Broadcom bcm283x based CPU.
func IsBCM283x() bool {
	if isArm {
		//_, err := os.Stat("/sys/bus/platform/drivers/bcm2835_thermal")
		//return err == nil
		hardware, ok := CPUInfo()["Hardware"]
		return ok && strings.HasPrefix(hardware, "BCM")
	}
	return false
}

// IsAllwinner returns true if running on an Allwinner based CPU.
//
// https://en.wikipedia.org/wiki/Allwinner_Technology
func IsAllwinner() bool {
	if isArm {
		// TODO(maruel): This is too vague.
		hardware, ok := CPUInfo()["Hardware"]
		return ok && strings.HasPrefix(hardware, "sun")
		// /sys/class/sunxi_info/sys_info
	}
	return false
}

// CPUInfo returns parsed data from /proc/cpuinfo.
func CPUInfo() map[string]string {
	if isLinux {
		return makeCPUInfoLinux()
	}
	return cpuInfo
}

// CPUInfo returns parsed data from /etc/os-release.
func OSRelease() map[string]string {
	if isLinux {
		return makeOSReleaseLinux()
	}
	return osRelease
}

// Board

// IsRaspberryPi returns true if running on a raspberry pi board.
//
// https://www.raspberrypi.org/
func IsRaspberryPi() bool {
	if isArm {
		// This is iffy at best.
		_, err := os.Stat("/sys/bus/platform/drivers/raspberrypi-firmware")
		return err == nil
	}
	return false
}

// IsPine64 returns true if running on a pine64 board.
//
// https://www.pine64.org/
func IsPine64() bool {
	if isArm {
		// This is iffy at best.
		_, err := os.Stat("/boot/pine64.dtb")
		return err == nil
	}
	return false
}

//

var (
	lock      sync.Mutex
	cpuInfo   map[string]string
	osRelease map[string]string
)

func splitSemiColon(content string) map[string]string {
	// Strictly speaking this format isn't ok, there can be multiple group.
	out := map[string]string{}
	for _, line := range strings.Split(content, "\n") {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		// This format may have space around the ':'.
		key := strings.TrimRightFunc(parts[0], unicode.IsSpace)
		if len(key) == 0 || key[0] == '#' {
			continue
		}
		// Ignore duplicate keys.
		// TODO(maruel): Keep them all.
		if _, ok := out[key]; !ok {
			// Trim on both side, trailing space was observed on "Features" value.
			out[key] = strings.TrimFunc(parts[1], unicode.IsSpace)
		}
	}
	return out
}

func splitStrict(content string) map[string]string {
	out := map[string]string{}
	for _, line := range strings.Split(content, "\n") {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := parts[0]
		if len(key) == 0 || key[0] == '#' {
			continue
		}
		// Overwrite previous key.
		value := parts[1]
		if len(value) > 2 && value[0] == '"' && value[len(value)-1] == '"' {
			// Not exactly 100% right but #closeenough. See for more details
			// https://www.freedesktop.org/software/systemd/man/os-release.html
			var err error
			value, err = strconv.Unquote(value)
			if err != nil {
				continue
			}
		}
		out[key] = value
	}
	return out
}

func makeCPUInfoLinux() map[string]string {
	lock.Lock()
	defer lock.Unlock()
	if cpuInfo == nil {
		if bytes, err := ioutil.ReadFile("/proc/cpuinfo"); err == nil {
			cpuInfo = splitSemiColon(string(bytes))
		} else {
			cpuInfo = map[string]string{}
		}
	}
	return cpuInfo
}

func makeOSReleaseLinux() map[string]string {
	lock.Lock()
	defer lock.Unlock()
	if osRelease == nil {
		if bytes, err := ioutil.ReadFile("/etc/os-release"); err == nil {
			osRelease = splitStrict(string(bytes))
		} else {
			osRelease = map[string]string{}
		}
	}
	return osRelease
}
