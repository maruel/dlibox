// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/maruel/anim1d"
	"github.com/maruel/dlibox/go/msgbus"
	"github.com/maruel/interrupt"
	"github.com/pkg/errors"
	"periph.io/x/periph/devices"
)

// LRU is the list of recent patterns. The first is the oldest.
type LRU struct {
	sync.Mutex
	Max      int
	Patterns []Pattern
}

func (l *LRU) ResetDefault() {
	l.Lock()
	defer l.Unlock()
	l.Max = 25
	l.Patterns = []Pattern{
		"{\"_type\":\"Aurore\"}",
		"{\"Child\":{\"Frame\":\"Lff0000ff0000ff0000ff0000ff0000ffffffffffffffffffffffffffffff\",\"_type\":\"Repeated\"},\"MovePerHour\":21600,\"_type\":\"Rotate\"}",
		"{\"Patterns\":[{\"_type\":\"Aurore\"},{\"C\":\"#ff9000\",\"_type\":\"NightStars\"},{\"AverageDelay\":0,\"Duration\":0,\"_type\":\"WishingStar\"}],\"_type\":\"Add\"}",
		"{\"Curve\":\"easeinout\",\"Patterns\":[\"#ff0000\",\"#00ff00\",\"#0000ff\"],\"ShowMS\":1000000,\"TransitionMS\":1000000,\"_type\":\"Loop\"}",
		"{\"Curve\":\"direct\",\"Left\":\"#000000\",\"Right\":\"#0000ff\",\"_type\":\"Gradient\"}",
		"{\"Curve\":\"direct\",\"Left\":\"#000000\",\"Right\":\"#ff0000\",\"_type\":\"Gradient\"}",
		"{\"Curve\":\"direct\",\"Left\":\"#000000\",\"Right\":\"#00ff00\",\"_type\":\"Gradient\"}",
		"{\"Curve\":\"direct\",\"Left\":\"#000000\",\"Right\":\"#ffffff\",\"_type\":\"Gradient\"}",
		"{\"Child\":\"Lff0000ff0000ee0000dd0000cc0000bb0000aa0000990000880000770000660000550000440000330000220000110000\",\"MovePerHour\":108000,\"_type\":\"PingPong\"}",
		"{\"After\":\"#000000\",\"Before\":{\"After\":\"#ffffff\",\"Before\":{\"After\":\"#ff7f00\",\"Before\":\"#000000\",\"Curve\":\"direct\",\"OffsetMS\":0,\"TransitionMS\":0,\"_type\":\"Transition\"},\"Curve\":\"direct\",\"OffsetMS\":0,\"TransitionMS\":0,\"_type\":\"Transition\"},\"Curve\":\"direct\",\"OffsetMS\":0,\"TransitionMS\":0,\"_type\":\"Transition\"}",
		"\"#000000\"",
		"{\"Child\":\"Lffffff\",\"MovePerHour\":108000,\"_type\":\"PingPong\"}",
		"{\"Curve\":\"easeinout\",\"Patterns\":[\"#ff0000\",\"#ff7f00\",\"#ffff00\",\"#00ff00\",\"#0000ff\",\"#4b0082\",\"#8b00ff\"],\"ShowMS\":1000000,\"TransitionMS\":10000000,\"_type\":\"Loop\"}",
		"\"Rainbow\"",
		"{\"C\":\"#ff9000\",\"_type\":\"NightStars\"}",
		"{\"Child\":\"L010001ff000000ff000000ff\",\"_type\":\"Chronometer\"}",
		morning,
	}
}

func (l *LRU) Validate() error {
	l.Lock()
	defer l.Unlock()
	for i, s := range l.Patterns {
		if err := s.Validate(); err != nil {
			return errors.Wrap(err, fmt.Sprintf("can't load recent pattern %d", i))
		}
	}
	return nil
}

// Inject moves the pattern at the top of LRU cache.
func (l *LRU) Inject(pattern string) {
	l.Lock()
	defer l.Unlock()
	if l.Max == 0 {
		// Always use a sain default value.
		l.Max = 25
	}
	for i, old := range l.Patterns {
		if old == Pattern(pattern) {
			copy(l.Patterns[i:], l.Patterns[i+1:])
			l.Patterns = l.Patterns[:len(l.Patterns)-1]
			break
		}
	}
	if len(l.Patterns) < l.Max {
		l.Patterns = append(l.Patterns, "")
	}
	if len(l.Patterns) > 1 {
		copy(l.Patterns[1:], l.Patterns)
	}
	l.Patterns[0] = Pattern(pattern)
}

