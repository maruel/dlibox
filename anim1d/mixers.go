// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package anim1d

import (
	"fmt"
	"time"
)

// TransitionType models visually pleasing transitions.
//
// They are modeled against CSS transitions.
// https://www.w3.org/TR/web-animations/#scaling-using-a-cubic-bezier-curve
type TransitionType string

const (
	TransitionEase      TransitionType = "ease"
	TransitionEaseIn    TransitionType = "ease-in"
	TransitionEaseInOut TransitionType = "ease-in-out"
	TransitionEaseOut   TransitionType = "ease-out" // Recommended and default value.
	TransitionLinear    TransitionType = "linear"
)

// scale scales input [0, 1] to output [0, 1] using the transition requested.
func (t TransitionType) scale(intensity float32) float32 {
	// TODO(maruel): Add support for arbitrary cubic-bezier().
	// TODO(maruel): Map ease-* to cubic-bezier().
	// TODO(maruel): Add support for steps() which is pretty cool.
	switch t {
	case TransitionEase:
		return cubicBezier(0.25, 0.1, 0.25, 1, intensity)
	case TransitionEaseIn:
		return cubicBezier(0.42, 0, 1, 1, intensity)
	case TransitionEaseInOut:
		return cubicBezier(0.42, 0, 0.58, 1, intensity)
	case TransitionEaseOut, "":
		fallthrough
	default:
		return cubicBezier(0, 0, 0.58, 1, intensity)
	case TransitionLinear:
		return intensity
	}
}

// ScalingType specifies a way to scales a pixel strip.
type ScalingType string

const (
	ScalingNearestSkip ScalingType = "nearestskip" // Selects the nearest pixel but when upscaling, skips on missing pixels.
	ScalingNearest     ScalingType = "nearest"     // Selects the nearest pixel, gives a blocky view.
	ScalingLinear      ScalingType = "linear"      // Linear interpolation, recommended and default value.
	ScalingBilinear    ScalingType = "bilinear"    // Bilinear interpolation, usually overkill for 1D.
)

func (s ScalingType) scale(in, out Frame) {
	// Use integer operations as much as possible for reasonable performance.
	li := len(in)
	lo := len(out)
	if li == 0 || lo == 0 {
		return
	}
	switch s {
	case ScalingNearestSkip:
		if li < lo {
			// Do not touch skipped pixels.
			for i, p := range in {
				out[(i*lo+lo/2)/li] = p
			}
			return
		}
		fallthrough
	case ScalingNearest, ScalingLinear, ScalingBilinear, "":
		fallthrough
	default:
		for i := range out {
			out[i] = in[(i*li+li/2)/lo]
		}
		/*
			case ScalingLinear:
				for i := range out {
					x := (i*li + li/2) / lo
					c := in[x]
					c.Add(in[x+1])
					out[i] = c
				}
		*/
	}
}

// Transition changes from In to Out over time.
//
// In gets sinceStart that is subtracted by Offset.
type Transition struct {
	Out        SPattern       // Old pattern that is disappearing
	In         SPattern       // New pattern to show
	Offset     time.Duration  // Offset at which the transiton from Out->In starts
	Duration   time.Duration  // Duration of the transition while both are rendered
	Transition TransitionType // Type of transition, defaults to EaseOut if not set
	buf        Frame
}

func (t *Transition) NextFrame(pixels Frame, sinceStart time.Duration) {
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
	buf                Frame
}

func (l *Loop) NextFrame(pixels Frame, sinceStart time.Duration) {
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

func (s *Crop) NextFrame(pixels Frame, sinceStart time.Duration) {
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
	bufs     []Frame
}

func (m *Mixer) NextFrame(pixels Frame, sinceStart time.Duration) {
	if len(m.Patterns) != len(m.Weights) {
		panic(fmt.Errorf("len(Patterns) (%d) != len(Weights) (%d)", len(m.Patterns), len(m.Weights)))
	}
	if len(m.bufs) != len(m.Patterns) {
		m.bufs = make([]Frame, len(m.Patterns))
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
		// TODO(maruel): Use uint32 calculation.
		var r, g, b float32
		for j := range m.bufs {
			c := m.bufs[j][i]
			w := m.Weights[j]
			r += float32(c.R) * w
			g += float32(c.G) * w
			b += float32(c.B) * w
		}
		pixels[i].R = FloatToUint8(r)
		pixels[i].G = FloatToUint8(g)
		pixels[i].B = FloatToUint8(b)
	}
}

// Scale adapts a larger or smaller patterns to the Strip size
//
// This is useful to create smoother animations or scale down images.
type Scale struct {
	Child  SPattern
	Scale  ScalingType // Defaults to ScalingLinear
	Length int         // A buffer of this length will be provided to Child and will be scaled to the actual pixels length
	Ratio  float32     // Scaling ratio to use, <1 means smaller, >1 means larger. Only one of Length or Ratio can be used
	buf    Frame
}

func (s *Scale) NextFrame(pixels Frame, sinceStart time.Duration) {
	if s.Child.Pattern == nil {
		return
	}
	l := s.Length
	if l == 0 {
		l = int(ceil(s.Ratio * float32(len(pixels))))
	}
	s.buf.reset(l)
	s.Child.NextFrame(s.buf, sinceStart)
	s.Scale.scale(s.buf, pixels)
}

// Private

func mix(intensity float32, a, b Frame) {
	for i := range a {
		c := b[i]
		t2 := 1 - intensity
		// TODO(maruel): Averaging colors in RGB space looks like hell.
		a[i].R = FloatToUint8(float32(a[i].R)*intensity + float32(c.R)*t2)
		a[i].G = FloatToUint8(float32(a[i].G)*intensity + float32(c.G)*t2)
		a[i].B = FloatToUint8(float32(a[i].B)*intensity + float32(c.B)*t2)
	}
}
