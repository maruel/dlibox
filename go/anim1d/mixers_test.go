// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package anim1d

import (
	"testing"

	"github.com/maruel/ut"
)

func TestTransitionType(t *testing.T) {
	for _, v := range []TransitionType{TransitionType(""), TransitionEase, TransitionEaseIn, TransitionEaseInOut, TransitionEaseOut, TransitionLinear} {
		ut.AssertEqual(t, float32(0.), v.scale(0.))
		ut.AssertEqual(t, float32(1.), v.scale(1.))
	}

	data := []struct {
		t        TransitionType
		i        float32
		expected float32
	}{
		{TransitionEase, 0.5, 0.8024033904075623},
		{TransitionEaseIn, 0.5, 0.3153568208217621},
		{TransitionEaseInOut, 0.5, 0.5},
		{TransitionEaseOut, 0.5, 0.6846432685852051},
		{TransitionType(""), 0.5, 0.6846432685852051},
		{TransitionLinear, 0.5, 0.5},
	}
	for i, line := range data {
		// TODO(maruel): Round a bit.
		ut.AssertEqualIndex(t, i, line.expected, line.t.scale(line.i))
	}
}

func TestScalingType(t *testing.T) {
	b := make(Frame, 1)
	for _, v := range []ScalingType{ScalingType(""), ScalingNearestSkip, ScalingNearest, ScalingLinear, ScalingBilinear} {
		v.scale(nil, nil)
		v.scale(nil, b)
		v.scale(b, nil)
	}

	// TODO(maruel): Add actual tests.
	red := Color{0xFF, 0x00, 0x00}
	blue := Color{0x00, 0x00, 0xFF}
	data := []struct {
		s        ScalingType
		i        Frame
		expected Frame
	}{
		{ScalingNearestSkip, Frame{red, blue}, Frame{red, blue}},
		{ScalingNearest, Frame{red, blue}, Frame{red, blue}},
		{ScalingLinear, Frame{red, blue}, Frame{red, blue}},
		{ScalingType(""), Frame{red, blue}, Frame{red, blue}},
		{ScalingBilinear, Frame{red, blue}, Frame{red, blue}},
	}
	for i, line := range data {
		// TODO(maruel): Round a bit.
		out := make(Frame, len(line.expected))
		line.s.scale(line.i, out)
		ut.AssertEqualIndex(t, i, line.expected, out)
	}
}

func TestGradient(t *testing.T) {
	a := &Color{0x10, 0x10, 0x10}
	b := &Color{0x20, 0x20, 0x20}
	testFrame(t, &Gradient{Left: SPattern{a}, Right: SPattern{b}, Transition: TransitionLinear}, expectation{0, Frame{{0x18, 0x18, 0x18}}})
	testFrame(t, &Gradient{Left: SPattern{a}, Right: SPattern{b}, Transition: TransitionLinear}, expectation{0, Frame{{0x10, 0x10, 0x10}, {0x20, 0x20, 0x20}}})
	testFrame(t, &Gradient{Left: SPattern{a}, Right: SPattern{b}, Transition: TransitionLinear}, expectation{0, Frame{{0x10, 0x10, 0x10}, {0x18, 0x18, 0x18}, {0x20, 0x20, 0x20}}})
}

func TestTransition(t *testing.T) {
	// TODO(maruel): Add.
}

func TestCycle(t *testing.T) {
	// TODO(maruel): Add.
}

func TestLoop(t *testing.T) {
	// TODO(maruel): Add.
}

func TestRotate(t *testing.T) {
	a := Color{10, 10, 10}
	b := Color{20, 20, 20}
	c := Color{30, 30, 30}
	p := &Rotate{Child: SPattern{Frame{a, b, c}}, MovesPerSec: 100}
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
	p := &Rotate{Child: SPattern{Frame{a, b, c}}, MovesPerSec: -100}
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

	p := &PingPong{Child: SPattern{Frame{a, b}}, MovesPerSec: 100}
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

	p = &PingPong{Child: SPattern{Frame{a, b, c, d, e, f}}, MovesPerSec: 1}
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
	// TODO(maruel): Add.
}

func TestMixer(t *testing.T) {
	// TODO(maruel): Add.
}

func TestScale(t *testing.T) {
	// TODO(maruel): Add.
}
