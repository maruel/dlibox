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

	"github.com/maruel/dlibox/go/pio"
	"github.com/maruel/dlibox/go/pio/host/sysfs"
	"github.com/maruel/dlibox/go/pio/protocols/i2c"
	"github.com/maruel/dlibox/go/pio/protocols/spi"
)

// Init calls pio.Init() and returns it as-is.
//
// The only difference is that by calling host.Init(), you are guaranteed to
// have all the drivers implemented in this library to be implicitly loaded.
func Init() (*pio.State, error) {
	return pio.Init()
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

// MaxSpeed returns the processor maximum speed in Hz.
//
// Returns 0 if it couldn't be calculated.
func MaxSpeed() int64 {
	if isLinux {
		return getMaxSpeedLinux()
	}
	return 0
}

// NewI2C opens an I²C bus using the most appropriate driver.
func NewI2C(busNumber int) (I2CCloser, error) {
	if isLinux {
		return sysfs.NewI2C(busNumber)
	}
	return nil, errors.New("no i²c driver found")
}

// NewSPI opens an SPI bus using the most appropriate driver.
func NewSPI(busNumber, cs int) (SPICloser, error) {
	if isLinux {
		return sysfs.NewSPI(busNumber, cs, 0)
	}
	return nil, errors.New("no spi driver found")
}

// NewI2CAuto opens the first available I²C bus.
//
// You can query the return value to determine which pins are being used.
func NewI2CAuto() (I2CCloser, error) {
	if isLinux {
		return newI2CAutoLinux()
	}
	return nil, errors.New("no i²c driver found")
}

// NewSPIAuto opens the first available SPI bus.
//
// You can query the return value to determine which pins are being used.
func NewSPIAuto() (SPICloser, error) {
	if isLinux {
		return newSPIAutoLinux()
	}
	return nil, errors.New("no spi driver found")
}

//

func newI2CAutoLinux() (I2CCloser, error) {
	if _, err := Init(); err != nil {
		return nil, err
	}
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
	if _, err := Init(); err != nil {
		return nil, err
	}
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
