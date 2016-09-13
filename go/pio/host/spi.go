// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// SPI

package host

// Mode determines how communication is done. The bits can be OR'ed to change
// the polarity and phase used for communication.
type Mode int

const (
	Mode0 Mode = 0x0
	Mode1 Mode = 0x1
	Mode2 Mode = 0x2
	Mode3 Mode = 0x3
)

// SPI defines the interface a concrete SPI driver must implement.
type SPI interface {
	Bus
	Configure(mode Mode, bits int) error
}
