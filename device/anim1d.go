// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package device

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/maruel/anim1d"
	"github.com/maruel/dlibox/nodes"
	"github.com/maruel/dlibox/shared"
	"github.com/maruel/interrupt"
	"github.com/maruel/msgbus"
	"periph.io/x/periph/conn/display"
	"periph.io/x/periph/conn/spi/spireg"
	"periph.io/x/periph/devices/apa102"
)

type anim1DDev struct {
	NodeBase
	Cfg *nodes.Anim1D
}

func (a *anim1DDev) init(b msgbus.Bus) error {
	s, err := spireg.Open(a.Cfg.SPI.ID)
	if err != nil {
		return err
	}
	if err = s.LimitSpeed(a.Cfg.SPI.Hz); err != nil {
		return err
	}
	opts := apa102.DefaultOpts
	opts.NumPixels = a.Cfg.NumberLights
	opts.Temperature = 6500
	apa, err := apa102.New(s, &opts)
	if err != nil {
		return err
	}
	str := &strip{Drawer: apa, b: b, s: s, fps: a.Cfg.FPS}
	if err != nil {
		str.Close()
		return err
	}
	/*
		if err := b.Publish(msgbus.Message{"$fake", fakeBytes}, msgbus.ExactlyOnce, true); err != nil {
			log.Printf("anim1d: publish failed: %v", err)
		}
	*/
	c, err := b.Subscribe("#", msgbus.ExactlyOnce)
	if err != nil {
		return err
	}
	shared.RetainedStr(b, "$fps", strconv.Itoa(str.fps))
	shared.RetainedStr(b, "$num", strconv.Itoa(a.Cfg.NumberLights))
	shared.RetainedStr(b, "intensity", "255")
	shared.RetainedStr(b, "temperature", "6500")

	p := newPainter(apa, str.fps)
	if err := p.SetPattern(`"#800000"`, 500*time.Millisecond); err != nil {
		return err
	}

	go func() {
		for msg := range c {
			str.onMsg(p, msg)
		}
	}()
	return nil
}

type strip struct {
	display.Drawer
	s   io.Closer
	fps int
	b   msgbus.Bus
}

func (l *strip) Close() error {
	l.b.Unsubscribe("#")
	if l.s != nil {
		return l.s.Close()
	}
	_, err := os.Stdout.Write([]byte("\033[0m\n"))
	return err
}

// Support both relative and absolute values.
func processRel(topic string, p []byte) (int, int, error) {
	if len(p) == 0 {
		return 0, 0, fmt.Errorf("anim1d: %s missing payload", topic)
	}
	s := string(p)
	op := 0
	if p[0] == '+' {
		op = 1
		s = s[1:]
	} else if p[0] == '-' {
		op = 2
		s = s[1:]
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0, 0, fmt.Errorf("anim1d: %s: %v", topic, err)
	}
	return v, op, nil
}

func (l *strip) onMsg(p *painterLoop, msg msgbus.Message) {
	switch msg.Topic {
	case "anim1d":
		s := string(msg.Payload)
		if err := p.SetPattern(s, 100*time.Millisecond); err != nil {
			log.Printf("painter.setautomated: invalid payload: %s", s)
		}
		break

	case "fake":
	case "fps":
	case "intensity":
		a, ok := l.Drawer.(*apa102.Dev)
		if !ok {
			log.Printf("anim1d: can't set intensity with fake LED")
			return
		}
		v, op, err := processRel(msg.Topic, msg.Payload)
		if err != nil {
			log.Print(err.Error())
			return
		}
		switch op {
		case 0:
		case 1:
			v = int(a.Intensity) + v
		case 2:
			v = int(a.Intensity) - v
		}
		if v < 0 {
			v = 0
		} else if v > 255 {
			v = 255
		}
		a.Intensity = uint8(v)
	case "num":
	case "temperature":
		a, ok := l.Drawer.(*apa102.Dev)
		if !ok {
			log.Printf("anim1d: can't set temperature with fake LED")
			return
		}
		v, op, err := processRel(msg.Topic, msg.Payload)
		if err != nil {
			log.Print(err.Error())
			return
		}
		switch op {
		case 0:
		case 1:
			v = int(a.Temperature) + v
		case 2:
			v = int(a.Temperature) - v
		}
		if v < 1000 {
			v = 1000
		} else if v > 35000 {
			v = 35000
		}
		a.Temperature = uint16(v)
	case "$fps":
	case "$num":
		break
	default:
		log.Printf("anim1d: unknown msg: %# v", msg)
	}
}

//

// painterLoop handles the "draw frame, write" loop.
type painterLoop struct {
	d             display.Drawer
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
func newPainter(d display.Drawer, fps int) *painterLoop {
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

	var root anim1d.Pattern = &anim1d.Color{}
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
					Before:       anim1d.SPattern{Pattern: root},
					After:        anim1d.SPattern{Pattern: newPat.p},
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
	// TODO(maruel): This is wrong.
	w := p.d.(io.Writer)
	for {
		pixels, ok := <-cWrite
		if pixels == nil || !ok {
			return
		}
		if err == nil {
			pixels.ToRGB(buf)
			if _, err = w.Write(buf); err != nil {
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
