// Copyright 2017 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Package button is a dlibox event source from a physical button.
package button

import "errors"

// Dev represents a physical GPIO input pin.
type Dev struct {
	Name string
	Pin  string
}

// Validate returns true if Dev is correctly initialized.
func (d *Dev) Validate() error {
	if len(d.Name) == 0 {
		return errors.New("button: Name is required")
	}
	if len(d.Pin) == 0 {
		return errors.New("button: Pin is required")
	}
	return nil
}
