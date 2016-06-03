// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package anim1d

import (
	"math/rand"
	"time"
)

// Color shows a single color on all lights. It knows how to renders itself
// into a frame.
//
// If you want a single dot, use a Frame of length one.
type Color struct {
	R, G, B uint8
}

// Add adds two color together with saturation, mixing according to the alpha
// channel.
func (c *Color) Add(d Color) {
	r := uint16(c.R) + uint16(d.R)
	if r > 255 {
		r = 255
	}
	c.R = uint8(r)
	g := uint16(c.G) + uint16(d.G)
	if g > 255 {
		g = 255
	}
	c.G = uint8(g)
	b := uint16(c.B) + uint16(d.B)
	if b > 255 {
		b = 255
	}
	c.B = uint8(b)
}

// Mix blends the second color with the first.
//
// gradient 0 means pure 'c', gradient 255 means pure 'd'.
func (c *Color) Mix(d Color, gradient uint8) {
	grad := uint16(gradient)
	grad1 := 255 - grad
	// unit test confirms the values cannot overflow.
	c.R = uint8(((uint16(c.R)+1)*grad1 + (uint16(d.R)+1)*grad) >> 8)
	c.G = uint8(((uint16(c.G)+1)*grad1 + (uint16(d.G)+1)*grad) >> 8)
	c.B = uint8(((uint16(c.B)+1)*grad1 + (uint16(d.B)+1)*grad) >> 8)
}

func (c *Color) NextFrame(pixels Frame, sinceStart time.Duration) {
	for i := range pixels {
		pixels[i] = *c
	}
}

func (c *Color) NativeDuration(pixels int) time.Duration {
	return 0
}

// Frame is a strip of colors. It knows how to renders itself into a frame
// (which is recursive).
type Frame []Color

func (f Frame) NextFrame(pixels Frame, sinceStart time.Duration) {
	copy(pixels, f)
}

func (f Frame) NativeDuration(pixels int) time.Duration {
	return 0
}

// Mix blends the second frame with the first.
//
// gradient 0 means pure 'f', gradient 255 means pure 'b'.
func (f Frame) Mix(b Frame, gradient uint8) {
	for i := range f {
		f[i].Mix(b[i], gradient)
	}
}

// reset() always resets the buffer to black.
func (f *Frame) reset(l int) {
	if len(*f) != l {
		*f = make(Frame, l)
	} else {
		s := *f
		for i := range s {
			s[i] = Color{}
		}
	}
}

// MakeRainbow renders rainbow colors.
type Rainbow struct {
}

func (r *Rainbow) NextFrame(pixels Frame, sinceStart time.Duration) {
	start := 380
	end := 781
	/*
		step := (end - start) / float32(len(pixels)-1)
			for i := range pixels {
				// TODO(maruel): Use log scale.
				pixels[i] = waveLength2RGB(start + step*float32(i))
			}
	*/

	// TODO(maruel): Still too much red not enough pink.
	delta := end - start
	scale := logn(2)
	step := 1. / float32(len(pixels))
	for i := range pixels {
		j := log1p(float32(len(pixels)-i-1)*step) / scale
		pixels[i] = waveLength2RGB(int(float32(start) + float32(delta)*(1-j)))
	}
}

func (r *Rainbow) NativeDuration(pixels int) time.Duration {
	return 0
}

// waveLengthToRGB returns a color over a rainbow.
//
// This code was inspired by public domain code on the internet.
func waveLength2RGB(w int) (c Color) {
	switch {
	case w < 380:
	case w < 420:
		// Red peaks at 1/3 at 420.
		c.R = uint8(196 - (170*(440-w))/(440-380))
		c.B = uint8(26 + (229*(w-380))/(420-380))
	case w < 440:
		c.R = uint8((0x89 * (440 - w)) / (440 - 420))
		c.B = 255
	case w < 490:
		c.G = uint8((255 * (w - 440)) / (490 - 440))
		c.B = 255
	case w < 510:
		c.G = 255
		c.B = uint8((255 * (510 - w)) / (510 - 490))
	case w < 580:
		c.R = uint8((255 * (w - 510)) / (580 - 510))
		c.G = 255
	case w < 645:
		c.R = 255
		c.G = uint8((255 * (645 - w)) / (645 - 580))
	case w < 700:
		c.R = 255
	case w < 781:
		c.R = uint8(26 + (229*(780-w))/(780-700))
	default:
	}
	return
}

// Repeated repeats a Frame to fill the pixels.
type Repeated struct {
	Frame Frame
}

