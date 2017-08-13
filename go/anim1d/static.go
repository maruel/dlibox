// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// static is for patterns that do not change over time.

package anim1d

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
)

// Color shows a single color on all lights. It knows how to renders itself
// into a frame.
//
// If you want a single dot, use a Frame of length one.
type Color struct {
	R, G, B uint8
}

func (c *Color) Render(pixels Frame, timeMS uint32) {
	for i := range pixels {
		pixels[i] = *c
	}
}

// Dim reduces the intensity of a color/pixel to scale it on intensity.
//
// 0 means completely dark, 255 the color c is unaffected.
func (c *Color) Dim(intensity uint8) {
	i := uint16(intensity)
	d := i >> 1
	c.R = uint8((uint16(c.R)*i + d) >> 8)
	c.G = uint8((uint16(c.G)*i + d) >> 8)
	c.B = uint8((uint16(c.B)*i + d) >> 8)
}

// Add adds two color together with saturation.
func (c *Color) Add(d Color) {
	r := uint16(c.R) + uint16(d.R)
	if r > 255 {
		r = 255
	}
	c.R = uint8(r)
	g := uint16(c.G) + uint16(d.G)
	if g > 255 {
		g = 255
	}
	c.G = uint8(g)
	b := uint16(c.B) + uint16(d.B)
	if b > 255 {
		b = 255
	}
	c.B = uint8(b)
}

// Mix blends the second color with the first.
//
// gradient 0 means pure 'c', gradient 255 means pure 'd'.
//
// It is the equivalent of:
//   c.Dim(255-gradient)
//   d.Dim(gradient)
//   c.Add(d)
// except that this function doesn't affect d.
func (c *Color) Mix(d Color, gradient uint8) {
	grad := uint16(gradient)
	grad1 := 255 - grad
	// unit test confirms the values cannot overflow.
	c.R = uint8(((uint16(c.R)+1)*grad1 + (uint16(d.R)+1)*grad) >> 8)
	c.G = uint8(((uint16(c.G)+1)*grad1 + (uint16(d.G)+1)*grad) >> 8)
	c.B = uint8(((uint16(c.B)+1)*grad1 + (uint16(d.B)+1)*grad) >> 8)
}

func (c *Color) String() string {
	return fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B)
}

// FromString converts a "#RRGGBB" encoded string to a Color.
//
// 'c' is untouched in case of error.
func (c *Color) FromString(s string) error {
	if len(s) != 7 || s[0] != '#' {
		return errors.New("invalid color string")
	}
	return c.FromRGBString(s[1:])
}

// FromRGBString converts a "RRGGBB" encoded string to a Color.
//
// 'c' is untouched in case of error.
func (c *Color) FromRGBString(s string) error {
	if len(s) != 6 {
		return errors.New("invalid color string")
	}
	r, err := strconv.ParseUint(s[0:2], 16, 8)
	if err != nil {
		return err
	}
	g, err := strconv.ParseUint(s[2:4], 16, 8)
	if err != nil {
		return err
	}
	b, err := strconv.ParseUint(s[4:6], 16, 8)
	if err != nil {
		return err
	}
	c.R = uint8(r)
	c.G = uint8(g)
	c.B = uint8(b)
	return nil
}

//

// Frame is a strip of colors. It knows how to renders itself into a frame
// (which is recursive).
type Frame []Color

func (f Frame) Render(pixels Frame, timeMS uint32) {
	copy(pixels, f)
}

// Dim reduces the intensity of a frame to scale it on intensity.
func (f Frame) Dim(intensity uint8) {
	for i := range f {
		f[i].Dim(intensity)
	}
}

// Add adds two frames together with saturation.
func (f Frame) Add(r Frame) {
	for i := range f {
		f[i].Add(r[i])
	}
}

// Mix blends the second frame with the first.
//
// gradient 0 means pure 'f', gradient 255 means pure 'b'.
//
// It is the equivalent of:
//   c.Dim(255-gradient)
//   d.Dim(gradient)
//   c.Add(d)
// except that this function doesn't affect d.
func (f Frame) Mix(b Frame, gradient uint8) {
	for i := range f {
		f[i].Mix(b[i], gradient)
	}
}

