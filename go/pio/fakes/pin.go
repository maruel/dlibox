// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package fakes

import (
	"sync"

	"github.com/maruel/dlibox/go/pio/host"
)

// Pin implements host.Pin.
type Pin struct {
	sync.Mutex
	host.Level
	Name      string
	Num       int
	EdgesChan chan host.Level
}

func (p *Pin) String() string {
	return p.Name
}

func (p *Pin) Number() int {
	return p.Num
}

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

func (p *Pin) Read() host.Level {
	p.Lock()
	defer p.Unlock()
	return p.Level
}

func (p *Pin) Edges() (chan host.Level, error) {
	p.Lock()
	defer p.Unlock()
	if p.EdgesChan == nil {
		p.EdgesChan = make(chan host.Level)
	}
	return p.EdgesChan, nil
}

func (p *Pin) Out() error {
	p.Lock()
	defer p.Unlock()
	return nil
}

func (p *Pin) Set(level host.Level) {
	p.Lock()
	defer p.Unlock()
	p.Level = level
}

var _ host.Pin = &Pin{}
