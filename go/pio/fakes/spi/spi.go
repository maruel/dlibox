// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Package spi implements a fake SPI bus.
package spi

import (
	"bytes"
	"errors"

	"github.com/maruel/dlibox/go/pio/buses"
)

// Bus registers everything written to it.
type Bus struct {
	Buf bytes.Buffer
}

// Close is a no-op.
func (b *Bus) Close() error {
	return nil
}

// Configure is a no-op.
func (b *Bus) Configure(mode buses.Mode, bits int) error {
	return nil
}

// Write accumulates all the bytes written.
func (b *Bus) Write(d []byte) (int, error) {
	return b.Buf.Write(d)
}

// Tx returns an error.
func (b *Bus) Tx(w, r []byte) error {
	return errors.New("not implemented")
}
