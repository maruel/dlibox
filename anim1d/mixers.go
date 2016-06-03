// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package anim1d

import (
	"log"
	"time"
)

// TransitionType models visually pleasing transitions.
//
// They are modeled against CSS transitions.
// https://www.w3.org/TR/web-animations/#scaling-using-a-cubic-bezier-curve
type TransitionType string

const (
	TransitionEase       TransitionType = "ease"
	TransitionEaseIn     TransitionType = "ease-in"
	TransitionEaseInOut  TransitionType = "ease-in-out"
	TransitionEaseOut    TransitionType = "ease-out" // Recommended and default value.
	TransitionLinear     TransitionType = "linear"
	TransitionStepStart  TransitionType = "steps(1,start)"
	TransitionStepMiddle TransitionType = "steps(1,middle)"
	TransitionStepEnd    TransitionType = "steps(1,end)"
)

const epsilon = 1e-7

// scale scales input [0, 1] to output [0, 1] using the transition requested.
//
// TODO(maruel): Implement a version that is integer based.
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
	case TransitionStepStart:
		if intensity < 0.+epsilon {
			return 0
		}
		return 1
	case TransitionStepMiddle:
		if intensity < 0.5 {
			return 0
		}
		return 1
	case TransitionStepEnd:
		if intensity > 1.-epsilon {
			return 1
		}
		return 0
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

// Gradient does a gradient between 2 patterns.
//
// A good example is using two colors but it can also be animations.
//
// TODO(maruel): Support N colors at M positions.
type Gradient struct {
	Left       SPattern
	Right      SPattern
	Transition TransitionType
	buf        Frame
}

func (g *Gradient) NextFrame(pixels Frame, sinceStart time.Duration) {
	if g.Left.Pattern == nil || g.Right.Pattern == nil {
		return
	}
	l := len(pixels) - 1
	g.buf.reset(len(pixels))
	g.Left.NextFrame(pixels, sinceStart)
	g.Right.NextFrame(g.buf, sinceStart)
	if l == 0 {
		pixels.Mix(g.buf, FloatToUint8(255.*g.Transition.scale(0.5)))
	} else {
		// TODO(maruel): Convert to integer calculation.
		max := float32(len(pixels) - 1)
		for i := range pixels {
			// [0, 1]
			intensity := float32(i) / max
			pixels[i].Mix(g.buf[i], FloatToUint8(255.*g.Transition.scale(intensity)))
		}
	}
}

// Transition changes from Before to After over time. It doesn't repeat.
//
// In gets sinceStart that is subtracted by Offset.
type Transition struct {
	Before     SPattern       // Old pattern that is disappearing
	After      SPattern       // New pattern to show
	Offset     time.Duration  // Offset at which the transiton from Before->In starts
	Duration   time.Duration  // Duration of the transition while both are rendered
	Transition TransitionType // Type of transition, defaults to EaseOut if not set
	buf        Frame
}

func (t *Transition) NextFrame(pixels Frame, sinceStart time.Duration) {
	if sinceStart <= t.Offset {
		// Before transition.
		if t.Before.Pattern != nil {
			t.Before.NextFrame(pixels, sinceStart)
		}
		return
	}
	if t.After.Pattern != nil {
		t.After.NextFrame(pixels, sinceStart-t.Offset)
	}
	if sinceStart >= t.Offset+t.Duration {
		// After transition.
		t.buf = nil
		return
	}
	t.buf.reset(len(pixels))

	// TODO(maruel): Add lateral animation and others.
	if t.Before.Pattern != nil {
		t.Before.NextFrame(t.buf, sinceStart)
	}
	pixels.Mix(t.buf, 255.-FloatToUint8(255.*t.Transition.scale(float32(sinceStart-t.Offset)/float32(t.Duration))))
}

// Cycle cycles between multiple patterns. It can be used as an animatable
// looping frame.
//
// TODO(maruel): Blend between frames with TransitionType, defaults to step.
// TODO(maruel): Merge with Loop.
type Cycle struct {
	Frames        []SPattern
	FrameDuration time.Duration
}

func (c *Cycle) NextFrame(pixels Frame, sinceStart time.Duration) {
	if len(c.Frames) == 0 {
		return
	}
	c.Frames[int(sinceStart/c.FrameDuration)%len(c.Frames)].NextFrame(pixels, sinceStart)
}

