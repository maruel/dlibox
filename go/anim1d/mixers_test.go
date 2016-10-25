// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package anim1d

import "testing"

func TestGradient(t *testing.T) {
	a := &Color{0x10, 0x10, 0x10}
	b := &Color{0x20, 0x20, 0x20}
	testFrame(t, &Gradient{Left: SPattern{a}, Right: SPattern{b}, Curve: Direct}, expectation{0, Frame{{0x18, 0x18, 0x18}}})
	testFrame(t, &Gradient{Left: SPattern{a}, Right: SPattern{b}, Curve: Direct}, expectation{0, Frame{{0x10, 0x10, 0x10}, {0x20, 0x20, 0x20}}})
	testFrame(t, &Gradient{Left: SPattern{a}, Right: SPattern{b}, Curve: Direct}, expectation{0, Frame{{0x10, 0x10, 0x10}, {0x18, 0x18, 0x18}, {0x20, 0x20, 0x20}}})
}

func TestTransition(t *testing.T) {
	// TODO(maruel): Add.
}

func TestLoop(t *testing.T) {
	// TODO(maruel): Add.
}

func TestRotate(t *testing.T) {
	a := Color{10, 10, 10}
	b := Color{20, 20, 20}
	c := Color{30, 30, 30}
	p := &Rotate{Child: SPattern{Frame{a, b, c}}, MovePerHour: MovePerHour{Const(360000)}}
	e := []expectation{
		{0, Frame{a, b, c}},
		{5, Frame{a, b, c}},
		{10, Frame{c, a, b}},
		{20, Frame{b, c, a}},
		{30, Frame{a, b, c}},
		{40, Frame{c, a, b}},
		{50, Frame{b, c, a}},
		{60, Frame{a, b, c}},
	}
	testFrames(t, p, e)
}

func TestRotateRev(t *testing.T) {
	// Works in reverse too.
	a := Color{10, 10, 10}
	b := Color{20, 20, 20}
	c := Color{30, 30, 30}
	p := &Rotate{Child: SPattern{Frame{a, b, c}}, MovePerHour: MovePerHour{Const(-360000)}}
	e := []expectation{
		{0, Frame{a, b, c}},
		{5, Frame{a, b, c}},
		{10, Frame{b, c, a}},
		{20, Frame{c, a, b}},
		{30, Frame{a, b, c}},
		{40, Frame{b, c, a}},
		{50, Frame{c, a, b}},
		{60, Frame{a, b, c}},
	}
	testFrames(t, p, e)
}

func TestChronometer(t *testing.T) {
	r := Color{0xff, 0x00, 0x00}
	g := Color{0x00, 0xff, 0x00}
	b := Color{0x00, 0x00, 0xff}
	p := &Chronometer{Child: SPattern{Frame{{}, r, g, b}}}
	exp := []expectation{
		{0, Frame{r, {}, {}, {}, {}, {}}},                        // 0:00:00
		{1000 * 10, Frame{g, r, {}, {}, {}, {}}},                 // 0:00:10
		{1000 * 20, Frame{g, {}, r, {}, {}, {}}},                 // 0:00:20
		{1000 * 60, Frame{r, {}, {}, {}, {}, {}}},                // 0:01:00
		{1000 * 600, Frame{r, g, {}, {}, {}, {}}},                // 0:10:00
		{1000 * 3600, Frame{r, b, {}, {}, {}, {}}},               // 1:00:00
		{1000 * (3600 + 20*60 + 30), Frame{{}, b, g, r, {}, {}}}, // 1:20:30
	}
	testFrames(t, p, exp)
}

func TestPingPong(t *testing.T) {
	a := Color{0x10, 0x10, 0x10}
	b := Color{0x20, 0x20, 0x20}
	c := Color{0x30, 0x30, 0x30}
	d := Color{0x40, 0x40, 0x40}
	e := Color{0x50, 0x50, 0x50}
	f := Color{0x60, 0x60, 0x60}

	p := &PingPong{Child: SPattern{Frame{a, b}}, MovePerHour: MovePerHour{Const(360000)}}
	exp := []expectation{
		{0, Frame{a, b, {}}},
		{5, Frame{a, b, {}}},
		{10, Frame{b, a, {}}},
		{20, Frame{{}, b, a}},
		{30, Frame{{}, a, b}},
		{40, Frame{a, b, {}}},
		{50, Frame{b, a, {}}},
		{60, Frame{{}, b, a}},
	}
	testFrames(t, p, exp)

	p = &PingPong{Child: SPattern{Frame{a, b, c, d, e, f}}, MovePerHour: MovePerHour{Const(3600)}}
	exp = []expectation{
		{0, Frame{a, b, c, d}},
		{500, Frame{a, b, c, d}},
		{1000, Frame{b, a, d, e}},
		{2000, Frame{c, b, a, f}},
		{3000, Frame{d, c, b, a}},
		{4000, Frame{e, d, a, b}},
		{5000, Frame{f, a, b, c}},
		{6000, Frame{a, b, c, d}},
	}
	testFrames(t, p, exp)
}

func TestCrop(t *testing.T) {
	// Crop skips the begining and the end of the source.
	f := Frame{
		{0x10, 0x10, 0x10},
		{0x20, 0x20, 0x20},
		{0x30, 0x30, 0x30},
	}
	p := &Crop{Child: SPattern{f}, Before: SValue{Const(1)}, After: SValue{Const(2)}}
	testFrame(t, p, expectation{0, f[1:3]})
}

func TestSubset(t *testing.T) {
	// Subset skips the begining and the end of the destination.
	f := Frame{
		{0x10, 0x10, 0x10},
		{0x20, 0x20, 0x20},
		{0x30, 0x30, 0x30},
	}
	p := &Subset{Child: SPattern{f}, Offset: SValue{Const(1)}, Length: SValue{Const(2)}}
	// Skip the begining and the end of the destination.
	expected := Frame{
		{},
		{0x10, 0x10, 0x10},
		{0x20, 0x20, 0x20},
		{},
	}
	testFrame(t, p, expectation{0, expected})
}

func TestDim(t *testing.T) {
	p := &Dim{Child: SPattern{&Color{0x60, 0x60, 0x60}}, Intensity: SValue{Const(127)}}
	testFrame(t, p, expectation{0, Frame{{0x2f, 0x2f, 0x2f}}})
}

func TestAdd(t *testing.T) {
	a := Color{0x60, 0x60, 0x60}
	b := Color{0x10, 0x20, 0x30}
	p := &Add{Patterns: []SPattern{{&a}, {&b}}}
	testFrame(t, p, expectation{0, Frame{{0x70, 0x80, 0x90}}})
}

func TestScale(t *testing.T) {
	f := Frame{{0x60, 0x60, 0x60}, {0x10, 0x20, 0x30}}
	p := &Scale{Child: SPattern{f}, Interpolation: NearestSkip, RatioMilli: SValue{Const(667)}}
	expected := Frame{{0x60, 0x60, 0x60}, {}, {0x10, 0x20, 0x30}}
	testFrame(t, p, expectation{0, expected})
}
