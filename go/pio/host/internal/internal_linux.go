// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package internal

import (
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func readAndSplit(path string) map[string]string {
	bytes, err := ioutil.ReadFile("/proc/cpuinfo")
	if err != nil {
		return nil
	}
	out := map[string]string{}
	for _, line := range strings.Split(string(bytes), "\n") {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		if len(key) == 0 || key[0] == '#' {
			continue
		}
		// Ignore duplicate keys
		if len(out[key]) == 0 {
			s := strings.TrimSpace(parts[1])
			if len(s) > 2 && s[0] == '"' && s[len(s)-1] == '"' {
				// Not exactly 100% right but #closeenough. See for more details
				// https://www.freedesktop.org/software/systemd/man/os-release.html
				s, err = strconv.Unquote(s[1 : len(s)-2])
				if err != nil {
					continue
				}
			}
			out[key] = s
		}
	}
	return out
}

func init() {
	// Technically speaking, cpuinfo doesn't contain quotes and os-release
	// doesn't contain duplicate keys. Make it more strictly correct if needed.
	if m := readAndSplit("/proc/cpuinfo"); m != nil {
		CPUInfo = m
	}
	if m := readAndSplit("/etc/os-release"); m != nil {
		OSRelease = m
	}
}

// OS

func IsArmbian() bool {
	// This is iffy at best.
	// Armbian presents itself as debian in /etc/os-release.
	_, err := os.Stat("/etc/armbian.txt")
	return err == nil
}

func IsRaspbian() bool {
	id, _ := OSRelease["ID"]
	return id == "raspbian"
}

// CPU

func IsBCM283x() bool {
	//_, err := os.Stat("/sys/bus/platform/drivers/bcm2835_thermal")
	//return err == nil
	hardware, ok := CPUInfo["Hardware"]
	return ok && strings.HasPrefix(hardware, "BCM")
}

func IsAllWinner() bool {
	// TODO(maruel): This is too vague.
	hardware, ok := CPUInfo["Hardware"]
	return ok && strings.HasPrefix(hardware, "sun")
	// /sys/class/sunxi_info/sys_info
}

// Board

func IsRaspberryPi() bool {
	// This is iffy at best.
	_, err := os.Stat("/sys/bus/platform/drivers/raspberrypi-firmware")
	return err == nil
}

func IsPine64() bool {
	// This is iffy at best.
	_, err := os.Stat("/boot/pine64.dtb")
	return err == nil
}
