// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package hosttest

import (
	"errors"
	"io"
	"sync"

	"github.com/maruel/dlibox/go/pio/protocols/gpio"
	"github.com/maruel/dlibox/go/pio/protocols/pins"
	"github.com/maruel/dlibox/go/pio/protocols/spi"
)

// SPI implements spi.Bus. It registers everything written to it.
//
// BUG(maruel): SPI does not support reading yet.
type SPI struct {
	sync.Mutex
	W io.Writer
}

// Close is a no-op.
func (s *SPI) Close() error {
	s.Lock()
	defer s.Unlock()
	return nil
}

// Speed is a no-op.
func (s *SPI) Speed(hz int64) error {
	s.Lock()
	defer s.Unlock()
	return nil
}

// Configure is a no-op.
func (s *SPI) Configure(mode spi.Mode, bits int) error {
	s.Lock()
	defer s.Unlock()
	return nil
}

// Write accumulates all the bytes written.
func (s *SPI) Write(d []byte) (int, error) {
	s.Lock()
	defer s.Unlock()
	return s.W.Write(d)
}

// Tx only support writes.
func (s *SPI) Tx(w, r []byte) error {
	if len(r) != 0 {
		return errors.New("not implemented")
	}
	_, err := s.Write(w)
	return err
}

func (s *SPI) CLK() gpio.PinOut {
	return pins.INVALID
}

func (s *SPI) MISO() gpio.PinIn {
	return pins.INVALID
}

func (s *SPI) MOSI() gpio.PinOut {
	return pins.INVALID
}

func (s *SPI) CS() gpio.PinOut {
	return pins.INVALID
}

var _ spi.Bus = &SPI{}
