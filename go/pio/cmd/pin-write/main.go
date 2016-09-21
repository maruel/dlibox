// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// pin-write is a small app to write a pin.
package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/maruel/dlibox/go/pio/host"
	"github.com/maruel/dlibox/go/pio/host/hal/pins"
)

func mainImpl() error {
	if len(os.Args) != 3 {
		return errors.New("specify pin to write to and its level (0 or 1)")
	}
	pin, err := strconv.Atoi(os.Args[1])
	if err != nil {
		return err
	}
	p := pins.ByNumber(pin)
	if p == nil {
		return errors.New("invalid pin number")
	}
	if err := p.Out(); err != nil {
		return err
	}
	switch os.Args[2] {
	case "0":
		p.Set(host.Low)
	case "1":
		p.Set(host.High)
	default:
		return errors.New("specify level as 0 or 1")
	}
	return nil
}

func main() {
	if err := mainImpl(); err != nil {
		fmt.Fprintf(os.Stderr, "pin-write: %s.\n", err)
		os.Exit(1)
	}
}
