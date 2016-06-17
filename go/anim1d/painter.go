// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package anim1d

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/maruel/dlibox/go/pio/devices"
	"github.com/maruel/interrupt"
)

// Pattern is a interface to draw an animated line.
type Pattern interface {
	// NextFrame fills the buffer with the image at this time frame.
	//
	// The image should be derived from timeMS, which is the time since this
	// pattern was started.
	//
	// Calling NextFrame() with a nil pattern is valid. Patterns should be
	// callable without crashing with an object initialized with default values.
	//
	// timeMS will cycle after 49.7 days. The reason it's not using time.Duration
	// is that int64 calculation on ARM is very slow and abysmal on xtensa, which
	// this code is transpiled to.
	NextFrame(pixels Frame, timeMS uint32)
}

// Painter handles the "draw frame, write" loop.
type Painter struct {
	d             devices.Display
	c             chan Pattern
	wg            sync.WaitGroup
	frameDuration time.Duration
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
	p.c <- pat.Pattern
	return nil
}

func (p *Painter) Close() error {
	p.c <- nil
	p.wg.Wait()
	p.c = nil
	return nil
}

// NewPainter returns a Painter that manages updating the Patterns to the
// strip.
//
// It Assumes the display uses native RGB packed pixels.
func NewPainter(d devices.Display, fps int) *Painter {
	p := &Painter{
		d:             d,
		c:             make(chan Pattern),
		frameDuration: time.Second / time.Duration(fps),
	}
	numLights := d.Bounds().Dx()
	// Tripple buffering.
	cGen := make(chan Frame, 3)
	cWrite := make(chan Frame, cap(cGen))
	for i := 0; i < cap(cGen); i++ {
		cGen <- make(Frame, numLights)
	}
	p.wg.Add(2)
	go p.runPattern(cGen, cWrite)
	go p.runWrite(cGen, cWrite, numLights)
	return p
}

// Private stuff.

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
		DurationMS: 500,
		Curve:      EaseOut,
	}
	var since time.Duration
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
			ease.OffsetMS = uint32(since / time.Millisecond)

		case pixels := <-cGen:
			for i := range pixels {
				pixels[i] = Color{}
			}
			ease.NextFrame(pixels, uint32(since/time.Millisecond))
			since += p.frameDuration
			cWrite <- pixels

		case <-interrupt.Channel:
			return
		}
	}
}

func (p *Painter) runWrite(cGen, cWrite chan Frame, numLights int) {
	defer p.wg.Done()
	tick := time.NewTicker(p.frameDuration)
	defer tick.Stop()
	var err error
	buf := make([]byte, numLights*3)
	for {
		pixels := <-cWrite
		if len(pixels) == 0 {
			return
		}
		if err == nil {
			pixels.ToRGB(buf)
			if _, err = p.d.Write(buf); err != nil {
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
