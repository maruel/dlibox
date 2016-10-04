// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package host

import (
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/maruel/dlibox/go/pio"
)

// Init calls pio.Init() and returns it as-is.
//
// The only difference is that by calling host.Init(), you are guaranteed to
// have all the drivers implemented in this library to be implicitly loaded.
func Init() (*pio.State, error) {
	return pio.Init()
}

// MaxSpeed returns the processor maximum speed in Hz.
//
// Returns 0 if it couldn't be calculated.
func MaxSpeed() int64 {
	if isLinux {
		return getMaxSpeedLinux()
	}
	return 0
}

func getMaxSpeedLinux() int64 {
	lock.Lock()
	defer lock.Unlock()
	if maxSpeed == -1 {
		if bytes, err := ioutil.ReadFile("/sys/devices/system/cpu/cpu0/cpufreq/scaling_max_freq"); err == nil {
			s := strings.TrimSpace(string(bytes))
			if i, err := strconv.ParseInt(s, 10, 64); err == nil {
				// Weirdly, the speed is listed as khz. :(
				maxSpeed = i * 1000
			} else {
				log.Printf("Failed to parse scaling_max_freq: %s", s)
				maxSpeed = 0
			}
		} else {
			log.Printf("Failed to read scaling_max_freq: %v", err)
			maxSpeed = 0
		}
	}
	return maxSpeed
}

var (
	lock     sync.Mutex
	maxSpeed int64 = -1
)
