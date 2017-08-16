// Copyright 2017 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Package leds is a dlibox sink.
package leds

import "errors"

// Dev is an APA102 LED strip.
//
// TODO(maruel): make it more generic so other kind of display are supported.
type Dev struct {
	Name string
	I2C  struct {
		ID string
	}
	SPI struct {
		ID string
		Hz int64
	}
	// NumberLights is the number of lights controlled by this device. If lower
	// than the actual number of lights, the remaining lights will flash oddly.
	NumberLights int
}

func (d *Dev) Validate() error {
	if len(d.Name) == 0 {
		return errors.New("leds: Name is required")
	}
	if len(d.I2C.ID) != 0 {
		if len(d.SPI.ID) != 0 || d.SPI.Hz != 0 {
			return errors.New("leds: can't use both I2C and SPI")
		}
	} else {
		if len(d.SPI.ID) == 0 {
			return errors.New("leds: SPI.ID is required")
		}
		if d.SPI.Hz < 1000 {
			return errors.New("leds: SPI.Hz is required")
		}
	}
	if d.NumberLights == 0 {
		return errors.New("leds: NumberLights is required")
	}
	return nil
}
