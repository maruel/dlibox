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

type TransitionType string

const (
	// TODO(maruel): FIX!!
	TransitionEaseInOut TransitionType = "easeinout" // Use the same Bezier curve than CSS3 transition "ease-in-out"; cubic-bezier(0, XXX, 0.58, 1).
	TransitionEaseOut   TransitionType = "easeout"   // Use the same Bezier curve than CSS3 transition "ease-out"; cubic-bezier(0, 0, 0.58, 1). It is the default.
	TransitionLinear    TransitionType = "linear"
)

func (t TransitionType) scale(intensity float32) float32 {
	if t == TransitionLinear {
		return intensity
	} else if t == TransitionEaseInOut {
		// TODO(maruel): FIX!!
		return cubicBezier(0, 0, 0.58, 1, intensity)
	}
	return cubicBezier(0, 0, 0.58, 1, intensity)
}

type ScalingType string

const (
	ScalingNearest  ScalingType = "nearest"
	ScalingBilinear ScalingType = "bilinear"
	ScalingLanczos3 ScalingType = "lanczos3"
)

// Transition changes from In to Out over time.
//
// In gets sinceStart that is subtracted by Offset.
type Transition struct {
	Out        SPattern       // Old pattern that is disappearing
	In         SPattern       // New pattern to show
	Offset     time.Duration  // Offset at which the transiton from Out->In starts
	Duration   time.Duration  // Duration of the transition while both are rendered
	Transition TransitionType // Type of transition, defaults to EaseOut if not set
	buf        buffer
}

func (t *Transition) NextFrame(pixels []color.NRGBA, sinceStart time.Duration) {
	if sinceStart <= t.Offset {
		// Before transition.
		if t.Out.Pattern != nil {
			t.Out.NextFrame(pixels, sinceStart)
		}
		return
	}
	if t.In.Pattern != nil {
		t.In.NextFrame(pixels, sinceStart-t.Offset)
	}
	if sinceStart >= t.Offset+t.Duration {
		// After transition.
		t.buf = nil
		return
	}
	t.buf.reset(len(pixels))

	// TODO(maruel): Add lateral animation and others.
	if t.Out.Pattern != nil {
		t.Out.NextFrame(t.buf, sinceStart)
	}
	intensity := t.Transition.scale(float32(sinceStart-t.Offset) / float32(t.Duration))
	mix(intensity, pixels, t.buf)
}

// Loop rotates between all the animations.
//
// Display starts with one DurationShow for Patterns[0], then starts looping.
// sinceStart is not modified so it's like as all animations continued
// animating behind.
type Loop struct {
	Patterns           []SPattern
	DurationShow       time.Duration  // Duration for each pattern to be shown as pure
	DurationTransition time.Duration  // Duration of the transition between two patterns
	Transition         TransitionType // Type of transition, defaults to EaseOut if not set
	buf                buffer
}

func (l *Loop) NextFrame(pixels []color.NRGBA, sinceStart time.Duration) {
	l.buf.reset(len(pixels))
	ds := float32(l.DurationShow.Seconds())
	dt := float32(l.DurationTransition.Seconds())
	cycleDuration := ds + dt
	cycles := float32(sinceStart.Seconds()) / cycleDuration
	baseIndex := int(cycles)
	lp := len(l.Patterns)
	if lp == 0 {
		return
	}
	a := l.Patterns[baseIndex%lp]
	a.NextFrame(pixels, sinceStart)
	// [0, 1[
	delta := (cycles - float32(baseIndex))
	offset := delta * cycleDuration
	if offset <= ds {
		return
	}
	b := l.Patterns[(baseIndex+1)%lp]
	// ]0, 1[
	intensity := 1. - (offset-ds)/dt
	// TODO(maruel): Add lateral animation and others.
	b.NextFrame(l.buf, sinceStart)
	mix(l.Transition.scale(intensity), pixels, l.buf)
}

// Crop draws a subset of a strip, not touching the rest.
type Crop struct {
	Child  SPattern
	Start  int // Starting pixels to skip
	Length int // Length of the pixels to affect
}

func (s *Crop) NextFrame(pixels []color.NRGBA, sinceStart time.Duration) {
	if s.Child.Pattern != nil {
		s.Child.NextFrame(pixels[s.Start:s.Length], sinceStart)
	}
}

