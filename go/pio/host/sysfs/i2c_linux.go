// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package sysfs

// MakeI2C opens an I²C bus via its sysfs interface as described at
// https://www.kernel.org/doc/Documentation/i2c/dev-interface.
//
// busNumber is the bus number as exported by sysfs. For example if the path is
// /dev/i2c-1, busNumber should be 1.
//
// The resulting object is safe for concurent use.
func MakeI2C(busNumber int) (*I2C, error) {
	return makeI2C(busNumber)
}
