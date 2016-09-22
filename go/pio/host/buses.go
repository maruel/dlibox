// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package host

import (
	"errors"
	"io"

	"github.com/maruel/dlibox/go/pio/host/drivers/sysfs"
	"github.com/maruel/dlibox/go/pio/protocols/i2c"
	"github.com/maruel/dlibox/go/pio/protocols/spi"
)

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

// NewI2C opens the first available I²C bus.
//
// You can query the return value to determine which pins are being used.
func NewI2C() (I2CCloser, error) {
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

// NewSPI opens the first available SPI bus.
//
// You can query the return value to determine which pins are being used.
func NewSPI() (SPICloser, error) {
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
