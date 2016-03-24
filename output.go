// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package dotstar

import (
	"bytes"
	"image/color"
	"io"
	"time"

	"github.com/maruel/ansi256"
	"github.com/mattn/go-colorable"
)

type dotStar struct {
	w          io.WriteCloser
	b          []byte
	brightness int // 0-31
}

func (d *dotStar) Close() error {
	return d.w.Close()
}

func (d *dotStar) Write(pixels []color.NRGBA) error {
	// https://cpldcpu.files.wordpress.com/2014/08/apa-102c-super-led-specifications-2014-en.pdf
	numLights := len(pixels)
	// End frames are needed to be able to push enough SPI clock signals due to
	// internal half-delay of data signal from each individual LED. See
	// https://cpldcpu.wordpress.com/2014/11/30/understanding-the-apa102-superled/
	l := 4*(numLights+1) + numLights/2/8 + 1
	if len(d.b) < l {
		d.b = make([]byte, l)
	}
	// Start frame is all zeros. Just skip it.
	s := d.b[4:]
	brightness := byte(0xE0 + d.brightness)
	for i := range pixels {
		r, g, b, _ := pixels[i].RGBA()
		// BGR.
		s[4*i] = brightness
		s[4*i+1] = byte(b >> 8)
		s[4*i+2] = byte(g >> 8)
		s[4*i+3] = byte(r >> 8)
	}
	// End frames
	s = s[4*numLights:]
	for i := range s {
		s[i] = 0xFF
	}
	_, err := d.w.Write(d.b)
	return err
}

func (d *dotStar) MinDelay() time.Duration {
	// As per APA102-C spec, it's max refresh rate is 400hz.
	// https://en.wikipedia.org/wiki/Flicker_fusion_threshold is a recommended
	// reading.
	return time.Second / 400
}

// MakeDotStar returns a stripe that communicates over SPI.
func MakeDotStar() (Strip, error) {
	// The speed must be high, as there's 32 bits sent per LED, creating a
	// staggered effect. See
	// https://cpldcpu.wordpress.com/2014/11/30/understanding-the-apa102-superled/
	w, err := makeSPI("", 20000000)
	if err != nil {
		return nil, err
	}
	return &dotStar{w: w, brightness: 31}, err
}

//

type screenStrip struct {
	w io.Writer
	b bytes.Buffer
}

func (s *screenStrip) Close() error {
	return nil
}

func (s *screenStrip) Write(pixels []color.NRGBA) error {
	// This code is designed to minimize the amount of memory allocated per call.
	s.b.Reset()
	_, _ = s.b.WriteString("\r\033[0m")
	for _, c := range pixels {
		_, _ = io.WriteString(&s.b, ansi256.Default.Block(c))
	}
	_, _ = s.b.WriteString("\033[0m ")
	_, err := s.b.WriteTo(s.w)
	return err
}

func (s *screenStrip) MinDelay() time.Duration {
	// Limit to 30hz, especially for ssh connections.
	return time.Second / 30
}

// MakeScreen returns a stripe that display at the console.
func MakeScreen() Strip {
	return &screenStrip{w: colorable.NewColorableStdout()}
}
