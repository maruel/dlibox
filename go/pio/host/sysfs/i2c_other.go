// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// +build !linux

package sysfs

import "errors"

// MakeI2C opens an I²C bus via its sysfs interface as described at
// https://www.kernel.org/doc/Documentation/i2c/dev-interface It is not
// Raspberry Pi specific.
//
// The resulting object is safe for concurent use.
func MakeI2C(bus int) (*I2C, error) {
	return nil, errors.New("sysfs.i2c is not supported on this platform")
}
