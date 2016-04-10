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

func doGPIO() error {
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

func writeColors() error {
	d, err := dotstar.MakeDotStar()
	if err != nil {
		return err
	}
	defer d.Close()

	pixels := make([]color.NRGBA, 150)
	a := color.NRGBA{0, 0, 0, 255}
	b := color.NRGBA{255, 0, 0, 255}
	for i := range pixels {
		intensity := 1. - float64(i)/float64(len(pixels)-1)
		pixels[i] = color.NRGBA{
			uint8((float64(a.R)*intensity + float64(b.R)*(1-intensity))),
			uint8((float64(a.G)*intensity + float64(b.G)*(1-intensity))),
			uint8((float64(a.B)*intensity + float64(b.B)*(1-intensity))),
			uint8((float64(a.A)*intensity + float64(b.A)*(1-intensity))),
		}
	}
	return d.Write(pixels)
}

func writeColorsRaw() error {
	// Uses the raw SPI protocol.
	w, err := dotstar.MakeSPI("", 10000000)
	if err != nil {
		return err
	}
	numLights := 150
	l := 4*(numLights+1) + numLights/2/8 + 1
	buf := make([]byte, l)
	// Set end frames right away.
	s := buf[4+4*numLights:]
	for i := range s {
		s[i] = 0xFF
	}
	// Start frame is all zeros. Just skip it.
	s = buf[4 : 4+4*numLights]
	for i := 0; i < 150; i++ {
		r := byte(255 - i)
		/*
				if i&1 == 0 {
					r = byte(i + 255 - 150)
				}
			r := byte(i)
		*/
		s[4*i+0], s[4*i+1], s[4*i+2], s[4*i+3] = dotstar.ColorToAPA102(color.NRGBA{r, 0, 0, 255})
	}
	/*
		i := 0
		s[4*i+0] = byte(0xE0 + 31)
		s[4*i+1] = byte(1)
		s[4*i+3] = byte(1)
		i++
		s[4*i+0] = byte(0xE0 + 1)
		s[4*i+1] = byte(31)
		s[4*i+3] = byte(31)
		i++
		s[4*i+0] = byte(0xE0 + 1)
		s[4*i+1] = byte(63)
		s[4*i+3] = byte(63)
		i++
		s[4*i+0] = byte(0xE0 + 31)
		s[4*i+1] = byte(1)
		s[4*i+3] = byte(1)
	*/
	/*
		for i := 0; i < numLights; i++ {
			// bBGR.
			s[4*i+0] = byte(0xE0 + 1)
			s[4*i+1] = byte(0)
			s[4*i+3] = byte(0)
		}
	*/
	_, err = w.Write(buf)
	return err

}

func mainImpl() error {
	//return doGPIO()
	return writeColorsRaw()
}

func main() {
	if err := mainImpl(); err != nil {
		fmt.Fprintf(os.Stderr, "playing\n: %s.\n", err)
		os.Exit(1)
	}
}
