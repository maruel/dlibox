// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/maruel/fastbezier"
)

func mainImpl() error {
	flag.Parse()

	if flag.NArg() != 5 {
		return errors.New("supply 5 values")
	}
	x0, err := strconv.ParseFloat(flag.Arg(0), 64)
	if err != nil {
		return err
	}
	y0, err := strconv.ParseFloat(flag.Arg(1), 64)
	if err != nil {
		return err
	}
	x1, err := strconv.ParseFloat(flag.Arg(2), 64)
	if err != nil {
		return err
	}
	y1, err := strconv.ParseFloat(flag.Arg(3), 64)
	if err != nil {
		return err
	}
	steps, err := strconv.Atoi(flag.Arg(4))
	if err != nil {
		return err
	}
	f := fastbezier.Make(float32(x0), float32(y0), float32(x1), float32(y1), uint16(steps))
	_, err = fmt.Printf("%s\n", f)
	return err
}

func main() {
	if err := mainImpl(); err != nil {
		fmt.Fprintf(os.Stderr, "usage: makebezier <x0> <y0> <x1> <y1> <steps>\nmakebezier: %s.\n", err)
		os.Exit(1)
	}
}
