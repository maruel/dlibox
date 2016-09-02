// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// values contains all kind of non pattern types usable as values.

package anim1d

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

//

const epsilon = 1e-7
