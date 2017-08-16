// Copyright 2017 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Package sound is a dlibox sink.
package sound

import "errors"

// Dev is a sound output device.
type Dev struct {
	Name     string
	DeviceID int // Non-zero when there's more than one sound card installed.
}

func (d *Dev) Validate() error {
	if len(d.Name) == 0 {
		return errors.New("sound: Name is required")
	}
	return nil
}
