// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// ssd1306 writes to a display driven by a ssd1306 controler.
package main

import (
	"errors"
	"flag"
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	_ "image/png"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/maruel/dlibox/go/bw2d"
	"github.com/maruel/dlibox/go/pio/devices/ssd1306"
	"github.com/maruel/dlibox/go/pio/host"
	"github.com/maruel/dlibox/go/pio/host/drivers/sysfs"
	"github.com/maruel/dlibox/go/psf"
	"github.com/nfnt/resize"
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
			if p2 := filepath.Join(p, "src/github.com/maruel/dlibox/go/pio/cmd/ssd1306", name); access(p2) {
				return p2
			}
		}
	}
	return ""
}

// loadImg loads an image from disk.
func loadImg(name string) (image.Image, *gif.GIF, error) {
	p := findFile(name)
	if len(p) == 0 {
		return nil, nil, fmt.Errorf("couldn't find file %s", name)
	}
	f, err := os.Open(p)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	// Try to decode as an animated GIF first, then fall back to generic decoding.
	if g, err := gif.DecodeAll(f); err == nil {
		if len(g.Image) > 1 {
			log.Printf("Image %s as animated GIF", name)
			return nil, g, nil
		}
		log.Printf("Image %s", name)
		return g.Image[0], nil, nil
	}
	if _, err = f.Seek(0, 0); err != nil {
		return nil, nil, err
	}
	img, _, err := image.Decode(f)
	if err != nil {
		return nil, nil, err
	}
	log.Printf("Image %s", name)
	return img, nil, nil
}

func demo(s *ssd1306.Dev) error {
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

func convert(s *ssd1306.Dev, src image.Image, f *psf.Font, text string) (*bw2d.Image, error) {
	// Resize automatically while keeping aspect ratio.
	src = resize.Thumbnail(uint(s.W), uint(s.H), src, resize.Lanczos3)
	img, err := bw2d.New(s.W, s.H)
	if err != nil {
		return nil, err
	}
	r := src.Bounds()
	// Center the image.
	r = r.Add(image.Point{(s.W - r.Max.X) / 2, (s.H - r.Max.Y) / 2})
	draw.Draw(img, r, src, image.Point{}, draw.Src)
	// Use nil instead of bw2d.Off to not print the black pixels, or reverse the
	// two argument for inverted text.
	f.Draw(img, 0, s.H-f.H-1, bw2d.On, bw2d.Off, text)
	return img, nil
}

func mainImpl() error {
	i2cId := flag.Int("i2c", -1, "specify I²C bus to use")
	spiId := flag.Int("spi", -1, "specify SPI bus to use")
	csId := flag.Int("cs", 0, "specify SPI chip select (CS) to use")
	speed := flag.Int("speed", 4000000, "specify SPI speed in Hz to use")
	fontName := flag.String("f", "VGA8", "PSF font to use; use psf -l to list them")
	h := flag.Int("h", 64, "display height")
	imgName := flag.String("i", "ballerine.gif", "image to load; try bunny.gif")
	text := flag.String("t", "pio is awesome", "text to display")
	w := flag.Int("w", 128, "display width")
	demoMode := flag.Bool("d", false, "demo scrolling")
	rotated := flag.Bool("r", false, "Rotate the display by 180°")
	verbose := flag.Bool("v", false, "verbose mode")
	flag.Parse()
	if !*verbose {
		log.SetOutput(ioutil.Discard)
	}
	log.SetFlags(log.Lmicroseconds)
	if flag.NArg() != 0 {
		return errors.New("unexpected argument, try -help")
	}

	// Open the device on the right bus.
	var s *ssd1306.Dev
	if *spiId >= 0 {
		bus, err := sysfs.NewSPI(*spiId, *csId, int64(*speed))
		if err != nil {
			return err
		}
		defer bus.Close()
		log.Printf("Using pins CLK: %s  MOSI: %s  CS: %s", bus.CLK(), bus.MOSI(), bus.CS())
		s, err = ssd1306.NewSPI(bus, *w, *h, *rotated)
		if err != nil {
			return err
		}
	} else if *i2cId >= 0 {
		bus, err := sysfs.NewI2C(*i2cId)
		if err != nil {
			return err
		}
		defer bus.Close()
		log.Printf("Using pins SCL: %s  SDA: %s", bus.SCL(), bus.SDA())
		s, err = ssd1306.NewI2C(bus, *w, *h, *rotated)
		if err != nil {
			return err
		}
	} else {
		// Get the first I²C bus available.
		bus, err := host.NewI2C()
		if err != nil {
			return err
		}
		defer bus.Close()
		log.Printf("Using pins SCL: %s  SDA: %s", bus.SCL(), bus.SDA())
		s, err = ssd1306.NewI2C(bus, *w, *h, *rotated)
		if err != nil {
			return err
		}
	}

	// Load console font.
	f, err := psf.Load(*fontName)
	if err != nil {
		return err
	}
	log.Printf("Font: %dx%d", f.W, f.H)

	// Load image.
	src, g, err := loadImg(*imgName)
	if err != nil {
		return err
	}
	// If an animated GIF, draw it in a loop.
	if g != nil {
		// Resize all the images up front to save on CPU processing.
		imgs := make([]*bw2d.Image, len(g.Image))
		for i := range g.Image {
			imgs[i], err = convert(s, g.Image[i], f, *text)
			if err != nil {
				return err
			}
		}
		for i := 0; g.LoopCount <= 0 || i < g.LoopCount*len(g.Image); i++ {
			index := i % len(g.Image)
			c := time.After(time.Duration(10*g.Delay[index]) * time.Millisecond)
			if _, err := s.Write(imgs[index].Buf); err != nil {
				return err
			}
			<-c
		}
		return nil
	}

	if src == nil {
		// Create a blank image.
		if src, err = bw2d.New(s.W, s.H); err != nil {
			return err
		}
	}

	img, err := convert(s, src, f, *text)
	if err != nil {
		return err
	}
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
