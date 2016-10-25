// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package anim1d

import (
	"encoding/json"
	"strconv"
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

func TestJSONPatterns(t *testing.T) {
	for _, p := range knownPatterns {
		p2 := &SPattern{p}
		b, err := json.Marshal(p2)
		ut.AssertEqual(t, nil, err)
		if isColorOrFrameOrRainbow(p) {
			ut.AssertEqual(t, uint8('"'), b[0])
		} else {
			ut.AssertEqual(t, uint8('{'), b[0])
		}
		// Must not crash on nil members and empty frame.
		p2.NextFrame(Frame{}, 0)
		p2.Pattern = nil
		ut.AssertEqualf(t, nil, json.Unmarshal(b, p2), "%s", b)
	}
}

func TestJSONPatternsSpotCheck(t *testing.T) {
	// Increase coverage of edge cases.
	serializePattern(t, &Color{1, 2, 3}, `"#010203"`)
	serializePattern(t, &Frame{}, `"L"`)
	serializePattern(t, &Frame{{1, 2, 3}, {4, 5, 6}}, `"L010203040506"`)
	serializePattern(t, &Rainbow{}, `"Rainbow"`)
	serializePattern(t, &PingPong{}, `{"Child":{},"MovePerHour":0,"_type":"PingPong"}`)
	serializePattern(t, &Chronometer{}, `{"Child":{},"_type":"Chronometer"}`)

	// Create one more complex. Assert that int64 is not mangled.
	p := &Transition{
		Before: SPattern{
			&Transition{
				After:      SPattern{&Color{255, 255, 255}},
				OffsetMS:   600000,
				DurationMS: 600000,
				Curve:      Direct,
			},
		},
		After:      SPattern{&Color{}},
		OffsetMS:   30 * 60000,
		DurationMS: 600000,
		Curve:      Direct,
	}
	expected := `{"After":"#000000","Before":{"After":"#ffffff","Before":{},"Curve":"direct","DurationMS":600000,"OffsetMS":600000,"_type":"Transition"},"Curve":"direct","DurationMS":600000,"OffsetMS":1800000,"_type":"Transition"}`
	serializePattern(t, p, expected)
}

func TestJSONValues(t *testing.T) {
	for _, v := range knownValues {
		v2 := &SValue{v}
		b, err := json.Marshal(v2)
		ut.AssertEqual(t, nil, err)
		if isConst(v) {
			if _, err := strconv.ParseInt(string(b), 10, 32); err != nil {
				t.Fatalf("%v", err)
			}
		} else if isRand(v) {
			if b[0] != '"' {
				t.Fatalf("%v", b)
			}
		} else {
			ut.AssertEqual(t, uint8('{'), b[0])
		}
		v2.Value = nil
		ut.AssertEqualf(t, nil, json.Unmarshal(b, v2), "%s", b)
	}
}

func TestJSONValuesSpotCheck(t *testing.T) {
	// Increase coverage of edge cases.
	serializeValue(t, &Rand{TickMS: 43}, `{"TickMS":43,"_type":"Rand"}`)
}

//

func serializePattern(t *testing.T, p Pattern, expected string) {
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

func serializeValue(t *testing.T, v Value, expected string) {
	v2 := &SValue{v}
	b, err := json.Marshal(v2)
	ut.AssertEqual(t, nil, err)
	ut.AssertEqual(t, expected, string(b))
	v2.Value = nil
	ut.AssertEqual(t, nil, json.Unmarshal(b, v2))
}

func isConst(v Value) bool {
	_, ok := v.(*Const)
	return ok
}

func isRand(v Value) bool {
	_, ok := v.(*Rand)
	return ok
}

// marshalPattern is a shorthand to JSON encode a pattern.
func marshalPattern(p Pattern) []byte {
	b, err := json.Marshal(&SPattern{p})
	if err != nil {
		panic(err)
	}
	return b
}
