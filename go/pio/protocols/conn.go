// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package protocols

import "io"

// Conn defines the interface for a connection on a point-to-point
// communication channel.
//
// The channel may either be write-only or read-write. It may either be
// half-duplex or full duplex.
//
// This is the lowest common denominator for I²C (when talking to a specific
// device), SPI, UART, etc.
//
// It is expected (but not enforced) that all Conn implement fmt.Stringer.
type Conn interface {
	// io.Writer can be used for a write-only device.
	io.Writer
	// Tx does a single transaction.
	//
	// For full duplex protocols (SPI, UART), the two buffers must have the same
	// length as both reading and writing happen simultaneously.
	//
	// For half duplex protocols (I²C), there is no restriction as reading
	// happens after writing, and r can be nil.
	Tx(w, r []byte) error
}
