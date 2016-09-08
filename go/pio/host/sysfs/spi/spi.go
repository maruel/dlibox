// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// +build !linux

package spi

import (
	"errors"

	"github.com/maruel/dlibox/go/pio/host"
)

var err = errors.New("not implemented on non-linux OSes")

// Bus is not implemented on non-linux OSes.
type Bus struct{}

// Make is not implemented on non-linux OSes.
func Make(bus int) (*Bus, error) {
	return nil, err
}

func (s *Bus) Close() error {
	return err
}

func (s *Bus) Configure(mode host.Mode, bits int) error {
	return err
}

func (s *Bus) Write(b []byte) (int, error) {
	return 0, err
}

func (s *Bus) Tx(w, r []byte) error {
	return err
}
