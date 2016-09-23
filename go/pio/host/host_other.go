// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// +build !linux

package host

import "errors"

// MaxSpeed returns the processor maximum speed in Hz.
//
// Returns 0 if it couldn't be calculated.
func MaxSpeed() int64 {
	return 0
}

// NewI2C opens an I²C bus using the most appropriate driver.
func NewI2C(busNumber int) (I2CCloser, error) {
	return nil, errors.New("no i²c driver found")
}

// NewSPI opens an SPI bus using the most appropriate driver.
func NewSPI(busNumber, cs int) (SPICloser, error) {
	return nil, errors.New("no spi driver found")
}

// NewI2CAuto opens the first available I²C bus.
//
// You can query the return value to determine which pins are being used.
func NewI2CAuto() (I2CCloser, error) {
	return nil, errors.New("no i²c driver found")
}

// NewSPIAuto opens the first available SPI bus.
//
// You can query the return value to determine which pins are being used.
func NewSPIAuto() (SPICloser, error) {
	return nil, errors.New("no spi driver found")
}
