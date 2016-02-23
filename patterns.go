// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"bytes"
	"image/color"
	"image/png"
	"time"
)

// Patterns is a map of nice predefined patterns.
var Patterns map[string]Pattern

func init() {
	red := color.NRGBA{255, 0, 0, 255}
	white := color.NRGBA{255, 255, 255, 255}
	Patterns = map[string]Pattern{
		"canne": &Repeated{[]color.NRGBA{red, red, red, red, white, white, white, white}, 6},
		"comète": &PingPong{
			Trail: []color.NRGBA{
				{0xff, 0xff, 255, 255},
				{0xD7, 0xD7, 255, 255},
				{0xAF, 0xAF, 255, 255},
				{0x87, 0x87, 255, 255},
				{0x5F, 0x5F, 255, 255},
			},
			MovesPerSec: 30,
		},
		"étoile floue":   LoadAnimate("étoile floue.png", 16*time.Millisecond, false),
		"K2000":          &PingPong{K2000, color.NRGBA{0, 0, 0, 255}, 30},
		"glow":           &Glow{[]color.NRGBA{{255, 255, 255, 255}, {0, 128, 0, 255}}, 1},
		"glow gris":      &Glow{[]color.NRGBA{{255, 255, 255, 255}, {}}, 0.33},
		"glow rainbow":   &Glow{RainbowColors, 1.},
		"pingpong":       &PingPong{Trail: []color.NRGBA{{255, 255, 255, 255}}, MovesPerSec: 30},
		"rainbow static": &Rainbow{},
		"red":            &StaticColor{color.NRGBA{255, 70, 70, 0}},
	}
}

// StaticColor shows a single color on all lights.
type StaticColor struct {
	C color.NRGBA
}

func (s *StaticColor) NextFrame(pixels []color.NRGBA, sinceStart time.Duration) {
	for i := range pixels {
		pixels[i] = s.C
	}
}

// Glow alternates betweens colors over time.
type Glow struct {
	Colors []color.NRGBA // Colors to cycle through. Use at least 2 colors.
	Hz     float64       // Color change rate per second. Should be below 0.1 for smooth change.
}

func (g *Glow) NextFrame(pixels []color.NRGBA, sinceStart time.Duration) {
	cycles := sinceStart.Seconds() * g.Hz
	baseIndex := int(cycles)
	// [0, 1]
	intensity := cycles - float64(baseIndex)
	//intensity := 0.5*math.Sin(sinceStart.Seconds()*g.Hz*math.Pi*2) + 0.5
	a := g.Colors[baseIndex%len(g.Colors)]
	b := g.Colors[(baseIndex+1)%len(g.Colors)]
	c := color.NRGBA{
		uint8((float64(a.R)*intensity + float64(b.R)*(1-intensity))),
		uint8((float64(a.G)*intensity + float64(b.G)*(1-intensity))),
		uint8((float64(a.B)*intensity + float64(b.B)*(1-intensity))),
		uint8((float64(a.A)*intensity + float64(b.A)*(1-intensity))),
	}
	for i := range pixels {
		pixels[i] = c
	}
}

// RainbowColors are approximate rainbox colors without alpha.
var RainbowColors = []color.NRGBA{
	{255, 0, 0, 255},
	{255, 127, 0, 255},
	{255, 255, 0, 255},
	{0, 255, 0, 255},
	{0, 0, 255, 255},
	{75, 0, 130, 255},
	{139, 0, 255, 255},
}

// PingPong shows a 'ball' with a trail that bounces from one side to
// the other.
//
// Can be used for a ball, a water wave or K2000 (Knight Rider) style light.
type PingPong struct {
	Trail       []color.NRGBA // [0] is the front pixel.
	Background  color.NRGBA
	MovesPerSec float64 // Expressed in number of light jumps per second.
}

// K2000 can be used with PingPong to look like Knight Rider.
// https://en.wikipedia.org/wiki/Knight_Rider_(1982_TV_series)
var K2000 = []color.NRGBA{
	{0xff, 0, 0, 255},
	{0xff, 0, 0, 255},
	{0xee, 0, 0, 255},
	{0xdd, 0, 0, 255},
	{0xcc, 0, 0, 255},
	{0xbb, 0, 0, 255},
	{0xaa, 0, 0, 255},
	{0x99, 0, 0, 255},
	{0x88, 0, 0, 255},
	{0x77, 0, 0, 255},
	{0x66, 0, 0, 255},
	{0x55, 0, 0, 255},
	{0x44, 0, 0, 255},
	{0x33, 0, 0, 255},
	{0x22, 0, 0, 255},
	{0x11, 0, 0, 255},
}

func (p *PingPong) NextFrame(pixels []color.NRGBA, sinceStart time.Duration) {
	for i := range pixels {
		pixels[i] = p.Background
	}
	// The last point of each extremity is only lit on one tick but every other
	// points are lit twice during a full cycle. This means the full cycle is
	// 2*(len(pixels)-1). For a 3 pixels line, the cycle is: x00, 0x0, 00x, 0x0.
	cycle := 2 * (len(pixels) - 1)
	moves := int(sinceStart.Seconds() * p.MovesPerSec)
	if moves < len(pixels) {
		// On the first cycle, the trail has not bounced yet.
		pos := moves % cycle
		for i := 0; i < len(p.Trail) && pos-i >= 0; i++ {
			pixels[pos-i] = p.Trail[i]
		}
	} else {
		for i := len(p.Trail) - 1; i >= 0; i-- {
			r := (moves - i) % cycle
			if r >= len(pixels) {
				r = cycle - r
			}
			pixels[r] = p.Trail[i]
		}
	}
}

