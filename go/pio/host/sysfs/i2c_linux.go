// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package sysfs

import "github.com/maruel/dlibox/go/pio/drivers"

// NewI2C opens an I²C bus via its sysfs interface as described at
// https://www.kernel.org/doc/Documentation/i2c/dev-interface.
//
// busNumber is the bus number as exported by sysfs. For example if the path is
// /dev/i2c-1, busNumber should be 1.
//
// The resulting object is safe for concurent use.
func NewI2C(busNumber int) (*I2C, error) {
	return newI2C(busNumber)
}

// EnumerateI2C returns the available I²C buses.
func EnumerateI2C() ([]int, error) {
	return enumerateI2C()
}

func init() {
	drivers.MustRegister(&driverI2C{})
}
