// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// pins prints out the function of each pin.
package main

import (
	"flag"
	"fmt"
	"os"
	"sort"

	"github.com/maruel/dlibox/go/pio/host/cpu"
	"github.com/maruel/dlibox/go/pio/host/headers"
	"github.com/maruel/dlibox/go/pio/host/pins"
	"github.com/maruel/dlibox/go/pio/protocols/gpio"
)

func printFunc(invalid bool) {
	max := 0
	functional := pins.Functional()
	funcs := make([]string, 0, len(functional))
	for f := range functional {
		if l := len(f); l > 0 && f[0] != '<' {
			funcs = append(funcs, f)
			if l > max {
				max = l
			}
		}
	}
	sort.Strings(funcs)
	for _, name := range funcs {
		pin := functional[name]
		if invalid || pin != gpio.INVALID {
			if pin == nil {
				fmt.Printf("%-*s: INVALID\n", max, name)
			} else {
				fmt.Printf("%-*s: %s\n", max, name, pin)
			}
		}
	}
}

func printGPIO(invalid bool) {
	maxName := 0
	maxFn := 0
	all := pins.All()
	for _, p := range all {
		if invalid || headers.IsConnected(p) {
			if l := len(p.String()); l > maxName {
				maxName = l
			}
			if l := len(p.Function()); l > maxFn {
				maxFn = l
			}
		}
	}
	for _, p := range all {
		if headers.IsConnected(p) {
			fmt.Printf("%-*s: %s\n", maxName, p, p.Function())
		} else if invalid {
			fmt.Printf("%-*s: %-*s (not connected)\n", maxName, p, maxFn, p.Function())
		}
	}
}

func printHardware(invalid bool) {
	all := headers.All()
	names := make([]string, 0, len(all))
	for name := range all {
		names = append(names, name)
	}
	sort.Strings(names)
	maxName := 0
	maxFn := 0
	for _, header := range all {
		if len(header) == 0 || len(header[0]) != 2 {
			continue
		}
		for _, line := range header {
			for _, p := range line {
				if l := len(p.String()); l > maxName {
					maxName = l
				}
				if l := len(p.Function()); l > maxFn {
					maxFn = l
				}
			}
		}
	}
	for i, name := range names {
		if i != 0 {
			fmt.Print("\n")
		}
		header := all[name]
		if len(header) == 0 {
			fmt.Printf("%s: No pin connected\n", name)
			continue
		}
		sum := 0
		for _, line := range header {
			sum += len(line)
		}
		fmt.Printf("%s: %d pins\n", name, sum)
		if len(header[0]) == 2 {
			fmt.Printf("  %*s  %*s  Pos  Pos  %-*s Func\n", maxFn, "Func", maxName, "Name", maxName, "Name")
			for i, line := range header {
				fmt.Printf("  %*s  %*s  %3d  %-3d  %-*s %s\n", maxFn, line[0].Function(), maxName, line[0], 2*i+1, 2*i+2, maxName, line[1], line[1].Function())
			}
			continue
		}
		fmt.Printf("  Pos  %-*s  Func\n", maxName, "Name")
		pos := 1
		for _, line := range header {
			for _, item := range line {
				fmt.Printf("  %-3d  %-*s  %s\n", pos, maxName, item, item.Function())
				pos++
			}
		}
	}
}

func mainImpl() error {
	all := flag.Bool("a", false, "print everything")
	fun := flag.Bool("f", false, "print functional pins (e.g. I2C1_SCL)")
	gpio := flag.Bool("g", false, "print GPIO pins (e.g. GPIO1) (default)")
	hardware := flag.Bool("h", false, "print hardware pins (e.g. P1_1)")
	info := flag.Bool("i", false, "show general information")
	invalid := flag.Bool("n", false, "show not connected/INVALID pins")
	flag.Parse()
	if *all {
		*fun = true
		*gpio = true
		*hardware = true
		*info = true
		*invalid = true
	} else if !*fun && !*gpio && !*hardware && !*info {
		*gpio = true
	}

	// Explicitly initialize to catch any error.
	subsystem, err := pins.Init(true)
	if err != nil {
		return err
	}
	if *info {
		fmt.Printf("Subsystem: %s\nMaxSpeed: %dMhz\n", subsystem, cpu.MaxSpeed()/1000000)
	}
	if *fun {
		printFunc(*invalid)
	}
	if *gpio {
		printGPIO(*invalid)
	}
	if *hardware {
		printHardware(*invalid)
	}
	return nil
}

func main() {
	if err := mainImpl(); err != nil {
		fmt.Fprintf(os.Stderr, "pins: %s.\n", err)
		os.Exit(1)
	}
}
