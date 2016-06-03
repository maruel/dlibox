// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package anim1d

import (
	"encoding/json"
	"io"
	"log"
	"sync"
	"time"

	"github.com/maruel/dotstar/rpi"
	"github.com/maruel/interrupt"
)

// Pattern is a interface to draw an animated line.
type Pattern interface {
	// NextFrame fills the buffer with the next image.
	//
	// The image should be derived from the duration since this pattern was
	// started.
	//
	// Calling NextFrame() with a nil pattern is valid. Patterns should be
	// callable without crashing with an object initialized with default values.
	//
	// First call is guaranteed to be called with sinceStart == 0.
	NextFrame(pixels Frame, sinceStart time.Duration)

	// TODO(maruel): Will have to think about it.
	// NativeDuration returns the looping duration, if any. It is used for
	// animated GIF generation.
	//NativeDuration(pixels int) time.Duration
}

// Strip is an 1D output device.
type Strip interface {
	io.Closer
	// Write writes a new frame.
	Write(pixels Frame) error
	// MinDelay returns the minimum delay between each draw refresh.
	MinDelay() time.Duration
}

// Painter handles the "draw frame, write" loop.
type Painter struct {
	s  Strip
	c  chan Pattern
	wg sync.WaitGroup
}

// SetPattern changes the current pattern to a new one.
//
// The pattern is in JSON encoded format. The function will return an error if
// the encoding is bad. The function is synchronous, it returns only after the
// pattern was effectively set.
func (p *Painter) SetPattern(s string) error {
	var pat SPattern
	if err := json.Unmarshal([]byte(s), &pat); err != nil {
		return err
	}
	p.c <- pat
	return nil
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
	cGen := make(chan Frame, 3)
	cWrite := make(chan Frame, cap(cGen))
	for i := 0; i < cap(cGen); i++ {
		cGen <- make(Frame, numLights)
	}
	p.wg.Add(2)
	go p.runPattern(cGen, cWrite)
	go p.runWrite(cGen, cWrite)
	return p
}

// Private stuff.

// d60Hz is the duration of one frame at 60Hz.
const d60Hz = 16666667 * time.Nanosecond
const d30Hz = 33333333 * time.Nanosecond

func getDelay(s Strip) time.Duration {
	delay := s.MinDelay()
	defaultHz := d60Hz
	if rpi.Version() == 1 {
		// Use 30Hz on a rPi1 because it is too slow.
		defaultHz = d30Hz
	}
	if delay < defaultHz {
		delay = defaultHz
	}
	return delay
}

var black = &Color{}

func (p *Painter) runPattern(cGen, cWrite chan Frame) {
	defer p.wg.Done()
	defer func() {
		// Tell runWrite() to quit.
		cWrite <- nil
	}()
	ease := Transition{
		Before:     SPattern{black},
		After:      SPattern{black},
		Duration:   500 * time.Millisecond,
		Transition: TransitionEaseOut,
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
			ease.Before = ease.After
			ease.After.Pattern = newPat
			ease.Offset = since

		case pixels := <-cGen:
			for i := range pixels {
				pixels[i] = Color{}
			}
			ease.NextFrame(pixels, since)
			since += delay
			cWrite <- pixels

		case <-interrupt.Channel:
			return
		}
	}
}

func (p *Painter) runWrite(cGen, cWrite chan Frame) {
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
