// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package apa102

import (
	"io"
	"time"

	"github.com/maruel/dlibox/go/anim1d"
	"github.com/maruel/temperature"
)

// maxOut is the maximum intensity of each channel on a APA102 LED.
const maxOut = 0x1EE1

// ramp converts input from [0, 0xFF] as intensity to lightness on a scale of
// [0, maxOut] or other desired range [0, max].
//
// It tries to use the same curve independent of the scale used. max can be
// changed to change the color temperature or to limit power dissipation.
//
// It's the reverse of lightness; https://en.wikipedia.org/wiki/Lightness
func ramp(l uint8, max uint16) uint16 {
	if l == 0 {
		// Make sure black is black.
		return 0
	}
	// linearCutOff defines the linear section of the curve. Inputs between
	// [0, linearCutOff] are mapped linearly to the output. It is 1% of maximum
	// output.
	linearCutOff := uint32((max + 50) / 100)
	l32 := uint32(l)
	if l32 < linearCutOff {
		return uint16(l32)
	}

	// Maps [linearCutOff, 255] to use [linearCutOff*max/255, max] using a x^3
	// ramp.
	// Realign input to [0, 255-linearCutOff]. It now maps to
	// [0, max-linearCutOff*max/255].
	//const inRange = 255
	l32 -= linearCutOff
	inRange := 255 - linearCutOff
	outRange := uint32(max) - linearCutOff
	offset := inRange >> 1
	y := (l32*l32*l32 + offset) / inRange
	return uint16((y*outRange+(offset*offset))/inRange/inRange + linearCutOff)
}

// rampTable is a lookup table generated by calling ramp() on each intensity.
type rampTable [256]uint16

// rampCache is actually a leak, it never gets cleaned.
var rampCache = map[uint16]*rampTable{}

// ensureRampCached makes sure the ramp LUT for 'max' is precalculated and
// returns it.
func ensureRampCached(max uint16) *rampTable {
	if r, ok := rampCache[max]; ok {
		return r
	}
	r := &rampTable{}
	for i := 0; i <= 255; i++ {
		r[uint8(i)] = ramp(uint8(i), max)
	}
	rampCache[max] = r
	return r
}

// Serializes converts a buffer of colors to the APA102 SPI format.
func raster(pixels anim1d.Frame, buf *[]byte, maxR, maxG, maxB uint16) {
	// https://cpldcpu.files.wordpress.com/2014/08/apa-102c-super-led-specifications-2014-en.pdf
	numLights := len(pixels)
	// End frames are needed to be able to push enough SPI clock signals due to
	// internal half-delay of data signal from each individual LED. See
	// https://cpldcpu.wordpress.com/2014/11/30/understanding-the-apa102-superled/
	l := 4*(numLights+1) + numLights/2/8 + 1
	if len(*buf) != l {
		*buf = make([]byte, l)
		// It is not necessary to set the end frames to 0xFFFFFFFF.
		// Set end frames right away.
		s := (*buf)[4+4*numLights:]
		for i := range s {
			s[i] = 0xFF
		}
	}
	// Make sure the ramps are cached.
	rampR := ensureRampCached(maxR)
	rampG := ensureRampCached(maxG)
	rampB := ensureRampCached(maxB)

	// Start frame is all zeros. Just skip it.
	s := (*buf)[4 : 4+4*numLights]
	for i, c := range pixels {
		// Converts a color into the 4 bytes needed to control an APA-102 LED.
		//
		// The response as seen by the human eye is very non-linear. The APA-102
		// provides an overall brightness PWM but it is relatively slower and
		// results in human visible flicker. On the other hand the minimal color
		// (1/255) is still too intense at full brightness, so for very dark color,
		// it is worth using the overall brightness PWM. The goal is to use
		// brightness!=31 as little as possible.
		//
		// Global brightness frequency is 580Hz and color frequency at 19.2kHz.
		// https://cpldcpu.wordpress.com/2014/08/27/apa102/
		// Both are multiplicative, so brightness@50% and color@50% means an
		// effective 25% duty cycle but it is not properly distributed, which is
		// the main problem.
		//
		// It is unclear to me if brightness is exactly in 1/31 increment as I don't
		// have an oscilloscope to confirm. Same for color in 1/255 increment.
		// TODO(maruel): I have one now!
		//
		// Each channel duty cycle ramps from 100% to 1/(31*255) == 1/7905.
		//
		// Computes brighness, blue, green, red.
		r := rampR[c.R]
		g := rampG[c.G]
		b := rampB[c.B]
		m := r | g | b
		if m <= 1023 {
			if m <= 255 {
				s[4*i], s[4*i+1], s[4*i+2], s[4*i+3] = byte(0xE0+1), byte(b), byte(g), byte(r)
			} else if m <= 511 {
				s[4*i], s[4*i+1], s[4*i+2], s[4*i+3] = byte(0xE0+2), byte(b>>1), byte(g>>1), byte(r>>1)
			} else {
				s[4*i], s[4*i+1], s[4*i+2], s[4*i+3] = byte(0xE0+4), byte((b+2)>>2), byte((g+2)>>2), byte((r+2)>>2)
			}
		} else {
			// In this case we need to use a ramp of 255-1 even for lower colors.
			s[4*i], s[4*i+1], s[4*i+2], s[4*i+3] = byte(0xE0+31), byte((b+15)/31), byte((g+15)/31), byte((r+15)/31)
		}
	}
}

type APA102 struct {
	Intensity   uint8  // Set an intensity between 0 (off) and 255 (full brightness).
	Temperature uint16 // In Kelvin.
	w           io.Writer
	buf         []byte
}

func (a *APA102) Write(pixels anim1d.Frame) error {
	tr, tg, tb := temperature.ToRGB(a.Temperature)
	r := uint16((uint32(maxOut)*uint32(a.Intensity)*uint32(tr) + 127*127) / 65025)
	g := uint16((uint32(maxOut)*uint32(a.Intensity)*uint32(tg) + 127*127) / 65025)
	b := uint16((uint32(maxOut)*uint32(a.Intensity)*uint32(tb) + 127*127) / 65025)
	raster(pixels, &a.buf, r, g, b)
	_, err := a.w.Write(a.buf)
	return err
}

func (a *APA102) MinDelay() time.Duration {
	// As per APA102-C spec, it's max refresh rate is 400hz.
	// https://en.wikipedia.org/wiki/Flicker_fusion_threshold is a recommended
	// reading.
	return time.Second / 400
}

// MakeAPA102 returns a strip that communicates over SPI to APA102 LEDs.
//
// w should be a *SPI as returned by rpi.MakeSPI. The speed must be high, as
// there's 32 bits sent per LED, creating a staggered effect. See
// https://cpldcpu.wordpress.com/2014/11/30/understanding-the-apa102-superled/
func MakeAPA102(w io.Writer) *APA102 {
	return &APA102{
		Intensity:   255,
		Temperature: 6500,
		w:           w,
	}
}