// Animate represents an animatable looping frame.
type Animate struct {
	Frames        [][]color.NRGBA
	FrameDuration time.Duration
}

// LoadAnimate loads an Animate from a PNG file.
//
// Returns nil if the file can't be found. If vertical is true, rotate the
// image by 90°.
func LoadAnimate(name string, frameDuration time.Duration, vertical bool) *Animate {
	content := read(name)
	if content == nil {
		return nil
	}
	img, err := png.Decode(bytes.NewReader(content))
	if err != nil {
		return nil
	}
	bounds := img.Bounds()
	maxY := bounds.Max.Y
	maxX := bounds.Max.X
	if vertical {
		// Invert axes.
		maxY, maxX = maxX, maxY
	}
	buf := make([][]color.NRGBA, maxY)
	for y := 0; y < maxY; y++ {
		buf[y] = make([]color.NRGBA, maxX)
		for x := 0; x < maxX; x++ {
			if vertical {
				buf[y][x] = color.NRGBAModel.Convert(img.At(y, x)).(color.NRGBA)
			} else {
				buf[y][x] = color.NRGBAModel.Convert(img.At(x, y)).(color.NRGBA)
			}
		}
	}
	return &Animate{buf, frameDuration}
}

func (a *Animate) NextFrame(pixels []color.NRGBA, sinceStart time.Duration) {
	copy(pixels, a.Frames[int(sinceStart/a.FrameDuration)%len(a.Frames)])
	for i := len(a.Frames[0]); i < len(pixels); i++ {
		pixels[i] = color.NRGBA{}
	}
}

// MakeRainbow returns rainbow colors including alpha.
type Rainbow struct {
}

func (r *Rainbow) NextFrame(pixels []color.NRGBA, sinceStart time.Duration) {
	// Add buffer both before and after.
	start := 370.
	end := 790.
	step := (end - start) / float64(len(pixels)-1)
	for i := range pixels {
		pixels[i] = waveLength2RGB(start + step*float64(i))
	}
}

// waveLengthToRGB returns a color over a rainbow, including alpha.
//
// This code was inspired by public domain code on the internet.
func waveLength2RGB(w float64) (c color.NRGBA) {
	switch {
	case 380 <= w && w < 440:
		c.R = byte(255. * (440 - w) / (440 - 380))
		c.B = 255
	case 440 <= w && w < 490:
		c.G = byte(255. * (w - 440) / (490 - 440))
		c.B = 255
	case 490 <= w && w < 510:
		c.G = 255
		c.B = byte(255. * (510 - w) / (510 - 490))
	case 510 <= w && w < 580:
		c.R = byte(255. * (w - 510) / (580 - 510))
		c.G = 255
	case 580 <= w && w < 645:
		c.R = 255
		c.G = byte(255. * (645 - w) / (645 - 580))
	case 645 <= w && w < 781:
		c.R = 255
	}
	switch {
	case 380 <= w && w < 420:
		c.A = byte(255. * (0.3 + 0.7*(w-380)/(420-380)))
	case 420 <= w && w < 701:
		c.A = 255
	case 701 <= w && w < 781:
		c.A = byte(255.*0.3 + 0.7*(780-w)/(780-700))
	}
	return
}

/*
type point struct {
	star  int
	start time.Time
}

type Stars struct {
	Stars     []Animate
	Frequency float64 // Number of explosions by second.
	points    []point
}

func (s *Stars) NextFrame(pixels []color.NRGBA,  sinceStart time.Duration) {
	// random
	// animate.
}
*/

// Repeated prints a repeated pattern that can also cycle either way.
//
// Use negative to go left. Can be used for 'candy bar'.
type Repeated struct {
	Points      []color.NRGBA
	MovesPerSec float64 // Expressed in number of light jumps per second.
}

func (r *Repeated) NextFrame(pixels []color.NRGBA, sinceStart time.Duration) {
	offset := len(r.Points) - int(sinceStart.Seconds()*r.MovesPerSec)%len(r.Points)
	for i := range pixels {
		pixels[i] = r.Points[(i+offset)%len(r.Points)]
	}
}

//

/*

// Mixer merges the output from multiple patterns.
type Mixer struct {
	Patterns []Pattern
	Weight   []float64
	buf      []color.NRGBA
}

func (m *Mixer) NextFrame(pixels []color.NRGBA,  sinceStart time.Duration) {
	if len(m.buf) != len(pixels) {
		m.buf = make([]color.NRGBA, len(pixels))
	}
	delay := time.Hour
	for i := range pixels {
		pixels[i] = color.NRGBA{}
	}
	for i := range m.Patterns {
		for i := range m.buf {
			m.buf[i] = color.NRGBA{}
		}
		d := m.Patterns[i].NextFrame(m.buf,  sinceStart)
		if d < delay {
			delay = d
		}
		for i := range pixels {
			m.buf[i].A = uint8(m.Weight[i] * 255)
			r, g, b, _ := m.buf[i].RGBA()
			pixels[i].R += uint8(r)
			pixels[i].G += uint8(g)
			pixels[i].B += uint8(b)
		}
	}
}

func roundF(x float64) float64 {
	if x < 0 {
		return math.Ceil(x - 0.5)
	}
	return math.Floor(x + 0.5)
}
*/
