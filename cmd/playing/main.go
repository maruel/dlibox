// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// playing is a small app to play with the pins, nothing more. You are not
// expected to use it as-is.
package main

import (
	"fmt"
	"os"
	"reflect"

	"github.com/maruel/dlibox-go/anim1d"
	"github.com/maruel/dlibox-go/apa102"
)

func printFrame(p anim1d.Pattern, l int) {
	// Generate a frame.
	pixels := make(anim1d.Frame, l)
	p.NextFrame(pixels, 0)

	// Convert to apa102 protocol.
	var d []byte
	apa102.Raster(pixels, &d)

	// Print it.
	const cols = 16
	fmt.Printf("uint8_t %s[] = {", reflect.TypeOf(p).Elem().Name())
	for i, b := range d {
		if i%cols == 0 {
			fmt.Printf("\n  ")
		}
		fmt.Printf("0x%02x,", b)
		if i%cols != cols-1 && i != len(d)-1 {
			fmt.Printf(" ")
		}
	}
	fmt.Printf("\n};\n")
}

func mainImpl() error {
	printFrame(&anim1d.Rainbow{}, 144)
	printFrame(&anim1d.Color{0x7f, 0x7f, 0x7f}, 144)
	return nil
}

func main() {
	if err := mainImpl(); err != nil {
		fmt.Fprintf(os.Stderr, "playing\n: %s.\n", err)
		os.Exit(1)
	}
}
