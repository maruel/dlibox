// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// pin-read reads a pin.
package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/maruel/dlibox/go/pio/host/pins"
	"github.com/maruel/dlibox/go/pio/protocols/gpio"
)

func printLevel(l gpio.Level) error {
	if l == gpio.Low {
		_, err := os.Stdout.Write([]byte{'0', '\n'})
		return err
	}
	_, err := os.Stdout.Write([]byte{'1', '\n'})
	return err
}

func mainImpl() error {
	pullUp := flag.Bool("u", false, "pull up")
	pullDown := flag.Bool("d", false, "pull down")
	edges := flag.Bool("e", false, "wait for edges")
	verbose := flag.Bool("v", false, "enable verbose logs")
	flag.Parse()

	if !*verbose {
		log.SetOutput(ioutil.Discard)
	}
	log.SetFlags(log.Lmicroseconds)

	//pull := gpio.PullNoChange
	pull := gpio.Float
	if *pullUp {
		if *pullDown {
			return errors.New("use only one of -d or -u")
		}
		pull = gpio.Up
	}
	if *pullDown {
		pull = gpio.Down
	}
	if flag.NArg() != 1 {
		return errors.New("specify pin to read")
	}

	pin, err := strconv.Atoi(flag.Args()[0])
	if err != nil {
		return err
	}
	p := pins.ByNumber(pin)
	if p == nil {
		return errors.New("specify a valid pin number")
	}
	if err = p.In(pull); err != nil {
		return err
	}
	if *edges {
		c, err := p.Edges()
		defer p.DisableEdges()
		if err != nil {
			return err
		}
		for {
			if err = printLevel(<-c); err != nil {
				return err
			}
		}
	} else {
		return printLevel(p.Read())
	}
	return nil
}

func main() {
	if err := mainImpl(); err != nil {
		fmt.Fprintf(os.Stderr, "pin-read: %s.\n", err)
		os.Exit(1)
	}
}
