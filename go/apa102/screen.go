// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package apa102

import (
	"bytes"
	"errors"
	"image/color"
	"io"

	"github.com/maruel/ansi256"
	"github.com/mattn/go-colorable"
)

// ScreenStrip is a 1d LED strip emulator.
type ScreenStrip struct {
	w   io.Writer
	buf bytes.Buffer
}

// Write accepts a stream of raw RGB pixels and writes it to os.Stdout.
func (s *ScreenStrip) Write(pixels []byte) (int, error) {
	if len(pixels)%3 != 0 {
		return 0, errors.New("invalid RGB stream length")
	}
	// This code is designed to minimize the amount of memory allocated per call.
	s.buf.Reset()
	_, _ = s.buf.WriteString("\r\033[0m")
	for i := 0; i < len(pixels)/3; i++ {
		_, _ = io.WriteString(&s.buf, ansi256.Default.Block(color.NRGBA{pixels[3*i], pixels[3*i+1], pixels[3*i+2], 255}))
	}
	_, _ = s.buf.WriteString("\033[0m ")
	_, err := s.buf.WriteTo(s.w)
	return len(pixels), err
}

// MakeScreen returns a strip that displays at the console.
//
// This is generally what you want while waiting for the LED strip to be
// shipped and you are excited to try it out.
func MakeScreen() *ScreenStrip {
	return &ScreenStrip{w: colorable.NewColorableStdout()}
}