// Painter contains settings about patterns.
type Painter struct {
	sync.Mutex
	Named   map[string]Pattern // Patterns that are 'named'.
	Startup Pattern            // Startup pattern to use. If not set, use Last.
	Last    Pattern            // Last pattern used.
}

func (p *Painter) ResetDefault() {
	p.Lock()
	defer p.Unlock()
	p.Named = map[string]Pattern{
		"Aurora":      "{\"_type\":\"Aurore\"}",
		"Candy":       "{\"Child\":{\"Frame\":\"Lff0000ff0000ff0000ff0000ff0000ffffffffffffffffffffffffffffff\",\"_type\":\"Repeated\"},\"MovePerHour\":21600,\"_type\":\"Rotate\"}",
		"Night":       "{\"Patterns\":[{\"_type\":\"Aurore\"},{\"C\":\"#ff9000\",\"_type\":\"NightStars\"},{\"AverageDelay\":0,\"Duration\":0,\"_type\":\"WishingStar\"}],\"_type\":\"Add\"}",
		"Colors":      "{\"Curve\":\"easeinout\",\"Patterns\":[\"#ff0000\",\"#00ff00\",\"#0000ff\"],\"ShowMS\":1000000,\"TransitionMS\":1000000,\"_type\":\"Loop\"}",
		"Blue":        "{\"Curve\":\"direct\",\"Left\":\"#000000\",\"Right\":\"#0000ff\",\"_type\":\"Gradient\"}",
		"Red":         "{\"Curve\":\"direct\",\"Left\":\"#000000\",\"Right\":\"#ff0000\",\"_type\":\"Gradient\"}",
		"Green":       "{\"Curve\":\"direct\",\"Left\":\"#000000\",\"Right\":\"#00ff00\",\"_type\":\"Gradient\"}",
		"Gray":        "{\"Curve\":\"direct\",\"Left\":\"#000000\",\"Right\":\"#ffffff\",\"_type\":\"Gradient\"}",
		"K2000":       "{\"Child\":\"Lff0000ff0000ee0000dd0000cc0000bb0000aa0000990000880000770000660000550000440000330000220000110000\",\"MovePerHour\":108000,\"_type\":\"PingPong\"}",
		"WakeUp":      "{\"After\":\"#000000\",\"Before\":{\"After\":\"#ffffff\",\"Before\":{\"After\":\"#ff7f00\",\"Before\":\"#000000\",\"Curve\":\"direct\",\"OffsetMS\":0,\"TransitionMS\":60000,\"_type\":\"Transition\"},\"Curve\":\"direct\",\"OffsetMS\":0,\"TransitionMS\":0,\"_type\":\"Transition\"},\"Curve\":\"direct\",\"OffsetMS\":0,\"TransitionMS\":0,\"_type\":\"Transition\"}",
		"Black":       "\"#000000\"",
		"Dot":         "{\"Child\":\"Lffffff\",\"MovePerHour\":108000,\"_type\":\"PingPong\"}",
		"":            "{\"Curve\":\"easeinout\",\"Patterns\":[\"#ff0000\",\"#ff7f00\",\"#ffff00\",\"#00ff00\",\"#0000ff\",\"#4b0082\",\"#8b00ff\"],\"ShowMS\":1000000,\"TransitionMS\":10000000,\"_type\":\"Loop\"}",
		"Rainbow":     "\"Rainbow\"",
		"NightStars":  "{\"C\":\"#ff9000\",\"_type\":\"NightStars\"}",
		"Chronometer": "{\"Child\":\"L0100010f0000000f0000000f\",\"_type\":\"Chronometer\"}",
		"Morning":     morning,
	}
	p.Startup = "\"#010001\""
	p.Last = "\"#010001\""
}

func (p *Painter) Validate() error {
	p.Lock()
	defer p.Unlock()
	for k, v := range p.Named {
		if err := v.Validate(); err != nil {
			return errors.Wrap(err, fmt.Sprintf("can't load pattern %s", k))
		}
	}
	if len(p.Startup) != 0 {
		if err := p.Startup.Validate(); err != nil {
			return errors.Wrap(err, "can't load pattern for Last")
		}
	}
	if len(p.Last) != 0 {
		return p.Last.Validate()
	}
	return nil
}

func initPainter(b msgbus.Bus, leds devices.Display, fps int, config *Painter, lru *LRU) (*painterNode, error) {
	config.Lock()
	defer config.Unlock()
	lru.Lock()
	defer lru.Unlock()
	p := newPainter(leds, fps)
	if len(config.Last) != 0 {
		if err := p.SetPattern(string(config.Last), 500*time.Millisecond); err != nil {
			return nil, err
		}
	} else if len(config.Startup) != 0 {
		if err := p.SetPattern(string(config.Startup), 500*time.Millisecond); err != nil {
			return nil, err
		}
	}
	c, err := b.Subscribe("painter/#", msgbus.BestEffort)
	if err != nil {
		return nil, err
	}
	pp := &painterNode{p, b, config, lru}
	go func() {
		for msg := range c {
			pp.onMsg(msg)
		}
	}()
	return pp, nil
}

