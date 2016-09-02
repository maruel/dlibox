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
		ut.AssertEqual(t, uint16(0), v.Scale(0))
		ut.AssertEqual(t, uint16(65535), v.Scale(65535))
	}

	half := uint16(65535 >> 1)
	data := []struct {
		t        Curve
		i        uint16
		expected uint16
	}{
		{Ease, half, 0xcd01},
		{EaseIn, half, 0x50df},
		{EaseInOut, half, 0x7ffe},
		{EaseOut, half, 0xaf1d},
		{Curve(""), half, 0xaf1d},
		{Direct, half, half},
	}
	for i, line := range data {
		ut.AssertEqualIndex(t, i, line.expected, line.t.Scale(line.i))
	}
}

func TestInterpolationEmpty(t *testing.T) {
	b := make(Frame, 1)
	for _, v := range []Interpolation{Interpolation(""), NearestSkip, Nearest, Linear} {
		v.Scale(nil, nil)
		v.Scale(nil, b)
		v.Scale(b, nil)
	}
}

func TestInterpolation(t *testing.T) {
	red := Color{0xFF, 0x00, 0x00}
	green := Color{0x00, 0xFF, 0x00}
	blue := Color{0x00, 0x00, 0xFF}
	yellow := Color{0xFF, 0xFF, 0x00}
	cyan := Color{0x00, 0xFF, 0xFF}
	magenta := Color{0xFF, 0x00, 0xFF}
	white := Color{0xFF, 0xFF, 0xFF}
	black := Color{}
	input := Frame{red, green, blue, yellow, cyan, magenta, white}
	data := []struct {
		s        Interpolation
		input    Frame
		expected Frame
	}{
		{
			NearestSkip,
			input,
			Frame{red, black, green, black, blue, black, yellow, black, cyan, black, magenta, black, white},
		},
		{NearestSkip, input, Frame{yellow}},
		{NearestSkip, input, Frame{green, magenta}},
		{NearestSkip, input, Frame{green, yellow, magenta}},
		{
			Nearest,
			input,
			Frame{red, red, green, green, blue, blue, yellow, yellow, cyan, cyan, magenta, magenta, white, white},
		},
		{Nearest, input, Frame{yellow}},
		{Nearest, input, Frame{green, magenta}},
		{Nearest, input, Frame{green, yellow, magenta}},
		// TODO(maruel): This is broken.
		/*{
			Linear,
			input,
			Frame{red, red, green, green, blue, blue, yellow, yellow, cyan, cyan, magenta, magenta, white, white},
		},*/
		{Linear, input, Frame{Color{0x80, 0xFF, 0x7F}}},
		{Linear, input, Frame{Color{0x00, 0x80, 0x7F}, Color{0xFF, 0x7F, 0xFF}}},
		{Linear, input, Frame{Color{0x0, 0x80, 0x7F}, Color{0x80, 0xFF, 0x7F}, Color{0xFF, 0x7F, 0xFF}}},
	}
	for i, line := range data {
		out := make(Frame, len(line.expected))
		line.s.Scale(line.input, out)
		ut.AssertEqualIndex(t, i, line.expected, out)
	}
}

func BenchmarkSetupCache(b *testing.B) {
	// Calculate how much this one-time initialization cost is.
	for i := 0; i < b.N; i++ {
		setupCache()
	}
}
