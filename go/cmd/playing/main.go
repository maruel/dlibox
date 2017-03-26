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

	"github.com/maruel/dlibox/go/psf"
	"github.com/maruel/interrupt"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/conn/ir"
	"periph.io/x/periph/devices"
	"periph.io/x/periph/devices/bme280"
	"periph.io/x/periph/devices/lirc"
	"periph.io/x/periph/devices/ssd1306"
	"periph.io/x/periph/devices/ssd1306/image1bit"
	"periph.io/x/periph/host"
)

func loadImg(path string) (*image1bit.Image, error) {
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
	img, err := image1bit.New(image.Rect(0, 0, r.Max.X, r.Max.Y))
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
	name := flag.String("i2c", "", "I²C bus to use")
	verbose := flag.Bool("v", false, "enable verbose logs")
	flag.Parse()

	if !*verbose {
		log.SetOutput(ioutil.Discard)
	}
	log.SetFlags(log.Lmicroseconds)

	// Initialize periph.
	if _, err := host.Init(); err != nil {
		return err
	}

	button := make(chan bool)
	motion := make(chan bool)
	keys := make(chan ir.Key)
	env := make(chan *devices.Environment)

	f8, err := psf.Load("VGA8")
	if err != nil {
		return err
	}
	src, err := loadImg("bunny.png")
	if err != nil {
		return err
	}

	i, err := i2c.OpenByName(*name)
	if err != nil {
		return err
	}

	// Display
	s, err := ssd1306.NewI2C(i, 128, 64, false)
	if err != nil {
		return err
	}
	src.Inverse()
	img, err := image1bit.New(image.Rect(0, 0, s.W, s.H))
	if err != nil {
		return err
	}
	r := src.Bounds()
	r = r.Add(image.Point{(img.W - r.Max.X), (img.H - r.Max.Y) / 2})
	draw.Draw(img, r, src, image.Point{}, draw.Src)
	f8.Draw(img, 0, 0, image1bit.On, nil, "dlibox!")
	f8.Draw(img, 0, s.H-f8.H-1, image1bit.On, nil, "is awesome")
	if _, err = s.Write(img.Buf); err != nil {
		return err
	}
	go displayLoop(s, f8, img, button, motion, env, keys)

	if useBME280 {
		b, err := bme280.NewI2C(i, nil)
		if err != nil {
			return err
		}
		defer b.Stop()
		go sensorLoop(b, env)
	}

	if useButton {
		p := gpio.ByNumber(24)
		if p == nil {
			return errors.New("no pin 24")
		}
		if err := p.In(gpio.PullUp, gpio.BothEdges); err != nil {
			return err
		}
		c := make(chan gpio.Level)
		go func() {
			p.WaitForEdge(-1)
			c <- p.Read()
		}()
		go buttonLoop(c, button)
	}

	/*
		// Relays
		p := gpio.ByNumber(17)
		if p == nil {
			return errors.New("no pin 17")
		}
		if err := .Out(); err != nil {
			return err
		}
		p.SetLow()
		p = gpio.ByNumber(27)
		if p == nil {
			return errors.New("no pin 27")
		}
		if err := p.Out(); err != nil {
			return err
		}
		p.SetLow()
	*/

	if usePir {
		p := gpio.ByNumber(19)
		if p == nil {
			return errors.New("no pin 19")
		}
		if err := p.In(gpio.PullDown, gpio.BothEdges); err != nil {
			return err
		}
		c := make(chan gpio.Level)
		go func() {
			p.WaitForEdge(-1)
			c <- p.Read()
		}()
		go pirLoop(c, motion)
	}

	if useIR {
		irBus, err := lirc.New()
		if err != nil {
			return err
		}
		go irLoop(irBus, keys)
	}

	interrupt.HandleCtrlC()
	<-interrupt.Channel

	return nil
}

func displayLoop(s *ssd1306.Dev, f *psf.Font, img *image1bit.Image, button, motion <-chan bool, env <-chan *devices.Environment, keys <-chan ir.Key) {
	tick := time.NewTicker(time.Second)
	defer tick.Stop()
	for {
		draw := false
		select {
		case b := <-button:
			if b {
				f.Draw(img, 0, f.H*4, image1bit.On, image1bit.Off, "Bouton!")
			} else {
				f.Draw(img, 0, f.H*4, image1bit.On, image1bit.Off, "          ")
			}
			draw = true

		case m := <-motion:
			if m {
				f.Draw(img, 0, f.H*5, image1bit.On, image1bit.Off, "Mouvement!")
			} else {
				f.Draw(img, 0, f.H*5, image1bit.On, image1bit.Off, "          ")
			}

		case t := <-env:
			f.Draw(img, 0, f.H, image1bit.On, image1bit.Off, fmt.Sprintf("%8s", t.Temperature))
			f.Draw(img, 0, f.H*2, image1bit.On, image1bit.Off, fmt.Sprintf("%9s", t.Pressure))
			f.Draw(img, 0, f.H*3, image1bit.On, image1bit.Off, fmt.Sprintf("%10s", t.Humidity))

		case s := <-keys:
			f.Draw(img, 0, f.H*6, image1bit.On, image1bit.Off, string(s))
			draw = true

		case <-tick.C:
			f.Draw(img, 0, 0, image1bit.On, image1bit.Off, time.Now().Format("15:04:05"))
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

func irLoop(irBus ir.Conn, keys chan<- ir.Key) {
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

func buttonLoop(b <-chan gpio.Level, c chan<- bool) {
	for {
		select {
		case l := <-b:
			log.Printf("Bouton: %s", l)
			c <- l == gpio.Low
		case <-interrupt.Channel:
			break
		}
	}
}

func pirLoop(b <-chan gpio.Level, c chan<- bool) {
	for {
		select {
		case l := <-b:
			log.Printf("PIR: %s", l)
			c <- l == gpio.High
		case <-interrupt.Channel:
			break
		}
	}
}

func sensorLoop(b *bme280.Dev, c chan<- *devices.Environment) {
	tick := time.NewTicker(500 * time.Millisecond)
	defer tick.Stop()
	for {
		env := &devices.Environment{}
		if err := b.Sense(env); err != nil {
			log.Fatal(err)
		}
		c <- env
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
