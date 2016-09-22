// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Package uart defines the UART protocol.
package uart

import (
	"github.com/maruel/dlibox/go/pio/protocols"
	"github.com/maruel/dlibox/go/pio/protocols/gpio"
)

// Parity determines the parity bit when transmitting, if any.
type Parity int8

const (
	None Parity = 0
	Odd  Parity = 1
	Even Parity = 2
)

// Bus defines the interface a concrete UART driver must implement.
type Bus interface {
	protocols.Bus
	// Speed changes the bus speed.
	Speed(baud int64) error
	// Configure changes the communication parameters of the bus.
	Configure(startBit, stopBit bool, parity Parity, bits int) error

	// RX returns the receive pin.
	RX() gpio.PinIn
	// TX returns the transmit pin.
	TX() gpio.PinOut
	// RTS returns the request to send pin.
	RTS() gpio.PinIO
	// CTS returns the clear to send pin.
	CTS() gpio.PinIO
}
