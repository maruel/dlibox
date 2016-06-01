// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package anim1d

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/maruel/ut"
)

func TestNilObject(t *testing.T) {
	c := []Color{{}}
	for _, p := range knownPatterns {
		p.NextFrame(nil, 0)
		p.NextFrame(c, 0)
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

func TestJSON(t *testing.T) {
	for _, p := range knownPatterns {
		p2 := &SPattern{p}
		b, err := json.Marshal(p2)
		ut.AssertEqual(t, nil, err)
		if _, ok := p.(*Color); ok {
			ut.AssertEqual(t, uint8('"'), b[0])
		} else {
			ut.AssertEqual(t, uint8('{'), b[0])
		}
		p2.Pattern = nil
		ut.AssertEqual(t, nil, json.Unmarshal(b, p2))
	}
	serialize(t, &Color{1, 2, 3, 4}, `"#01020304"`)
	serialize(t, &PingPong{}, `{"Background":"#00000000","MovesPerSec":0,"Trail":null,"_type":"PingPong"}`)
	serialize(t, &Animation{}, `{"FrameDuration":0,"Frames":null,"_type":"Animation"}`)

	// Create one more complex. Assert that int64 is not mangled.
	p := &Transition{
		Out: SPattern{
			&Transition{
				In:         SPattern{&Color{255, 255, 255, 255}},
				Offset:     10 * time.Minute,
				Duration:   10 * time.Minute,
				Transition: TransitionLinear,
			},
		},
		In:         SPattern{&Color{}},
		Offset:     30 * time.Minute,
		Duration:   10 * time.Minute,
		Transition: TransitionLinear,
	}
	expected := `{"Duration":600000000000,"In":"#00000000","Offset":1800000000000,"Out":{"Duration":600000000000,"In":"#ffffffff","Offset":600000000000,"Out":{},"Transition":"linear","_type":"Transition"},"Transition":"linear","_type":"Transition"}`
	serialize(t, p, expected)
}
