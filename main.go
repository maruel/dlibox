// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Packages the static files in a .go file.
//go:generate go run package/main.go -out static_files_gen.go web/static images

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

	"golang.org/x/exp/inotify"

	"github.com/maruel/interrupt"
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
	cpuprofile := flag.String("cpuprofile", "", "dump CPU profile in file")
	port := flag.Int("port", 8010, "http port to listen on")
	verbose := flag.Bool("verbose", false, "enable log output")
	fake := flag.Bool("fake", false, "use a fake camera mock, useful to test without the hardware")
	colorTest := flag.Bool("color-test", false, "prints all term-256 colors and exit")
	flag.Parse()
	if flag.NArg() != 0 {
		return errors.New("unknown arguments")
	}
	if !*verbose {
		log.SetOutput(ioutil.Discard)
	}

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			return err
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	interrupt.HandleCtrlC()
	defer interrupt.Set()

	if len(flag.Args()) != 0 {
		return fmt.Errorf("unexpected argument: %s", flag.Args())
	}

	if *colorTest {
		for y := 0; y < 32; y++ {
			extra := "  "
			for x := 0; x < 8; x++ {
				i := x + 8*y
				if x == 7 {
					extra = ""
				}
				fmt.Printf("\033[48;5;%dm   %3d   \033[0m%s", i, i, extra)
			}
			fmt.Printf("\n")
		}
		return nil
	}

	ws := StartWebServer(*port)
	defer ws.Close()

	var err error
	var s Strip
	if *fake {
		s = MakeScreen()
	} else {
		s, err = MakeDotStar()
		if err != nil {
			return err
		}
	}
	p := MakePainter(s, 80)

	go func() {
		patterns := []struct {
			d int
			p Pattern
		}{
			{3, Patterns["rainbow static"]},
			{10, Patterns["glow rainbow"]},
			{10, Patterns["étoile floue"]},
			{7, Patterns["canne"]},
			{7, Patterns["K2000"]},
			{7, Patterns["comète"]},
			{5, Patterns["pingpong"]},
			{5, Patterns["glow"]},
			{5, Patterns["glow gris"]},
			{3, Patterns["red"]},
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
	defer fmt.Printf("\033[0m\n")
	return watchFile(os.Args[0])
}

func main() {
	if err := mainImpl(); err != nil {
		fmt.Fprintf(os.Stderr, "\ndotstar: %s.\n", err)
		os.Exit(1)
	}
}
