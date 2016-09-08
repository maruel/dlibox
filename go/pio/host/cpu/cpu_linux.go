// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package cpu

import (
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

func init() {
	if bytes, err := ioutil.ReadFile("/sys/devices/system/cpu/cpu0/cpufreq/scaling_max_freq"); err == nil {
		s := strings.TrimSpace(string(bytes))
		if i, err := strconv.ParseInt(s, 10, 64); err == nil {
			// Weirdly, the speed is listed as khz. :(
			MaxSpeed = i * 1000
		} else {
			log.Printf("Failed to parse scaling_max_freq: %s", s)
		}
	} else {
		log.Printf("Failed to read scaling_max_freq: %v", err)
	}
}
