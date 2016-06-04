// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// playing is a small app to play with the pins, nothing more. You are not
// expected to use it as-is.
package main

import (
	"fmt"
	"os"

	"github.com/kr/pretty"
	"github.com/maruel/dlibox-go/anim1d"
	"github.com/maruel/dlibox-go/apa102"
)

func mainImpl() error {
	pixels := make(anim1d.Frame, 150)
	var p anim1d.Rainbow
	p.NextFrame(pixels, 0)
	var d []byte
	apa102.Raster(pixels, &d)
	pretty.Print(d)
	return nil
}

func main() {
	if err := mainImpl(); err != nil {
		fmt.Fprintf(os.Stderr, "playing\n: %s.\n", err)
		os.Exit(1)
	}
}
