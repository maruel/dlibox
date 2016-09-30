// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// pio-info prints out information about the loaded pio drivers.
package main

import (
	"fmt"
	"os"

	"github.com/maruel/dlibox/go/pio/host"
)

func mainImpl() error {
	state, err := host.Init()
	if err != nil {
		return err
	}
	if len(state.Failed) != 0 {
		fmt.Printf("Drivers failed to load:\n")
		for _, f := range state.Failed {
			fmt.Printf("  - %s: %v\n", f.D, f.Err)
			err = f.Err
		}
	}
	fmt.Printf("Using drivers:\n")
	for _, driver := range state.Loaded {
		fmt.Printf("  - %s\n", driver.String())
	}
	fmt.Printf("Drivers skipped:\n")
	for _, driver := range state.Skipped {
		fmt.Printf("  - %s\n", driver.String())
	}
	return err
}

func main() {
	if err := mainImpl(); err != nil {
		fmt.Fprintf(os.Stderr, "pio-info: %s.\n", err)
		os.Exit(1)
	}
}
