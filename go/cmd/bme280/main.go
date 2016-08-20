// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// bme280 is a small app to read from a BME280.
package main

import (
	"fmt"
	"os"

	"github.com/maruel/dlibox/go/bme280"
	"github.com/maruel/dlibox/go/rpi"
)

func mainImpl() error {
	i2c, err := rpi.MakeI2C(1)
	if err != nil {
		return err
	}
	b, err := bme280.MakeBME280(i2c)
	if err != nil {
		return err
	}
	fmt.Printf("%v", b.ChipID())
	return nil
}

func main() {
	if err := mainImpl(); err != nil {
		fmt.Fprintf(os.Stderr, "bme280: %s.\n", err)
		os.Exit(1)
	}
}
