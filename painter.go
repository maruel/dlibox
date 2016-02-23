// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"image/color"
	"io"
	"sync"
	"time"

	"github.com/maruel/interrupt"
)

// Pattern is a interface to draw an animated line.
type Pattern interface {
	// NextFrame fills the buffer with the next image.
	//
	// The image should be derived from the duration since this pattern was
	// started.
	//
	// First call is guaranteed to be called with sinceStart == 0.
	NextFrame(pixels []color.NRGBA, sinceStart time.Duration)

	// TODO(maruel): Add generic function to increase/decrease the speed.
}

// Strip is an 1D output device.
type Strip interface {
	io.Closer
	// Write writes a new frame.
	Write(pixels []color.NRGBA) error
	// MinDelay returns the minimum delay between each draw refresh.
	MinDelay() time.Duration
}

type Painter struct {
	s  Strip
	c  chan Pattern
	wg sync.WaitGroup
}

func (p *Painter) SetPattern(pat Pattern) {
	p.c <- pat
}

func (p *Painter) Close() error {
	p.c <- nil
	p.wg.Wait()
	return p.s.Close()
}

// MakePainter returns a Painter that manages updating the Patterns to the
// Strip.
func MakePainter(s Strip, numLights int) *Painter {
	p := &Painter{s: s, c: make(chan Pattern)}
	p.wg.Add(1)
	go p.runPattern(numLights)
	return p
}

// Private stuff.

// d60Hz is the duration of one frame at 60Hz.
const d60Hz = 16666666 * time.Nanosecond

func (p *Painter) runPattern(numLights int) {
	defer p.wg.Done()
	var pat Pattern = &StaticColor{color.NRGBA{}}
	pixels := make([]color.NRGBA, numLights)
	start := time.Now()
	var since time.Duration
	delay := p.s.MinDelay()
	if delay < d60Hz {
		delay = d60Hz
	}
	timer := time.After(delay)
	for {
		// TODO(maruel): If ever become CPU bound or SPI I/O bound, call
		// NextFrame() in one goroutine and Write() in another one. This is
		// especially useful in multicore systems like rPi2.
		pat.NextFrame(pixels, since)
		if err := p.s.Write(pixels); err != nil {
			return
		}
		select {
		case pat = <-p.c:
			if pat == nil {
				// Request to terminate.
				return
			}
			// New pattern.
			// TODO(maruel): Use a fade-in, fade-out of 500ms.
			start = time.Now()
			since = 0
			timer = time.After(delay)
		case <-timer:
			since = time.Since(start)
			timer = time.After(delay)
		case <-interrupt.Channel:
			return
		}
	}
}

/*
type subset struct {
	s     Strip
	start int
	end   int
}

func (s*subset) Write(pixels[]color.NRGBA) error {
	// Argh.
	return s.s.Write(pixels[:end])
}

// Subset returns a new Strip that only affects the subset of the original
// Strip.
func Subset(s Strip, start, end int) Strip {
	return &subset{s, start, end}
}
*/
