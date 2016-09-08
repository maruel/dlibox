// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package fakes

import (
	"errors"

	"github.com/maruel/dlibox/go/pio/host"
)

// Pin implements host.Pin.
type Pin struct {
	L host.Level
}

func (p *Pin) In(pull host.Pull, edge host.Edge) error {
	if pull == host.Down {
		p.L = host.Low
	} else if pull == host.Up {
		p.L = host.High
	}
	if edge != host.EdgeNone {
		return errors.New("not implemented")
	}
	return nil
}

func (p *Pin) ReadInstant() host.Level {
	return p.L
}

func (p *Pin) ReadEdge() host.Level {
	return p.L
}

func (p *Pin) Out() error {
	return nil
}

func (p *Pin) SetLow() {
	p.L = host.Low
}

func (p *Pin) SetHigh() {
	p.L = host.High
}

var _ host.Pin = &Pin{}
