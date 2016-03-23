// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Packages the static files in a .go file.
//go:generate go run ../package/main.go -out static_files_gen.go web/static images

package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime/pprof"
	"time"

	"github.com/kardianos/osext"
	"github.com/maruel/dotstar"
	"github.com/maruel/interrupt"
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

func mainImpl() error {
	thisFile, err := osext.Executable()
	if err != nil {
		return err
	}

	cpuprofile := flag.String("cpuprofile", "", "dump CPU profile in file")
	port := flag.Int("port", 8010, "http port to listen on")
	verbose := flag.Bool("verbose", false, "enable log output")
	fake := flag.Bool("fake", false, "use a fake camera mock, useful to test without the hardware")
	cycle := flag.Bool("cycle", false, "enable cycling through a few animations")
	flag.Parse()
	if flag.NArg() != 0 {
		return errors.New("unknown arguments")
	}
	if len(flag.Args()) != 0 {
		return fmt.Errorf("unexpected argument: %s", flag.Args())
	}

	if !*verbose {
		log.SetOutput(ioutil.Discard)
	}

	interrupt.HandleCtrlC()
	defer interrupt.Set()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			return err
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	dotstar.Patterns["étoile floue"] = dotstar.LoadAnimate(mustRead("étoile floue.png"), 16*time.Millisecond, false)

	var s dotstar.Strip
	if *fake {
		s = dotstar.MakeScreen()
	} else {
		s, err = dotstar.MakeDotStar()
		if err != nil {
			return err
		}
	}
	p := dotstar.MakePainter(s, 80)

	StartWebServer(p, *port)

	if *cycle {
		go func() {
			patterns := []struct {
				d int
				p dotstar.Pattern
			}{
				{3, dotstar.Patterns["rainbow static"]},
				{10, dotstar.Patterns["glow rainbow"]},
				{10, dotstar.Patterns["étoile floue"]},
				{7, dotstar.Patterns["canne"]},
				{7, dotstar.Patterns["K2000"]},
				{7, dotstar.Patterns["comète"]},
				{5, dotstar.Patterns["pingpong"]},
				{5, dotstar.Patterns["glow"]},
				{5, dotstar.Patterns["glow gris"]},
				{3, dotstar.Patterns["red"]},
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

	defer fmt.Printf("\033[0m\n")
	return watchFile(thisFile)
}

func main() {
	if err := mainImpl(); err != nil {
		fmt.Fprintf(os.Stderr, "\ndotstar: %s.\n", err)
		os.Exit(1)
	}
}
