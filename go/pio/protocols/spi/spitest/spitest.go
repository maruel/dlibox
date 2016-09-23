// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package spitest

import (
	"errors"
	"io"
	"sync"

	"github.com/maruel/dlibox/go/pio/protocols/gpio"
	"github.com/maruel/dlibox/go/pio/protocols/pins"
	"github.com/maruel/dlibox/go/pio/protocols/spi"
)

// Record implements spi.Conn. It registers everything written to it.
type Record struct {
	sync.Mutex
	W io.Writer
}

// Close is a no-op.
func (s *Record) Close() error {
	s.Lock()
	defer s.Unlock()
	return nil
}

// Speed is a no-op.
func (s *Record) Speed(hz int64) error {
	return nil
}

// Configure is a no-op.
func (s *Record) Configure(mode spi.Mode, bits int) error {
	return nil
}

// Write accumulates all the bytes written.
func (s *Record) Write(d []byte) (int, error) {
	s.Lock()
	defer s.Unlock()
	return s.W.Write(d)
}

// Tx only support writes.
func (s *Record) Tx(w, r []byte) error {
	if len(r) != 0 {
		return errors.New("not implemented")
	}
	_, err := s.Write(w)
	return err
}

func (s *Record) CLK() gpio.PinOut {
	return pins.INVALID
}

func (s *Record) MISO() gpio.PinIn {
	return pins.INVALID
}

func (s *Record) MOSI() gpio.PinOut {
	return pins.INVALID
}

func (s *Record) CS() gpio.PinOut {
	return pins.INVALID
}

var _ spi.Conn = &Record{}
