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
	"sync"

	"github.com/maruel/dlibox/go/modules/nodes/leds"
	"github.com/maruel/dlibox/go/msgbus"
	"periph.io/x/periph/conn/spi/spireg"
	"periph.io/x/periph/devices"
	"periph.io/x/periph/devices/apa102"
)

type LEDs struct {
	sync.Mutex
	strips []leds.Dev
}

func fake() {
	// Output (terminal with ANSI codes or APA102).
	// Hardcode to 100 characters when using a terminal output.
	// TODO(maruel): Query the terminal and use its width.
	/*
		num = 100
		fakeBytes = []byte("1")
		l = &strip{Display: screen.New(num), b: b, fps: 30}
	*/
}

func (l *LEDs) init(b msgbus.Bus) error {
	for _, cfg := range l.strips {
		fps := 30
		s, err := spireg.Open(cfg.SPI.ID)
		if err != nil {
			return err
		}
		if err = s.LimitSpeed(cfg.SPI.Hz); err != nil {
			return err
		}
		a, err := apa102.New(s, cfg.NumberLights, 255, 6500)
		if err != nil {
			return err
		}
		str := &strip{Display: a, b: b, s: s, fps: fps}
		c, err := b.Subscribe("leds/#", msgbus.BestEffort)
		if err != nil {
			str.Close()
			return err
		}
		//if err := b.Publish(msgbus.Message{"leds/fake", fakeBytes}, msgbus.MinOnce, true); err != nil {
		//	log.Printf("leds: publish failed: %v", err)
		//}
		if err := b.Publish(msgbus.Message{"leds/fps", []byte(strconv.Itoa(str.fps))}, msgbus.MinOnce, true); err != nil {
			log.Printf("leds: publish failed: %v", err)
		}
		if err := b.Publish(msgbus.Message{"leds/num", []byte(strconv.Itoa(cfg.NumberLights))}, msgbus.MinOnce, true); err != nil {
			log.Printf("leds: publish failed: %v", err)
		}
		if err := b.Publish(msgbus.Message{"leds/intensity", []byte("255")}, msgbus.MinOnce, true); err != nil {
			log.Printf("leds: publish failed: %v", err)
		}
		if err := b.Publish(msgbus.Message{"leds/temperature", []byte("6500")}, msgbus.MinOnce, true); err != nil {
			log.Printf("leds: publish failed: %v", err)
		}
		go func() {
			for msg := range c {
				str.onMsg(msg)
			}
		}()
	}
	return nil
}

type strip struct {
	devices.Display
	s   io.Closer
	fps int
	b   msgbus.Bus
}

func (l *strip) Close() error {
	l.b.Unsubscribe("leds/#")
	if l.s != nil {
		return l.s.Close()
	}
	_, err := os.Stdout.Write([]byte("\033[0m\n"))
	return err
}

// Support both relative and absolute values.
func processRel(topic string, p []byte) (int, int, error) {
	if len(p) == 0 {
		return 0, 0, fmt.Errorf("leds: %s missing payload", topic)
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
		return 0, 0, fmt.Errorf("leds: %s: %v", topic, err)
	}
	return v, op, nil
}

func (l *strip) onMsg(msg msgbus.Message) {
	switch msg.Topic {
	case "leds/fake":
	case "leds/fps":
	case "leds/intensity":
		a, ok := l.Display.(*apa102.Dev)
		if !ok {
			log.Printf("leds: can't set intensity with fake LED")
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
	case "leds/num":
	case "leds/temperature":
		a, ok := l.Display.(*apa102.Dev)
		if !ok {
			log.Printf("leds: can't set temperature with fake LED")
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
		log.Printf("leds: unknown msg: %# v", msg)
	}
}
