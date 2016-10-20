// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/google/pio/devices"
	"github.com/maruel/dlibox/go/anim1d"
	"github.com/maruel/dlibox/go/modules"
	"github.com/pkg/errors"
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
		"{\"Patterns\":[{\"_type\":\"Aurore\"},{\"Seed\":0,\"Stars\":null,\"_type\":\"NightStars\"},{\"AverageDelay\":0,\"Duration\":0,\"_type\":\"WishingStar\"}],\"_type\":\"Add\"}",
		"{\"Curve\":\"easeinout\",\"DurationShowMS\":1000000,\"DurationTransitionMS\":1000000,\"Patterns\":[\"#ff0000\",\"#00ff00\",\"#0000ff\"],\"_type\":\"Loop\"}",
		"{\"Curve\":\"direct\",\"Left\":\"#000000\",\"Right\":\"#0000ff\",\"_type\":\"Gradient\"}",
		"{\"Curve\":\"direct\",\"Left\":\"#000000\",\"Right\":\"#ff0000\",\"_type\":\"Gradient\"}",
		"{\"Curve\":\"direct\",\"Left\":\"#000000\",\"Right\":\"#00ff00\",\"_type\":\"Gradient\"}",
		"{\"Curve\":\"direct\",\"Left\":\"#000000\",\"Right\":\"#ffffff\",\"_type\":\"Gradient\"}",
		"{\"Child\":\"Lff0000ff0000ee0000dd0000cc0000bb0000aa0000990000880000770000660000550000440000330000220000110000\",\"MovePerHour\":108000,\"_type\":\"PingPong\"}",
		"{\"After\":\"#000000\",\"Before\":{\"After\":\"#ffffff\",\"Before\":{\"After\":\"#ff7f00\",\"Before\":\"#000000\",\"Curve\":\"direct\",\"DurationMS\":0,\"OffsetMS\":0,\"_type\":\"Transition\"},\"Curve\":\"direct\",\"DurationMS\":0,\"OffsetMS\":0,\"_type\":\"Transition\"},\"Curve\":\"direct\",\"DurationMS\":0,\"OffsetMS\":0,\"_type\":\"Transition\"}",
		"\"#000000\"",
		"{\"Child\":\"Lffffff\",\"MovePerHour\":108000,\"_type\":\"PingPong\"}",
		"{\"Curve\":\"easeinout\",\"DurationShowMS\":1000000,\"DurationTransitionMS\":10000000,\"Patterns\":[\"#ff0000\",\"#ff7f00\",\"#ffff00\",\"#00ff00\",\"#0000ff\",\"#4b0082\",\"#8b00ff\"],\"_type\":\"Loop\"}",
		"\"Rainbow\"",
		"{\"Seed\":0,\"Stars\":null,\"_type\":\"NightStars\"}",
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
	p.Named = map[string]Pattern{}
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

func initPainter(b modules.Bus, leds devices.Display, fps int, config *Painter, lru *LRU) (*painter, error) {
	config.Lock()
	defer config.Unlock()
	lru.Lock()
	defer lru.Unlock()
	p := anim1d.NewPainter(leds, fps)
	if len(config.Last) != 0 {
		if err := p.SetPattern(string(config.Last)); err != nil {
			return nil, err
		}
	} else if len(config.Startup) != 0 {
		if err := p.SetPattern(string(config.Startup)); err != nil {
			return nil, err
		}
	}
	c, err := b.Subscribe("painter/+", modules.ExactlyOnce)
	if err != nil {
		return nil, err
	}
	pp := &painter{p, b, config, lru}
	go func() {
		for msg := range c {
			pp.onMsg(msg)
		}
	}()
	return pp, nil
}

type painter struct {
	p      *anim1d.Painter
	b      modules.Bus
	config *Painter
	lru    *LRU
}

// Temporary.
func (p *painter) SetPattern(s string) error {
	return p.p.SetPattern(s)
}

func (p *painter) Close() error {
	p.b.Unsubscribe("painter")
	return p.p.Close()
}

func (p *painter) onMsg(msg modules.Message) {
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

func (p *painter) setautomated(payload []byte) {
	// Skip the LRU.
	s := string(payload)
	if err := p.p.SetPattern(s); err != nil {
		log.Printf("painter.setautomated: invalid payload: %s", s)
	}
}

func (p *painter) setnow(payload []byte) {
	// Skip the 500ms ease-out.
	s := string(payload)
	if err := p.p.SetPattern(s); err != nil {
		log.Printf("painter.setnow: invalid payload: %s", s)
	}
}

func (p *painter) setuser(payload []byte) {
	// Add it to the LRU.
	s := string(payload)
	if err := p.p.SetPattern(s); err != nil {
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
