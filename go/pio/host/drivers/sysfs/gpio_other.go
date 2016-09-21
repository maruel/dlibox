// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// +build !linux

package sysfs

import (
	"errors"
	"os"
)

// Init initializes GPIO sysfs handling code.
//
// Uses gpio sysfs as described at
// https://www.kernel.org/doc/Documentation/gpio/sysfs.txt
//
// GPIO sysfs is often the only way to do edge triggered interrupts. Doing this
// requires cooperation from a driver in the kernel.
//
// The main drawback of GPIO sysfs is that it doesn't expose internal pull
// resistor and it is much slower than using memory mapped hardware registers.
//
// Init returns an error no non-Linux OS.
func Init() error {
	return errors.New("gpio sysfs is not implemented on non-linux OSes")
}

//

type event struct{}

func (e event) wait(ep int) (int, error) {
	return 0, errors.New("unreachable code")
}

func (e event) makeEvent(f *os.File) (int, error) {
	return 0, errors.New("unreachable code")
}
