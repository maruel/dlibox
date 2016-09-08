// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// I²C bus.

package host

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
