// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package rpi

import (
	"io/ioutil"
	"strconv"
	"strings"
)

var version int

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

func maxCPUSpeed() int64 {
	bytes, err := ioutil.ReadFile("/sys/devices/system/cpu/cpu0/cpufreq/scaling_max_freq")
	if err != nil {
		return 0
	}
	i, _ := strconv.ParseInt(string(bytes), 10, 64)
	return i
}

// Version returns the Raspberry Pi version 1, 2 or 3.
//
// This function is not futureproof, it will return 0 on a Raspberry Pi 4
// whenever it comes out.
func Version() int {
	if version == 0 {
		i, err := strconv.Atoi(loadCPUInfo()["Revision"])
		if err == nil {
			if i < 0x20 {
				version = 1
			} else if i == 0xa01041 || i == 0xa21041 {
				version = 2
			} else if i == 0xa02082 {
				version = 3
			}
		}
	}
	return version
}
