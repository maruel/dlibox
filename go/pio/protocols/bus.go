// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package protocols

import "io"

// Conn defines the interface for a generic bus.
//
// This is the lowest common denominator for I²C (when talking to a specific
// device), SPI, UART, etc.
type Conn interface {
	// io.Writer can be used for a write-only device.
	io.Writer
	// Tx does a single transaction.
	//
	// For full duplex protocols (SPI, UART), the two buffers must have the same
	// length.
	//
	// For half duplex protocols (I²C), there is no restriction and r can be nil.
	//
	// Warning: 'write' is the first argument.
	Tx(w, r []byte) error
}