type painterNode struct {
	p      *painterLoop
	b      msgbus.Bus
	config *Painter
	lru    *LRU
}

func (p *painterNode) Close() error {
	p.b.Unsubscribe("painter/#")
	return p.p.Close()
}

func (p *painterNode) onMsg(msg msgbus.Message) {
	switch msg.Topic {
	case "painter/setautomated":
		p.setautomated(msg.Payload)
	case "painter/setnow":
		p.setnow(msg.Payload)
	case "painter/setuser":
		p.setuser(msg.Payload)
	default:
		log.Printf("painter unknown msg: %# v", msg)
	}
}

func (p *painterNode) setautomated(payload []byte) {
	// Skip the LRU.
	s := string(payload)
	if err := p.p.SetPattern(s, 500*time.Millisecond); err != nil {
		log.Printf("painter.setautomated: invalid payload: %s", s)
	}
}

func (p *painterNode) setnow(payload []byte) {
	// Skip the 500ms ease-out.
	s := string(payload)
	if err := p.p.SetPattern(s, 0); err != nil {
		log.Printf("painter.setnow: invalid payload: %s", s)
	}
}

func (p *painterNode) setuser(payload []byte) {
	// Add it to the LRU.
	s := string(payload)
	if err := p.p.SetPattern(s, 500*time.Millisecond); err != nil {
		log.Printf("painter.setuser: invalid payload: %s", s)
		return
	}
	var pat anim1d.SPattern
	if err := pat.UnmarshalJSON(payload); err != nil {
		log.Printf("painter.setuser: internal error: %s", s)
		return
	}
	if c, ok := pat.Pattern.(*anim1d.Color); !ok || (c.R == 0 && c.G == 0 && c.B == 0) {
		r, _ := pat.MarshalJSON()
		p.lru.Inject(string(r))
	}
	p.config.Lock()
	defer p.config.Unlock()
	p.config.Last = Pattern(s)
}

//

// painterLoop handles the "draw frame, write" loop.
type painterLoop struct {
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
func (p *painterLoop) SetPattern(s string, transition time.Duration) error {
	var pat anim1d.SPattern
	if err := json.Unmarshal([]byte(s), &pat); err != nil {
		return err
	}
	p.c <- newPattern{pat.Pattern, transition}
	return nil
}

func (p *painterLoop) Close() error {
	select {
	case p.c <- newPattern{}:
	default:
	}
	close(p.c)
	p.wg.Wait()
	return nil
}

// newPainter returns a painterLoop that manages updating the Patterns to the
// strip.
//
// It Assumes the display uses native RGB packed pixels.
func newPainter(d devices.Display, fps int) *painterLoop {
	p := &painterLoop{
		d:             d,
		c:             make(chan newPattern),
		frameDuration: time.Second / time.Duration(fps),
	}
	numLights := d.Bounds().Dx()
	// Tripple buffering.
	cGen := make(chan anim1d.Frame, 3)
	cWrite := make(chan anim1d.Frame, cap(cGen))
	for i := 0; i < cap(cGen); i++ {
		cGen <- make(anim1d.Frame, numLights)
	}
	p.wg.Add(2)
	go p.runPattern(cGen, cWrite)
	go p.runWrite(cGen, cWrite, numLights)
	return p
}

// Private stuff.

var black = &anim1d.Color{}

type newPattern struct {
	p anim1d.Pattern
	d time.Duration
}

func (p *painterLoop) runPattern(cGen <-chan anim1d.Frame, cWrite chan<- anim1d.Frame) {
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

	var root anim1d.Pattern = black
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
				root = &anim1d.Transition{
					Before:       anim1d.SPattern{root},
					After:        anim1d.SPattern{newPat.p},
					OffsetMS:     uint32(since / time.Millisecond),
					TransitionMS: uint32(newPat.d / time.Millisecond),
					Curve:        anim1d.EaseOut,
				}
			}

		case pixels, ok := <-cGen:
			if !ok {
				return
			}
			for i := range pixels {
				pixels[i] = anim1d.Color{}
			}
			timeMS := uint32(since / time.Millisecond)
			root.Render(pixels, timeMS)
			since += p.frameDuration
			cWrite <- pixels
			if t, ok := root.(*anim1d.Transition); ok {
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

func (p *painterLoop) runWrite(cGen chan<- anim1d.Frame, cWrite <-chan anim1d.Frame, numLights int) {
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
