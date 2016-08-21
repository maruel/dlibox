// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// ssd1306 is a small app to write to a display driven by a ssd1306 controler.
package main

import (
	"flag"
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"os"
	"time"

	"image/draw"
	_ "image/gif"
	_ "image/png"

	"github.com/maruel/dlibox/go/bw2d"
	"github.com/maruel/dlibox/go/rpi"
	"github.com/maruel/dlibox/go/ssd1306"
)

// loadImg loads an image as black and white.
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
	img := bw2d.Make(r.Max.X, r.Max.Y)
	draw.Draw(img, r, src, image.Point{}, draw.Src)
	return img, nil
}

func demo(s *ssd1306.SSD1306) error {
	if err := s.Scroll(ssd1306.Left, ssd1306.FrameRate2); err != nil {
		return err
	}
	time.Sleep(2 * time.Second)
	if err := s.Scroll(ssd1306.Right, ssd1306.FrameRate2); err != nil {
		return err
	}
	time.Sleep(2 * time.Second)
	if err := s.Scroll(ssd1306.UpLeft, ssd1306.FrameRate2); err != nil {
		return err
	}
	time.Sleep(2 * time.Second)
	if err := s.Scroll(ssd1306.UpRight, ssd1306.FrameRate2); err != nil {
		return err
	}
	time.Sleep(2 * time.Second)
	if err := s.StopScroll(); err != nil {
		return err
	}
	//if _, err := s.Write(img.Buf); err != nil {
	//	return err
	//}
	if err := s.SetContrast(0); err != nil {
		return err
	}
	time.Sleep(2 * time.Second)
	if err := s.SetContrast(0xFF); err != nil {
		return err
	}
	time.Sleep(2 * time.Second)
	return nil
}

func mainImpl() error {
	bus := flag.Int("b", 1, "I²C bus to use")
	demoMode := flag.Bool("d", false, "demo scrolling")
	rotated := flag.Bool("r", false, "Rotate the display by 180°")
	verbose := flag.Bool("v", false, "verbose mode")
	flag.Parse()
	if !*verbose {
		log.SetOutput(ioutil.Discard)
	}
	log.SetFlags(log.Lmicroseconds)

	// Open the device
	i2c, err := rpi.MakeI2C(*bus)
	if err != nil {
		return err
	}
	s, err := ssd1306.MakeSSD1306(i2c, 128, 64, *rotated)
	if err != nil {
		return err
	}

	f, err := bw2d.BasicFont(8)
	if err != nil {
		return err
	}
	log.Printf("Font: %dx%d", f.W, f.H)

	src, err := loadImg("bunny.png")
	if err != nil {
		return err
	}
	src.Inverse()
	img := bw2d.Make(128, 64)
	r := src.Bounds()
	r = r.Add(image.Point{(img.W - r.Max.X) / 2, (img.H - r.Max.Y) / 2})
	draw.Draw(img, r, src, image.Point{}, draw.Src)
	img.Text(0, 0, f, "dlibox!")
	img.Text(0, s.H-f.H-1, f, "is awesome")
	if _, err = s.Write(img.Buf); err != nil {
		return err
	}

	if *demoMode {
		if err := demo(s); err != nil {
			return err
		}
	}
	if err := s.Enable(false); err != nil {
		return err
	}
	return err
}

func main() {
	if err := mainImpl(); err != nil {
		fmt.Fprintf(os.Stderr, "ssd1306: %s.\n", err)
		os.Exit(1)
	}
}