func (c *Cycle) NativeDuration(pixels int) time.Duration {
	return c.FrameDuration * time.Duration(len(c.Frames))
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
	pixels.Mix(l.buf, FloatToUint8(255.*l.Transition.scale(intensity)))
}

// Rotate rotates a pattern that can also cycle either way.
//
// Use negative to go left. Can be used for 'candy bar'.
//
// Similar to PingPong{} except that it doesn't bounce.
//
// Use 5x oversampling with Scale{} to create smoother animation.
type Rotate struct {
	Child       SPattern
	MovesPerSec float32 // Expressed in number of light jumps per second.
	buf         Frame
}

func (r *Rotate) NextFrame(pixels Frame, sinceStart time.Duration) {
	l := len(pixels)
	if l == 0 || r.Child.Pattern == nil {
		return
	}
	r.buf.reset(l)
	r.Child.NextFrame(r.buf, sinceStart)
	offset := int(float32(sinceStart.Seconds())*r.MovesPerSec) % l
	if offset < 0 {
		offset = l + offset
	}
	copy(pixels[offset:], r.buf)
	copy(pixels[:offset], r.buf[l-offset:])
}

func (r *Rotate) NativeDuration(pixels int) time.Duration {
	return time.Duration(float32(pixels)/r.MovesPerSec) * time.Second
}

// PingPong shows a 'ball' with a trail that bounces from one side to
// the other.
//
// Can be used for a ball, a water wave or K2000 (Knight Rider) style light.
// The trail can be a Frame or a dynamic pattern.
//
// To get smoothed movement, use Scale{} with a 5x factor or so.
// TODO(maruel): That's a bit inefficient, enable ScalingType here.
type PingPong struct {
	Child       SPattern // [0] is the front pixel so the pixels are effectively drawn in reverse order.
	MovesPerSec float32  // Expressed in number of light jumps per second.
	buf         Frame
}

func (p *PingPong) NextFrame(pixels Frame, sinceStart time.Duration) {
	if len(pixels) == 0 || p.Child.Pattern == nil {
		return
	}
	p.buf.reset(len(pixels)*2 - 1)
	p.Child.NextFrame(p.buf, sinceStart)
	// The last point of each extremity is only lit on one tick but every other
	// points are lit twice during a full cycle. This means the full cycle is
	// 2*(len(pixels)-1). For a 3 pixels line, the cycle is: x00, 0x0, 00x, 0x0.
	//
	// For Child being Frame "01234567":
	//   move == 0  -> "01234567"
	//   move == 2  -> "21056789"
	//   move == 5  -> "543210ab"
	//   move == 7  -> "76543210"
	//   move == 9  -> "98765012"
	//   move == 11 -> "ba901234"
	//   move == 13 -> "d0123456"
	//   move 14 -> move 0; "2*(8-1)"
	cycle := 2 * (len(pixels) - 1)
	// TODO(maruel): Smoothing with TransitionType, defaults to Step.
	pos := int(float32(sinceStart.Seconds())*p.MovesPerSec) % cycle

	// Once it works the following code looks trivial but everytime it takes me
	// an absurd amount of time to rewrite it.
	if pos >= len(pixels)-1 {
		// Head runs left.
		// pos2 is the position from the right.
		pos2 := pos + 1 - len(pixels)
		// limit is the offset at which order change.
		limit := len(pixels) - pos2 - 1
		for i := range pixels {
			if i < limit {
				// Going right.
				pixels[i] = p.buf[len(pixels)-i+pos2-1]
			} else {
				// Going left.
				pixels[i] = p.buf[i-limit]
			}
		}
	} else {
		// Head runs right.
		for i := range pixels {
			if i <= pos {
				// Going right.
				pixels[i] = p.buf[pos-i]
			} else {
				// Going left.
				pixels[i] = p.buf[pos+i]
			}
		}
	}
}

func (p *PingPong) NativeDuration(pixels int) time.Duration {
	cycle := 2 * (pixels - 1)
	return time.Duration(p.MovesPerSec*float32(cycle)) * time.Second
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
		log.Printf("len(Patterns) (%d) != len(Weights) (%d)", len(m.Patterns), len(m.Weights))
		return
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
