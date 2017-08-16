// Copyright 2017 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Package pir is a dlibox event source from a motion detector.
package pir

import "errors"

type Dev struct {
	Name string
	Pin  string
}

func (d *Dev) Validate() error {
	if len(d.Name) == 0 {
		return errors.New("pir: Name is required")
	}
	if len(d.Pin) == 0 {
		return errors.New("pir: Pin is required")
	}
	return nil
}
