// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package anim1d

import (
	"encoding/json"
	"testing"

	"github.com/maruel/ut"
)

func TestNilObject(t *testing.T) {
	c := Frame{}
	d := Frame{{}}
	for _, p := range knownPatterns {
		p.NextFrame(nil, 0)
		p.NextFrame(c, 0)
		p.NextFrame(d, 0)
	}
}

func serialize(t *testing.T, p Pattern, expected string) {
	p2 := &SPattern{p}
	b, err := json.Marshal(p2)
	ut.AssertEqual(t, nil, err)
	ut.AssertEqual(t, expected, string(b))
	p2.Pattern = nil
	ut.AssertEqual(t, nil, json.Unmarshal(b, p2))
}

func isColorOrFrameOrRainbow(p Pattern) bool {
	if _, ok := p.(*Color); ok {
		return ok
	}
	if _, ok := p.(*Frame); ok {
		return ok
	}
	_, ok := p.(*Rainbow)
	return ok
}

func TestJSON(t *testing.T) {
	for _, p := range knownPatterns {
		p2 := &SPattern{p}
		b, err := json.Marshal(p2)
		ut.AssertEqual(t, nil, err)
		if isColorOrFrameOrRainbow(p) {
			ut.AssertEqual(t, uint8('"'), b[0])
		} else {
			ut.AssertEqual(t, uint8('{'), b[0])
		}
		p2.Pattern = nil
		ut.AssertEqualf(t, nil, json.Unmarshal(b, p2), "%s", b)
	}
	serialize(t, &Color{1, 2, 3}, `"#010203"`)
	serialize(t, &Frame{}, `"L"`)
	serialize(t, &Frame{{1, 2, 3}, {4, 5, 6}}, `"L010203040506"`)
	serialize(t, &Rainbow{}, `"Rainbow"`)
	serialize(t, &PingPong{}, `{"Child":{},"MovesPerSec":0,"_type":"PingPong"}`)
	serialize(t, &Chronometer{}, `{"Child":{},"_type":"Chronometer"}`)
	serialize(t, &Cycle{}, `{"FrameDurationMS":0,"Frames":null,"_type":"Cycle"}`)

	// Create one more complex. Assert that int64 is not mangled.
	p := &Transition{
		Before: SPattern{
			&Transition{
				After:      SPattern{&Color{255, 255, 255}},
				OffsetMS:   600000,
				DurationMS: 600000,
				Transition: TransitionLinear,
			},
		},
		After:      SPattern{&Color{}},
		OffsetMS:   30 * 60000,
		DurationMS: 600000,
		Transition: TransitionLinear,
	}
	expected := `{"After":"#000000","Before":{"After":"#ffffff","Before":{},"DurationMS":600000,"OffsetMS":600000,"Transition":"linear","_type":"Transition"},"DurationMS":600000,"OffsetMS":1800000,"Transition":"linear","_type":"Transition"}`
	serialize(t, p, expected)
}
