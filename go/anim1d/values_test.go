// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package anim1d

import (
	"testing"

	"github.com/maruel/ut"
)

func TestCurve(t *testing.T) {
	for _, v := range []Curve{Curve(""), Ease, EaseIn, EaseInOut, EaseOut, Direct} {
		ut.AssertEqual(t, float32(0.), v.Scale(0.))
		ut.AssertEqual(t, float32(1.), v.Scale(1.))
	}

	data := []struct {
		t        Curve
		i        float32
		expected float32
	}{
		{Ease, 0.5, 0.8024033904075623},
		{EaseIn, 0.5, 0.3153568208217621},
		{EaseInOut, 0.5, 0.5},
		{EaseOut, 0.5, 0.6846432685852051},
		{Curve(""), 0.5, 0.6846432685852051},
		{Direct, 0.5, 0.5},
	}
	for i, line := range data {
		// TODO(maruel): Round a bit.
		ut.AssertEqualIndex(t, i, line.expected, line.t.Scale(line.i))
	}
}

func TestInterpolation(t *testing.T) {
	b := make(Frame, 1)
	for _, v := range []Interpolation{Interpolation(""), NearestSkip, Nearest, Linear, Bilinear} {
		v.Scale(nil, nil)
		v.Scale(nil, b)
		v.Scale(b, nil)
	}

	// TODO(maruel): Add actual tests.
	red := Color{0xFF, 0x00, 0x00}
	blue := Color{0x00, 0x00, 0xFF}
	data := []struct {
		s        Interpolation
		i        Frame
		expected Frame
	}{
		{NearestSkip, Frame{red, blue}, Frame{red, blue}},
		{Nearest, Frame{red, blue}, Frame{red, blue}},
		{Linear, Frame{red, blue}, Frame{red, blue}},
		{Interpolation(""), Frame{red, blue}, Frame{red, blue}},
		{Bilinear, Frame{red, blue}, Frame{red, blue}},
	}
	for i, line := range data {
		// TODO(maruel): Round a bit.
		out := make(Frame, len(line.expected))
		line.s.Scale(line.i, out)
		ut.AssertEqualIndex(t, i, line.expected, out)
	}
}
