// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"sync"

	"github.com/maruel/dlibox/go/donotuse/conn/gpio"
	"github.com/maruel/dlibox/go/modules"
	"github.com/pkg/errors"
)

// PIR contains a motion detection behavior.
type PIR struct {
	sync.Mutex
	Pin int
	Cmd Command
}

func (p *PIR) ResetDefault() {
	p.Lock()
	defer p.Unlock()
	p.Pin = -1
	p.Cmd = Command{"painter/setnow", []byte("\"#ffffff\"")}
}

func (p *PIR) Validate() error {
	p.Lock()
	defer p.Unlock()
	if err := p.Cmd.Validate(); err != nil {
		return errors.Wrap(err, "can't load command for PIR")
	}
	return nil
}

func initPIR(b modules.Bus, config *PIR) error {
	if config.Pin == -1 {
		return nil
	}
	p := gpio.ByNumber(config.Pin)
	if p == nil {
		return nil
	}
	if err := p.In(gpio.Down, gpio.Both); err != nil {
		return err
	}
	go func() {
		for {
			p.WaitForEdge(-1)
			if p.Read() == gpio.High {
				config.Lock()
				b.Publish(modules.Message(config.Cmd), modules.ExactlyOnce, false)
				config.Unlock()
			}
		}
	}()
	return nil
}
