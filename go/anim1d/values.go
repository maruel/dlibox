// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// values contains all kind of non pattern types usable as values.

package anim1d

// Curve models visually pleasing curves between 0 and 1.
//
// They are modeled against CSS transitions.
// https://www.w3.org/TR/web-animations/#scaling-using-a-cubic-bezier-curve
type Curve string

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

// Scale scales input [0, 1] to output [0, 1] using the transformation curve
// requested.
//
// TODO(maruel): Implement a version that is integer based.
func (c Curve) Scale(intensity float32) float32 {
	// TODO(maruel): Add support for arbitrary cubic-bezier().
	// TODO(maruel): Map ease-* to cubic-bezier().
	// TODO(maruel): Add support for steps() which is pretty cool.
	switch c {
	case Ease:
		return cubicBezier(0.25, 0.1, 0.25, 1, intensity)
	case EaseIn:
		return cubicBezier(0.42, 0, 1, 1, intensity)
	case EaseInOut:
		return cubicBezier(0.42, 0, 0.58, 1, intensity)
	case EaseOut, "":
		fallthrough
	default:
		return cubicBezier(0, 0, 0.58, 1, intensity)
	case Direct:
		return intensity
	case StepStart:
		if intensity < 0.+epsilon {
			return 0
		}
		return 1
	case StepMiddle:
		if intensity < 0.5 {
			return 0
		}
		return 1
	case StepEnd:
		if intensity > 1.-epsilon {
			return 1
		}
		return 0
	}
}

// Interpolation specifies a way to scales a pixel strip.
type Interpolation string

const (
	NearestSkip Interpolation = "nearestskip" // Selects the nearest pixel but when upscaling, skips on missing pixels.
	Nearest     Interpolation = "nearest"     // Selects the nearest pixel, gives a blocky view.
	Linear      Interpolation = "linear"      // Linear interpolation, recommended and default value.
	Bilinear    Interpolation = "bilinear"    // Bilinear interpolation, usually overkill for 1D.
)

func (i Interpolation) Scale(in, out Frame) {
	// Use integer operations as much as possible for reasonable performance.
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
		fallthrough
	case Nearest, Linear, Bilinear, "":
		fallthrough
	default:
		for i := range out {
			out[i] = in[(i*li+li/2)/lo]
		}
		/*
			case Linear:
				for i := range out {
					x := (i*li + li/2) / lo
					c := in[x]
					c.Add(in[x+1])
					out[i] = c
				}
		*/
	}
}

//

const epsilon = 1e-7
