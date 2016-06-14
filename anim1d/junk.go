// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// These are incomplete and will either be removed or fixed.

package anim1d

import (
	"math/rand"
	"time"
)

type point struct {
	star  int
	start time.Time
}

// NightSky has:
//    - Stars
//    - WishingStar
//    - Aurores
//    - Super nova.
//    - Rotation de la terre?
//    - Station Internationale?
type NightSky struct {
	Stars     []Cycle
	Frequency float32 // Number of explosions by second.
	points    []point
}

func (c *NightSky) NextFrame(pixels Frame, timeMS uint32) {
	// random
	// animate.
}

// Aurore commence lentement, se transforme lentement et éventuellement
// disparait.
type Aurore struct {
}

func (a *Aurore) NextFrame(pixels Frame, timeMS uint32) {
	// TODO(maruel): Redo.
	y := float32(timeMS) * .01
	for i := range pixels {
		x := float32(i)
		//a := 32 + 31*sin(x/(37.+15*cos(y/74)))*cos(y/(31+11*sin(x/57)))
		b := (32 + 31*(sin(hypot(200-y, 320-x)/16))) * (0.5 + 0.5*sin(y*0.1))
		pixels[i].R = 0
		//pixels[i].G = uint8(a + b)
		pixels[i].G = uint8(b)
		pixels[i].B = 0
	}
}

type NightStar struct {
	Intensity uint8
	Type      int
}

type NightStars struct {
	Stars []NightStar
	Seed  int // Change it to create a different pseudo-random animation.
	r     *rand.Rand
}

func (e *NightStars) NextFrame(pixels Frame, timeMS uint32) {
	if e.r == nil {
		e.r = rand.New(rand.NewSource(int64(e.Seed)))
	}
	if len(e.Stars) != len(pixels) {
		e.Stars = make([]NightStar, len(pixels))
		for i := 0; i < len(pixels); {
			// Add a star. Decide it's relative position, intensity and type.
			// ExpFloat64() ?
			f := abs(3 * float32(e.r.NormFloat64()))
			if f < 1 {
				continue
			}
			i += int(roundF(f))
			if i >= len(pixels) {
				break
			}
			// e.r.Intn(255)
			intensity := abs(float32(e.r.NormFloat64()))
			if intensity > 255 {
				intensity = 0
			}
			e.Stars[i].Intensity = FloatToUint8(intensity)
		}
	}
	for i := range e.Stars {
		if j := e.Stars[i].Intensity; j != 0 {
			// TODO(maruel): Type, oscillation.
			if j != 0 {
				f := FloatToUint8(float32(e.r.NormFloat64())*4 + float32(j))
				pixels[i] = Color{f, f, f}
			}
		}
	}
}

// WishingStar draws a wishing star from time to time.
//
// It will only draw one star at a time. To increase the likelihood of getting
// many simultaneously, create multiple instances and use Mixer with Weights of
// 1.
type WishingStar struct {
	Duration     time.Duration // Average duration of a star.
	AverageDelay time.Duration // Average delay between each wishing star.
}

func (w *WishingStar) NextFrame(pixels Frame, timeMS uint32) {
	/*
		// Create a deterministic replay by using the current number of
		// the wishing star as the seed for the current flow. Make it independent of
		// any other non-deterministic source.
		i := timeMS / w.AverageDelay
		r := rand.New(rand.NewSource(int64(i)))
		// Always calculate things in the same order to keep the calculation
		// deterministic.
		startOffset := r.Int63()
		startPos := r.Int63()
		intensity := r.Int63()
		orientation := r.Intn(2)
		// Draw according to these parameters.
		// - Trail
		// - Observed speed based on orientation
	*/
}
