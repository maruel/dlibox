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
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"

	"github.com/kardianos/osext"
	"github.com/maruel/dlibox/go/donotuse/host"
	"github.com/maruel/dlibox/go/modules"
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

	// Initialize pio.
	if _, err := host.Init(); err != nil {
		return err
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

	// Initialize modules.

	bus := &modules.LocalBus{}
	local := modules.Rebase(bus, "dlibox/")
	_, err = initDisplay(local, &config.Settings.Display)
	if err != nil {
		// Non-fatal.
		log.Printf("Display not connected: %v", err)
	}

	leds, end, properties2, fps, err := initLEDs(*fake, &config.Settings.APA102)
	if err != nil {
		return err
	}
	defer end()
	properties = append(properties, properties2...)

	p, err := initPainter(local, leds, fps, &config.Settings.Painter)
	if err != nil {
		return err
	}
	defer p.Close()
	if err := config.Init(p); err != nil {
		return err
	}
	startWebServer(*port, p, &config.Config)

	if err = initButton(p, nil, &config.Settings.Button); err != nil {
		// Non-fatal.
		log.Printf("Button not connected: %v", err)
	}

	if err = initIR(p, &config.Settings.IR); err != nil {
		// Non-fatal.
		log.Printf("IR not connected: %v", err)
	}

	if err = initPIR(p, &config.Settings.PIR); err != nil {
		// Non-fatal.
		log.Printf("PIR not connected: %v", err)
	}

	//service, err := initmDNS(*port, properties)
	//if err != nil {
	//	return err
	//}
	//defer service.Close()

	return watchFile(thisFile)
}

func main() {
	if err := mainImpl(); err != nil {
		fmt.Fprintf(os.Stderr, "\ndlibox: %s.\n", err)
		os.Exit(1)
	}
}
