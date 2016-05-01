// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Packages the static files in a .go file.
//go:generate go run ../package/main.go -out static_files_gen.go web/static images

// dotstar drives the dotstar LED strip on a Raspberry Pi. It runs a web server
// for remote control.
package main

import (
	"flag"
	"fmt"
	"image/color"
	"io/ioutil"
	"log"
	"os"
	"runtime/pprof"
	"sort"
	"time"

	"github.com/hashicorp/mdns"
	"github.com/kardianos/osext"
	"github.com/maruel/dotstar"
	"github.com/maruel/interrupt"
	"github.com/stianeikeland/go-rpio"
	"golang.org/x/exp/inotify"
)

func watchFile(fileName string) error {
	fi, err := os.Stat(fileName)
	if err != nil {
		return err
	}
	mod0 := fi.ModTime()
	watcher, err := inotify.NewWatcher()
	if err != nil {
		return err
	}
	if err = watcher.Watch(fileName); err != nil {
		return err
	}
	for {
		select {
		case <-interrupt.Channel:
			return err
		case err = <-watcher.Error:
			return err
		case <-watcher.Event:
			if fi, err = os.Stat(fileName); err != nil || !fi.ModTime().Equal(mod0) {
				return err
			}
		}
	}
}

func listenToPin(pinNumber int, p *dotstar.Painter, r *dotstar.PatternRegistry) {
	pin := rpio.Pin(pinNumber)
	pin.Input()
	pin.PullUp()
	last := rpio.High
	names := make([]string, 0, len(r.Patterns))
	for n := range r.Patterns {
		names = append(names, n)
	}
	sort.Strings(names)
	index := 0
	for {
		// Types of press:
		// - Short press (<2s)
		// - 2s press
		// - 4s press
		// - double-click (incompatible with repeated short press)
		//
		// Functions:
		// - Bonne nuit
		// - Next / Prev
		// - Éteindre (longer press après bonne nuit?)
		if state := pin.Read(); state != last {
			last = state
			if state == rpio.Low {
				index = (index + 1) % len(names)
				p.SetPattern(r.Patterns[names[index]])
			}
		}
		select {
		case <-interrupt.Channel:
			return
		case <-time.After(time.Millisecond):
		}
	}
}

type Weekday int

func (w Weekday) IsEnabledFor(d time.Weekday) bool {
	return false
}

