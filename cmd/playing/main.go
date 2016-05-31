// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// playing is a small app to play with the pins, nothing more. You are not
// expected to use it as-is.
package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"time"

	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/plotutil"
	"github.com/gonum/plot/vg"
	"github.com/maruel/dotstar/anim1d"
	"github.com/maruel/dotstar/apa102"
	"github.com/maruel/dotstar/rpi"
	"github.com/stianeikeland/go-rpio"
)

func getPatterns() error {
	var red = anim1d.Color{255, 0, 0, 255}
	var white = anim1d.Color{255, 255, 255, 255}
	rainbowColors := make([]anim1d.SPattern, len(anim1d.RainbowColors))
	for i, c := range anim1d.RainbowColors {
		rainbowColors[i].Pattern = &c
	}
	patterns := map[string]anim1d.Pattern{
		"Canne de Noël": &anim1d.Repeated{
			[]anim1d.Color{red, red, red, red, red, white, white, white, white, white},
			6,
		},
		"K2000": &anim1d.PingPong{anim1d.K2000Colors, anim1d.Color{0, 0, 0, 255}, 30},
		"Noir":  &anim1d.Color{},
		"Ping pong": &anim1d.PingPong{
			Trail:       []anim1d.Color{{255, 255, 255, 255}},
			MovesPerSec: 30,
		},
		"Rainbow cycle": &anim1d.Loop{
			Patterns:           rainbowColors,
			DurationShow:       1 * time.Second,
			DurationTransition: 10 * time.Second,
			Transition:         anim1d.TransitionEaseInOut,
		},
		"Rainbow static":       &anim1d.Rainbow{},
		"Étoiles cintillantes": &anim1d.NightStars{},
		"Ciel étoilé": &anim1d.Mixer{
			Patterns: []anim1d.SPattern{
				{&anim1d.Aurore{}},
				{&anim1d.NightStars{}},
				{&anim1d.WishingStar{}},
			},
			Weights: []float32{1, 1, 1},
		},
		"Aurores": &anim1d.Aurore{},
		// Transition from black to orange to white then to black.
		"Morning alarm": &anim1d.Transition{
			Out: anim1d.SPattern{
				&anim1d.Transition{
					Out: anim1d.SPattern{
						&anim1d.Transition{
							Out:        anim1d.SPattern{&anim1d.Color{}},
							In:         anim1d.SPattern{&anim1d.Color{255, 127, 0, 255}},
							Duration:   10 * time.Minute,
							Transition: anim1d.TransitionLinear,
						},
					},
					In:         anim1d.SPattern{&anim1d.Color{255, 255, 255, 255}},
					Offset:     10 * time.Minute,
					Duration:   10 * time.Minute,
					Transition: anim1d.TransitionLinear,
				},
			},
			In:         anim1d.SPattern{&anim1d.Color{}},
			Offset:     30 * time.Minute,
			Duration:   10 * time.Minute,
			Transition: anim1d.TransitionLinear,
		},
		// Test de couleurs:
		"Cycle RGB": &anim1d.Loop{
			Patterns: []anim1d.SPattern{
				{&anim1d.Color{255, 0, 0, 255}},
				{&anim1d.Color{0, 255, 0, 255}},
				{&anim1d.Color{0, 0, 255, 255}},
			},
			DurationShow:       1000 * time.Millisecond,
			DurationTransition: 1000 * time.Millisecond,
			Transition:         anim1d.TransitionEaseInOut,
		},
		"Dégradé":       &anim1d.Gradient{anim1d.Color{0, 0, 0, 255}, anim1d.Color{255, 255, 255, 255}},
		"Dégradé rouge": &anim1d.Gradient{anim1d.Color{0, 0, 0, 255}, anim1d.Color{255, 0, 0, 255}},
		"Dégradé vert":  &anim1d.Gradient{anim1d.Color{0, 0, 0, 255}, anim1d.Color{0, 255, 0, 255}},
		"Dégradé bleu":  &anim1d.Gradient{anim1d.Color{0, 0, 0, 255}, anim1d.Color{0, 0, 255, 255}},
	}

	for k, v := range patterns {
		p := anim1d.SPattern{v}
		b, err := json.Marshal(&p)
		if err != nil {
			return err
		}
		fmt.Printf("%q: %q\n", k, string(b))
	}
	return nil
}

func set(c anim1d.Color) {
	rpi.SetPinPWM(rpi.GPIO4, float32(c.R)/255.)
	rpi.SetPinPWM(rpi.GPIO22, float32(c.G)/255.)
	rpi.SetPinPWM(rpi.GPIO18, float32(c.B)/255.)
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
	c := []anim1d.Color{
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
		rpi.SetPinPWM(rpi.GPIO17, float32(j)/255.)
		j = (j + 1) % 256
	}
	return nil
}

func writeColors() error {
	d, err := apa102.MakeDotStar()
	if err != nil {
		return err
	}
	defer d.Close()

	pixels := make([]anim1d.Color, 150)
	a := anim1d.Color{0, 0, 0, 255}
	b := anim1d.Color{255, 0, 0, 255}
	for i := range pixels {
		intensity := 1. - float32(i)/float32(len(pixels)-1)
		pixels[i] = anim1d.Color{
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
	w, err := rpi.MakeSPI("", 10000000)
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
		s[4*i+0], s[4*i+1], s[4*i+2], s[4*i+3] = apa102.ColorToAPA102(anim1d.Color{r, 0, 0, 255})
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
	fmt.Printf("Yo\n")
	return getPatterns()
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
