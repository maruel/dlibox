// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package dotstar

import (
	"fmt"
	"io"
	"os"
)

type Pin int

const (
	P1_7  Pin = 4
	P1_11 Pin = 17
	P1_12 Pin = 18
	P1_13 Pin = 21
	P1_15 Pin = 22
	P1_16 Pin = 23
	P1_18 Pin = 24
	P1_22 Pin = 25

	GPIO4  Pin = 4
	GPIO17 Pin = 17
	GPIO18 Pin = 18
	GPIO21 Pin = 21
	GPIO22 Pin = 22
	GPIO23 Pin = 23
	GPIO24 Pin = 24
	GPIO25 Pin = 25
)

// Close the handle implicitly open by either SetPinPWM or ReleasePinPWM.
//
// It's not required to call it.
func Close() error {
	if handle != nil {
		w := handle
		handle = nil
		return w.Close()
	}
	return nil
}

// SetPinPWM enables and sets the PWM duty on a GPIO output pin.
//
// duty must be [0, 1].
//
// See https://github.com/sarfata/pi-blaster for more details. It relies on
// pi-blaster being installed and enabled.
func SetPinPWM(pin Pin, duty float32) error {
	if duty < 0 || duty > 1 {
		return fmt.Errorf("duty %f is invalid for blaster", duty)
	}
	err := open()
	if err == nil {
		_, err = io.WriteString(handle, fmt.Sprintf("%d=%f\n", pin, duty))
	}
	return err
}

// ReleasePin releases a GPIO output and leave it floating.
//
// This function must be called on process exit for each activated pin
// otherwise the pin will stay in the state.
func ReleasePin(pin Pin) error {
	err := open()
	if err == nil {
		_, err = io.WriteString(handle, fmt.Sprintf("release %d\n", pin))
	}
	return err
}

//

var handle io.WriteCloser

func open() error {
	if handle == nil {
		f, err := os.OpenFile("/dev/pi-blaster", os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return err
		}
		handle = f
	}
	return nil
}
