// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// playing is a small app to play with the pins, nothing more. You are not
// expected to use it as-is.
package main

import (
	"fmt"
	"image/color"
	"os"
	"time"

	"github.com/maruel/dotstar"
	"github.com/stianeikeland/go-rpio"
)

func set(c color.NRGBA) {
	dotstar.SetPinPWM(dotstar.GPIO4, float64(c.R)/255.)
	dotstar.SetPinPWM(dotstar.GPIO22, float64(c.G)/255.)
	dotstar.SetPinPWM(dotstar.GPIO18, float64(c.B)/255.)
}

func mainImpl() error {
	// Open and map memory to access gpio, check for errors
	if err := rpio.Open(); err != nil {
		return err
	}
	defer rpio.Close()

	pin := rpio.Pin(25)
	pin.Input()
	pin.PullUp()
	last := rpio.High
	c := []color.NRGBA{
		{255, 0, 0, 255},
		{0, 255, 0, 255},
		{0, 0, 255, 255},
		{255, 255, 0, 255},
	}
	i := 0
	set(c[i])
	j := 0
	for {
		if state := pin.Read(); state != last {
			last = state
			//fmt.Printf("%d\n", state)
			if state == rpio.Low {
				i = (i + 1) % len(c)
				set(c[i])
			}
		}
		time.Sleep(time.Millisecond)
		dotstar.SetPinPWM(dotstar.GPIO17, float64(j)/255.)
		j = (j + 1) % 256
	}
	return nil
}

func main() {
	if err := mainImpl(); err != nil {
		fmt.Fprintf(os.Stderr, "playing\n: %s.\n", err)
		os.Exit(1)
	}
}
