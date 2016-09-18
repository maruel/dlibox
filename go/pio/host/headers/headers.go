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
	"github.com/maruel/dlibox/go/pio/host/internal"
	"github.com/maruel/dlibox/go/pio/host/internal/pine64"
	"github.com/maruel/dlibox/go/pio/host/internal/rpi"
)

// All contains all the on-board headers on a micro computer. The map key is
// the header name, e.g. "P1" or "EULER" and the value is a slice of slice of
// pins. For a 2x20 header, it's going to be a slice of [20][2]host.Pin.
func All() map[string][][]host.Pin {
	lock.Lock()
	defer lock.Unlock()
	initAll()
	return all
}

// IsConnected returns true if the pin is on a header.
func IsConnected(p host.Pin) bool {
	lock.Lock()
	defer lock.Unlock()
	// Populate the map on first use.
	if reverse == nil {
		initAll()
		reverse = map[string]bool{}
		for name, header := range all {
			for i, line := range header {
				for j, item := range line {
					if item == nil || len(item.String()) == 0 {
						fmt.Printf("%s[%d][%d]\n", name, i, j)
					}
					reverse[item.String()] = true
				}
			}
		}
	}
	b, _ := reverse[p.String()]
	return b
}

//

var (
	lock    sync.Mutex
	all     map[string][][]host.Pin
	reverse map[string]bool
)

func initAll() {
	if all == nil {
		if internal.IsRaspberryPi() {
			all = rpi.Headers
		} else if internal.IsPine64() {
			all = pine64.Headers
		} else {
			all = map[string][][]host.Pin{}
		}
	}
}
