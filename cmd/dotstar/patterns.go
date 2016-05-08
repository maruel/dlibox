// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"image/color"

	"github.com/maruel/dotstar/anim1d"
)

var red = color.NRGBA{255, 0, 0, 255}
var white = color.NRGBA{255, 255, 255, 255}

func getRegistry() *anim1d.PatternRegistry {
	return &anim1d.PatternRegistry{
		Patterns: map[string]anim1d.Pattern{
			"Canne de Noël":        &anim1d.Repeated{[]color.NRGBA{red, red, red, red, white, white, white, white}, 6},
			"K2000":                &anim1d.PingPong{anim1d.K2000, color.NRGBA{0, 0, 0, 255}, 30},
			"Glow rainbow":         &anim1d.Glow{anim1d.RainbowColors, 1. / 3.},
			"Ping pong":            &anim1d.PingPong{Trail: []color.NRGBA{{255, 255, 255, 255}}, MovesPerSec: 30},
			"Rainbow static":       &anim1d.Rainbow{},
			"Étoiles cintillantes": &anim1d.ÉtoilesCintillantes{},
			"Ciel étoilé": &anim1d.Mixer{
				Patterns: []anim1d.Pattern{
					&anim1d.Aurore{},
					&anim1d.ÉtoilesCintillantes{},
					&anim1d.ÉtoileFilante{},
				},
				Weights: []float32{1, 1, 1},
			},
			"Dégradé":       &anim1d.Dégradé{color.NRGBA{0, 0, 0, 255}, color.NRGBA{255, 255, 255, 255}},
			"Dégradé rouge": &anim1d.Dégradé{color.NRGBA{0, 0, 0, 255}, color.NRGBA{255, 0, 0, 255}},
			"Dégradé vert":  &anim1d.Dégradé{color.NRGBA{0, 0, 0, 255}, color.NRGBA{0, 255, 0, 255}},
			"Dégradé bleu":  &anim1d.Dégradé{color.NRGBA{0, 0, 0, 255}, color.NRGBA{0, 0, 255, 255}},
			"Aurores":       &anim1d.Aurore{},
		},
		NumberLEDs:       150,
		ThumbnailHz:      10,
		ThumbnailSeconds: 10,
	}
}

//"Dégradé png": anim1d.LoadAnimate(mustRead("dégradé.png"), 16*time.Millisecond, false),
//"Étoile floue":         anim1d.LoadAnimate(mustRead("étoile floue.png"), 16*time.Millisecond, false),
//"Glow":                 &anim1d.Glow{[]color.NRGBA{{255, 255, 255, 255}, {0, 128, 0, 255}}, 1},
//"Glow gris":            &anim1d.Glow{[]color.NRGBA{{255, 255, 255, 255}, {}}, 0.33},
//"Ping pong bleue": &anim1d.PingPong{
//	Trail: []color.NRGBA{
//		{0xff, 0xff, 255, 255},
//		{0xD7, 0xD7, 255, 255},
//		{0xAF, 0xAF, 255, 255},
//		{0x87, 0x87, 255, 255},
//		{0x5F, 0x5F, 255, 255},
//	},
//	MovesPerSec: 30,
//},
