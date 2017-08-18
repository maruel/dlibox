// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package device

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/maruel/dlibox/nodes"
	"github.com/maruel/msgbus"
	"periph.io/x/periph/conn/spi/spireg"
	"periph.io/x/periph/devices"
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
	apa, err := apa102.New(s, a.Cfg.NumberLights, 255, 6500)
	if err != nil {
		return err
	}
	str := &strip{Display: apa, b: b, s: s, fps: a.Cfg.FPS}
	c, err := b.Subscribe("anim1d/#", msgbus.BestEffort)
	if err != nil {
		str.Close()
		return err
	}
	//if err := b.Publish(msgbus.Message{"$fake", fakeBytes}, msgbus.MinOnce, true); err != nil {
	//	log.Printf("anim1d: publish failed: %v", err)
	//}
	if err := b.Publish(msgbus.Message{"$fps", []byte(strconv.Itoa(str.fps))}, msgbus.MinOnce, true); err != nil {
		log.Printf("anim1d: publish failed: %v", err)
	}
	if err := b.Publish(msgbus.Message{"$num", []byte(strconv.Itoa(a.Cfg.NumberLights))}, msgbus.MinOnce, true); err != nil {
		log.Printf("anim1d: publish failed: %v", err)
	}
	if err := b.Publish(msgbus.Message{"intensity", []byte("255")}, msgbus.MinOnce, true); err != nil {
		log.Printf("anim1d: publish failed: %v", err)
	}
	if err := b.Publish(msgbus.Message{"temperature", []byte("6500")}, msgbus.MinOnce, true); err != nil {
		log.Printf("anim1d: publish failed: %v", err)
	}
	go func() {
		for msg := range c {
			str.onMsg(msg)
		}
	}()
	return nil
}

type strip struct {
	devices.Display
	s   io.Closer
	fps int
	b   msgbus.Bus
}

func (l *strip) Close() error {
	l.b.Unsubscribe("anim1d/#")
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

func (l *strip) onMsg(msg msgbus.Message) {
	switch msg.Topic {
	case "anim1d/fake":
	case "anim1d/fps":
	case "anim1d/intensity":
		a, ok := l.Display.(*apa102.Dev)
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
	case "anim1d/num":
	case "anim1d/temperature":
		a, ok := l.Display.(*apa102.Dev)
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
	default:
		log.Printf("anim1d: unknown msg: %# v", msg)
	}
}
