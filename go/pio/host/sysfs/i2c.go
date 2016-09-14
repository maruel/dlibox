// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// +build !linux

package sysfs

import "errors"

var err = errors.New("not implemented on non-linux OSes")

// I2C is not implemented on non-linux OSes.
type I2C struct{}

// MakeI2C is not implemented on non-linux OSes.
func MakeI2C(bus int) (*I2C, error) {
	return nil, err
}

func (i *I2C) Close() error {
	return err
}

func (i *I2C) Tx(addr uint16, w, r []byte) error {
	return err
}
