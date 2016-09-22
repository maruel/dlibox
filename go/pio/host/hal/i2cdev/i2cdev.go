// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// package i2cdev is an adapter code to directly address an I²C device on a I²C
// bus without having to continuously specify the address when doing I/O.
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

// Tx does a transaction by adding the device's address to each command.
//
// It's a wrapper for Dev.Bus.Tx().
func (d *Dev) Tx(w, r []byte) error {
	return d.Bus.Tx(d.Addr, w, r)
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

// Shortcuts

// ReadRegUint8 reads a 8 bit register.
func (d *Dev) ReadRegUint8(reg byte) (uint8, error) {
	var v [1]byte
	err := d.ReadReg(reg, v[:])
	return v[0], err
}

// ReadRegUint16 reads a 16 bit register as big endian.
func (d *Dev) ReadRegUint16(reg byte) (uint16, error) {
	var v [2]byte
	err := d.ReadReg(reg, v[:])
	return uint16(v[0])<<8 | uint16(v[1]), err
}

// WriteRegUint8 writes a 8 bit register.
func (d *Dev) WriteRegUint8(reg byte, v uint8) error {
	_, err := d.Write([]byte{reg, v})
	return err
}

// WriteRegUint16 writes a 16 bit register.
func (d *Dev) WriteRegUint16(reg byte, v uint16) error {
	_, err := d.Write([]byte{reg, byte(v >> 16), byte(v)})
	return err
}
