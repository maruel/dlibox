// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"image/color"
	"time"

	"github.com/maruel/dotstar"
)

func getRegistry() *dotstar.PatternRegistry {
	red := color.NRGBA{255, 0, 0, 255}
	white := color.NRGBA{255, 255, 255, 255}
	return &dotstar.PatternRegistry{
		Patterns: map[string]dotstar.Pattern{
			"Canne de Noël": &dotstar.Repeated{[]color.NRGBA{red, red, red, red, white, white, white, white}, 6},
			"Ping pong bleue": &dotstar.PingPong{
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
			"Glow":                 &dotstar.Glow{[]color.NRGBA{{255, 255, 255, 255}, {0, 128, 0, 255}}, 1},
			"Glow gris":            &dotstar.Glow{[]color.NRGBA{{255, 255, 255, 255}, {}}, 0.33},
			"Glow rainbow":         &dotstar.Glow{dotstar.RainbowColors, 1. / 3.},
			"Ping pong":            &dotstar.PingPong{Trail: []color.NRGBA{{255, 255, 255, 255}}, MovesPerSec: 30},
			"Rainbow static":       &dotstar.Rainbow{},
			"Étoiles cintillantes": &dotstar.ÉtoilesCintillantes{},
			"Étoile floue":         dotstar.LoadAnimate(mustRead("étoile floue.png"), 16*time.Millisecond, false),
			"Ciel étoilé": &dotstar.Mixer{
				Patterns: []dotstar.Pattern{
					&dotstar.Aurore{},
					&dotstar.ÉtoilesCintillantes{},
					&dotstar.ÉtoileFilante{},
				},
				Weights: []float64{1, 1, 1},
			},
			//"Dégradé png": dotstar.LoadAnimate(mustRead("dégradé.png"), 16*time.Millisecond, false),
			"Dégradé": &dotstar.Dégradé{color.NRGBA{0, 0, 0, 255}, color.NRGBA{255, 255, 255, 255}},
		},
		NumberLEDs:       150,
		ThumbnailHz:      10,
		ThumbnailSeconds: 10,
	}
}
