// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"image/color"

	"github.com/maruel/dotstar"
)

var Registry = dotstar.PatternRegistry{
	NumberLEDs:       150,
	ThumbnailHz:      10,
	ThumbnailSeconds: 10,
}

func init() {
	red := color.NRGBA{255, 0, 0, 255}
	white := color.NRGBA{255, 255, 255, 255}
	Registry.Patterns = map[string]dotstar.Pattern{
		"black": &dotstar.StaticColor{},
		"canne": &dotstar.Repeated{[]color.NRGBA{red, red, red, red, white, white, white, white}, 6},
		"pingpong bleue": &dotstar.PingPong{
			Trail: []color.NRGBA{
				{0xff, 0xff, 255, 255},
				{0xD7, 0xD7, 255, 255},
				{0xAF, 0xAF, 255, 255},
				{0x87, 0x87, 255, 255},
				{0x5F, 0x5F, 255, 255},
			},
			MovesPerSec: 30,
		},
		"K2000":                &dotstar.PingPong{dotstar.K2000, color.NRGBA{0, 0, 0, 255}, 30},
		"glow":                 &dotstar.Glow{[]color.NRGBA{{255, 255, 255, 255}, {0, 128, 0, 255}}, 1},
		"glow gris":            &dotstar.Glow{[]color.NRGBA{{255, 255, 255, 255}, {}}, 0.33},
		"glow rainbow":         &dotstar.Glow{dotstar.RainbowColors, 1.},
		"pingpong":             &dotstar.PingPong{Trail: []color.NRGBA{{255, 255, 255, 255}}, MovesPerSec: 30},
		"rainbow static":       &dotstar.Rainbow{},
		"étoiles cintillantes": &dotstar.ÉtoilesCintillantes{},
		"ciel étoilé": &dotstar.Mixer{
			Patterns: []dotstar.Pattern{
				&dotstar.Aurore{},
				&dotstar.ÉtoilesCintillantes{},
				&dotstar.ÉtoileFilante{},
			},
			Weights: []float64{1, 1, 1},
		},
	}
}
