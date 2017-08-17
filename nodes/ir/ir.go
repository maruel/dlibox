// Copyright 2017 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Package ir is a dlibox event source from a InfraRed remote.
package ir

import "errors"

type Dev struct {
	Name string
}

func (d *Dev) Validate() error {
	if len(d.Name) == 0 {
		return errors.New("ir: Name is required")
	}
	return nil
}
