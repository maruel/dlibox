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

func TestTransition(t *testing.T) {
	// TODO(maruel): Add.
}

func TestLoop(t *testing.T) {
	// TODO(maruel): Add.
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
