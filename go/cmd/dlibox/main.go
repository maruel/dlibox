// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Packages the static files in a .go file.
//go:generate go run ../package/main.go -out static_files_gen.go images ../../../web

// dlibox drives the dlibox LED strip on a Raspberry Pi. It runs a web server
// for remote control.
package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"runtime/pprof"

	"github.com/kardianos/osext"
	"github.com/maruel/dlibox/go/anim1d"
	"github.com/maruel/dlibox/go/apa102"
	"github.com/maruel/dlibox/go/bw2d"
	"github.com/maruel/dlibox/go/psf"
	"github.com/maruel/dlibox/go/rpi"
	"github.com/maruel/dlibox/go/ssd1306"
	"github.com/maruel/interrupt"
)

func mainImpl() error {
	thisFile, err := osext.Executable()
	if err != nil {
		return err
	}

	cpuprofile := flag.String("cpuprofile", "", "dump CPU profile in file")
	port := flag.Int("port", 8010, "http port to listen on")
	verbose := flag.Bool("verbose", false, "enable log output")
	fake := flag.Bool("fake", false, "use a terminal mock, useful to test without the hardware")
	flag.Parse()
	if flag.NArg() != 0 {
		return fmt.Errorf("unexpected argument: %s", flag.Args())
	}

	if !*verbose {
		log.SetOutput(ioutil.Discard)
	}

	interrupt.HandleCtrlC()
	defer interrupt.Set()

	var properties []string
	if *cpuprofile != "" {
		// Run with cpuprofile, then use 'go tool pprof' to analyze it. See
		// http://blog.golang.org/profiling-go-programs for more details.
		f, err := os.Create(*cpuprofile)
		if err != nil {
			return err
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
		properties = append(properties, "profiled=1")
	}

	// Config.
	config := ConfigMgr{}
	config.ResetDefault()
	if err := config.Load(); err != nil {
		return err
	}
	defer config.Close()

	b, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	log.Printf("Config:\n%s", string(b))

	fps := 60
	if rpi.MaxSpeed < 900000 {
		// Use 30Hz on slower devices because it is too slow.
		fps = 30
	}

	// Output (screen or APA102).
	var s io.Writer
	if *fake {
		s = apa102.MakeScreen()
		// Use lower refresh rate too.
		fps = 30
		properties = append(properties, "fake=1")
	} else {
		// Verify the pinout is as expected.
		if rpi.IR_OUT != rpi.GPIO5 || rpi.IR_IN != rpi.GPIO13 {
			return errors.New("configure lirc for out=5, in=13")
		}
		spi, err := rpi.MakeSPI(0, 0, config.APA102.SPIspeed)
		if err != nil {
			return err
		}
		s = apa102.MakeAPA102(spi)

		i2c, err := rpi.MakeI2C(1)
		s, err := ssd1306.MakeSSD1306(i2c, 128, 64, false)
		if err != nil {
			return err
		}
		f12, err := psf.Load("Terminus12x6")
		if err != nil {
			return err
		}
		f20, err := psf.Load("Terminus20x10")
		if err != nil {
			return err
		}
		// TODO(maruel): Leverage bme280 while at it but don't fail if not
		// connected.
		img := bw2d.Make(s.W, s.H)
		f20.Draw(img, 0, 0, bw2d.On, nil, "dlibox!")
		f12.Draw(img, 0, s.H-f12.H-1, bw2d.On, nil, "is awesome")
		if _, err = s.Write(img.Buf); err != nil {
			return err
		}

		properties = append(properties, fmt.Sprintf("APA102=%d", config.APA102.NumberLights))
	}

	// Painter.
	numLights := config.APA102.NumberLights
	if *fake {
		// Hardcode to 100 characters when using a terminal output.
		numLights = 100
	}
	p := anim1d.MakePainter(s, numLights, fps)
	if err := config.Init(p); err != nil {
		return err
	}
	startWebServer(*port, p, &config.Config)

	service, err := initmDNS(*port, properties)
	if err != nil {
		return err
	}
	defer service.Close()

	defer fmt.Printf("\033[0m\n")
	return watchFile(thisFile)
}

func main() {
	if err := mainImpl(); err != nil {
		fmt.Fprintf(os.Stderr, "\ndlibox: %s.\n", err)
		os.Exit(1)
	}
}
