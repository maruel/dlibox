// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// bme280 is a small app to read from a BME280.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/maruel/dlibox/go/bme280"
	"github.com/maruel/dlibox/go/rpi"
)

func mainImpl() error {
	bus := flag.Int("b", 1, "I²C bus to use")
	sample1x := flag.Bool("1", false, "sample at 1x")
	sample2x := flag.Bool("2", false, "sample at 2x")
	sample4x := flag.Bool("4", false, "sample at 4x")
	sample8x := flag.Bool("8", false, "sample at 8x")
	sample16x := flag.Bool("16", false, "sample at 16x")
	loop := flag.Bool("l", false, "loop every 100ms")
	verbose := flag.Bool("v", false, "verbose mode")
	flag.Parse()
	if !*verbose {
		log.SetOutput(ioutil.Discard)
	}
	log.SetFlags(log.Lmicroseconds)

	i2c, err := rpi.MakeI2C(*bus)
	if err != nil {
		return err
	}
	s := bme280.O4x
	if *sample1x {
		s = bme280.O1x
	} else if *sample2x {
		s = bme280.O2x
	} else if *sample4x {
		s = bme280.O4x
	} else if *sample8x {
		s = bme280.O8x
	} else if *sample16x {
		s = bme280.O16x
	}
	b, err := bme280.MakeBME280(i2c, s, s, s, bme280.S20ms, bme280.FOff)
	if err != nil {
		return err
	}
	for {
		t, p, h, err := b.Read()
		if err != nil {
			return err
		}
		fmt.Printf("%f°C %fkPa %f%%rH\n", t, p, h)
		if !*loop {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	return nil
}

func main() {
	if err := mainImpl(); err != nil {
		fmt.Fprintf(os.Stderr, "bme280: %s.\n", err)
		os.Exit(1)
	}
}
