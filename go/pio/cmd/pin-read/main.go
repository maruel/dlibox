// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// pin-read is a small app to read a pin.
package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/maruel/dlibox/go/pio/buses"
	"github.com/maruel/dlibox/go/pio/buses/bcm283x"
)

func read(p buses.Pin, edge buses.Edge) {
	if p.ReadEdge() == buses.Low {
		os.Stdout.Write([]byte{'0', '\n'})
	} else {
		os.Stdout.Write([]byte{'1', '\n'})
	}
}

func mainImpl() error {
	pullUp := flag.Bool("u", false, "pull up")
	pullDown := flag.Bool("d", false, "pull down")
	edgeRising := flag.Bool("r", false, "wait for rising edge; can be used along -f")
	edgeFalling := flag.Bool("f", false, "wait for falling edge; can be used along -r")
	loop := flag.Bool("l", false, "loop")
	verbose := flag.Bool("v", false, "enable verbose logs")
	flag.Parse()

	if !*verbose {
		log.SetOutput(ioutil.Discard)
	}
	log.SetFlags(log.Lmicroseconds)

	//pull := buses.PullNoChange
	pull := buses.Float
	if *pullUp {
		if *pullDown {
			return errors.New("use only one of -d or -u")
		}
		pull = buses.Up
	}
	if *pullDown {
		pull = buses.Down
	}
	if flag.NArg() != 1 {
		return errors.New("specify pin to read")
	}

	edge := buses.EdgeNone
	if *edgeRising {
		edge |= buses.Rising
	}
	if *edgeFalling {
		edge |= buses.Falling
	}

	pin, err := strconv.Atoi(flag.Args()[0])
	if err != nil {
		return err
	}
	if pin > 53 || pin < 0 {
		return errors.New("specify pin between 0 and 53")
	}
	p := bcm283x.Pin(pin)

	if err = p.In(pull, edge); err != nil {
		return err
	}
	read(p, edge)
	for *loop {
		read(p, edge)
	}
	return nil
}

func main() {
	if err := mainImpl(); err != nil {
		fmt.Fprintf(os.Stderr, "pin-read: %s.\n", err)
		os.Exit(1)
	}
}
