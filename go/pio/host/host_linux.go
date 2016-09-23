// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package host

import "github.com/maruel/dlibox/go/pio/host/sysfs"

// MaxSpeed returns the processor maximum speed in Hz.
//
// Returns 0 if it couldn't be calculated.
func MaxSpeed() int64 {
	return getMaxSpeedLinux()
}

// NewI2C opens an I²C bus using the most appropriate driver.
func NewI2C(busNumber int) (I2CCloser, error) {
	return sysfs.NewI2C(busNumber)
}

// NewSPI opens an SPI bus using the most appropriate driver.
func NewSPI(busNumber, cs int) (SPICloser, error) {
	return sysfs.NewSPI(busNumber, cs, 0)
}

// NewI2CAuto opens the first available I²C bus.
//
// You can query the return value to determine which pins are being used.
func NewI2CAuto() (I2CCloser, error) {
	return newI2CAutoLinux()
}

// NewSPIAuto opens the first available SPI bus.
//
// You can query the return value to determine which pins are being used.
func NewSPIAuto() (SPICloser, error) {
	return newSPIAutoLinux()
}
