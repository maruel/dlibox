// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// led reads the state of a LED or change it.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/maruel/dlibox/go/pio/host"
	"github.com/maruel/dlibox/go/pio/host/sysfs"
)

func mainImpl() error {
	verbose := flag.Bool("v", false, "enable verbose logs")
	flag.Parse()

	if !*verbose {
		log.SetOutput(ioutil.Discard)
	}
	log.SetFlags(0)

	state, err := host.Init()
	if err != nil {
		return err
	}

	if len(state.Failed) != 0 {
		fmt.Printf("Got the following errors:\n")
		for _, f := range state.Failed {
			fmt.Printf("  - %s: %v\n", f.D, f.Err)
		}
	}
	log.Printf("Using drivers:")
	for _, d := range state.Loaded {
		log.Printf("  - %s", d)
	}
	log.Printf("Skipped drivers:")
	for _, d := range state.Skipped {
		log.Printf("  - %s", d)
	}
	for _, led := range sysfs.LEDs {
		fmt.Printf("%s: %s\n", led, led.Function())
	}
	return nil
}

func main() {
	if err := mainImpl(); err != nil {
		fmt.Fprintf(os.Stderr, "led: %s.\n", err)
		os.Exit(1)
	}
}
