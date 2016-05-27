// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"image/color"
	"time"

	"github.com/maruel/dotstar/anim1d"
	"github.com/maruel/dotstar/anim1d/animio"
)

var red = color.NRGBA{255, 0, 0, 255}
var white = color.NRGBA{255, 255, 255, 255}

var rainbowColors []anim1d.Pattern

func init() {
	rainbowColors = make([]anim1d.Pattern, len(anim1d.RainbowColors))
	for i, c := range anim1d.RainbowColors {
		rainbowColors[i] = &anim1d.StaticColor{c}
	}
}

func getRegistry() *animio.PatternRegistry {
	return &animio.PatternRegistry{
		Patterns: map[string]anim1d.Pattern{
			"Canne de Noël": &anim1d.Repeated{
				[]color.NRGBA{red, red, red, red, red, white, white, white, white, white},
				6,
			},
			"K2000": &anim1d.PingPong{anim1d.K2000Colors, color.NRGBA{0, 0, 0, 255}, 30},
			"Ping pong": &anim1d.PingPong{
				Trail:       []color.NRGBA{{255, 255, 255, 255}},
				MovesPerSec: 30,
			},
			"Rainbow cycle": &anim1d.Loop{
				Patterns:           rainbowColors,
				DurationShow:       1 * time.Second,
				DurationTransition: 10 * time.Second,
				Transition:         anim1d.TransitionEaseInOut,
			},
			"Rainbow static":       &anim1d.Rainbow{},
			"Étoiles cintillantes": &anim1d.NightStars{},
			"Ciel étoilé": &anim1d.Mixer{
				Patterns: []anim1d.Pattern{
					&anim1d.Aurore{},
					&anim1d.NightStars{},
					&anim1d.WishingStar{},
				},
				Weights: []float32{1, 1, 1},
			},
			"Aurores": &anim1d.Aurore{},
			// Transition from black to orange to white then to black.
			"Morning alarm": &anim1d.Transition{
				Out: &anim1d.Transition{
					Out: &anim1d.Transition{
						Out:        &anim1d.StaticColor{},
						In:         &anim1d.StaticColor{color.NRGBA{255, 127, 0, 255}},
						Duration:   10 * time.Minute,
						Transition: anim1d.TransitionLinear,
					},
					In:         &anim1d.StaticColor{color.NRGBA{255, 255, 255, 255}},
					Offset:     10 * time.Minute,
					Duration:   10 * time.Minute,
					Transition: anim1d.TransitionLinear,
				},
				In:         &anim1d.StaticColor{},
				Offset:     30 * time.Minute,
				Duration:   10 * time.Minute,
				Transition: anim1d.TransitionLinear,
			},
			// Test de couleurs:
			"Cycle RGB": &anim1d.Loop{
				Patterns: []anim1d.Pattern{
					&anim1d.StaticColor{color.NRGBA{255, 0, 0, 255}},
					&anim1d.StaticColor{color.NRGBA{0, 255, 0, 255}},
					&anim1d.StaticColor{color.NRGBA{0, 0, 255, 255}},
				},
				DurationShow:       1000 * time.Millisecond,
				DurationTransition: 1000 * time.Millisecond,
				Transition:         anim1d.TransitionEaseInOut,
			},
			"Dégradé":       &anim1d.Gradient{color.NRGBA{0, 0, 0, 255}, color.NRGBA{255, 255, 255, 255}},
			"Dégradé rouge": &anim1d.Gradient{color.NRGBA{0, 0, 0, 255}, color.NRGBA{255, 0, 0, 255}},
			"Dégradé vert":  &anim1d.Gradient{color.NRGBA{0, 0, 0, 255}, color.NRGBA{0, 255, 0, 255}},
			"Dégradé bleu":  &anim1d.Gradient{color.NRGBA{0, 0, 0, 255}, color.NRGBA{0, 0, 255, 255}},
		},
		NumberLEDs:       100,
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
