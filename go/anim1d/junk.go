// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// These are incomplete and will either be removed or fixed.

package anim1d

import (
	"math/rand"
	"time"
)

// TODO(maruel): Create NightSky with:
//    - Stars
//    - WishingStar
//    - Aurores
//    - Super nova.
//    - Rotation de la terre?
//    - Station Internationale?

// Aurore commence lentement, se transforme lentement et éventuellement
// disparait.
type Aurore struct {
}

func (a *Aurore) Render(pixels Frame, timeMS uint32) {
	// TODO(maruel): Redo.
	y := float32(timeMS) * .01
	for i := range pixels {
		x := float32(i)
		//a := 32 + 31*sin(x/(37.+15*cos(y/74)))*cos(y/(31+11*sin(x/57)))
		b := (128 + 127*(sin(hypot(200-y, 320-x)/16))) * (0.5 + 0.5*sin(y*0.1))
		pixels[i].R = 0
		//pixels[i].G = uint8(a + b)
		pixels[i].G = uint8(b)
		pixels[i].B = 0
	}
}

type NightStars struct {
	C     Color
	stars Frame
}

func (n *NightStars) Render(pixels Frame, timeMS uint32) {
	if len(n.stars) != len(pixels) {
		r := rand.NewSource(0)
		n.stars = make(Frame, len(pixels))
		for i := range n.stars {
			j := int32(r.Int63())
			// Cut off at 25%.
			if j&0x30000 != 0x30000 {
				continue
			}
			// Use gamma == 2 and limit intensity at 50%.
			d := int(j&0xff+1) * int((j>>8)&0xff+1)
			n.stars[i] = n.C
			n.stars[i].Dim(uint8((d-1)>>8) / 2)
		}
	}

	r := rand.NewSource(int64((&Rand{}).Eval(timeMS, len(pixels))))
	copy(pixels, n.stars)
	for i := range n.stars {
		j := uint8(r.Int63())
		// Use gamma == 2.
		d := int32(j&0xf+1) * int32((j>>4)+1)
		pixels[i].Dim(255 - uint8((d-1)>>4))
	}
}

type Lightning struct {
	Center    SValue // offset of the center, from the left
	HalfWidth SValue // in pixels
	Intensity int    // the maximum intensity
	StartMS   SValue // when it started
}

var lightningCycle = []struct {
	offsetMS  uint32
	intensity uint8
}{
	{0, 0},
	{150, 255},
	{300, 0},
	{450, 255},
	{600, 0},
	{750, 255},
	{900, 0},
	{1050, 76},
	{1200, 51},
	{1350, 26},
	{1500, 0},
}

func (l *Lightning) Render(pixels Frame, timeMS uint32) {
	// Will fail after 25 days.
	offset := timeMS - uint32(l.StartMS.Eval(timeMS, len(pixels)))
	intensity := uint8(0)
	for i := 0; i < len(lightningCycle); i++ {
		if lightningCycle[i].offsetMS > offset {
			intensity = lightningCycle[i-1].intensity
			break
		}
	}
	if intensity == 0 {
		return
	}
	center := l.Center.Eval(timeMS, len(pixels))
	halfWidth := l.HalfWidth.Eval(timeMS, len(pixels))
	left := center - halfWidth
	right := center + halfWidth
	width := left - right
	min := MinMax32(left, 0, int32(len(pixels)-1))
	max := MinMax32(right, 0, int32(len(pixels)-1))
	b := Bell{}
	for i := min; i < max; i++ {
		x := (i - left) * 65535 / width
		pixels[i] = Color{intensity, intensity, intensity}
		pixels[i].Dim(uint8(b.Scale(uint16(x)) >> 8))
	}
}

// Thunderstorm creates strobe-like lightning.
type Thunderstorm struct {
	AvgMS   int // Average between lightning strikes
	current []Lightning
	nextMS  uint32
}

func (t *Thunderstorm) Render(pixels Frame, timeMS uint32) {
	/*
		//freq := 3
		if t.current == nil {
			if timeMS != 0 {
				// TODO(maruel): Backfill for determinism. Will skip for now.
			}
			t.current = []Lightning{}
			r := rand.NewSource(0)
			t.nextMS = uint32((&Bell{}).Eval(uint16(r.Int63())))
		}
		for timeMS > t.nextMS {
			t.current = append(t.current, Lightning{})
			r := rand.NewSource(0)
			t.nextMS = uint32((&Bell{}).Eval(uint16(r.Int63()))) + t.nextMS
		}
		// Calculate all triggers up to now.
		// Calculate location up to now.
		// Create one random
		r := rand.NewSource(int64((&Rand{}).Eval(timeMS)))
		// TODO(maruel): Slight coloring?
		for i := range pixels {
			pixels[i] = Color{}
		}
	*/
}

// WishingStar draws a wishing star from time to time.
//
// It will only draw one star at a time. To increase the likelihood of getting
// many simultaneously, create multiple instances.
type WishingStar struct {
	Duration     time.Duration // Average duration of a star.
	AverageDelay time.Duration // Average delay between each wishing star.
}

func (w *WishingStar) Render(pixels Frame, timeMS uint32) {
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
