// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Packages the static files in a .go file.
//go:generate go run ../package/main.go -out static_files_gen.go images web

// dotstar drives the dotstar LED strip on a Raspberry Pi. It runs a web server
// for remote control.
package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"runtime/pprof"

	"github.com/kardianos/osext"
	"github.com/maruel/dotstar/anim1d"
	"github.com/maruel/dotstar/apa102"
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

type Config struct {
	// Alarm clock.
	Alarms
	// Number of lights to display. If lower than the actual number of lights,
	// the remaining lights will flash oddly.
	NumberLights int
}

func (c *Config) Load(n string) {
	if f, err := os.Open(n); err == nil {
		defer f.Close()
		_ = json.NewDecoder(f).Decode(c)
	}
}

func (c *Config) Save(n string) {
	if b, err := json.MarshalIndent(c, "", "  "); err == nil {
		if f, err := os.Create(n); err == nil {
			defer f.Close()
			_, _ = f.Write(append(b, '\n'))
		}
	}
}

// TODO(maruel): Make it configurable via the web UI.
var config = Config{
	Alarms: Alarms{
		{
			Enabled: true,
			Hour:    6,
			Minute:  30,
			Days:    Monday | Tuesday | Wednesday | Thursday | Friday,
			Pattern: "#FFFFFFFF",
			/*
				Pattern: anim1d.Marshal(&anim1d.EaseOut{
					After:       &anim1d.Color{},
					Before:      &anim1d.Repeated{[]color.NRGBA{red, red, red, red, white, white, white, white}, 6},
					Duration: 20 * time.Minute,
				}),
			*/
		},
		//"{\"Duration\":600000000000,\"After\":\"#00000000\",\"Offset\":1800000000000,\"Before\":{\"Duration\":600000000000,\"After\":\"#ffffffff\",\"Offset\":600000000000,\"Before\":{\"Duration\":600000000000,\"After\":\"#ff7f00ff\",\"Offset\":0,\"Before\":\"#00000000\",\"Transition\":\"linear\",\"_type\":\"Transition\"},\"Transition\":\"linear\",\"_type\":\"Transition\"},\"Transition\":\"linear\",\"_type\":\"Transition\"}",
	},
	NumberLights: 150,
}

// getHome returns the home directory even when cross compiled.
//
// When cross compiling, user.Current() fails.
func getHome() (string, error) {
	if u, err := user.Current(); err == nil {
		return u.HomeDir, nil
	}
	if s := os.Getenv("HOME"); len(s) != 0 {
		return s, nil
	}
	return "", errors.New("can't find HOME")
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
	numLights := flag.Int("n", 0, "number of lights to display. If lower than the actual number of lights, the remaining lights will flash oddly. When combined with -fake, number of characters to display on the line.")
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
		properties = append(properties, "profiled")
	}

	var s anim1d.Strip
	if *fake {
		s = apa102.MakeScreen()
		properties = append(properties, "fake")
	} else {
		s, err = apa102.MakeDotStar()
		if err != nil {
			return err
		}
		properties = append(properties, "APA102")
	}
	if *numLights == 0 {
		*numLights = config.NumberLights
	}
	p := anim1d.MakePainter(s, *numLights)

	home, err := getHome()
	if err != nil {
		return err
	}
	configPath := filepath.Join(home, "dotstar.json")
	config.Load(configPath)
	defer config.Save(configPath)
	config.Alarms.Reset(p)
	startWebServer(*port, p)

	/* TODO(maruel): Make this work.
	service, err := initmDNS(properties)
	if err != nil {
		return err
	}
	defer service.Close()
	*/

	defer fmt.Printf("\033[0m\n")
	return watchFile(thisFile)
}

func main() {
	if err := mainImpl(); err != nil {
		fmt.Fprintf(os.Stderr, "\ndotstar: %s.\n", err)
		os.Exit(1)
	}
}
