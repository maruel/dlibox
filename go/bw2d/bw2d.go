// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Package bw2d implements black and white (1 bit per pixel) 2D graphics.
//
// It is compatible with package image/draw.
package bw2d

import (
	"compress/gzip"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"os"
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

// Text draw text in the image.
//
// Returns the end point.
func (i *Image) Text(x, y int, f *Font, text string) (int, int) {
	// Slow path.
	// TODO(maruel): Handle overflow by returning an error.
	// TODO(maruel): Operations like rotation.
	for _, r := range text {
		c := f.Letters[r]
		for yL := 0; yL < f.H; yL++ {
			for xL := 0; xL < f.W; xL++ {
				v := c[xL/8+yL*f.W/8]
				b := byte(1 << (7 - byte(xL&7)))
				i.Set(x+xL, y+yL, bit(v&b != 0))
			}
		}
		x += f.W
	}
	return x, y + f.H
}

type Font struct {
	H       int
	W       int
	Letters [][]byte
}

// BasicFont is a Debian provided simple font.
//
// Valid values are 8, 14 and 16
func BasicFont(h int) (*Font, error) {
	// http://unix.stackexchange.com/questions/216184/what-is-the-difference-between-uni1-uni2-and-uni3-terminal-font-codesets
	// /usr/share/consolefonts/*-VGA8.psf.gz

	// TODO(maruel): Load multiple code pages.
	return loadPSF(fmt.Sprintf("/usr/share/consolefonts/Uni1-VGA%d.psf.gz", h))
}

func loadPSF(path string) (*Font, error) {
	fd, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	g, err := gzip.NewReader(fd)
	if err != nil {
		return nil, err
	}
	defer g.Close()
	b, err := ioutil.ReadAll(g)
	if err != nil {
		return nil, err
	}
	f := &Font{}
	if b[0] == 0x36 {
		if b[1] != 0x04 {
			return nil, errors.New("invalid file: invalid header")
		}
		switch b[2] {
		case 0, 2:
			f.Letters = make([][]byte, 256)
		case 1, 3:
			f.Letters = make([][]byte, 512)
		default:
			return nil, errors.New("invalid file: number of glyphs")
		}
		f.H = int(b[3])
		f.W = 8
		b = b[4:]
	} else if b[0] == 0x72 {
		return nil, errors.New("not implemented: mode 0x72")
	} else {
		return nil, errors.New("invalid file: invalid magic byte")
	}

	// Keep references to the original memory block. This wastes the header but
	// it's still better than fragmenting the heap.
	l := f.H * f.W / 8
	for c := range f.Letters {
		f.Letters[c] = b[:l]
		b = b[l:]
	}
	return f, nil
}

/*
func LoadTTF(b []byte) (*Font, error) {
	//f := basicfont.Face7x13
	//_, _, _, _, _ = f.Glyph(fixed.Point26_6{1, 1}, 'a')
	//
	// http://www.lowing.org/fonts/
	// http://www.alvit.de/blog/article/25-best-license-free-pixelfonts
	// http://www.dafont.com/bitmap.php
	//
	// http://www.kottke.org/plus/type/silkscreen/
	// http://www.kottke.org/plus/type/silkscreen/download/silkscreen.zip
	// License: custom "free for personal and corporate use"
	//
	// https://www.gnome.org/fonts/
	// http://ftp.gnome.org/pub/GNOME/sources/ttf-bitstream-vera/1.10/ttf-bitstream-vera-1.10.tar.bz2
	// License: custom; "free"
	//
	// Only installed if X is installed:
	// /usr/share/fonts/liberation/LiberationSerif-Regular.ttf

	// https://developer.apple.com/fonts/TrueType-Reference-Manual/
	f, err := truetype.Parse(b)
	//t, err := freetype.ParseFont(b)
	if err != nil {
		return nil, err
	}
	log.Printf("%s", f.Name(truetype.NameIDCopyright))
	g := truetype.GlyphBuf{}
	a := f.Index('a')
	if err := g.Load(f, fixed.Int26_6(1<<6), a, font.HintingNone); err != nil {
		return nil, err
	}
	log.Printf("%# v", pretty.Formatter(g))
	return &Font{}, nil
}
*/

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
