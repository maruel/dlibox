// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Package bw2d implements black and white (1 bit per pixel) 2D graphics.
//
// It is compatible with package image/draw.
package bw2d

import (
	"image"
	"image/color"
	"image/draw"
)

// Whenever a bit is set or not.
const On = bit(true)
const Off = bit(false)

// Image is a 1bit image.
type Image struct {
	W   int
	H   int
	Buf []byte
}

func Make(w, h int) *Image {
	return &Image{w, h, make([]byte, w*h/8)}
}

func (i *Image) SetAll() {
	for j := range i.Buf {
		i.Buf[j] = 0xFF
	}
}

func (i *Image) Clear() {
	for j := range i.Buf {
		i.Buf[j] = 0
	}
}

func (i *Image) Inverse() {
	for j := range i.Buf {
		i.Buf[j] ^= 0xFF
	}
}

// ColorModel implements image.Image.
func (i *Image) ColorModel() color.Model {
	return color.ModelFunc(convert)
}

// Bounds implements image.Image.
func (i *Image) Bounds() image.Rectangle {
	return image.Rectangle{Max: image.Point{X: i.W, Y: i.H}}
}

// At implements image.Image.
func (i *Image) At(x, y int) color.Color {
	// Addressing is a bit odd, each byte is 8 vertical bits.
	o := x + y/8*i.W
	b := byte(1 << byte(y&7))
	return bit(i.Buf[o]&b != 0)
}

// Set implements draw.Image
func (i *Image) Set(x, y int, c color.Color) {
	if x >= i.W {
		panic("out of bound")
	}
	if y >= i.H {
		panic("out of bound")
	}
	o := x + y/8*i.W
	b := byte(1 << byte(y&7))
	if convertBit(c) {
		i.Buf[o] |= b
	} else {
		i.Buf[o] &^= b
	}
}

// Private stuff.

var _ draw.Image = &Image{}

// Anything not transparent and not pure black is white.
func convert(c color.Color) color.Color {
	return convertBit(c)
}

// Anything not transparent and not pure black is white.
func convertBit(c color.Color) bit {
	switch t := c.(type) {
	case bit:
		return t
	default:
		// Values are on 16 bits.
		r, g, b, a := c.RGBA()
		return bit((r+g+b) > 0x10000 && a >= 0x4000)
	}
}

type bit bool

func (b bit) RGBA() (uint32, uint32, uint32, uint32) {
	if b {
		return 255, 255, 255, 255
	}
	return 0, 0, 0, 0
}
