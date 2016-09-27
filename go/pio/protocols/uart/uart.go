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
type Parity byte

const (
	None  Parity = 'N'
	Odd   Parity = 'O'
	Even  Parity = 'E'
	Mark  Parity = 'M' // always 1
	Space Parity = 'S' // always 0
)

// Stop determines what stop bit to use.
type Stop int8

const (
	One     Stop = 0 // 1 stop bit
	OneHalf Stop = 1 // 1.5 stop bits
	Two     Stop = 2 // 2 stop bits
)

// conn defines the interface a concrete UART driver must implement.
//
// BUG(maruel): Export once the interface is used and functional.
type conn interface {
	protocols.Conn
	// Speed changes the bus speed.
	Speed(baud int64) error
	// Configure changes the communication parameters of the bus.
	Configure(stopBit Stop, parity Parity, bits int) error
}

// Pins defines the pins that an UART bus interconnect is using on the host.
//
// It is expected that a implementer of Conn also implement Pins but this is
// not a requirement.
type Pins interface {
	// RX returns the receive pin.
	RX() gpio.PinIn
	// TX returns the transmit pin.
	TX() gpio.PinOut
	// RTS returns the request to send pin.
	RTS() gpio.PinIO
	// CTS returns the clear to send pin.
	CTS() gpio.PinIO
}
