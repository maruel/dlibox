// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"image/color"
	"io"
	"log"
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
}

// Strip is an 1D output device.
type Strip interface {
	io.Closer
	// Write writes a new frame.
	Write(pixels []color.NRGBA) error
	// MinDelay returns the minimum delay between each draw refresh.
	MinDelay() time.Duration
}

// Painter handles the "draw frame, write" loop.
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
	// Tripple buffering.
	cGen := make(chan []color.NRGBA, 3)
	cWrite := make(chan []color.NRGBA, cap(cGen))
	for i := 0; i < cap(cGen); i++ {
		cGen <- make([]color.NRGBA, numLights)
	}
	p.wg.Add(2)
	go p.runPattern(cGen, cWrite)
	go p.runWrite(cGen, cWrite)
	return p
}

// Private stuff.

// d60Hz is the duration of one frame at 60Hz.
const d60Hz = 16666666 * time.Nanosecond

func getDelay(s Strip) time.Duration {
	delay := s.MinDelay()
	if delay < d60Hz {
		delay = d60Hz
	}
	return delay
}

func (p *Painter) runPattern(cGen, cWrite chan []color.NRGBA) {
	defer p.wg.Done()
	defer func() {
		cWrite <- nil
	}()
	ease := EaseOut{
		In:       &StaticColor{color.NRGBA{}},
		Out:      &StaticColor{color.NRGBA{}},
		Duration: 500 * time.Millisecond,
	}
	var since time.Duration
	delay := getDelay(p.s)
	for {
		select {
		case newPat := <-p.c:
			if newPat == nil {
				// Request to terminate.
				return
			}

			// New pattern.
			ease.Out = ease.In
			ease.In = newPat
			ease.Offset = since
			since = 0

		case pixels := <-cGen:
			for i := range pixels {
				pixels[i] = color.NRGBA{}
			}
			ease.NextFrame(pixels, since)
			since += delay
			cWrite <- pixels

		case <-interrupt.Channel:
			return
		}
	}
}

func (p *Painter) runWrite(cGen, cWrite chan []color.NRGBA) {
	defer p.wg.Done()
	delay := getDelay(p.s)
	tick := time.NewTicker(delay)
	defer tick.Stop()
	var err error
	for {
		pixels := <-cWrite
		if len(pixels) == 0 {
			return
		}
		if err == nil {
			if err = p.s.Write(pixels); err != nil {
				log.Printf("Writing failed: %s", err)
			}
		}
		cGen <- pixels

		select {
		case <-tick.C:
		case <-interrupt.Channel:
			return
		}
	}
}
