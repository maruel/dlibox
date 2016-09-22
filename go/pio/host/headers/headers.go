// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package headers

import (
	"fmt"
	"sync"

	"github.com/maruel/dlibox/go/pio/protocols/pins"
)

// All contains all the on-board headers on a micro computer.
//
// The map key is the header name, e.g. "P1" or "EULER" and the value is a
// slice of slice of pins. For a 2x20 header, it's going to be a slice of
// [20][2]pins.Pin.
func All() map[string][][]pins.Pin {
	lock.Lock()
	defer lock.Unlock()
	// TODO(maruel): Return a copy?
	return allHeaders
}

// IsConnected returns true if the pin is on a header.
func IsConnected(p pins.Pin) bool {
	lock.Lock()
	defer lock.Unlock()
	// Populate the map on first use.
	if connectedPins == nil {
		connectedPins = map[string]bool{}
		for name, header := range allHeaders {
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

// Register registers a physical header.
func Register(name string, pins [][]pins.Pin) {
	lock.Lock()
	defer lock.Unlock()
	// TODO(maruel): Copy the slices?
	allHeaders[name] = pins
}

//

var (
	lock          sync.Mutex
	allHeaders    = map[string][][]pins.Pin{} // every known headers as per internal lookup table
	connectedPins map[string]bool             // GPIO pin name to bool
)
