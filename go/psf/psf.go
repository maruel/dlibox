// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Package psf loads PSF bitmap fonts as installed via consolefonts on Debian
// derived linux distributions.
//
// They can be used to draw on a image.Image.
package psf

import (
	"compress/gzip"
	"encoding/binary"
	"errors"
	"image/color"
	"image/draw"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

const root = "/usr/share/consolefonts/"

// Choices of family and sizes on a Jessie system:
// Fixed13
// Fixed14
// Fixed15
// Fixed16
// Fixed18
// Terminus12x6
// Terminus14
// Terminus16
// Terminus20x10
// Terminus22x11
// Terminus24x12
// Terminus28x14
// Terminus32x16
// TerminusBold14
// TerminusBold16
// TerminusBold20x10
// TerminusBold22x11
// TerminusBold24x12
// TerminusBold28x14
// TerminusBold32x16
// TerminusBoldVGA14
// TerminusBoldVGA16
// VGA8
// VGA14
// VGA16
// VGA28x16
// VGA32x16

// Enumerate returns the font families accessible.
func Enumerate() ([]string, error) {
	m, err := filepath.Glob(root + "Uni2-*")
	if err != nil {
		return nil, err
	}
	out := make([]string, 0, len(m))
	for _, i := range m {
		parts := strings.SplitN(i, "-", 2)
		if len(parts) == 2 {
			out = append(out, strings.SplitN(parts[1], ".", 2)[0])
		}
	}
	return out, nil
}

// Load loads a font from Debian provided consolefonts.
//
// Call Enumerate to enumerate the font families available.
func Load(name string) (*Font, error) {
	// It opens all the psf for the font family and merges all of them.
	f := &Font{Letters: map[rune][]byte{}}
	m, err := filepath.Glob(root + "*-" + name + ".psf.gz")
	if err != nil {
		return nil, err
	}
	for _, n := range m {
		if err := f.loadPSF(n); err != nil {
			return nil, err
		}
	}
	return f, nil
}

// Draw draws text on an image.
//
// Use nil as the color to use a transparent foreground or background.
//
// Returns the end point.
func (f *Font) Draw(dst draw.Image, x, y int, fore, back color.Color, text string) (int, int) {
	for _, r := range text {
		bitmap := f.Letters[r]
		for yL := 0; yL < f.H; yL++ {
			for xL := 0; xL < f.W; xL++ {
				v := bitmap[xL/8+yL*f.W/8]
				b := byte(1 << (7 - byte(xL&7)))
				if v&b != 0 {
					if fore != nil {
						dst.Set(x+xL, y+yL, fore)
					}
				} else {
					if back != nil {
						dst.Set(x+xL, y+yL, back)
					}
				}
			}
		}
		x += f.W
	}
	return x, y + f.H
}

// Font is a rasterized font to be used on low resolution display.
type Font struct {
	Version int
	H       int
	W       int
	Letters map[rune][]byte
}

// loadPSf loads a font according to the specification at
// https://www.win.tue.nl/~aeb/linux/kbd/font-formats-1.html
//
// It tries to be as strict as possible.
func (f *Font) loadPSF(path string) error {
	b, err := readCompressed(path)
	if err != nil {
		return err
	}
	version, nbGlyphs, mapping, height, width, charSize, data, err := parseHeader(b)

	// Update Font accordingly.
	if f.Version == 0 {
		f.Version = version
	}
	if f.Version != version {
		return errors.New("mixed versions")
	}
	if f.H == 0 {
		f.H = height
	}
	if f.H != height {
		return errors.New("unexpected different font size")
	}
	if f.W == 0 {
		f.W = width
	}
	if f.W != width {
		return errors.New("unexpected different font size")
	}

	// Grab the glyph bitmaps.
	for c := 0; c < nbGlyphs; c++ {
		n := make([]byte, charSize)
		copy(n, data)
		for _, i := range mapping[c] {
			f.Letters[i] = n
		}
		data = data[charSize:]
	}
	if len(data) != 0 {
		return errors.New("unconsumed data bytes")
	}
	return nil
}

func readCompressed(path string) ([]byte, error) {
	// The psf files are all small, so it's faster to just read the whole file in memory.
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
	return ioutil.ReadAll(g)
}

func parseHeader(b []byte) (int, int, map[int][]rune, int, int, int, []byte, error) {
	if len(b) < 16 {
		return 0, 0, nil, 0, 0, 0, nil, errors.New("invalid file: too short")
	}
	mapping := map[int][]rune{}
	if b[0] == 0x36 && b[1] == 0x04 {
		mode := b[2]
		nbGlyphs := 256
		if mode&1 != 0 {
			nbGlyphs = 512
		}
		charSize := int(b[3])

		if mode&2 != 0 {
			// Create unicode mapping as a stream uint16.
			uni := b[charSize*nbGlyphs:]
			for i := 0; len(uni) != 0; i++ {
				for {
					uc := binary.LittleEndian.Uint16(uni)
					uni = uni[2:]
					if uc == 0xFFFF {
						break
					}
					if uc == 0xFFFE {
						for binary.LittleEndian.Uint16(uni) != 0xFFFF {
							uni = uni[2:]
						}
						break
					}
					// psf1 only supports UCS2
					mapping[i] = append(mapping[i], rune(uc))
				}
			}
		}
		// For psf1, the width is hard coded at 8 pixels.
		return 1, nbGlyphs, mapping, charSize, 8, charSize, b[4 : 4+charSize*nbGlyphs], nil
	}

	if b[0] == 0x72 && b[1] == 0xB5 && b[2] == 0x4A && b[3] == 0x86 && binary.LittleEndian.Uint32(b[4:]) == 0 {
		hdrSize := int(binary.LittleEndian.Uint32(b[8:]))
		flags := binary.LittleEndian.Uint32(b[0xC:])
		nbGlyphs := int(binary.LittleEndian.Uint32(b[0x10:]))
		charSize := int(binary.LittleEndian.Uint32(b[0x14:]))
		height := int(binary.LittleEndian.Uint32(b[0x18:]))
		width := int(binary.LittleEndian.Uint32(b[0x1C:]))
		if flags&1 == 0 {
			return 0, 0, nil, 0, 0, 0, nil, errors.New("invalid file: no unicode lookup table")
		}
		off := hdrSize + charSize*nbGlyphs
		uni := b[off:]
		for i := 0; len(uni) != 0; i++ {
			for {
				if uni[0] == 0xFF {
					uni = uni[1:]
					break
				}
				if uni[0] == 0xFE {
					uni = uni[1:]
					continue
				}
				r, s := utf8.DecodeRune(uni)
				if s == 0 {
					return 0, 0, nil, 0, 0, 0, nil, errors.New("invalid file: bad unicode table")
				}
				uni = uni[s:]
				mapping[i] = append(mapping[i], r)
			}
		}
		return 2, nbGlyphs, mapping, height, width, charSize, b[hdrSize:off], nil
	}
	return 0, 0, nil, 0, 0, 0, nil, errors.New("invalid file: invalid magic byte")
}
