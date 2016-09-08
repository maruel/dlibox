// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// I²C device.

package devices

import "github.com/maruel/dlibox/go/pio/host"

// IO represents one IO to do inside a transaction, with the device address
// already implied.
type IO struct {
	Op  host.Op
	Buf []byte
}

// I2C is a device on a I²C bus.
//
// It saves from repeatedly specifying the device address and implements
// utility functions.
type I2C struct {
	Bus  host.I2C
	Addr uint16
}

// Write writes to the I²C bus without reading, implementing io.Writer.
//
// It's a wrapper for Tx()
func (i *I2C) Write(b []byte) (int, error) {
	if err := i.Tx([]IO{{host.WriteStop, b}}); err != nil {
		return 0, err
	}
	return len(b), nil
}

// ReadReg writes the register number to the I²C bus, then reads data.
//
// It's a wrapper for Tx()
func (i *I2C) ReadReg(reg byte, b []byte) error {
	return i.Tx([]IO{{host.Write, []byte{reg}}, {host.ReadStop, b}})
}

// Tx does a transaction by adding the device's address to each command.
//
// It's a wrapper for I2C.Tx() and causes a memory allocation. Use Bus.Tx()
// directly if memory constrained.
func (i *I2C) Tx(ios []IO) error {
	full := make([]host.IOFull, 0, len(ios))
	for x := range ios {
		full = append(full, host.IOFull{i.Addr, ios[x].Op, ios[x].Buf})
	}
	return i.Bus.Tx(full)
}
