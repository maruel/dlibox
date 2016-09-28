// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package host

import (
	"errors"
	"io"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/maruel/dlibox/go/pio/drivers"
	"github.com/maruel/dlibox/go/pio/host/sysfs"
	"github.com/maruel/dlibox/go/pio/protocols/i2c"
	"github.com/maruel/dlibox/go/pio/protocols/spi"
)

// Init calls drivers.Init() and returns it as-is.
//
// The only difference is that by calling host.Init(), you are guaranteed to
// have all the drivers implemented in this library to be implicitly loaded.
func Init() (*drivers.State, error) {
	return drivers.Init()
}

// I2CCloser is a generic I2C bus that can be closed.
type I2CCloser interface {
	io.Closer
	i2c.Conn
}

// SPICloser is a generic SPI bus that can be closed.
type SPICloser interface {
	io.Closer
	spi.Conn
}

//

func newI2CAutoLinux() (I2CCloser, error) {
	Init()
	buses, err := sysfs.EnumerateI2C()
	if err != nil {
		return nil, err
	}
	if len(buses) == 0 {
		return nil, errors.New("no I²C bus found")
	}
	return sysfs.NewI2C(buses[0])
	// TODO(maruel): Fallback with bitbang.NewI2C(). Find two pins available and
	// use them.
}

func newSPIAutoLinux() (SPICloser, error) {
	Init()
	buses, err := sysfs.EnumerateSPI()
	if err != nil {
		return nil, err
	}
	if len(buses) == 0 {
		return nil, errors.New("no SPI bus found")
	}
	return sysfs.NewSPI(buses[0][0], buses[0][1], 0)
	// TODO(maruel): Fallback with bitbang.NewSPI(). Find 4 pins available and
	// use them.
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
