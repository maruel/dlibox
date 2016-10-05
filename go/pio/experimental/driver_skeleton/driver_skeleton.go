// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package driver_skeleton

import (
	"errors"

	"github.com/maruel/dlibox/go/pio"
	"github.com/maruel/dlibox/go/pio/protocols"
	"github.com/maruel/dlibox/go/pio/protocols/i2c"
)

// FIXME: Expose public symbols as relevant. Do not export more than needed!
// See
// https://github.com/maruel/dlibox/blob/master/go/pio/DESIGN.md#requirements
// for the expectations.

// Dev is a handle to the device. FIXME.
type Dev struct {
	c protocols.Conn
}

// New opens a handle to the device. FIXME.
func New(i i2c.Conn) (*Dev, error) {
	d := &Dev{&i2c.Dev{i, 42}}
	// FIXME: Simulate a setup dance.
	var b [2]byte
	if err := d.c.Tx([]byte("in"), b[:]); err != nil {
		return nil, err
	}
	if b[0] != 'I' || b[1] != 'N' {
		return nil, errors.New("unexpected reply")
	}
	return d, nil
}

// Read is a method on your device. FIXME.
func (d *Dev) Read() string {
	var b [12]byte
	if err := d.c.Tx([]byte("what"), b[:]); err != nil {
		return err.Error()
	}
	return string(b[:])
}

// FIXME: A driver is generally only needed for host drivers. If you implement
// a device driver, delete the remainder of this file.

// driver implements pio.Driver.
type driver struct {
}

func (d *driver) String() string {
	// FIXME: Change this string to be unique. It must match the directory name.
	return "driver_skeleton"
}

func (d *driver) Type() pio.Type {
	// FIXME: Change this to be the type of driver.
	return pio.Device
}

func (d *driver) Prerequisites() []string {
	// FIXME: Declare prerequisites drivers if relevant.
	return nil
}

func (d *driver) Init() (bool, error) {
	// FIXME: If the driver is not needed, do the following:
	// return false, errors.New("not running on a skeleton")

	// FIXME: Add implementation.

	return true, errors.New("not implemented")
}

// FIXME: This verifies that the driver implements all the required methods.
var _ pio.Driver = &driver{}
