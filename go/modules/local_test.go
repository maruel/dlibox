// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package modules

import "testing"

func TestParseTopicGood(t *testing.T) {
	data := []string{
		"//",
		"+/+/+",
		"+/+/#",
		"sport/tennis/+",
		"sport/tennis/#",
		"sport/tennis/player1",
		"sport/tennis/player1/ranking",
		"sport/tennis/+/score/wimbledon",
	}
	for i, line := range data {
		if parseTopic(line) == nil {
			t.Fatalf("%d: parseTopic(%#v) returned nil", i, line)
		}
	}
}

func TestParseTopicBad(t *testing.T) {
	data := []string{
		"",
		"sport/tennis#",
		"sport/tennis/#/ranking",
		"sport/tennis+",
	}
	for i, line := range data {
		if parseTopic(line) != nil {
			t.Fatalf("%d: parseTopic(%#v) returned non nil", i, line)
		}
	}
}

func TestMatchSuccess(t *testing.T) {
	data := [][2]string{
		{"sport/tennis/#", "sport/tennis"},
		{"sport/tennis/#", "sport/tennis/player1"},
		{"sport/tennis/#", "sport/tennis/player1/ranking"},
		{"sport/tennis/+", "sport/tennis"},
		{"sport/tennis/+", "sport/tennis/player1"},
		{"sport/+/player1", "sport/tennis/player1"},
	}
	for i, line := range data {
		q := parseTopic(line[0])
		if !q.match(line[1]) {
			t.Fatalf("%d: %#v.match(%#v) returned false", i, line[0], line[1])
		}
	}
}

func TestMatchFail(t *testing.T) {
	data := [][2]string{
		{"sport/tennis/#", "sport"},
		{"sport/tennis/#", "sport/badminton"},
		{"sport/tennis/+", "sport"},
		{"sport/tennis/+", "sport/badminton"},
		{"sport/+/player1", "sport/tennis/player2"},
	}
	for i, line := range data {
		q := parseTopic(line[0])
		if q.match(line[1]) {
			t.Fatalf("%d: %#v.match(%#v) returned true", i, line[0], line[1])
		}
	}
}
