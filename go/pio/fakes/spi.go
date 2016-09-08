// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package fakes

import (
	"errors"
	"io"

	"github.com/maruel/dlibox/go/pio/host"
)

// SPI implements host.SPI. It registers everything written to it.
type SPI struct {
	W io.Writer
}

// Close is a no-op.
func (s *SPI) Close() error {
	return nil
}

// Configure is a no-op.
func (s *SPI) Configure(mode host.Mode, bits int) error {
	return nil
}

// Write accumulates all the bytes written.
func (s *SPI) Write(d []byte) (int, error) {
	return s.W.Write(d)
}

// Tx returns an error.
func (s *SPI) Tx(w, r []byte) error {
	return errors.New("not implemented")
}

var _ host.SPI = &SPI{}
