// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// i2c is a small app to communicate an i²c device.
package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/maruel/dlibox/go/rpi"
)

func mainImpl() error {
	if len(os.Args) < 4 {
		return errors.New("specify bus, address and byte(s) to write")
	}
	bus, err := strconv.Atoi(os.Args[1])
	if err != nil {
		return err
	}
	addr, err := strconv.Atoi(os.Args[2])
	if err != nil {
		return err
	}
	if addr < 0 || addr > 65535 {
		return errors.New("invalid address")
	}
	buf := make([]byte, 0, len(os.Args)-3)
	for _, a := range os.Args[3:] {
		b, err := strconv.Atoi(a)
		if err != nil {
			return err
		}
		if b < 0 || b > 255 {
			return errors.New("invalid byte")
		}
		buf = append(buf, byte(b))
	}

	i, err := rpi.MakeI2C(bus)
	if err != nil {
		return err
	}
	if err := i.Address(uint16(addr)); err != nil {
		return err
	}
	_, err = i.Write(buf)
	return err
}

func main() {
	if err := mainImpl(); err != nil {
		fmt.Fprintf(os.Stderr, "i2c: %s.\n", err)
		os.Exit(1)
	}
}
