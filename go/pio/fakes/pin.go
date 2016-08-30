// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package fakes

import (
	"errors"

	"github.com/maruel/dlibox/go/pio/buses"
)

// Pin implements buses.Pin.
type Pin struct {
	L buses.Level
}

func (p *Pin) In(pull buses.Pull, edge buses.Edge) error {
	if pull == buses.Down {
		p.L = buses.Low
	} else if pull == buses.Up {
		p.L = buses.High
	}
	if edge != buses.EdgeNone {
		return errors.New("not implemented")
	}
	return nil
}

func (p *Pin) ReadInstant() buses.Level {
	return p.L
}

func (p *Pin) ReadEdge() buses.Level {
	return p.L
}

func (p *Pin) Out() error {
	return nil
}

func (p *Pin) SetLow() {
	p.L = buses.Low
}

func (p *Pin) SetHigh() {
	p.L = buses.High
}

var _ buses.Pin = &Pin{}
