// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Package pins2 exists to remove an import cycle.
//
// Both host and drivers need to reference it.
package pins2

import (
	"sync"

	"github.com/maruel/dlibox/go/pio/protocols/gpio"
)

var (
	Lock       sync.Mutex
	All        PinList
	ByNumber   map[int]gpio.PinIO
	ByName     map[string]gpio.PinIO
	ByFunction map[string]gpio.PinIO
)

type PinList []gpio.PinIO

func (p PinList) Len() int           { return len(p) }
func (p PinList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PinList) Less(i, j int) bool { return p[i].Number() < p[j].Number() }
