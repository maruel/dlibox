// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// values contains all kind of non pattern types usable as values.

package anim1d

import (
	"math/rand"

	"github.com/maruel/fastbezier"
)

// MinMax limits the value between a min and a max
func MinMax(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

// MinMax32 limits the value between a min and a max
func MinMax32(v, min, max int32) int32 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

// Values

// Value defines a value that may be constant or that may evolve over time.
type Value interface {
	Eval(timeMS uint32, l int) int32
}

// Const is a constant value.
type Const int32

// Eval implements Value.
func (c Const) Eval(timeMS uint32, l int) int32 {
	return int32(c)
}

// Percent is a percentage of the length. It is stored as a 16.16 fixed point.
type Percent int32

// Eval implements Value.
func (p Percent) Eval(timeMS uint32, l int) int32 {
	return int32(int64(l) * int64(p) / 65536)
}

// OpAdd adds a constant to timeMS.
type OpAdd struct {
	AddMS int32
}

// Eval implements Value.
func (o *OpAdd) Eval(timeMS uint32, l int) int32 {
	return int32(timeMS) + o.AddMS
}

// OpMod is a value that is cycling downward.
type OpMod struct {
	TickMS int32 // The cycling time. Maximum is ~25 days.
}

// Eval implements Value.
func (o *OpMod) Eval(timeMS uint32, l int) int32 {
	return int32(timeMS % uint32(o.TickMS))
}

// OpStep is a value that is cycling upward.
//
// It is useful for offsets that are increasing as a stepping function.
type OpStep struct {
	TickMS int32 // The cycling time. Maximum is ~25 days.
}

// Eval implements Value.
func (o *OpStep) Eval(timeMS uint32, l int) int32 {
	return int32(timeMS / uint32(o.TickMS) * uint32(o.TickMS))
}

// Rand is a value that pseudo-randomly changes every TickMS millisecond. If
// unspecified, changes every 60fps.
type Rand struct {
	TickMS int32 // The resolution at which the random value changes.
}

// Eval implements Value.
func (r *Rand) Eval(timeMS uint32, l int) int32 {
	m := uint32(r.TickMS)
	if m == 0 {
		m = 16
	}
	return int32(rand.NewSource(int64(timeMS / m)).Int63())
}

// MovePerHour is the number of movement per hour.
//
// Can be either positive or negative. Maximum supported value is ±3600000, 1000
// move/sec.
//
// Sample values:
//   - 1: one move per hour
//   - 60: one move per minute
//   - 3600: one move per second
//   - 216000: 60 move per second
type MovePerHour SValue

// Eval is not a Value implementation but it leverages an inner one.
func (m *MovePerHour) Eval(timeMS uint32, l int, cycle int) int {
	s := SValue(*m)
	// Prevent overflows.
	v := MinMax32(s.Eval(timeMS, l), -3600000, 3600000)
	// TODO(maruel): Reduce the amount of int64 code in there yet keeping it from
	// overflowing.
	// offset ranges [0, 3599999]
	offset := timeMS % 3600000
	// (1<<32)/3600000 = 1193 is too low. Temporarily upgrade to int64 to
	// calculate the value.
	low := int64(offset) * int64(v) / 3600000
	hour := timeMS / 3600000
	high := int64(hour) * int64(v)
	if cycle != 0 {
		return int((low + high) % int64(cycle))
	}
	return int(low + high)
}

/*
// Equation evaluate an equation at every call.
type Equation struct {
	V string
	f func(timeMS uint32) int32
}

// Eval implements Value.
func (e *Equation) Eval(timeMS uint32) int32 {
	// Compiles the equation to an actual value and precompile it.
	if e.f == nil {
		e.f = func(timeMS uint32) int32 {
			return 0
		}
	}
	return e.f(timeMS)
}
*/

// Scalers

// Bell is a "good enough" approximation of a gaussian curve by using 2
// symmetrical ease-in-out bezier curves.
//
// It is not named Gaussian since it is not a gaussian curve; it really is a
// bell.
type Bell struct{}

// Scale scales input [0, 65535] to output [0, 65535] as a bell curve.
func (b *Bell) Scale(v uint16) uint16 {
	switch {
	case v == 0:
		return 0
	case v == 65535:
		return 0
	case v == 32767:
		return 65535

	case v < 32767:
		return EaseInOut.Scale(v * 2)
	default:
		return EaseInOut.Scale(65535 - v*2)
	}
}

// Curve models visually pleasing curves.
//
// They are modeled against CSS transitions.
// https://www.w3.org/TR/web-animations/#scaling-using-a-cubic-bezier-curve
type Curve string

// All the kind of known curves.
const (
	Ease       Curve = "ease"
	EaseIn     Curve = "ease-in"
	EaseInOut  Curve = "ease-in-out"
	EaseOut    Curve = "ease-out" // Recommended and default value.
	Direct     Curve = "direct"   // linear mapping
	StepStart  Curve = "steps(1,start)"
	StepMiddle Curve = "steps(1,middle)"
	StepEnd    Curve = "steps(1,end)"
)

var lutCache map[Curve]fastbezier.LUT

func setupCache() map[Curve]fastbezier.LUT {
	cache := map[Curve]fastbezier.LUT{
		Ease:      fastbezier.Make(0.25, 0.1, 0.25, 1, 18),
		EaseIn:    fastbezier.Make(0.42, 0, 1, 1, 18),
		EaseInOut: fastbezier.Make(0.42, 0, 0.58, 1, 18),
		EaseOut:   fastbezier.Make(0, 0, 0.58, 1, 18),
	}
	cache[""] = cache[EaseOut]
	return cache
}

func init() {
	lutCache = setupCache()
}

// Scale scales input [0, 65535] to output [0, 65535] using the curve
// requested.
func (c Curve) Scale(intensity uint16) uint16 {
	switch c {
	case Ease, EaseIn, EaseInOut, EaseOut, "":
		return lutCache[c].Eval(intensity)
	default:
		return lutCache[""].Eval(intensity)
	case Direct:
		return intensity
	case StepStart:
		if intensity < 256 {
			return 0
		}
		return 65535
	case StepMiddle:
		if intensity < 32768 {
			return 0
		}
		return 65535
	case StepEnd:
		if intensity >= 65535-256 {
			return 65535
		}
		return 0
	}
}

// Scale8 saves on casting.
func (c Curve) Scale8(intensity uint16) uint8 {
	return uint8(c.Scale(intensity) >> 8)
}

// Interpolation specifies a way to scales a pixel strip.
type Interpolation string

// All the kinds of interpolations.
const (
	NearestSkip Interpolation = "nearestskip" // Selects the nearest pixel but when upscaling, skips on missing pixels.
	Nearest     Interpolation = "nearest"     // Selects the nearest pixel, gives a blocky view.
	Linear      Interpolation = "linear"      // Linear interpolation, recommended and default value.
)

// Scale interpolates a frame into another using integers as much as possible
// for reasonable performance.
func (i Interpolation) Scale(in, out Frame) {
	li := len(in)
	lo := len(out)
	if li == 0 || lo == 0 {
		return
	}
	switch i {
	case NearestSkip:
		if li < lo {
			// Do not touch skipped pixels.
			for i, p := range in {
				out[(i*lo+lo/2)/li] = p
			}
			return
		}
		// When the destination is smaller than the source, Nearest and NearestSkip
		// have the same behavior.
		fallthrough
	case Nearest, "":
		fallthrough
	default:
		for i := range out {
			out[i] = in[(i*li+li/2)/lo]
		}
	case Linear:
		for i := range out {
			x := (i*li + li/2) / lo
			c := in[x]
			if x < li-1 {
				gradient := uint8(127)
				c.Mix(in[x+1], gradient)
			}
			out[i] = c
			//a := in[(i*li+li/2)/lo]
			//b := in[(i*li+li/2)/lo]
			//out[i] = (a + b) / 2
		}
	}
}

//

const epsilon = 1e-7
