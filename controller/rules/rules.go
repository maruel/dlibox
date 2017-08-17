// Copyright 2017 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Package rules defines the generic rules engine.
package rules

import (
	"periph.io/x/periph/conn/ir"
)

// Signal defines what signal triggers this rule.
//
// TODO(maruel): Define the language.
type Signal string

// Eval evaluates if a signal is a trigger for this rule.
func (s Signal) Eval(input string) bool {
	return string(s) == input
}

// Rule defines a signal that triggers a command.
type Rule struct {
	Signal Signal
	Cmd    Command
}

// Rules is named rules.
type Rules map[string]Rule

// Default returns default rules that can be set on a fresh instance.
func Default() []Rule {
	return []Rule{
		{Signal("ir/" + string(ir.KEY_CHANNELDOWN)), Command{"leds/temperature", "-500"}},
		{Signal("ir/" + string(ir.KEY_CHANNEL)), Command{"leds/temperature", "5000"}},
		{Signal("ir/" + string(ir.KEY_CHANNELUP)), Command{"leds/temperature", "+500"}},
		{Signal("ir/" + string(ir.KEY_PREVIOUS)), Command{"leds/temperature", "3000"}},
		{Signal("ir/" + string(ir.KEY_NEXT)), Command{"leds/temperature", "5000"}},
		{Signal("ir/" + string(ir.KEY_PLAYPAUSE)), Command{"leds/temperature", "6500"}},
		{Signal("ir/" + string(ir.KEY_VOLUMEDOWN)), Command{"leds/intensity", "-15"}},
		{Signal("ir/" + string(ir.KEY_VOLUMEUP)), Command{"leds/intensity", "+15"}},
		{Signal("ir/" + string(ir.KEY_EQ)), Command{"leds/intensity", "128"}},
		{Signal("ir/" + string(ir.KEY_NUMERIC_0)), Command{"leds/intensity", "0"}},
		{Signal("ir/" + string(ir.KEY_100PLUS)), Command{"painter/setuser", "\"#ffffff\""}},
		{Signal("ir/" + string(ir.KEY_200PLUS)), Command{"leds/intensity", "255"}},
		{Signal("ir/" + string(ir.KEY_NUMERIC_1)), Command{"painter/setuser", "\"Rainbow\""}},
		{Signal("ir/" + string(ir.KEY_NUMERIC_2)), Command{"painter/setuser", "{\"Child\":\"Rainbow\",\"MovePerHour\":108000,\"_type\":\"Rotate\"}"}},
		{Signal("ir/" + string(ir.KEY_NUMERIC_3)), Command{"painter/setuser", "{\"Child\":{\"Frame\":\"Lff0000ff0000ff0000ff0000ff0000ffffffffffffffffffffffffffffff\",\"_type\":\"Repeated\"},\"MovePerHour\":21600,\"_type\":\"Rotate\"}"}},
		{Signal("ir/" + string(ir.KEY_NUMERIC_4)), Command{"painter/setuser", "{\"Child\":\"L0100010f0000000f0000000f\",\"_type\":\"Chronometer\"}"}},
		{Signal("ir/" + string(ir.KEY_NUMERIC_5)), Command{"painter/setuser", "{\"Child\":\"Lff0000ff0000ee0000dd0000cc0000bb0000aa0000990000880000770000660000550000440000330000220000110000\",\"MovePerHour\":108000,\"_type\":\"PingPong\"}"}},
		{Signal("ir/" + string(ir.KEY_NUMERIC_6)), Command{"painter/setuser", "{\"C\":\"#ff9000\",\"_type\":\"NightStars\"}"}},
		{Signal("ir/" + string(ir.KEY_NUMERIC_7)), Command{"painter/setuser", "{\"Curve\":\"ease-out\",\"Patterns\":[{\"Patterns\":[{\"Child\":\"Lff0000ff0000ee0000dd0000cc0000bb0000aa0000990000880000770000660000550000440000330000220000110000\",\"MovePerHour\":108000,\"_type\":\"Rotate\"},{\"_type\":\"Aurore\"}],\"_type\":\"Add\"},{\"Patterns\":[{\"_type\":\"Aurore\"},{\"C\":\"#ffffff\",\"_type\":\"NightStars\"}],\"_type\":\"Add\"}],\"ShowMS\":10000,\"TransitionMS\":5000,\"_type\":\"Loop\"}"}},
		{Signal("ir/" + string(ir.KEY_NUMERIC_8)), Command{"painter/setuser", "{\"Left\":{\"Curve\":\"ease-out\",\"Patterns\":[\"#000f00\",\"#00ff00\",\"#1f0f00\",\"#ffa900\"],\"ShowMS\":100,\"TransitionMS\":700,\"_type\":\"Loop\"},\"Offset\":\"50%\",\"Right\":{\"Curve\":\"ease-out\",\"Patterns\":[\"#1f0f00\",\"#ffa900\",\"#000f00\",\"#00ff00\"],\"ShowMS\":100,\"TransitionMS\":700,\"_type\":\"Loop\"},\"_type\":\"Split\"}"}},
		{Signal("ir/" + string(ir.KEY_NUMERIC_9)), Command{"painter/setuser", "{\"Child\":\"Lffffff\",\"MovePerHour\":108000,\"_type\":\"PingPong\"}"}},
	}
}
