// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package apa102

import (
	"io"
	"time"

	"github.com/maruel/dlibox-go/anim1d"
	"github.com/maruel/dlibox-go/rpi"
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

// colorToAPA102 converts a color into the 4 bytes needed to control an APA-102
// LED.
//
// The response as seen by the human eye is very non-linear. The APA-102
// provides an overall brightness PWM but it is relatively slower and results
// in human visible flicker. On the other hand the minimal color (1/255) is
// still too intense at full brightness, so for very dark color, it is worth
// using the overall brightness PWM. The goal is to use brightness!=31 as
// little as possible.
//
// Global brightness frequency is 580Hz and color frequency at 19.2kHz.
// https://cpldcpu.wordpress.com/2014/08/27/apa102/
// Both are multiplicative, so brightness@50% and color@50% means an effective
// 25% duty cycle but it is not properly distributed, which is the main problem.
//
// It is unclear to me if brightness is exactly in 1/31 increment as I don't
// have an oscilloscope to confirm. Same for color in 1/255 increment.
//
// Each channel duty cycle ramps from 100% to 1/(31*255) == 1/7905.
//
// Return brighness, blue, green, red.
func colorToAPA102(c anim1d.Color, max uint16) (byte, byte, byte, byte) {
	r := ramp(c.R, max)
	g := ramp(c.G, max)
	b := ramp(c.B, max)
	if r <= 255 && g <= 255 && b <= 255 {
		return byte(0xE0 + 1), byte(b), byte(g), byte(r)
	} else if r <= 511 && g <= 511 && b <= 511 {
		return byte(0xE0 + 2), byte(b >> 1), byte(g >> 1), byte(r >> 1)
	} else if r <= 1023 && g <= 1023 && b <= 1023 {
		return byte(0xE0 + 4), byte((b + 2) >> 2), byte((g + 2) >> 2), byte((r + 2) >> 2)
	} else {
		// In this case we need to use a ramp of 255-1 even for lower colors.
		return byte(0xE0 + 31), byte((b + 15) / 31), byte((g + 15) / 31), byte((r + 15) / 31)
	}
}

// Serializes converts a buffer of colors to the APA102 SPI format.
func raster(pixels anim1d.Frame, buf *[]byte, max uint16) {
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
	// Start frame is all zeros. Just skip it.
	s := (*buf)[4 : 4+4*numLights]
	for i := range pixels {
		s[4*i], s[4*i+1], s[4*i+2], s[4*i+3] = colorToAPA102(pixels[i], max)
	}
}

type APA102 struct {
	Intensity   uint8  // Set an intensity between 0 (off) and 255 (full brightness).
	Temperature uint16 // In Kelvin.
	w           io.WriteCloser
	buf         []byte
}

func (a *APA102) Close() error {
	w := a.w
	a.w = nil
	return w.Close()
}

func (a *APA102) Write(pixels anim1d.Frame) error {
	// TODO(maruel): Calculate power in duty cycle of each channel.
	raster(pixels, &a.buf, uint16((uint32(maxOut)*uint32(a.Intensity)+127)/255))
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
// This is generally what you want once the hardware is connected.
func MakeAPA102(speed int64) (*APA102, error) {
	// The speed must be high, as there's 32 bits sent per LED, creating a
	// staggered effect. See
	// https://cpldcpu.wordpress.com/2014/11/30/understanding-the-apa102-superled/
	w, err := rpi.MakeSPI("", speed)
	if err != nil {
		return nil, err
	}
	return &APA102{
		Intensity:   255,
		Temperature: 7000,
		w:           w,
	}, err
}
