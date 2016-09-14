// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// I²C bus.

package host

// I2C defines the function a concrete I²C driver must implement.
//
// This interface doesn't implement Bus since a device address must be
// specified. Use i2cdev.Dev as an adapter to get a Bus compatible object.
type I2C interface {
	Tx(addr uint16, w, r []byte) error
}
