// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// playing is a small app to play with the pins, nothing more. You are not
// expected to use it as-is.
package main

import (
	"errors"
	"flag"
	"fmt"
	"image"
	"image/draw"
	"io/ioutil"
	"log"
	"os"
	"time"

	_ "image/png"

	"github.com/maruel/dlibox/go/bw2d"
	"github.com/maruel/dlibox/go/pio/devices/bme280"
	"github.com/maruel/dlibox/go/pio/devices/ssd1306"
	"github.com/maruel/dlibox/go/pio/host"
	"github.com/maruel/dlibox/go/pio/host/ir"
	"github.com/maruel/dlibox/go/pio/host/pins"
	"github.com/maruel/dlibox/go/pio/host/sysfs"
	"github.com/maruel/dlibox/go/psf"
	"github.com/maruel/interrupt"
)

func loadImg(path string) (*bw2d.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	src, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	r := src.Bounds()
	img, err := bw2d.Make(r.Max.X, r.Max.Y)
	if err != nil {
		return nil, err
	}
	draw.Draw(img, r, src, image.Point{}, draw.Src)
	return img, nil
}

func mainImpl() error {
	useBME280 := true
	useButton := true
	useIR := true
	usePir := true
	verbose := flag.Bool("v", false, "enable verbose logs")
	flag.Parse()

	if !*verbose {
		log.SetOutput(ioutil.Discard)
	}
	log.SetFlags(log.Lmicroseconds)

	button := make(chan bool)
	motion := make(chan bool)
	keys := make(chan host.Key)
	bme := make(chan env)

	f8, err := psf.Load("VGA8")
	if err != nil {
		return err
	}
	src, err := loadImg("bunny.png")
	if err != nil {
		return err
	}

	i, err := sysfs.MakeI2C(1)
	if err != nil {
		return err
	}

	// Display
	s, err := ssd1306.MakeI2C(i, 128, 64, false)
	if err != nil {
		return err
	}
	src.Inverse()
	img, err := bw2d.Make(s.W, s.H)
	if err != nil {
		return err
	}
	r := src.Bounds()
	r = r.Add(image.Point{(img.W - r.Max.X), (img.H - r.Max.Y) / 2})
	draw.Draw(img, r, src, image.Point{}, draw.Src)
	f8.Draw(img, 0, 0, bw2d.On, nil, "dlibox!")
	f8.Draw(img, 0, s.H-f8.H-1, bw2d.On, nil, "is awesome")
	if _, err = s.Write(img.Buf); err != nil {
		return err
	}
	go displayLoop(s, f8, img, button, motion, bme, keys)

	if useBME280 {
		b, err := bme280.Make(i, bme280.O4x, bme280.O4x, bme280.O4x, bme280.S20ms, bme280.F4)
		if err != nil {
			return err
		}
		defer b.Stop()
		go sensorLoop(b, bme)
	}

	if useButton {
		p := pins.ByNumber(24)
		if p == nil {
			return errors.New("no pin 24")
		}
		if err := p.In(host.Up); err != nil {
			return err
		}
		c, err := p.Edges()
		if err != nil {
			return err
		}
		go buttonLoop(c, button)
	}

	/*
		// Relays
		p := pins.ByNumber(17)
		if p == nil {
			return errors.New("no pin 17")
		}
		if err := .Out(); err != nil {
			return err
		}
		p.SetLow()
		p = pins.ByNumber(27)
		if p == nil {
			return errors.New("no pin 27")
		}
		if err := p.Out(); err != nil {
			return err
		}
		p.SetLow()
	*/

	if usePir {
		p := pins.ByNumber(19)
		if p == nil {
			return errors.New("no pin 19")
		}
		if err := p.In(host.Down); err != nil {
			return err
		}
		c, err := p.Edges()
		if err != nil {
			return err
		}
		go pirLoop(c, motion)
	}

	if useIR {
		irBus, err := ir.Make()
		if err != nil {
			return err
		}
		go irLoop(irBus, keys)
	}

	interrupt.HandleCtrlC()
	<-interrupt.Channel

	return nil
}

func displayLoop(s *ssd1306.Dev, f *psf.Font, img *bw2d.Image, button, motion <-chan bool, bme <-chan env, keys <-chan host.Key) {
	tick := time.NewTicker(time.Second)
	defer tick.Stop()
	for {
		draw := false
		select {
		case b := <-button:
			if b {
				f.Draw(img, 0, f.H*4, bw2d.On, bw2d.Off, "Bouton!")
			} else {
				f.Draw(img, 0, f.H*4, bw2d.On, bw2d.Off, "          ")
			}
			draw = true

		case m := <-motion:
			if m {
				f.Draw(img, 0, f.H*5, bw2d.On, bw2d.Off, "Mouvement!")
			} else {
				f.Draw(img, 0, f.H*5, bw2d.On, bw2d.Off, "          ")
			}

		case t := <-bme:
			f.Draw(img, 0, f.H, bw2d.On, bw2d.Off, fmt.Sprintf("%.2fC ", t.t))
			f.Draw(img, 0, f.H*2, bw2d.On, bw2d.Off, fmt.Sprintf("%.2fkPa ", t.p))
			f.Draw(img, 0, f.H*3, bw2d.On, bw2d.Off, fmt.Sprintf("%.2f%% ", t.h))

		case s := <-keys:
			f.Draw(img, 0, f.H*6, bw2d.On, bw2d.Off, string(s))
			draw = true

		case <-tick.C:
			f.Draw(img, 0, 0, bw2d.On, bw2d.Off, time.Now().Format("15:04:05"))
			draw = true

		case <-interrupt.Channel:
			break
		}
		if draw {
			if _, err := s.Write(img.Buf); err != nil {
				return
			}
		}
	}
}

func irLoop(irBus host.IR, keys chan<- host.Key) {
	c := irBus.Channel()
	for {
		select {
		case <-interrupt.Channel:
			break
		case msg := <-c:
			log.Printf("IR: %#v", msg)
			keys <- msg.Key
		}
	}
}

func buttonLoop(b chan host.Level, c chan<- bool) {
	for {
		select {
		case l := <-b:
			log.Printf("Bouton: %s", l)
			c <- l == host.Low
		case <-interrupt.Channel:
			break
		}
	}
}

func pirLoop(b chan host.Level, c chan<- bool) {
	for {
		select {
		case l := <-b:
			log.Printf("PIR: %s", l)
			c <- l == host.High
		case <-interrupt.Channel:
			break
		}
	}
}

func sensorLoop(b *bme280.Dev, c chan<- env) {
	tick := time.NewTicker(500 * time.Millisecond)
	defer tick.Stop()
	for {
		t, p, h, err := b.Read()
		if err != nil {
			log.Fatal(err)
		}
		//log.Printf("%6.3fC  %7.4fkPa  %7.4f%%", t, p, h)
		if err == nil {
			c <- env{t, p, h}
		}
		<-tick.C
	}
}

type env struct {
	t, p, h float32
}

func main() {
	if err := mainImpl(); err != nil {
		fmt.Fprintf(os.Stderr, "playing: %s.\n", err)
		os.Exit(1)
	}
}
