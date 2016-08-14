// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// pin-read is a small app to read a pin.
package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/maruel/dlibox/go/rpi"
)

func mainImpl() error {
	pull := rpi.PullNoChange
	if len(os.Args) == 3 {
		i, err := strconv.Atoi(os.Args[2])
		if err != nil {
			return err
		}
		switch i {
		case 0:
			pull = rpi.Float
		case 1:
			pull = rpi.Down
		case 2:
			pull = rpi.Up
		default:
			return errors.New("specify pull value as 0 for float, 1 for down, 2 for up, or omit for no change")
		}
	} else if len(os.Args) != 2 {
		return errors.New("specify pin to read")
	}
	pin, err := strconv.Atoi(os.Args[1])
	if err != nil {
		return err
	}
	if pin > 53 || pin < 0 {
		return errors.New("specify pin between 0 and 53")
	}
	p := rpi.Pin(pin)
	p.In(pull)
	if p.Read() == rpi.Low {
		os.Stdout.Write([]byte{'0', '\n'})
	} else {
		os.Stdout.Write([]byte{'1', '\n'})
	}
	return nil
}

func main() {
	if err := mainImpl(); err != nil {
		fmt.Fprintf(os.Stderr, "pin-read: %s.\n", err)
		os.Exit(1)
	}
}
