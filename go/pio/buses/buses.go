// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Package buses defines generic interfaces for buses.
//
// Subpackages contain the concrete implementations. Devices accept interface,
// constructors return concrete type.
package buses

import "io"

// SPI

// Mode determines how communication is done. The bits can be OR'ed to change
// the polarity and phase used for communication.
type Mode int

const (
	Mode0 Mode = 0x0
	Mode1 Mode = 0x1
	Mode2 Mode = 0x2
	Mode3 Mode = 0x3
)

// SPI defines the functions a contrete SPI driver must implement.
type SPI interface {
	io.Writer
	Configure(mode Mode, bits int) error
	Tx(r, w []byte) error
}

// I²C

// Op is an operation to do in one IO operation.
type Op uint8

const (
	Write     Op = 0 // Write operation that should be followed by another one.
	WriteStop Op = 1 // Write operation then sends a STOP.
	Read      Op = 2 // Read operation that should be following by another one.
	ReadStop  Op = 3 // Read operation then a sends a STOP.
)

// IOFull is one I/O to do inside a transaction.
//
// This is useful if a client needs to do I/O to multiple devices as a single
// transaction.
type IOFull struct {
	Addr uint16
	Op   Op
	Buf  []byte
}

// I2C defines the function a concrete I²C driver must implement.
type I2C interface {
	Tx([]IOFull) error
}

// Dev is a device on a I²C bus.
//
// It saves from repeatedly specifying the device address and implements
// utility functions.
type Dev struct {
	Bus  I2C
	Addr uint16
}

// Write writes to the I²C bus without reading, implementing io.Writer.
//
// It's a wrapper for Tx()
func (d *Dev) Write(b []byte) (int, error) {
	if err := d.Tx([]IO{{WriteStop, b}}); err != nil {
		return 0, err
	}
	return len(b), nil
}

// WriteBytes writes to the I²C bus without reading.
//
// This is a shorter form than Write(). It's a wrapper for Tx()
func (d *Dev) WriteBytes(b ...byte) error {
	return d.Tx([]IO{{WriteStop, b}})
}

// ReadReg writes the register number to the I²C bus, then reads data.
//
// It's a wrapper for Tx()
func (d *Dev) ReadReg(reg byte, b []byte) error {
	return d.Tx([]IO{{Write, []byte{reg}}, {ReadStop, b}})
}

// IO represents one IO to do inside a transaction, with the device address
// already implied.
type IO struct {
	Op  Op
	Buf []byte
}

// Tx does a transaction by adding the device's address to each command.
//
// It's a wrapper for I2C.Tx() and causes a memory allocation. Use Bus.Tx()
// directly if memory constrained.
func (d *Dev) Tx(ios []IO) error {
	full := make([]IOFull, 0, len(ios))
	for i := range ios {
		full = append(full, IOFull{d.Addr, ios[i].Op, ios[i].Buf})
	}
	return d.Bus.Tx(full)
}
