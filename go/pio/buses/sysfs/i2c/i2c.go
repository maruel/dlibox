// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// +build !linux

package i2c

import (
	"errors"

	"github.com/maruel/dlibox/go/pio/buses"
)

var err = errors.New("not implemented on non-linux OSes")

// Bus is not implemented on non-linux OSes.
type Bus struct{}

// Make is not implemented on non-linux OSes.
func Make(bus int) (*Bus, error) {
	return nil, err
}

func (b *Bus) Close() error {
	return err
}

func (b *Bus) Tx(ios []buses.IOFull) error {
	return err
}
