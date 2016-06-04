// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package apa102

import (
	"bytes"
	"image/color"
	"io"
	"time"

	"github.com/maruel/ansi256"
	"github.com/maruel/dlibox-go/anim1d"
	"github.com/mattn/go-colorable"
)

type screenStrip struct {
	w   io.Writer
	buf bytes.Buffer
}

func (s *screenStrip) Close() error {
	return nil
}

func (s *screenStrip) Write(pixels anim1d.Frame) error {
	// This code is designed to minimize the amount of memory allocated per call.
	s.buf.Reset()
	_, _ = s.buf.WriteString("\r\033[0m")
	for _, c := range pixels {
		_, _ = io.WriteString(&s.buf, ansi256.Default.Block(color.NRGBA{c.R, c.G, c.B, 255}))
	}
	_, _ = s.buf.WriteString("\033[0m ")
	_, err := s.buf.WriteTo(s.w)
	return err
}

func (s *screenStrip) MinDelay() time.Duration {
	// Limit to 30hz, especially for ssh connections.
	return time.Second / 30
}

// MakeScreen returns a strip that displays at the console.
//
// This is generally what you want while waiting for the LED strip to be
// shipped and you are excited to try it out.
func MakeScreen() anim1d.Strip {
	return &screenStrip{w: colorable.NewColorableStdout()}
}
