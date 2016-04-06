// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package dotstar

import (
	"image/color"
	"testing"
	"time"
)

func TestStaticColor(t *testing.T) {
	p := &StaticColor{color.NRGBA{255, 255, 255, 255}}
	e := []expectation{{3 * time.Second, []color.NRGBA{{255, 255, 255, 255}}}}
	frames(t, p, e)
}

func TestGlow1(t *testing.T) {
	p := &Glow{[]color.NRGBA{{255, 255, 255, 255}, {0, 0, 0, 255}}, 1}
	e := []expectation{
		{0, []color.NRGBA{{0xFF, 0xFF, 0xFF, 0xFF}}},
		{250 * time.Millisecond, []color.NRGBA{{0xBF, 0xBF, 0xBF, 0xFF}}},
		{500 * time.Millisecond, []color.NRGBA{{0x7F, 0x7F, 0x7F, 0xFF}}},
		{750 * time.Millisecond, []color.NRGBA{{0x3F, 0x3F, 0x3F, 0xFF}}},
		{1000 * time.Millisecond, []color.NRGBA{{0x00, 0x00, 0x00, 0xFF}}},
	}
	frames(t, p, e)
}

func TestGlow2(t *testing.T) {
	p := &Glow{[]color.NRGBA{{255, 255, 255, 255}, {0, 0, 0, 255}}, 0.1}
	e := []expectation{
		{0, []color.NRGBA{{0xFF, 0xFF, 0xFF, 0xFF}}},
		{2500 * time.Millisecond, []color.NRGBA{{0xBF, 0xBF, 0xBF, 0xFF}}},
		{5000 * time.Millisecond, []color.NRGBA{{0x80, 0x80, 0x80, 0xFF}}},
		{7500 * time.Millisecond, []color.NRGBA{{0x3F, 0x3F, 0x3F, 0xFF}}},
		{10000 * time.Millisecond, []color.NRGBA{{0x00, 0x00, 0x00, 0xFF}}},
	}
	frames(t, p, e)
}

func TestPingPong(t *testing.T) {
	a := color.NRGBA{10, 10, 10, 10}
	b := color.NRGBA{20, 20, 20, 20}
	p := &PingPong{[]color.NRGBA{a, b}, color.NRGBA{}, 1000}
	e := []expectation{
		{0, []color.NRGBA{a, {}, {}}},
		{500 * time.Microsecond, []color.NRGBA{a, {}, {}}},
		{1 * time.Millisecond, []color.NRGBA{b, a, {}}},
		{2 * time.Millisecond, []color.NRGBA{{}, b, a}},
		{3 * time.Millisecond, []color.NRGBA{{}, a, b}},
		{4 * time.Millisecond, []color.NRGBA{a, b, {}}},
		{5 * time.Millisecond, []color.NRGBA{b, a, {}}},
		{6 * time.Millisecond, []color.NRGBA{{}, b, a}},
	}
	frames(t, p, e)
}

func TestRepeated(t *testing.T) {
	a := color.NRGBA{10, 10, 10, 10}
	b := color.NRGBA{20, 20, 20, 20}
	c := color.NRGBA{30, 30, 30, 30}
	p := &Repeated{[]color.NRGBA{a, b, c}, 1000}
	e := []expectation{
		{0, []color.NRGBA{a, b, c, a, b}},
		{500 * time.Microsecond, []color.NRGBA{a, b, c, a, b}},
		{1 * time.Millisecond, []color.NRGBA{c, a, b, c, a}},
		{2 * time.Millisecond, []color.NRGBA{b, c, a, b, c}},
		{3 * time.Millisecond, []color.NRGBA{a, b, c, a, b}},
		{4 * time.Millisecond, []color.NRGBA{c, a, b, c, a}},
		{5 * time.Millisecond, []color.NRGBA{b, c, a, b, c}},
		{6 * time.Millisecond, []color.NRGBA{a, b, c, a, b}},
	}
	frames(t, p, e)
}

func TestRepeatedRev(t *testing.T) {
	// Works in reverse too.
	a := color.NRGBA{10, 10, 10, 10}
	b := color.NRGBA{20, 20, 20, 20}
	c := color.NRGBA{30, 30, 30, 30}
	p := &Repeated{[]color.NRGBA{a, b, c}, -1000}
	e := []expectation{
		{0, []color.NRGBA{a, b, c, a, b}},
		{500 * time.Microsecond, []color.NRGBA{a, b, c, a, b}},
		{1 * time.Millisecond, []color.NRGBA{b, c, a, b, c}},
		{2 * time.Millisecond, []color.NRGBA{c, a, b, c, a}},
		{3 * time.Millisecond, []color.NRGBA{a, b, c, a, b}},
		{4 * time.Millisecond, []color.NRGBA{b, c, a, b, c}},
		{5 * time.Millisecond, []color.NRGBA{c, a, b, c, a}},
		{6 * time.Millisecond, []color.NRGBA{a, b, c, a, b}},
	}
	frames(t, p, e)
}

//

type expectation struct {
	offset time.Duration
	colors []color.NRGBA
}

func frames(t *testing.T, p Pattern, expectations []expectation) {
	pixels := make([]color.NRGBA, len(expectations[0].colors))
	for frame, e := range expectations {
		p.NextFrame(pixels, e.offset)
		for j := range e.colors {
			a := e.colors[j]
			b := pixels[j]
			dR := int(a.R) - int(b.R)
			dG := int(a.G) - int(b.G)
			dB := int(a.B) - int(b.B)
			dA := int(a.A) - int(b.A)
			if dR > 1 || dR < -1 || dG > 1 || dG < -1 || dB > 1 || dB < -1 || dA > 1 || dA < -1 {
				t.Fatalf("frame=%d; pixel=%d; %v != %v", frame, j, a, b)
			}
		}
	}
}
