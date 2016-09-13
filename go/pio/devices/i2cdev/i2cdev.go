// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// package i2cdev is an adapter code to directly address an I²C device on a I²C
// bus without having to continuously specify the address when doing I/O.
//
// BUG(maruel): Where does this package/struct belong?
package i2cdev

import "github.com/maruel/dlibox/go/pio/host"

// Dev is a device on a I²C bus.
//
// It implements host.Bus.
//
// It saves from repeatedly specifying the device address and implements
// utility functions.
type Dev struct {
	Bus  host.I2C
	Addr uint16
}

// Write writes to the I²C bus without reading, implementing io.Writer.
//
// It's a wrapper for Tx()
func (d *Dev) Write(b []byte) (int, error) {
	if err := d.Tx(b, nil); err != nil {
		return 0, err
	}
	return len(b), nil
}

// ReadReg writes the register number to the I²C bus, then reads data.
//
// It's a wrapper for Tx()
func (d *Dev) ReadReg(reg byte, b []byte) error {
	return d.Tx([]byte{reg}, b)
}

// Tx does a transaction by adding the device's address to each command.
//
// It's a wrapper for Dev.Bus.Tx().
func (d *Dev) Tx(w, r []byte) error {
	return d.Bus.Tx(d.Addr, w, r)
}
