// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"log"

	"github.com/google/pio/devices"
	"github.com/maruel/dlibox/go/anim1d"
	"github.com/maruel/dlibox/go/modules"
)

func initPainter(bus modules.Bus, leds devices.Display, fps int, config *Painter) (*painter, error) {
	p := anim1d.NewPainter(leds, fps)
	if len(config.Last) != 0 {
		if err := p.SetPattern(string(config.Last)); err != nil {
			return nil, err
		}
	}
	if len(config.Startup) != 0 {
		if err := p.SetPattern(string(config.Startup)); err != nil {
			return nil, err
		}
	}
	c, err := bus.Subscribe("painter", modules.ExactlyOnce)
	if err != nil {
		return nil, err
	}
	pp := &painter{p, bus}
	go func() {
		for msg := range c {
			pp.onMsg(msg)
		}
	}()
	return pp, nil
}

type painter struct {
	p   *anim1d.Painter
	bus modules.Bus
}

// Temporary.
func (p *painter) SetPattern(s string) error {
	return p.p.SetPattern(s)
}

func (p *painter) Close() error {
	p.bus.Unsubscribe("painter")
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
		log.Printf("painter.setautomated: invalid payload: %s", s)
	}
}

func (p *painter) setuser(payload []byte) {
	// Add it to the LRU.
	s := string(payload)
	if err := p.p.SetPattern(s); err != nil {
		log.Printf("painter.setautomated: invalid payload: %s", s)
	}
}
