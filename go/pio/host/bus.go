// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Generic bus.

package host

import "io"

// Bus defines the interface for a generic bus. This is the lowest common
// denominator for I²C (when talking to a specific device) and SPI.
type Bus interface {
	io.Writer
	// Tx does a single transaction.
	//
	// For SPI, the two buffers must have the same length. For I²C, there is no
	// restriction and r can be nil.
	//
	// Warning: 'write' is the first argument.
	Tx(w, r []byte) error
}
