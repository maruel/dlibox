// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// playing is a small app to play with the pins, nothing more. You are not
// expected to use it as-is.
package main

import (
	"fmt"
	"image/color"
	"math"
	"os"
	"time"

	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/plotutil"
	"github.com/gonum/plot/vg"
	"github.com/maruel/dotstar"
	"github.com/stianeikeland/go-rpio"
)

func set(c color.NRGBA) {
	dotstar.SetPinPWM(dotstar.GPIO4, float32(c.R)/255.)
	dotstar.SetPinPWM(dotstar.GPIO22, float32(c.G)/255.)
	dotstar.SetPinPWM(dotstar.GPIO18, float32(c.B)/255.)
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
		dotstar.SetPinPWM(dotstar.GPIO17, float32(j)/255.)
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
		intensity := 1. - float32(i)/float32(len(pixels)-1)
		pixels[i] = color.NRGBA{
			uint8((float32(a.R)*intensity + float32(b.R)*(1-intensity))),
			uint8((float32(a.G)*intensity + float32(b.G)*(1-intensity))),
			uint8((float32(a.B)*intensity + float32(b.B)*(1-intensity))),
			uint8((float32(a.A)*intensity + float32(b.A)*(1-intensity))),
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

const maxIn = float32(0xFFFF)
const maxOut = float32(0x1EE1)
const lowCut = 30 * 255
const lowCutf = float32(lowCut)
const kin = lowCutf / maxIn
const kout = lowCutf / 255. / maxOut

func processRampCubeRoot(l uint32) uint32 {
	return uint32(float32(math.Pow(float64(float32(l)/maxIn), 3))*maxOut + 0.4)
}

func processRampFullRanges(l uint32) uint32 {
	// Linear [0->0] to [30*255->30].
	if l < lowCut {
		return uint32(float32(l)/255. + 0.4)
	}
	// Make sure the line cuts at lowCut starting with y==lowCut equals fY/255.
	// Put l1 on adapted linear basis [0, 1].
	l1 := (float32(l)/maxIn - kin) / (1. - kin)
	y := float32(math.Pow(float64(l1), 3))
	y2 := (y/(1+kout) + kout) * maxOut
	return uint32(y2 + 0.4)
}

func processRamp(l uint32) uint32 {
	// Linear [0->0] to [30*255->30].
	if l < lowCut {
		return uint32(float32(l)/255. + 0.4)
	}
	// Range [lowCut/maxIn, 1]
	y := float32(math.Pow(float64(float32(l)/maxIn), 3))
	// Range [(lowCut/maxIn)^3, 1]. We need to realign to [lowCut/255, 1]
	klow := float32(math.Pow(float64(lowCut/maxIn), 3)) + (lowCutf/255.+10.)/maxIn
	y2 := ((y + klow) / (1 + klow)) * maxOut
	return uint32(y2 + 0.4)
}

func doPlot() error {
	p, err := plot.New()
	if err != nil {
		return err
	}

	p.Title.Text = "Plotutil example"
	p.X.Label.Text = "Input range"
	p.Y.Label.Text = "Effective PWM"
	pts := make(plotter.XYs, 70*256)
	for i := range pts {
		pts[i].X = float64(i)
		pts[i].Y = float64(processRamp(uint32(i)))
	}
	/*
		pts := make(plotter.XYs, 0xffff+1)
		for i := range pts {
			pts[i].X = float64(i)
			pts[i].Y = float64(processRamp(uint32(i)))
		}
	*/
	if err = plotutil.AddLinePoints(p, "PWM on 0x1EE1 cycle", pts); err != nil {
		return err
	}

	return p.Save(32*vg.Inch, 32*vg.Inch, "points.png")
}

func mainImpl() error {
	return doPlot()
	//return doGPIO()
	return writeColorsRaw()
}

func main() {
	if err := mainImpl(); err != nil {
		fmt.Fprintf(os.Stderr, "playing\n: %s.\n", err)
		os.Exit(1)
	}
}
