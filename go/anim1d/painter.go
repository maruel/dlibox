// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package anim1d

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/maruel/interrupt"
	"periph.io/x/periph/devices"
)

// Pattern is a interface to draw an animated line.
type Pattern interface {
	// Render fills the buffer with the image at this time frame.
	//
	// The image should be derived from timeMS, which is the time since this
	// pattern was started.
	//
	// Calling Render() with a nil pattern is valid. Patterns should be callable
	// without crashing with an object initialized with default values.
	//
	// timeMS will cycle after 49.7 days. The reason it's not using time.Duration
	// is that int64 calculation on ARM is very slow and abysmal on xtensa, which
	// this code is transpiled to.
	Render(pixels Frame, timeMS uint32)
}

// Painter handles the "draw frame, write" loop.
type Painter struct {
	d             devices.Display
	c             chan newPattern
	wg            sync.WaitGroup
	frameDuration time.Duration
}

// SetPattern changes the current pattern to a new one.
//
// The pattern is in JSON encoded format. The function will return an error if
// the encoding is bad. The function is synchronous, it returns only after the
// pattern was effectively set.
func (p *Painter) SetPattern(s string, transition time.Duration) error {
	var pat SPattern
	if err := json.Unmarshal([]byte(s), &pat); err != nil {
		return err
	}
	p.c <- newPattern{pat.Pattern, transition}
	return nil
}

func (p *Painter) Close() error {
	select {
	case p.c <- newPattern{}:
	default:
	}
	close(p.c)
	p.wg.Wait()
	return nil
}

// NewPainter returns a Painter that manages updating the Patterns to the
// strip.
//
// It Assumes the display uses native RGB packed pixels.
func NewPainter(d devices.Display, fps int) *Painter {
	p := &Painter{
		d:             d,
		c:             make(chan newPattern),
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

type newPattern struct {
	p Pattern
	d time.Duration
}

func (p *Painter) runPattern(cGen <-chan Frame, cWrite chan<- Frame) {
	defer func() {
		// Tell runWrite() to quit.
		for loop := true; loop; {
			select {
			case _, loop = <-cGen:
			default:
				loop = false
			}
		}
		select {
		case cWrite <- nil:
		default:
		}
		close(cWrite)
		p.wg.Done()
	}()

	var root Pattern = black
	var since time.Duration
	for {
		select {
		case newPat, ok := <-p.c:
			if newPat.p == nil || !ok {
				// Request to terminate.
				return
			}

			// New pattern.
			if newPat.d == 0 {
				root = newPat.p
			} else {
				root = &Transition{
					Before:       SPattern{root},
					After:        SPattern{newPat.p},
					OffsetMS:     uint32(since / time.Millisecond),
					TransitionMS: uint32(newPat.d / time.Millisecond),
					Curve:        EaseOut,
				}
			}

		case pixels, ok := <-cGen:
			if !ok {
				return
			}
			for i := range pixels {
				pixels[i] = Color{}
			}
			timeMS := uint32(since / time.Millisecond)
			root.Render(pixels, timeMS)
			since += p.frameDuration
			cWrite <- pixels
			if t, ok := root.(*Transition); ok {
				if t.OffsetMS+t.TransitionMS < timeMS {
					root = t.After.Pattern
					since -= time.Duration(t.OffsetMS) * time.Millisecond
				}
			}

		case <-interrupt.Channel:
			return
		}
	}
}

func (p *Painter) runWrite(cGen chan<- Frame, cWrite <-chan Frame, numLights int) {
	defer func() {
		// Tell runPattern() to quit.
		for loop := true; loop; {
			select {
			case _, loop = <-cWrite:
			default:
				loop = false
			}
		}
		select {
		case cGen <- nil:
		default:
		}
		close(cGen)
		p.wg.Done()
	}()

	tick := time.NewTicker(p.frameDuration)
	defer tick.Stop()
	var err error
	buf := make([]byte, numLights*3)
	for {
		pixels, ok := <-cWrite
		if pixels == nil || !ok {
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
