// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"image"
	"image/color"
	"math"
	"time"

	"github.com/nfnt/resize"
)

// EaseOut changes from In to Out over time using the same Bezier curve than
// CSS3 transition "ease-out"; cubic-bezier(0, 0, 0.58, 1).
type EaseOut struct {
	In, Out  Pattern
	Duration time.Duration
	Offset   time.Duration
	buf      []color.NRGBA
}

func (e *EaseOut) NextFrame(pixels []color.NRGBA, sinceStart time.Duration) {
	e.In.NextFrame(pixels, sinceStart)
	if sinceStart >= e.Duration {
		return
	}
	if len(e.buf) != len(pixels) {
		e.buf = make([]color.NRGBA, len(pixels))
	} else {
		for i := range e.buf {
			e.buf[i] = color.NRGBA{}
		}
	}

	x := float64(sinceStart) / float64(e.Duration)
	e.Out.NextFrame(e.buf, sinceStart+e.Offset)
	t := cubicBezier(0, 0, 0.58, 1, x)
	for i := range pixels {
		c := e.buf[i]
		t2 := 1 - t
		pixels[i].R = floatToUint8(float64(pixels[i].R)*t + float64(c.R)*t2)
		pixels[i].G = floatToUint8(float64(pixels[i].G)*t + float64(c.G)*t2)
		pixels[i].B = floatToUint8(float64(pixels[i].B)*t + float64(c.B)*t2)
		pixels[i].A = floatToUint8(float64(pixels[i].A)*t + float64(c.A)*t2)
	}
}

// Subset draws a subpart of a strip, not touching the rest.
type Subset struct {
	Child Pattern
	Start int
	End   int
}

func (s *Subset) NextFrame(pixels []color.NRGBA, sinceStart time.Duration) {
	s.Child.NextFrame(pixels[s.Start:s.End], sinceStart)
}

// Mixer is a generic mixer that merges the output from multiple patterns.
type Mixer struct {
	Patterns []Pattern
	Weights  []float64 // In theory Sum(Weights) should be 1 but it doesn't need to.
	bufs     [][]color.NRGBA
}

func (m *Mixer) NextFrame(pixels []color.NRGBA, sinceStart time.Duration) {
	if len(m.Patterns) != len(m.Weights) {
		panic("Invalid Weights")
	}
	if len(m.bufs) != len(m.Patterns) {
		m.bufs = make([][]color.NRGBA, len(m.Patterns))
	}
	for i := range m.bufs {
		if len(m.bufs[i]) != len(pixels) {
			m.bufs[i] = make([]color.NRGBA, len(pixels))
		}
	}

	// Draw each pattern.
	for i := range m.Patterns {
		b := m.bufs[i]
		for j := range b {
			b[j] = color.NRGBA{}
		}
		m.Patterns[i].NextFrame(b, sinceStart)
	}

	// Merge patterns.
	for i := range pixels {
		var r, g, b, a float64
		for j := range m.bufs {
			c := m.bufs[j][i]
			w := m.Weights[j]
			r += float64(c.R) * w
			g += float64(c.G) * w
			b += float64(c.B) * w
			a += float64(c.A) * w
		}
		pixels[i].R = floatToUint8(r)
		pixels[i].G = floatToUint8(g)
		pixels[i].B = floatToUint8(b)
		pixels[i].A = floatToUint8(a)
	}
}

// Interpolate creates a N times larger striped then scale down. This is useful
// to create smoother animations.
type Interpolate struct {
	Child Pattern
	X     float64
	buf   []color.NRGBA
	img   image.NRGBA
}

func (i *Interpolate) NextFrame(pixels []color.NRGBA, sinceStart time.Duration) {
	l := int(math.Ceil(i.X * float64(len(pixels))))
	if l == 0 {
		return
	}
	if len(i.buf) != l {
		i.buf = make([]color.NRGBA, l)
	}
	for j := range i.buf {
		i.buf[j] = color.NRGBA{}
	}
	i.Child.NextFrame(i.buf, sinceStart)
	// TODO(maruel): Find a way to not have to double-buffer, e.g. alias i.buf
	// and i.img.Pix.
	if len(i.img.Pix) != 4*l {
		i.img = *image.NewNRGBA(image.Rect(0, 0, l, 1))
	}
	for j := range i.buf {
		i.img.SetNRGBA(j, 0, i.buf[j])
	}
	// TODO(maruel): Switch to code that doesn't allocate memory and doesn't
	// split the image to do concurrent processing. It's probably 10x slower than
	// it needs to be, and this is a concern on a rPi1.
	n := resize.Resize(uint(len(pixels)), 1, &i.img, resize.Lanczos3).(*image.NRGBA)
	for j := range pixels {
		pixels[j] = n.NRGBAAt(j, 0)
	}
}

// Skip creates a N times larger striped then strip down by omitting pixels.
// This creates an illusion of seeing through a screen window. Can be coupled
// with Interpolate to give a fuzzy view.
type Skip struct {
	Child Pattern
	X     int
	buf   []color.NRGBA
}

func (s *Skip) NextFrame(pixels []color.NRGBA, sinceStart time.Duration) {
	if len(s.buf) != s.X*len(pixels) {
		s.buf = make([]color.NRGBA, s.X*len(pixels))
	}
	for i := range s.buf {
		s.buf[i] = color.NRGBA{}
	}
	s.Child.NextFrame(s.buf, sinceStart)
	for j := range pixels {
		pixels[j] = s.buf[j*s.X]
	}
}

// Limit the maximum intensity. Does this by scaling the alpha channel.
type IntensityLimiter struct {
	Child Pattern
	Max   int // Maximum value between 0 (off) to 255 (full intensity).
}

func (i *IntensityLimiter) NextFrame(pixels []color.NRGBA, sinceStart time.Duration) {
	i.Child.NextFrame(pixels, sinceStart)
	for j := range pixels {
		pixels[j].A = uint8((int(pixels[j].A) + i.Max - 1) * 255 / i.Max)
	}
}

//

func roundF(x float64) float64 {
	if x < 0 {
		return math.Ceil(x - 0.5)
	}
	return math.Floor(x + 0.5)
}

func floatToUint8(x float64) uint8 {
	if x >= 255. {
		return 255
	}
	if x <= 0. {
		return 0
	}
	return uint8(roundF(x))
}

// cubicBezier returns [0, 1] for input `t` based on the cubic bezier curve
// (x0,y0), (x1, y1).
// Inspired by https://github.com/golang/mobile/blob/master/exp/sprite/clock/tween.go.
func cubicBezier(x0, y0, x1, y1, x float64) float64 {
	t := x
	for i := 0; i < 5; i++ {
		t2 := t * t
		t3 := t2 * t
		d := 1 - t
		d2 := d * d

		nx := 3*d2*t*x0 + 3*d*t2*x1 + t3
		dxdt := 3*d2*x0 + 6*d*t*(x1-x0) + 3*t2*(1-x1)
		if dxdt == 0 {
			break
		}

		t -= (nx - x) / dxdt
		if t <= 0 || t >= 1 {
			break
		}
	}
	if t < 0 {
		t = 0
	}
	if t > 1 {
		t = 1
	}

	// Solve for y using t.
	t2 := t * t
	t3 := t2 * t
	d := 1 - t
	d2 := d * d
	y := 3*d2*t*y0 + 3*d*t2*y1 + t3

	return y
}
