// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// apa102 is a small app to write to a strip of APA102 LED.
package main

import (
	"errors"
	"flag"
	"fmt"
	"image"
	_ "image/png"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/maruel/dlibox/go/anim1d"
	"github.com/maruel/dlibox/go/pio/buses/spi"
	"github.com/maruel/dlibox/go/pio/devices/apa102"
)

func access(name string) bool {
	_, err := os.Stat(name)
	return err == nil
}

func findFile(name string) string {
	if access(name) {
		return name
	}
	for _, p := range strings.Split(os.Getenv("GOPATH"), ":") {
		if len(p) != 0 {
			if p2 := filepath.Join(p, "src/github.com/maruel/dlibox/go/pio/cmd/apa102", name); access(p2) {
				return p2
			}
		}
	}
	return ""
}

// loadImg loads an image from disk.
func loadImg(name string) (image.Image, error) {
	p := findFile(name)
	if len(p) == 0 {
		return nil, fmt.Errorf("couldn't find file %s", name)
	}
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	log.Printf("Image %s", name)
	return img, nil
}

func mainImpl() error {
	bus := flag.Int("b", 0, "SPI bus to use")
	numLights := flag.Int("n", 150, "number of lights on the strip")
	intensity := flag.Int("l", 127, "light intensity [1-255]")
	temperature := flag.Int("t", 5000, "light temperature in °Kelvin [3500-7500]")
	imgName := flag.String("i", "", "image to load")
	speed := flag.Int("s", 4000000, "speed in Hz")
	verbose := flag.Bool("v", false, "verbose mode")
	flag.Parse()
	if !*verbose {
		log.SetOutput(ioutil.Discard)
	}
	log.SetFlags(log.Lmicroseconds)
	if flag.NArg() != 0 {
		return errors.New("unexpected argument, try -help")
	}
	if *intensity > 255 {
		return errors.New("max intensity is 255")
	}
	if *temperature > 65535 {
		return errors.New("max temperature is 65535")
	}

	// Open the device
	s, err := spi.Make(*bus, 0, int64(*speed))
	if err != nil {
		return err
	}
	a, err := apa102.Make(s, uint8(*intensity), uint16(*temperature))
	if err != nil {
		return err
	}

	if len(*imgName) != 0 {
		// Load an image and make it loop through the pixels.
		/*
			src, err := loadImg(*imgName)
			if err != nil {
				return err
			}
			if _, err := a.Write(nil); err != nil {
				return err
			}
		*/
	}

	// Draw a rainbow.
	r := anim1d.Rainbow{}
	pixels := make(anim1d.Frame, *numLights)
	r.NextFrame(pixels, 0)
	buf := make([]byte, *numLights*3)
	pixels.ToRGB(buf)
	_, err = a.Write(buf)
	return err
}

func main() {
	if err := mainImpl(); err != nil {
		fmt.Fprintf(os.Stderr, "apa102: %s.\n", err)
		os.Exit(1)
	}
}
