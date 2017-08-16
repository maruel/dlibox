// Copyright 2017 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package rules

import (
	"github.com/maruel/dlibox/modules"
	"periph.io/x/periph/conn/ir"
)

// Signal defines what signal triggers this rule.
//
// TODO(maruel): Define the language.
type Signal string

func (s Signal) Eval(input string) bool {
	return string(s) == input
}

type Rule struct {
	Signal Signal
	Cmd    modules.Command
}

type Rules map[string]Rule

func Default() []Rule {
	return []Rule{
		{Signal("ir/" + string(ir.KEY_CHANNELDOWN)), modules.Command{"leds/temperature", "-500"}},
		{Signal("ir/" + string(ir.KEY_CHANNEL)), modules.Command{"leds/temperature", "5000"}},
		{Signal("ir/" + string(ir.KEY_CHANNELUP)), modules.Command{"leds/temperature", "+500"}},
		{Signal("ir/" + string(ir.KEY_PREVIOUS)), modules.Command{"leds/temperature", "3000"}},
		{Signal("ir/" + string(ir.KEY_NEXT)), modules.Command{"leds/temperature", "5000"}},
		{Signal("ir/" + string(ir.KEY_PLAYPAUSE)), modules.Command{"leds/temperature", "6500"}},
		{Signal("ir/" + string(ir.KEY_VOLUMEDOWN)), modules.Command{"leds/intensity", "-15"}},
		{Signal("ir/" + string(ir.KEY_VOLUMEUP)), modules.Command{"leds/intensity", "+15"}},
		{Signal("ir/" + string(ir.KEY_EQ)), modules.Command{"leds/intensity", "128"}},
		{Signal("ir/" + string(ir.KEY_NUMERIC_0)), modules.Command{"leds/intensity", "0"}},
		{Signal("ir/" + string(ir.KEY_100PLUS)), modules.Command{"painter/setuser", "\"#ffffff\""}},
		{Signal("ir/" + string(ir.KEY_200PLUS)), modules.Command{"leds/intensity", "255"}},
		{Signal("ir/" + string(ir.KEY_NUMERIC_1)), modules.Command{"painter/setuser", "\"Rainbow\""}},
		{Signal("ir/" + string(ir.KEY_NUMERIC_2)), modules.Command{"painter/setuser", "{\"Child\":\"Rainbow\",\"MovePerHour\":108000,\"_type\":\"Rotate\"}"}},
		{Signal("ir/" + string(ir.KEY_NUMERIC_3)), modules.Command{"painter/setuser", "{\"Child\":{\"Frame\":\"Lff0000ff0000ff0000ff0000ff0000ffffffffffffffffffffffffffffff\",\"_type\":\"Repeated\"},\"MovePerHour\":21600,\"_type\":\"Rotate\"}"}},
		{Signal("ir/" + string(ir.KEY_NUMERIC_4)), modules.Command{"painter/setuser", "{\"Child\":\"L0100010f0000000f0000000f\",\"_type\":\"Chronometer\"}"}},
		{Signal("ir/" + string(ir.KEY_NUMERIC_5)), modules.Command{"painter/setuser", "{\"Child\":\"Lff0000ff0000ee0000dd0000cc0000bb0000aa0000990000880000770000660000550000440000330000220000110000\",\"MovePerHour\":108000,\"_type\":\"PingPong\"}"}},
		{Signal("ir/" + string(ir.KEY_NUMERIC_6)), modules.Command{"painter/setuser", "{\"C\":\"#ff9000\",\"_type\":\"NightStars\"}"}},
		{Signal("ir/" + string(ir.KEY_NUMERIC_7)), modules.Command{"painter/setuser", "{\"Curve\":\"ease-out\",\"Patterns\":[{\"Patterns\":[{\"Child\":\"Lff0000ff0000ee0000dd0000cc0000bb0000aa0000990000880000770000660000550000440000330000220000110000\",\"MovePerHour\":108000,\"_type\":\"Rotate\"},{\"_type\":\"Aurore\"}],\"_type\":\"Add\"},{\"Patterns\":[{\"_type\":\"Aurore\"},{\"C\":\"#ffffff\",\"_type\":\"NightStars\"}],\"_type\":\"Add\"}],\"ShowMS\":10000,\"TransitionMS\":5000,\"_type\":\"Loop\"}"}},
		{Signal("ir/" + string(ir.KEY_NUMERIC_8)), modules.Command{"painter/setuser", "{\"Left\":{\"Curve\":\"ease-out\",\"Patterns\":[\"#000f00\",\"#00ff00\",\"#1f0f00\",\"#ffa900\"],\"ShowMS\":100,\"TransitionMS\":700,\"_type\":\"Loop\"},\"Offset\":\"50%\",\"Right\":{\"Curve\":\"ease-out\",\"Patterns\":[\"#1f0f00\",\"#ffa900\",\"#000f00\",\"#00ff00\"],\"ShowMS\":100,\"TransitionMS\":700,\"_type\":\"Loop\"},\"_type\":\"Split\"}"}},
		{Signal("ir/" + string(ir.KEY_NUMERIC_9)), modules.Command{"painter/setuser", "{\"Child\":\"Lffffff\",\"MovePerHour\":108000,\"_type\":\"PingPong\"}"}},
	}
}
