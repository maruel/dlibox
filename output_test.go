// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package dotstar

import (
	"bytes"
	"image/color"
	"io"
	"testing"

	"github.com/maruel/ut"
)

func TestColorToAPA102(t *testing.T) {
	type col struct {
		b, B, G, R byte
	}
	data := []struct {
		c        color.NRGBA
		expected col
	}{
		{color.NRGBA{0xFF, 0xFF, 0xFF, 0xFF}, col{0xFF, 0xFF, 0xFF, 0xFF}},
		{color.NRGBA{0x00, 0x00, 0x00, 0xFF}, col{0xE1, 0x00, 0x00, 0x00}},
		{color.NRGBA{0xFF, 0xFF, 0xFF, 0x00}, col{0xE1, 0x00, 0x00, 0x00}},
		{color.NRGBA{0xFF, 0x00, 0x00, 0xFF}, col{0xFF, 0x00, 0x00, 0xFF}},
		{color.NRGBA{0x00, 0xFF, 0x00, 0xFF}, col{0xFF, 0x00, 0xFF, 0x00}},
		{color.NRGBA{0x00, 0x00, 0xFF, 0xFF}, col{0xFF, 0xFF, 0x00, 0x00}},
	}
	for i, line := range data {
		var actual col
		actual.b, actual.B, actual.G, actual.R = colorToAPA102(line.c)
		ut.AssertEqualIndex(t, i, line.expected, actual)
	}
}

func TestDotStarEmpty(t *testing.T) {
	b := &bytes.Buffer{}
	d := &DotStar{
		RedGamma:   1.,
		RedMax:     1.,
		GreenGamma: 1.,
		GreenMax:   1.,
		BlueGamma:  1.,
		BlueMax:    1.,
		w:          nopCloser{b},
	}
	ut.AssertEqual(t, nil, d.Write([]color.NRGBA{}))
	ut.AssertEqual(t, []byte{0x0, 0x0, 0x0, 0x0, 0xff}, b.Bytes())
}

func TestDotStar(t *testing.T) {
	b := &bytes.Buffer{}
	d := &DotStar{
		RedGamma:   1.,
		RedMax:     1.,
		GreenGamma: 1.,
		GreenMax:   1.,
		BlueGamma:  1.,
		BlueMax:    1.,
		w:          nopCloser{b},
	}
	colors := []color.NRGBA{
		{0xFE, 0xFE, 0xFE, 0xFF},
		{0xFE, 0xFE, 0xFE, 0x00},
		{0xF0, 0xF0, 0xF0, 0xFF},
		{0x80, 0x80, 0x80, 0xFF},
		{0x80, 0x00, 0x00, 0xFF},
		{0x00, 0x80, 0x00, 0xFF},
		{0x00, 0x00, 0x80, 0xFF},
	}
	ut.AssertEqual(t, nil, d.Write(colors))
	/*
		// TODO(maruel): Resolution loss.
		expected := []byte{
			0x00, 0x00, 0x00, 0x00,
			0xFF, 0xFF, 0xFF, 0xFF,
			0xFF, 0x00, 0x00, 0x00,
			0xFF, 0xF1, 0xF1, 0xF1,
			0xFF, 0x81, 0x81, 0x81,
			0xFF, 0x00, 0x00, 0x81,
			0xFF, 0x00, 0x81, 0x00,
			0xFF, 0x81, 0x00, 0x00,
			0xFF,
		}
		ut.AssertEqual(t, expected, b.Bytes())
	*/
}

func TestDotStarPowerLimited(t *testing.T) {
	b := &bytes.Buffer{}
	d := &DotStar{
		RedGamma:   1.,
		RedMax:     1.,
		GreenGamma: 1.,
		GreenMax:   1.,
		BlueGamma:  1.,
		BlueMax:    1.,
		AmpPerLED:  .02,
		AmpBudget:  0.1,
		w:          nopCloser{b},
	}
	colors := []color.NRGBA{
		{0xFE, 0xFE, 0xFE, 0xFF},
		{0xFE, 0xFE, 0xFE, 0x00},
		{0xF0, 0xF0, 0xF0, 0xFF},
		{0x80, 0x80, 0x80, 0xFF},
		{0x80, 0x00, 0x00, 0xFF},
		{0x00, 0x80, 0x00, 0xFF},
		{0x00, 0x00, 0x80, 0xFF},
	}
	ut.AssertEqual(t, nil, d.Write(colors))
	/*
		expected := []byte{
			0x00, 0x00, 0x00, 0x00,
			0xFF, 0x90, 0x90, 0x90,
			0xFF, 0x00, 0x00, 0x00,
			0xFF, 0x88, 0x88, 0x88,
			0xFF, 0x49, 0x49, 0x49,
			0xFF, 0x00, 0x00, 0x49,
			0xFF, 0x00, 0x49, 0x00,
			0xFF, 0x49, 0x00, 0x00,
			0xFF,
		}
		ut.AssertEqual(t, expected, b.Bytes())
	*/
}

func TestDotStarLong(t *testing.T) {
	b := &bytes.Buffer{}
	d := &DotStar{
		RedGamma:   1.,
		RedMax:     1.,
		GreenGamma: 1.,
		GreenMax:   1.,
		BlueGamma:  1.,
		BlueMax:    1.,
		w:          nopCloser{b},
	}
	colors := make([]color.NRGBA, 256)
	ut.AssertEqual(t, nil, d.Write(colors))
	expected := make([]byte, 4*(256+1)+17)
	for i := 0; i < 256; i++ {
		expected[4+4*i] = 0xFF
	}
	trailer := expected[4*257:]
	for i := range trailer {
		trailer[i] = 0xFF
	}
	//ut.AssertEqual(t, expected, b.Bytes())
}

//

type nopCloser struct {
	io.Writer
}

func (nopCloser) Close() error {
	return nil
}
