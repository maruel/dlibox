// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// piblaster specific code. Not thread safe.

package rpi

import (
	"fmt"
	"io"
	"os"
)

// IsPiblaster returns true if the pin can be used as a PWM source via
// piblaster.
//
// https://github.com/sarfata/pi-blaster#options
//
// TODO(maruel): "dtoverlay=pwm" or "dtoverlay=pwm-2chan" works without having
// to install anything, albeit with less pins supported.
func (p Pin) IsPiblaster() bool {
	switch p {
	case GPIO4, GPIO17, GPIO18, GPIO21, GPIO22, GPIO23, GPIO24, GPIO25, GPIO27:
		// TODO(maruel): Detect at runtime if enabled. This requires upstreaming a
		// change to https://github.com/sarfata/pi-blaster/blob/master/pi-blaster.c
		return true
	default:
		return false
	}
}

// SetPWM enables and sets the PWM duty on a GPIO output pin via piblaster.
//
// duty must be [0, 1].
//
// See https://github.com/sarfata/pi-blaster for more details. It relies on
// pi-blaster being installed and enabled.
func (p Pin) SetPWM(duty float32) error {
	if duty < 0 || duty > 1 {
		return fmt.Errorf("duty %f is invalid for blaster", duty)
	}
	err := openPiblaster()
	if err == nil {
		_, err = io.WriteString(piblasterHandle, fmt.Sprintf("%d=%f\n", p, duty))
	}
	return err
}

// ReleasePWM releases a GPIO output and leave it floating.
//
// This function must be called on process exit for each activated pin
// otherwise the pin will stay in the state.
func (p Pin) ReleasePWM() error {
	err := openPiblaster()
	if err == nil {
		_, err = io.WriteString(piblasterHandle, fmt.Sprintf("release %d\n", p))
	}
	return err
}

//

var piblasterHandle io.WriteCloser

func openPiblaster() error {
	if piblasterHandle == nil {
		f, err := os.OpenFile("/dev/pi-blaster", os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return err
		}
		piblasterHandle = f
	}
	return nil
}

func closePiblaster() error {
	if piblasterHandle != nil {
		w := piblasterHandle
		piblasterHandle = nil
		return w.Close()
	}
	return nil
}
