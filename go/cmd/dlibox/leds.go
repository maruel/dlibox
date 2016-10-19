// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Packages the static files in a .go file.
//go:generate go run ../package/main.go -out static_files_gen.go ../../../web

// dlibox drives the dlibox LED strip on a Raspberry Pi. It runs a web server
// for remote control.
package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/maruel/dlibox/go/donotuse/conn/spi"
	"github.com/maruel/dlibox/go/donotuse/devices"
	"github.com/maruel/dlibox/go/donotuse/devices/apa102"
	"github.com/maruel/dlibox/go/donotuse/host"
	"github.com/maruel/dlibox/go/screen"
)

// initLEDs initializes the LED strip.
//
// TODO(maruel): This function is horrible.
func initLEDs(fake bool, config *APA102) (devices.Display, func(), []string, int, error) {
	if fake {
		// Output (terminal with ANSI codes or APA102).
		// Hardcode to 100 characters when using a terminal output.
		// TODO(maruel): Query the terminal and use its width.
		leds := screen.New(100)
		end := func() { os.Stdout.Write([]byte("\033[0m\n")) }
		// Use lower refresh rate too.
		return leds, end, []string{"fake=1"}, 30, nil
	}

	fps := 60
	if host.MaxSpeed() < 900000 || runtime.NumCPU() < 4 {
		// Use 30Hz on slower devices because it is too slow.
		fps = 30
	}
	spiBus, err := spi.New(-1, -1)
	if err != nil {
		return nil, nil, nil, 0, err
	}
	if err = spiBus.Speed(config.SPIspeed); err != nil {
		return nil, nil, nil, 0, err
	}
	end := func() { spiBus.Close() }
	leds, err := apa102.New(spiBus, config.NumberLights, 255, 6500)
	if err != nil {
		return nil, end, nil, 0, err
	}
	return leds, end, []string{fmt.Sprintf("APA102=%d", config.NumberLights)}, fps, nil
}