func (r *Repeated) NextFrame(pixels Frame, sinceStart time.Duration) {
	if len(pixels) == 0 || len(r.Frame) == 0 {
		return
	}
	for i := 0; i < len(pixels); i += len(r.Frame) {
		copy(pixels[i:], r.Frame)
	}
}

func (r *Repeated) NativeDuration(pixels int) time.Duration {
	return 0
}

type point struct {
	star  int
	start time.Time
}

// NightSky has:
//    - Stars
//    - WishingStar
//    - Aurores
//    - Super nova.
//    - Rotation de la terre?
//    - Station Internationale?
type NightSky struct {
	Stars     []Cycle
	Frequency float32 // Number of explosions by second.
	points    []point
}

func (c *NightSky) NextFrame(pixels Frame, sinceStart time.Duration) {
	// random
	// animate.
}

// Aurore commence lentement, se transforme lentement et éventuellement
// disparait.
type Aurore struct {
}

func (a *Aurore) NextFrame(pixels Frame, sinceStart time.Duration) {
	// TODO(maruel): Redo.
	y := float32(sinceStart.Seconds()) * 10.
	for i := range pixels {
		x := float32(i)
		//a := 32 + 31*sin(x/(37.+15*cos(y/74)))*cos(y/(31+11*sin(x/57)))
		b := (32 + 31*(sin(hypot(200-y, 320-x)/16))) * (0.5 + 0.5*sin(y*0.1))
		pixels[i].R = 0
		//pixels[i].G = uint8(a + b)
		pixels[i].G = uint8(b)
		pixels[i].B = 0
	}
}

type NightStar struct {
	Intensity uint8
	Type      int
}

type NightStars struct {
	Stars []NightStar
	Seed  int // Change it to create a different pseudo-random animation.
	r     *rand.Rand
}

func (e *NightStars) NextFrame(pixels Frame, sinceStart time.Duration) {
	if e.r == nil {
		e.r = rand.New(rand.NewSource(int64(e.Seed)))
	}
	if len(e.Stars) != len(pixels) {
		e.Stars = make([]NightStar, len(pixels))
		for i := 0; i < len(pixels); {
			// Add a star. Decide it's relative position, intensity and type.
			// ExpFloat64() ?
			f := abs(3 * float32(e.r.NormFloat64()))
			if f < 1 {
				continue
			}
			i += int(roundF(f))
			if i >= len(pixels) {
				break
			}
			// e.r.Intn(255)
			intensity := abs(float32(e.r.NormFloat64()))
			if intensity > 255 {
				intensity = 0
			}
			e.Stars[i].Intensity = FloatToUint8(intensity)
		}
	}
	for i := range e.Stars {
		if j := e.Stars[i].Intensity; j != 0 {
			// TODO(maruel): Type, oscillation.
			if j != 0 {
				f := FloatToUint8(float32(e.r.NormFloat64())*4 + float32(j))
				pixels[i] = Color{f, f, f}
			}
		}
	}
}

// WishingStar draws a wishing star from time to time.
//
// It will only draw one star at a time. To increase the likelihood of getting
// many simultaneously, create multiple instances and use Mixer with Weights of
// 1.
type WishingStar struct {
	Duration     time.Duration // Average duration of a star.
	AverageDelay time.Duration // Average delay between each wishing star.
}

func (w *WishingStar) NextFrame(pixels Frame, sinceStart time.Duration) {
	/*
		// Create a deterministic replay by using the current number of
		// the wishing star as the seed for the current flow. Make it independent of
		// any other non-deterministic source.
		i := sinceStart / w.AverageDelay
		r := rand.New(rand.NewSource(int64(i)))
		// Always calculate things in the same order to keep the calculation
		// deterministic.
		startOffset := r.Int63()
		startPos := r.Int63()
		intensity := r.Int63()
		orientation := r.Intn(2)
		// Draw according to these parameters.
		// - Trail
		// - Observed speed based on orientation
	*/
}

// Gradient does a gradient between 2 colors as a static image.
//
// TODO(maruel): Support N colors at M positions.
// TODO(maruel): Blend arbitrary SPattern with different curves.
type Gradient struct {
	A, B       Color
	Transition TransitionType
}

func (g *Gradient) NextFrame(pixels Frame, sinceStart time.Duration) {
	l := len(pixels) - 1
	if l == 0 {
		pixels[0] = g.A
		pixels[0].Mix(g.B, FloatToUint8(255.*g.Transition.scale(0.5)))
		return
	}
	// TODO(maruel): Convert to integer calculation.
	max := float32(len(pixels) - 1)
	for i := range pixels {
		// [0, 1]
		intensity := float32(i) / max
		pixels[i] = g.A
		pixels[i].Mix(g.B, FloatToUint8(255.*g.Transition.scale(intensity)))
	}
}
