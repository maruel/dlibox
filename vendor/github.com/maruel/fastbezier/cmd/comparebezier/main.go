// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/maruel/fastbezier"
	"github.com/maruel/fastbezier/internal/rejected"
)

type eval struct {
	f          rejected.Evaluator
	name       string
	y          [65536]uint16
	deltas     [65536]int
	relDelta   [65536]string
	totalDelta int
}

func genEval(f rejected.Evaluator, name string, ref *eval) *eval {
	e := &eval{f: f, name: name}
	for i := 0; i < 65536; i++ {
		e.y[i] = f.Eval(uint16(i))
	}
	if ref != nil {
		for i := 0; i < 65536; i++ {
			d := int(e.y[i]) - int(ref.y[i])
			e.deltas[i] = d
			if d == 0 {
				e.relDelta[i] = "0"
			} else {
				e.relDelta[i] = fmt.Sprintf("%2.2f%%", 100.*float32(d)/65535.)
			}
			if d < 0 {
				d = -d
			}
			e.totalDelta += d
		}
	}
	return e
}

func mainImpl() error {
	state := flag.Bool("state", false, "Print the internal state of each evaluator")
	steps := flag.Int("steps", 50, "Number of steps to print")
	flag.Parse()
	if flag.NArg() != 0 {
		return errors.New("do not supply args")
	}
	values := []*eval{}
	values = append(values, genEval(rejected.MakePrecise(0.42, 0, 0.58, 1), "Precs", nil))
	values = append(values, genEval(fastbezier.Make(0.42, 0, 0.58, 1, 0), "LUT", values[0]))
	values = append(values, genEval(fastbezier.MakeFast(0.42, 0, 0.58, 1, 0), "LUTf", values[0]))
	values = append(values, genEval(rejected.MakePointsTrimmed(0.42, 0, 0.58, 1, 0), "PtsT", values[0]))
	values = append(values, genEval(rejected.MakePointsFull(0.42, 0, 0.58, 1, 0), "PtsF", values[0]))
	values = append(values, genEval(rejected.MakeTableTrimmed(0.42, 0, 0.58, 1, 0), "TblT", values[0]))
	values = append(values, genEval(rejected.MakeTableFull(0.42, 0, 0.58, 1, 0), "TblF", values[0]))
	if *state {
		for _, e := range values {
			fmt.Printf("%s\n", e.f)
		}
	}
	fmt.Printf("     x ")
	for _, v := range values {
		fmt.Printf("%6s", v.name)
	}
	fmt.Printf("  ")
	for i := 1; i < len(values); i++ {
		fmt.Printf("%6s", values[i].name)
	}
	fmt.Printf("  ")
	for i := 1; i < len(values); i++ {
		fmt.Printf("%9s", values[i].name)
	}
	fmt.Printf("\n")

	y := make([]uint16, len(values))
	delta := make([]int, len(values)-1)
	relDelta := make([]string, len(values)-1)
	for i := 0; i <= *steps; i++ {
		x := i * 65535 / *steps
		y[0] = values[0].y[x]
		for j := 1; j < len(values); j++ {
			y[j] = values[j].y[x]
			delta[j-1] = values[j].deltas[x]
			relDelta[j-1] = values[j].relDelta[x]
		}
		fmt.Printf("%6d %5v %5v %8v\n", x, y, delta, relDelta)
	}
	return nil
}

func main() {
	if err := mainImpl(); err != nil {
		fmt.Fprintf(os.Stderr, "comparebezier: %s.\n", err)
		os.Exit(1)
	}
}
