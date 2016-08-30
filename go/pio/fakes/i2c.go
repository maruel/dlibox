// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package fakes

import (
	"errors"

	"github.com/maruel/dlibox/go/pio/buses"
)

// I2C implements buses.I2C. It registers everything written to it.
type I2C struct {
	IO []buses.IOFull
}

// Tx currently only support writes.
func (i *I2C) Tx(ios []buses.IOFull) error {
	for i := range ios {
		if o := ios[i].Op; o != buses.Write && o != buses.WriteStop {
			return errors.New("not implemented")
		}
	}
	i.IO = append(i.IO, ios...)
	return nil
}

var _ buses.I2C = &I2C{}