// ToRGB converts the Frame to a raw RGB stream.
func (f Frame) ToRGB(b []byte) {
	for i := range f {
		b[3*i] = f[i].R
		b[3*i+1] = f[i].G
		b[3*i+2] = f[i].B
	}
}

// reset() always resets the buffer to black.
func (f *Frame) reset(l int) {
	if len(*f) != l {
		*f = make(Frame, l)
	} else {
		s := *f
		for i := range s {
			s[i] = Color{}
		}
	}
}

func (f Frame) isEqual(rhs Frame) bool {
	if len(f) != len(rhs) {
		return false
	}
	for j, p := range f {
		if p != rhs[j] {
			return false
		}
	}
	return true
}

// FromString converts a "LRRGGBB..." encoded string to a Frame.
//
// 'f' is untouched in case of error.
func (f *Frame) FromString(s string) error {
	if len(s) == 0 || (len(s)-1)%6 != 0 || s[0] != 'L' {
		return errors.New("invalid frame string")
	}
	l := (len(s) - 1) / 6
	f2 := make(Frame, l)
	for i := 0; i < l; i++ {
		if err := f2[i].FromRGBString(s[1+i*6 : 1+(i+1)*6]); err != nil {
			return err
		}
	}
	*f = f2
	return nil
}

func (f Frame) String() string {
	out := bytes.Buffer{}
	out.Grow(1 + 6*len(f))
	out.WriteByte('L')
	for _, c := range f {
		fmt.Fprintf(&out, "%02x%02x%02x", c.R, c.G, c.B)
	}
	return out.String()
}

//

// Rainbow renders rainbow colors.
type Rainbow struct {
	buf Frame
}

func (r *Rainbow) Render(pixels Frame, timeMS uint32) {
	if len(r.buf) != len(pixels) {
		r.buf.reset(len(pixels))
		const start = 380
		const end = 781
		const delta = end - start
		// TODO(maruel): Use integer arithmetic.
		scale := logn(2)
		step := 1. / float32(len(pixels))
		for i := range pixels {
			j := log1p(float32(len(pixels)-i-1)*step) / scale
			r.buf[i] = waveLength2RGB(int(start + delta*(1-j)))
		}
	}
	copy(pixels, r.buf)
}

func (r *Rainbow) String() string {
	return rainbowKey
}

// waveLengthToRGB returns a color over a rainbow.
//
// This code was inspired by public domain code on the internet.
func waveLength2RGB(w int) (c Color) {
	switch {
	case w < 380:
	case w < 420:
		// Red peaks at 1/3 at 420.
		c.R = uint8(196 - (170*(440-w))/(440-380))
		c.B = uint8(26 + (229*(w-380))/(420-380))
	case w < 440:
		c.R = uint8((0x89 * (440 - w)) / (440 - 420))
		c.B = 255
	case w < 490:
		c.G = uint8((255 * (w - 440)) / (490 - 440))
		c.B = 255
	case w < 510:
		c.G = 255
		c.B = uint8((255 * (510 - w)) / (510 - 490))
	case w < 580:
		c.R = uint8((255 * (w - 510)) / (580 - 510))
		c.G = 255
	case w < 645:
		c.R = 255
		c.G = uint8((255 * (645 - w)) / (645 - 580))
	case w < 700:
		c.R = 255
	case w < 781:
		c.R = uint8(26 + (229*(780-w))/(780-700))
	default:
	}
	return
}

// Repeated repeats a Frame to fill the pixels.
type Repeated struct {
	Frame Frame
}

func (r *Repeated) Render(pixels Frame, timeMS uint32) {
	if len(pixels) == 0 || len(r.Frame) == 0 {
		return
	}
	for i := 0; i < len(pixels); i += len(r.Frame) {
		copy(pixels[i:], r.Frame)
	}
}
