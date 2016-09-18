// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// pins is a small app to read the function of each pin.
package main

import (
	"flag"
	"fmt"
	"os"
	"sort"

	"github.com/maruel/dlibox/go/pio/host"
	"github.com/maruel/dlibox/go/pio/host/cpu"
	"github.com/maruel/dlibox/go/pio/host/headers"
	"github.com/maruel/dlibox/go/pio/host/pins"
)

func getMaxName() int {
	max := 0
	for _, p := range pins.All {
		if l := len(p.String()); l > max {
			max = l
		}
	}
	return max
}

func getMaxFn() int {
	max := 0
	for _, p := range pins.All {
		if l := len(p.Function()); l > max {
			max = l
		}
	}
	return max
}

func printFunc(invalid bool) {
	max := 0
	funcs := make([]string, 0, len(pins.Functional))
	for f := range pins.Functional {
		if l := len(f); l > 0 && f[0] != '<' {
			funcs = append(funcs, f)
			if l > max {
				max = l
			}
		}
	}
	sort.Strings(funcs)
	for _, name := range funcs {
		pin := pins.Functional[name]
		if invalid || pin != host.INVALID {
			if pin == nil {
				fmt.Printf("%-*s: INVALID\n", max, name)
			} else {
				fmt.Printf("%-*s: %s\n", max, name, pin)
			}
		}
	}
}

func printGPIO(invalid bool, maxName, maxFn int) {
	ids := make([]int, 0, len(pins.All))
	for i := range pins.All {
		ids = append(ids, i)
	}
	sort.Ints(ids)
	for _, id := range ids {
		p := pins.All[id]
		if headers.IsConnected(p) {
			fmt.Printf("%-*s: %s\n", maxName, p, p.Function())
		} else if invalid {
			fmt.Printf("%-*s: %-*s (not connected)\n", maxName, p, maxFn, p.Function())
		}
	}
}

func printHardware(invalid bool, maxName, maxFn int) {
	names := make([]string, 0, len(headers.All))
	for name := range headers.All {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		header := headers.All[name]
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
				fmt.Printf("  %*s  %*s  %3d  %-3d  %-*s %s\n", maxFn, line[0].Function(), maxName, line[0], i+1, i+2, maxName, line[1], line[1].Function())
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
	fun := flag.Bool("f", false, "print functional pins (e.g. I2C_SCL1)")
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
	if err := pins.Init(); err != nil {
		return err
	}
	if *info {
		fmt.Printf("MaxSpeed: %dMhz\n", cpu.MaxSpeed/1000000)
	}
	maxName := getMaxName()
	maxFn := getMaxFn()
	if *fun {
		printFunc(*invalid)
	}
	if *gpio {
		printGPIO(*invalid, maxName, maxFn)
	}
	if *hardware {
		printHardware(*invalid, maxName, maxFn)
	}
	return nil
}

func main() {
	if err := mainImpl(); err != nil {
		fmt.Fprintf(os.Stderr, "pins: %s.\n", err)
		os.Exit(1)
	}
}