// Mixer is a generic mixer that merges the output from multiple patterns.
//
// It doesn't animate.
type Mixer struct {
	Patterns []SPattern
	Weights  []float32 // In theory Sum(Weights) should be 1 but it doesn't need to. For example, mixing a night sky will likely have all of the Weights set to 1.
	bufs     []buffer
}

func (m *Mixer) NextFrame(pixels []color.NRGBA, sinceStart time.Duration) {
	if len(m.Patterns) != len(m.Weights) {
		panic(fmt.Errorf("len(Patterns) (%d) != len(Weights) (%d)", len(m.Patterns), len(m.Weights)))
	}
	if len(m.bufs) != len(m.Patterns) {
		m.bufs = make([]buffer, len(m.Patterns))
	}
	for i := range m.bufs {
		m.bufs[i].reset(len(pixels))
	}

	// Draw each pattern.
	for i := range m.Patterns {
		m.Patterns[i].NextFrame(m.bufs[i], sinceStart)
	}

	// Merge patterns.
	for i := range pixels {
		// TODO(maruel): Averaging colors in RGB space looks like hell.
		var r, g, b, a float32
		for j := range m.bufs {
			c := m.bufs[j][i]
			w := m.Weights[j]
			// TODO(maruel): Pre-multiply alpha.
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

// Scale adapts a larger or smaller patterns to the Strip size
//
// This is useful to create smoother animations or scale down images.
type Scale struct {
	Child  SPattern
	Scale  ScalingType // Defaults to ScalingLanczos3
	Length int         // A buffer of this length will be provided to Child and will be scaled to the actual pixels length
	Ratio  float32     // The scaling ratio to use, <1 means smaller, >1 means larger. Only one of Length or Ratio can be used
	buf    buffer
	img    image.NRGBA
}

func (s *Scale) NextFrame(pixels []color.NRGBA, sinceStart time.Duration) {
	if s.Child.Pattern == nil {
		return
	}
	l := s.Length
	r := s.Ratio
	if l == 0 {
		l = int(ceil(s.Ratio * float32(len(pixels))))
	} else {
		r = float32(l) / float32(len(pixels))
	}
	s.buf.reset(l)
	s.Child.NextFrame(s.buf, sinceStart)
	if s.Scale == ScalingNearest {
		for i := range pixels {
			pixels[i] = s.buf[int(roundF(float32(i)*r))]
		}
	} else {
		// TODO(maruel): Find a way to not have to double-buffer, e.g. alias s.buf
		// and s.img.Pix.
		if len(s.img.Pix) != 4*l {
			s.img = *image.NewNRGBA(image.Rect(0, 0, l, 1))
		}
		for i := range s.buf {
			s.img.SetNRGBA(i, 0, s.buf[i])
		}
		// TODO(maruel): Switch to code that doesn't allocate memory and doesn't
		// split the image to do concurrent processing. It's probably 10x slower
		// than it needs to be, and this is a concern on a rPi1.
		scale := resize.Lanczos3
		if s.Scale == ScalingBilinear {
			scale = resize.Bilinear
		}
		n := resize.Resize(uint(len(pixels)), 1, &s.img, scale).(*image.NRGBA)
		for i := range pixels {
			pixels[i] = n.NRGBAAt(i, 0)
		}
	}
}

// Private

func mix(intensity float32, a, b []color.NRGBA) {
	for i := range a {
		c := b[i]
		t2 := 1 - intensity
		// TODO(maruel): Averaging colors in RGB space looks like hell.
		a[i].R = FloatToUint8(float32(a[i].R)*intensity + float32(c.R)*t2)
		a[i].G = FloatToUint8(float32(a[i].G)*intensity + float32(c.G)*t2)
		a[i].B = FloatToUint8(float32(a[i].B)*intensity + float32(c.B)*t2)
		a[i].A = FloatToUint8(float32(a[i].A)*intensity + float32(c.A)*t2)
	}
}

type buffer []color.NRGBA

func (b *buffer) reset(l int) {
	if len(*b) != l {
		*b = make([]color.NRGBA, l)
	} else {
		s := *b
		for i := range s {
			s[i] = color.NRGBA{}
		}
	}
}
