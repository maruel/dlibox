// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// psf prints out available font family and runs supported.
package main

import (
	"flag"
	"fmt"
	"os"
	"sort"

	"github.com/maruel/psf"
)

func mainImpl() error {
	showRunes := flag.Bool("r", false, "show all runes supported")
	flag.Parse()

	var fonts []string
	if flag.NArg() > 0 {
		fonts = flag.Args()
	} else {
		var err error
		fonts, err = psf.Enumerate()
		if err != nil {
			return err
		}
	}
	sort.Strings(fonts)
	max := 0
	for _, name := range fonts {
		if len(name) > max {
			max = len(name)
		}
	}

	for _, name := range fonts {
		f, err := psf.Load(name)
		if err != nil {
			return err
		}
		size := fmt.Sprintf("%dx%d", f.W, f.H)
		fmt.Printf("%-*s: %6s %5d runes  (v%d)\n", max, name, size, len(f.Letters), f.Version)
		if *showRunes {
			var runes runesList
			for r, bitmap := range f.Letters {
				// If the bitmap is all background color (i.e. whitespace), skip it.
				found := false
				for _, b := range bitmap {
					if b != 0 {
						found = true
						break
					}
				}
				if found {
					runes = append(runes, r)
				}
			}
			sort.Sort(runes)
			for len(runes) != 0 {
				col := 78
				if len(runes) < col {
					col = len(runes)
				}
				fmt.Printf("  %s\n", string(runes[:col]))
				runes = runes[col:]
			}
		}
	}
	return nil
}

type runesList []rune

func (r runesList) Len() int           { return len(r) }
func (r runesList) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r runesList) Less(i, j int) bool { return r[i] < r[j] }

func main() {
	if err := mainImpl(); err != nil {
		fmt.Fprintf(os.Stderr, "psf: %s.\n", err)
		os.Exit(1)
	}
}
