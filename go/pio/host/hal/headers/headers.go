// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// package headers contains a table to represent the physical headers found on
// micro computers.
package headers

import (
	"fmt"
	"sync"

	"github.com/maruel/dlibox/go/pio/host"
	"github.com/maruel/dlibox/go/pio/host/hal/pine64"
	"github.com/maruel/dlibox/go/pio/host/hal/rpi"
	"github.com/maruel/dlibox/go/pio/host/internal"
)

// All contains all the on-board headers on a micro computer. The map key is
// the header name, e.g. "P1" or "EULER" and the value is a slice of slice of
// pins. For a 2x20 header, it's going to be a slice of [20][2]host.PinIO.
func All() map[string][][]host.PinIO {
	lock.Lock()
	defer lock.Unlock()
	initAll()
	return all
}

// IsConnected returns true if the pin is on a header.
func IsConnected(p host.PinIO) bool {
	lock.Lock()
	defer lock.Unlock()
	// Populate the map on first use.
	if connectedPins == nil {
		initAll()
		connectedPins = map[string]bool{}
		for name, header := range all {
			for i, line := range header {
				for j, pin := range line {
					if pin == nil || len(pin.String()) == 0 {
						fmt.Printf("%s[%d][%d]\n", name, i, j)
					}
					connectedPins[pin.String()] = true
				}
			}
		}
	}
	b, _ := connectedPins[p.String()]
	return b
}

//

var (
	lock          sync.Mutex
	all           map[string][][]host.PinIO // every known headers as per internal lookup table
	connectedPins map[string]bool           // GPIO pin name to bool
)

func initAll() {
	if all == nil {
		if internal.IsRaspberryPi() {
			all = rpi.Headers
		} else if internal.IsPine64() {
			all = pine64.Headers
		} else {
			// Implement!
			all = map[string][][]host.PinIO{}
		}
	}
}
