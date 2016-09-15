// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// +build !linux

package sysfs

import (
	"errors"

	"github.com/maruel/dlibox/go/pio/host"
)

var errNotImpl = errors.New("not implemented on non-linux OSes")

// SPI is not implemented on non-linux OSes.
type SPI struct{}

// MakeSPI is not implemented on non-linux OSes.
func MakeSPI(bus int) (*SPI, error) {
	return nil, err
}

func (s *SPI) Close() error {
	return errNotImpl
}

func (s *SPI) Configure(mode host.Mode, bits int) error {
	return errNotImpl
}

func (s *SPI) Write(b []byte) (int, error) {
	return 0, errNotImpl
}

func (s *SPI) Tx(w, r []byte) error {
	return errNotImpl
}
