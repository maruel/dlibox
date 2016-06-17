// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"errors"
	"fmt"
	"image/gif"
	"os"
)

func mainImpl() error {
	if len(os.Args) != 1 {
		return errors.New("usage: stream GIF stdin, get C out")
	}
	img, err := gif.Decode(os.Stdin)
	if err != nil {
		return err
	}
	size := img.Bounds().Max
	if _, err := fmt.Printf("const uint8_t data[] = {"); err != nil {
		return err
	}
	if size.X%8 != 0 {
		return errors.New("width must be multiple of 8")
	}
	n := 0
	for y := 0; y < size.Y; y++ {
		for x := 0; x < size.X; x += 8 {
			b := 0
			for i := 0; i < 8; i++ {
				b <<= 1
				r, _, _, _ := img.At(x+i, y).RGBA()
				if r != 0 {
					b |= 1
				}
			}
			if n%16 == 0 {
				if _, err := fmt.Printf("\n "); err != nil {
					return err
				}
			}
			n++
			if _, err := fmt.Printf(" 0x%02X,", b); err != nil {
				return err
			}
		}
	}
	if _, err := fmt.Printf("\n};\n"); err != nil {
		return err
	}
	return nil
}

func main() {
	if err := mainImpl(); err != nil {
		fmt.Fprintf(os.Stderr, "gif2c: %s\n", err)
		os.Exit(1)
	}
}