const (
	Sunday Weekday = 1 << iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

type Alarm struct {
	Enabled bool
	Hour    int
	Minute  int
	Days    Weekday
	Pattern dotstar.Pattern
}

func (a *Alarm) Next(now time.Time) time.Duration {
	if a.Enabled && a.Days != 0 {
		// TODO(maruel): Figure out next day it will happen.
		//w := now.Weekday()
		//return time.Date(now.Year).Sub(now)
	}
	return 0
}

func (a *Alarm) SetupTimer(now time.Time, p *dotstar.Painter) *time.Timer {
	if next := a.Next(now); next != 0 {
		return time.AfterFunc(next, func() {
			p.SetPattern(a.Pattern)
			// Rearm on the next event.
		})
	}
	return nil
}

type Config struct {
	Alarms []Alarm
}

func (c *Config) alarmLoop(p *dotstar.Painter) {
	now := time.Now()
	var timers []*time.Timer
	for _, a := range c.Alarms {
		if t := a.SetupTimer(now, p); t != nil {
			timers = append(timers, t)
		}
	}
	for _, t := range timers {
		t.Stop()
	}
}

// TODO(maruel): Save this in a file and make it configurable via the web UI.
var config = Config{
	Alarms: []Alarm{
		{
			Enabled: true,
			Hour:    6,
			Minute:  30,
			Days:    Monday | Tuesday | Wednesday | Thursday | Friday,
			Pattern: &dotstar.EaseOut{
				In:       &dotstar.StaticColor{},
				Out:      &dotstar.Repeated{[]color.NRGBA{red, red, red, red, white, white, white, white}, 6},
				Duration: 20 * time.Minute,
			},
		},
	},
}

func mainImpl() error {
	thisFile, err := osext.Executable()
	if err != nil {
		return err
	}

	cpuprofile := flag.String("cpuprofile", "", "dump CPU profile in file")
	port := flag.Int("port", 8010, "http port to listen on")
	verbose := flag.Bool("verbose", false, "enable log output")
	fake := flag.Bool("fake", false, "use a fake camera mock, useful to test without the hardware")
	demoMode := flag.Bool("demo", false, "enable cycling through a few animations as a demo")
	pinNumber := flag.Int("pin", 0, "pin to listen to")
	numLights := flag.Int("n", 150, "number of lights to display. If lower than the actual number of lights, the remaining lights will flash oddly. When combined with -fake, number of characters to display on the line.")
	flag.Parse()
	if flag.NArg() != 0 {
		return fmt.Errorf("unexpected argument: %s", flag.Args())
	}
	if *demoMode && *pinNumber != 0 {
		return fmt.Errorf("use only one of -demo or -pin")
	}

	if !*verbose {
		log.SetOutput(ioutil.Discard)
	}

	interrupt.HandleCtrlC()
	defer interrupt.Set()

	if *cpuprofile != "" {
		// Run with cpuprofile, then use 'go tool pprof' to analyze it. See
		// http://blog.golang.org/profiling-go-programs for more details.
		f, err := os.Create(*cpuprofile)
		if err != nil {
			return err
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	var s dotstar.Strip
	if *fake {
		s = dotstar.MakeScreen()
	} else {
		s, err = dotstar.MakeDotStar()
		if err != nil {
			return err
		}
	}
	p := dotstar.MakePainter(s, *numLights)

	registry := getRegistry()
	startWebServer(*port, p, registry)

	hostName, err := os.Hostname()
	if err != nil {
		return err
	}
	// "_http._tcp."
	// http://www.dns-sd.org/servicetypes.html
	// http://www.iana.org/form/ports-services
	service, err := mdns.NewMDNSService(hostName, "dlibox", "", "", 3611, nil, []string{"dlibox"})
	if err != nil {
		return err
	}
	server, err := mdns.NewServer(&mdns.Config{Zone: service})
	if err != nil {
		return err
	}
	defer server.Shutdown()
	entries := make(chan *mdns.ServiceEntry)
	go func() {
		// TODO(maruel): When another device polls for services, immediately
		// register the device too.
		for _ = range entries {
			//fmt.Printf("%s\n", c)
		}
	}()
	if err = mdns.Lookup("dlibox", entries); err != nil {
		close(entries)
		return err
	}
	close(entries)

	if *demoMode {
		go func() {
			patterns := []struct {
				d int
				p dotstar.Pattern
			}{
				{3, registry.Patterns["Rainbow static"]},
				{10, registry.Patterns["Glow rainbow"]},
				{10, registry.Patterns["Étoile floue"]},
				{7, registry.Patterns["Canne de Noël"]},
				{7, registry.Patterns["K2000"]},
				{5, registry.Patterns["Ping pong"]},
				{5, registry.Patterns["Glow"]},
				{5, registry.Patterns["Glow gris"]},
			}
			i := 0
			p.SetPattern(patterns[i].p)
			delay := time.Duration(patterns[i].d) * time.Second
			for {
				select {
				case <-time.After(delay):
					i = (i + 1) % len(patterns)
					p.SetPattern(patterns[i].p)
					delay = time.Duration(patterns[i].d) * time.Second
				case <-interrupt.Channel:
					return
				}
			}
		}()
	}

	if *pinNumber != 0 {
		// Open and map memory to access gpio, check for errors
		if err := rpio.Open(); err != nil {
			return err
		}
		defer rpio.Close()
		go listenToPin(*pinNumber, p, registry)
	}

	defer fmt.Printf("\033[0m\n")
	return watchFile(thisFile)
}

func main() {
	if err := mainImpl(); err != nil {
		fmt.Fprintf(os.Stderr, "\ndotstar: %s.\n", err)
		os.Exit(1)
	}
}
