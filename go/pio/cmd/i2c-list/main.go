// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// i2c-list lists all I²C buses.
package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/maruel/dlibox/go/pio/host"
	"github.com/maruel/dlibox/go/pio/protocols/i2c"
)

func mainImpl() error {
	if _, err := host.Init(); err != nil {
		return err
	}
	all := i2c.All()
	names := make([]string, 0, len(all))
	for name := range all {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		fmt.Printf("%s:\n", name)
		bus, err := all[name]()
		if err != nil {
			fmt.Printf("  Failed to open: %v\n", err)
			continue
		}
		if p, ok := bus.(i2c.Pins); ok {
			fmt.Printf("  SCL: %s\n", p.SCL())
			fmt.Printf("  SDA: %s\n", p.SDA())
		}
		bus.Close()
	}
	return nil
}

func main() {
	if err := mainImpl(); err != nil {
		fmt.Fprintf(os.Stderr, "i2c-list: %s.\n", err)
		os.Exit(1)
	}
}
