// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package anim1d

import (
	"fmt"
	"image"
	"image/color"
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

	x := float32(sinceStart) / float32(e.Duration)
	e.Out.NextFrame(e.buf, sinceStart+e.Offset)
	t := cubicBezier(0, 0, 0.58, 1, x)
	for i := range pixels {
		c := e.buf[i]
		t2 := 1 - t
		pixels[i].R = FloatToUint8(float32(pixels[i].R)*t + float32(c.R)*t2)
		pixels[i].G = FloatToUint8(float32(pixels[i].G)*t + float32(c.G)*t2)
		pixels[i].B = FloatToUint8(float32(pixels[i].B)*t + float32(c.B)*t2)
		pixels[i].A = FloatToUint8(float32(pixels[i].A)*t + float32(c.A)*t2)
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
	Weights  []float32 // In theory Sum(Weights) should be 1 but it doesn't need to. For example, mixing a night sky will likely have all of the Weights set to 1.
	bufs     [][]color.NRGBA
}

func (m *Mixer) NextFrame(pixels []color.NRGBA, sinceStart time.Duration) {
	if len(m.Patterns) != len(m.Weights) {
		panic(fmt.Errorf("len(Patterns) (%d) != len(Weights) (%d)", len(m.Patterns), len(m.Weights)))
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
		var r, g, b, a float32
		for j := range m.bufs {
			c := m.bufs[j][i]
			w := m.Weights[j]
			r += float32(c.R) * w
			g += float32(c.G) * w
			b += float32(c.B) * w
			a += float32(c.A) * w
		}
		pixels[i].R = FloatToUint8(r)
		pixels[i].G = FloatToUint8(g)
		pixels[i].B = FloatToUint8(b)
		pixels[i].A = FloatToUint8(a)
	}
}

// Interpolate creates a N times larger striped then scale down.
//
// This is useful to create smoother animations or scale down images for
// example. It always use Lanczos3.
type Interpolate struct {
	Child Pattern
	X     float32
	buf   []color.NRGBA
	img   image.NRGBA
}

func (i *Interpolate) NextFrame(pixels []color.NRGBA, sinceStart time.Duration) {
	l := int(ceil(i.X * float32(len(pixels))))
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
