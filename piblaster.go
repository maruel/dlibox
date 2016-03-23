// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package dotstar

import (
	"fmt"
	"os"
)

// SetPinPWM enables and sets the PWM duty on a GPIO output pin.
//
// duty must be [0, 1].
//
// See https://github.com/sarfata/pi-blaster for more details. It relies on
// pi-blaster being installed and enabled.
//
//  GPIO number   Pin in P1 header
//      4              P1-7
//      17             P1-11
//      18             P1-12
//      21             P1-13
//      22             P1-15
//      23             P1-16
//      24             P1-18
//      25             P1-22
func SetPinPWM(pin int, duty float64) error {
	if err := isPiBlasterPinValid(pin); err != nil {
		return err
	}
	if duty < 0 || duty > 1 {
		return fmt.Errorf("duty %f is invalid for blaster", duty)
	}
	return sendCommandToPiBlaster(fmt.Sprintf("%d=%f\n", pin, duty))
}

// ReleasePinPWM releases (disable) a GPIO output and leave it floating.
//
// This function must be called on process exit for each pin activated.
func ReleasePinPWM(pin int) error {
	if err := isPiBlasterPinValid(pin); err != nil {
		return err
	}
	return sendCommandToPiBlaster(fmt.Sprintf("release %d\n", pin))

}

// Whitelist the valid ports.
func isPiBlasterPinValid(pin int) error {
	if pin != 4 && pin != 17 && pin != 18 && !(21 <= pin && pin <= 25) {
		return fmt.Errorf("pin %d is invalid for blaster", pin)
	}
	return nil
}

func sendCommandToPiBlaster(command string) error {
	f, err := os.OpenFile("/dev/pi-blaster", os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(command)
	return err
}
