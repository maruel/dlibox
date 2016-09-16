// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package hosttest

import (
	"sync"

	"github.com/maruel/dlibox/go/pio/host"
)

// Pin implements host.Pin.
//
// Modify its members to simulate hardware events.
type Pin struct {
	Name string // Should be immutable
	Num  int    // Should be immutable
	Fn   string // Should be immutable

	sync.Mutex                 // Grab the Mutex before modifying the members to keep it concurrent safe
	host.Level                 // Used for both input and output
	EdgesChan  chan host.Level // Use it to fake edges
}

func (p *Pin) String() string {
	return p.Name
}

func (p *Pin) Number() int {
	return p.Num
}

func (p *Pin) Function() string {
	return p.Fn
}

// In is concurrent safe.
func (p *Pin) In(pull host.Pull) error {
	p.Lock()
	defer p.Unlock()
	if pull == host.Down {
		p.Level = host.Low
	} else if pull == host.Up {
		p.Level = host.High
	}
	return nil
}

// Read is concurrent safe.
func (p *Pin) Read() host.Level {
	p.Lock()
	defer p.Unlock()
	return p.Level
}

// Edges is concurrent safe.
func (p *Pin) Edges() (chan host.Level, error) {
	p.Lock()
	defer p.Unlock()
	if p.EdgesChan == nil {
		p.EdgesChan = make(chan host.Level)
	}
	return p.EdgesChan, nil
}

// Out is concurrent safe.
func (p *Pin) Out() error {
	p.Lock()
	defer p.Unlock()
	return nil
}

// Set is concurrent safe.
func (p *Pin) Set(level host.Level) {
	p.Lock()
	defer p.Unlock()
	p.Level = level
}

var _ host.Pin = &Pin{}
