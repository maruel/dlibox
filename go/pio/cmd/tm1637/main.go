// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// tm1637 is a small app to write to a digits LED display.
package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/maruel/dlibox/go/pio/devices/tm1637"
	"github.com/maruel/dlibox/go/pio/host"
	// TODO(maruel): Make this unneeded.
	_ "github.com/maruel/dlibox/go/pio/host/bcm283x"
)

func mainImpl() error {
	clk := flag.Int("c", 4, "CLK pin number")
	data := flag.Int("d", 5, "DIO pin number")
	off := flag.Bool("o", false, "set display as off")
	b1 := flag.Bool("b1", false, "set PWM to 1/16")
	b2 := flag.Bool("b2", false, "set PWM to 2/16")
	b4 := flag.Bool("b4", false, "set PWM to 4/16")
	b10 := flag.Bool("b10", false, "set PWM to 10/16 (default)")
	b12 := flag.Bool("b12", false, "set PWM to 12/16")
	b13 := flag.Bool("b13", false, "set PWM to 13/16")
	b14 := flag.Bool("b14", false, "set PWM to 14/16")
	verbose := flag.Bool("v", false, "verbose mode")
	asSeg := flag.Bool("s", false, "use hex encoded segments instead of numbers")
	asTime := flag.Bool("t", false, "expect two numbers representing time")
	showDot := flag.Bool("dot", false, "when -t is used, show dots")
	flag.Parse()
	if !*verbose {
		log.SetOutput(ioutil.Discard)
	}
	log.SetFlags(log.Lmicroseconds)

	b := tm1637.Brightness10
	switch {
	case *off:
		b = tm1637.Off
	case *b1:
		b = tm1637.Brightness1
	case *b2:
		b = tm1637.Brightness2
	case *b4:
		b = tm1637.Brightness4
	case *b10:
		b = tm1637.Brightness10
	case *b12:
		b = tm1637.Brightness12
	case *b13:
		b = tm1637.Brightness13
	case *b14:
		b = tm1637.Brightness14
	}
	if flag.NArg() > 6 {
		return errors.New("too many digits")
	}
	if b != tm1637.Off && flag.NArg() == 0 {
		// Turn it off
		b = tm1637.Off
	}
	var hour, minute int
	var digits []int
	var segments []byte
	if *asTime {
		if flag.NArg() != 2 {
			return errors.New("provide hh and mm")
		}
		x, err := strconv.ParseUint(flag.Arg(0), 10, 8)
		if err != nil {
			return err
		}
		hour = int(x)
		x, err = strconv.ParseUint(flag.Arg(1), 10, 8)
		if err != nil {
			return err
		}
		minute = int(x)
	} else if *asSeg {
		segments = make([]byte, flag.NArg())
		for i, d := range flag.Args() {
			x, err := strconv.ParseUint(d, 16, 8)
			if err != nil {
				return err
			}
			segments[i] = byte(x)
		}
	} else {
		digits = make([]int, flag.NArg())
		for i, d := range flag.Args() {
			x, err := strconv.ParseUint(d, 16, 8)
			if err != nil {
				return err
			}
			digits[i] = int(x)
		}
	}

	pClk := host.GetPinByNumber(*clk)
	if pClk == nil {
		return errors.New("specify a valid pin for clock")
	}
	pData := host.GetPinByNumber(*data)
	if pData == nil {
		return errors.New("specify a valid pin for data")
	}
	d, err := tm1637.Make(pClk, pData)
	if err != nil {
		return err
	}
	if err = d.SetBrightness(b); err != nil {
		return err
	}
	if len(segments) != 0 {
		return d.Segments(segments...)
	} else if len(digits) != 0 {
		return d.Digits(digits...)
	} else if *asTime {
		return d.Clock(hour, minute, *showDot)
	}
	return nil
}

func main() {
	if err := mainImpl(); err != nil {
		fmt.Fprintf(os.Stderr, "tm1637: %s.\n", err)
		os.Exit(1)
	}
}
