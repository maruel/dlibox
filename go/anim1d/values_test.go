// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package anim1d

import (
	"testing"

	"github.com/maruel/ut"
)

func TestMinMax(t *testing.T) {
	if MinMax(2, 0, 3) != 2 {
		t.Fail()
	}
	if MinMax(-2, 0, 3) != 0 {
		t.Fail()
	}
	if MinMax(4, 0, 3) != 3 {
		t.Fail()
	}
}

func TestMinMax32(t *testing.T) {
	if MinMax32(2, 0, 3) != 2 {
		t.Fail()
	}
	if MinMax32(-2, 0, 3) != 0 {
		t.Fail()
	}
	if MinMax32(4, 0, 3) != 3 {
		t.Fail()
	}
}

// Values

func TestSValue_Eval(t *testing.T) {
	var s SValue
	if s.Eval(23, 0) != 0 {
		t.Fail()
	}
}

func TestConst(t *testing.T) {
	if Const(2).Eval(23, 0) != 2 {
		t.Fail()
	}
}

func TestPercent(t *testing.T) {
	data := []struct {
		p        int32
		timeMS   uint32
		l        int
		expected int32
	}{
		{0, 0, 0, 0},
		{65536, 0, 0, 0},
		{65536, 1000, 0, 0},
		{65536, 0, 1000, 1000},
		{6554, 0, 1000, 100},
		{-65536, 0, 1000, -1000},
		{-6554, 0, 1000, -100},
	}
	for i, line := range data {
		ut.AssertEqualIndex(t, i, line.expected, Percent(line.p).Eval(line.timeMS, line.l))
	}
}

func TestOpAdd(t *testing.T) {
	if (&OpAdd{AddMS: 2}).Eval(23, 0) != 25 {
		t.Fail()
	}
	if (&OpAdd{AddMS: -2}).Eval(23, 0) != 21 {
		t.Fail()
	}
}

func TestOpMod(t *testing.T) {
	if (&OpMod{TickMS: 2}).Eval(23, 0) != 1 {
		t.Fail()
	}
}

func TestOpStep(t *testing.T) {
	if (&OpStep{TickMS: 21}).Eval(23, 0) != 21 {
		t.Fail()
	}
}

func TestRand(t *testing.T) {
	r1 := Rand{0}
	r2 := Rand{16}
	if r1.Eval(0, 0) != r2.Eval(15, 0) {
		t.Fail()
	}
	if r1.Eval(15, 0) == r2.Eval(16, 0) {
		t.Fail()
	}
	if r1.Eval(23, 0) != r2.Eval(23, 0) {
		t.Fail()
	}
}

// Scalers

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

func TestMovePerHour(t *testing.T) {
	data := []struct {
		mps      int32
		timeMS   uint32
		cycle    int
		expected int
	}{
		{1, 0, 10, 0},
		{1, 3600000, 10, 1},
		{1, 2 * 3600000, 10, 2},
		{1, 3 * 3600000, 10, 3},
		{1, 4 * 3600000, 10, 4},
		{1, 5 * 3600000, 10, 5},
		{1, 6 * 3600000, 10, 6},
		{1, 7 * 3600000, 10, 7},
		{1, 8 * 3600000, 10, 8},
		{1, 9 * 3600000, 10, 9},
		{1, 10 * 3600000, 10, 0},
		{1, 10 * 3600000, 11, 10},
		{60, 16, 10, 0},
		{60, 1000, 9, 0},
		{60, 1000, 10, 0},
		{60, 3600000, 10, 0},
		{3600, 3600000, 10, 0},
		{3600000, 0, 10, 0},
		{3600000, 1, 10, 1},
		{3600000, 2, 10, 2},
		{2 * 3600000, 1, 10, 1},
		{2 * 3600000, 2, 10, 2},
	}
	for i, line := range data {
		m := MovePerHour{Const(line.mps)}
		if actual := m.Eval(line.timeMS, 0, line.cycle); actual != line.expected {
			t.Fatalf("%d: %d.Eval(%d, %d) = %d != %d", i, line.mps, line.timeMS, line.cycle, actual, line.expected)
		}
	}
}

func BenchmarkSetupCache(b *testing.B) {
	// Calculate how much this one-time initialization cost is.
	for i := 0; i < b.N; i++ {
		setupCache()
	}
}
