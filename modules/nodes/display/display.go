// Copyright 2017 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Package display is a dlibox sink.
package display

import "errors"

// Dev is an ssd1306 display.
//
// TODO(maruel): make it more generic so other kind of display are supported.
type Dev struct {
	Name string
	I2C  struct {
		ID string
	}
	W, H int
}

func (d *Dev) Validate() error {
	if len(d.Name) == 0 {
		return errors.New("display: Name is required")
	}
	if len(d.I2C.ID) == 0 {
		return errors.New("display: I2C.ID is required")
	}
	if d.W == 0 {
		return errors.New("display: W is required")
	}
	if d.H == 0 {
		return errors.New("display: H is required")
	}
	return nil
}
