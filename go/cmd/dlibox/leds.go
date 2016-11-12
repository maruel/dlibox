// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Packages the static files in a .go file.
//go:generate go run ../package/main.go -out static_files_gen.go ../../../web

// dlibox drives the dlibox LED strip on a Raspberry Pi. It runs a web server
// for remote control.
package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strconv"
	"sync"

	"github.com/maruel/dlibox/go/donotuse/conn/spi"
	"github.com/maruel/dlibox/go/donotuse/devices"
	"github.com/maruel/dlibox/go/donotuse/devices/apa102"
	"github.com/maruel/dlibox/go/donotuse/host"
	"github.com/maruel/dlibox/go/modules"
	"github.com/maruel/dlibox/go/screen"
)

// APA102 contains light specific settings.
type APA102 struct {
	sync.Mutex
	// BusNumber is the SPI bus number to use, defaults to -1.
	BusNumber int
	// Speed of the transfer.
	SPIspeed int64
	// Number of lights controlled by this device. If lower than the actual
	// number of lights, the remaining lights will flash oddly.
	NumberLights int
}

func (a *APA102) ResetDefault() {
	a.Lock()
	defer a.Unlock()
	a.BusNumber = -1
	a.SPIspeed = 10000000
	a.NumberLights = 150
}

func (a *APA102) Validate() error {
	a.Lock()
	defer a.Unlock()
	return nil
}

// initLEDs initializes the LED strip.
func initLEDs(b modules.Bus, fake bool, config *APA102) (*leds, error) {
	if config.NumberLights == 0 {
		return nil, nil
	}
	var l *leds
	var fakeBytes []byte
	num := config.NumberLights
	if fake {
		// Output (terminal with ANSI codes or APA102).
		// Hardcode to 100 characters when using a terminal output.
		// TODO(maruel): Query the terminal and use its width.
		num = 100
		fakeBytes = []byte("1")
		l = &leds{Display: screen.New(num), b: b, fps: 30}
	} else {
		fps := 60
		if host.MaxSpeed() < 900000 || runtime.NumCPU() < 4 {
			// Use 30Hz on slower devices because it is too slow.
			fps = 30
		}
		s, err := spi.New(config.BusNumber, 0)
		if err != nil {
			return nil, err
		}
		if err = s.Speed(config.SPIspeed); err != nil {
			return nil, err
		}
		a, err := apa102.New(s, config.NumberLights, 255, 6500)
		if err != nil {
			return nil, err
		}
		l = &leds{Display: a, b: b, s: s, fps: fps}
	}
	c, err := b.Subscribe("leds/#", modules.BestEffort)
	if err != nil {
		l.Close()
		return nil, err
	}
	if err := b.Publish(modules.Message{"leds/fake", fakeBytes}, modules.MinOnce, true); err != nil {
		log.Printf("leds: publish failed: %v", err)
	}
	if err := b.Publish(modules.Message{"leds/fps", []byte(strconv.Itoa(l.fps))}, modules.MinOnce, true); err != nil {
		log.Printf("leds: publish failed: %v", err)
	}
	if err := b.Publish(modules.Message{"leds/num", []byte(strconv.Itoa(num))}, modules.MinOnce, true); err != nil {
		log.Printf("leds: publish failed: %v", err)
	}
	if err := b.Publish(modules.Message{"leds/intensity", []byte("255")}, modules.MinOnce, true); err != nil {
		log.Printf("leds: publish failed: %v", err)
	}
	if err := b.Publish(modules.Message{"leds/temperature", []byte("6500")}, modules.MinOnce, true); err != nil {
		log.Printf("leds: publish failed: %v", err)
	}
	go func() {
		for msg := range c {
			l.onMsg(msg)
		}
	}()
	return l, nil
}

type leds struct {
	devices.Display
	s   io.Closer
	fps int
	b   modules.Bus
}

func (l *leds) Close() error {
	err := l.b.Unsubscribe("leds/#")
	if err != nil {
		log.Printf("leds: failed to unsubscribe: leds/#: %v", err)
	}
	if l.s != nil {
		if err1 := l.s.Close(); err1 != nil {
			err = err1
		}
	} else {
		if _, err1 := os.Stdout.Write([]byte("\033[0m\n")); err1 != nil {
			err = err1
		}
	}
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

func (l *leds) onMsg(msg modules.Message) {
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
