// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Packages the static files in a .go file.
//go:generate go run ../package/main.go -out static_files_gen.go ../../../web

// dlibox drives the dlibox LED strip on a Raspberry Pi. It runs a web server
// for remote control.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"syscall"

	"github.com/kardianos/osext"
	"github.com/maruel/dlibox/go/anim1d"
	"github.com/maruel/dlibox/go/bw2d"
	"github.com/maruel/dlibox/go/pio/devices"
	"github.com/maruel/dlibox/go/pio/devices/apa102"
	"github.com/maruel/dlibox/go/pio/devices/ssd1306"
	"github.com/maruel/dlibox/go/pio/fakes"
	"github.com/maruel/dlibox/go/pio/fakes/screen"
	"github.com/maruel/dlibox/go/pio/host"
	"github.com/maruel/dlibox/go/pio/host/cpu"
	"github.com/maruel/dlibox/go/pio/host/ir"
	"github.com/maruel/dlibox/go/pio/host/sysfs/i2c"
	"github.com/maruel/dlibox/go/pio/host/sysfs/spi"
	"github.com/maruel/dlibox/go/psf"
	"github.com/maruel/interrupt"
)

func initDisplay() (devices.Display, error) {
	i2cBus, err := i2c.Make(1)
	if err != nil {
		return nil, err
	}
	display, err := ssd1306.MakeI2C(i2cBus, 128, 64, false)
	if err != nil {
		return nil, err
	}
	f12, err := psf.Load("Terminus12x6")
	if err != nil {
		return nil, err
	}
	f20, err := psf.Load("Terminus20x10")
	if err != nil {
		return nil, err
	}
	// TODO(maruel): Leverage bme280 while at it but don't fail if not
	// connected.
	img, err := bw2d.Make(display.W, display.H)
	if err != nil {
		return nil, err
	}
	f20.Draw(img, 0, 0, bw2d.On, nil, "dlibox!")
	f12.Draw(img, 0, display.H-f12.H-1, bw2d.On, nil, "is awesome")
	if _, err = display.Write(img.Buf); err != nil {
		return nil, err
	}
	return display, nil
}

func initIR(painter *anim1d.Painter, config *IR) error {
	bus, err := ir.Make()
	if err != nil {
		return err
	}
	go func() {
		c := bus.Channel()
		for {
			select {
			case msg, ok := <-c:
				if !ok {
					break
				}
				if !msg.Repeat {
					// TODO(maruel): Locking.
					if pat := config.Mapping[msg.Key]; len(pat) != 0 {
						painter.SetPattern(string(pat))
					}
				}
			}
		}
	}()
	return nil
}

func initPIR(painter *anim1d.Painter, config *PIR) error {
	p := host.GetPinByNumber(config.Pin)
	if p == nil {
		return nil
	}
	if err := p.In(host.Down); err != nil {
		return err
	}
	c, err := p.Edges()
	if err != nil {
		return err
	}
	go func() {
		for {
			if l := <-c; l == host.High {
				// TODO(maruel): Locking.
				painter.SetPattern(string(config.Pattern))
			}
		}
	}()
	return nil
}

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
	chanSignal := make(chan os.Signal)
	go func() {
		<-chanSignal
		interrupt.Set()
	}()
	signal.Notify(chanSignal, syscall.SIGTERM)

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
		log.Printf("Loading config failed: %v", err)
	}
	defer config.Close()

	b, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	log.Printf("Config:\n%s", string(b))

	fps := 60
	if cpu.MaxSpeed < 900000 || runtime.NumCPU() < 4 {
		// Use 30Hz on slower devices because it is too slow.
		fps = 30
	}

	// Output (terminal with ANSI codes or APA102).
	var leds devices.Display
	if *fake {
		// Hardcode to 100 characters when using a terminal output.
		// TODO(maruel): Query the terminal and use its width.
		leds = screen.Make(100)
		defer os.Stdout.Write([]byte("\033[0m\n"))
		// Use lower refresh rate too.
		fps = 30
		properties = append(properties, "fake=1")
	} else {
		spiBus, err := spi.Make(0, 0, config.Settings.APA102.SPIspeed)
		if err != nil {
			log.Printf("SPI failed: %v", err)
			leds = &fakes.Display{image.NewNRGBA(image.Rect(0, 0, config.Settings.APA102.NumberLights, 1))}
		} else {
			defer spiBus.Close()
			if leds, err = apa102.Make(spiBus, config.Settings.APA102.NumberLights, 255, 6500); err != nil {
				return err
			}
			properties = append(properties, fmt.Sprintf("APA102=%d", config.Settings.APA102.NumberLights))
		}
	}

	// Try to initialize the display.
	if _, err = initDisplay(); err != nil {
		log.Printf("Display not connected")
	}

	// Painter.
	p := anim1d.MakePainter(leds, fps)
	if err := config.Init(p); err != nil {
		return err
	}
	startWebServer(*port, p, &config.Config)

	if err = initIR(p, &config.Settings.IR); err != nil {
		log.Printf("IR not connected: %v", err)
	}

	if err = initPIR(p, &config.Settings.PIR); err != nil {
		log.Printf("PIR not connected: %v", err)
	}

	/*
		service, err := initmDNS(*port, properties)
		if err != nil {
			return err
		}
		defer service.Close()
	*/

	return watchFile(thisFile)
}

func main() {
	if err := mainImpl(); err != nil {
		fmt.Fprintf(os.Stderr, "\ndlibox: %s.\n", err)
		os.Exit(1)
	}
}
