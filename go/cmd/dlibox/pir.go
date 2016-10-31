// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/maruel/dlibox/go/donotuse/conn/gpio"
	"github.com/maruel/dlibox/go/modules"
	"github.com/pkg/errors"
)

// PIR contains a motion detection behavior.
type PIR struct {
	sync.Mutex
	PinNumber int
	Cmd       Command
}

func (p *PIR) ResetDefault() {
	p.Lock()
	defer p.Unlock()
	p.PinNumber = -1
	p.Cmd = Command{"painter/setnow", "\"#ffffff\""}
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
	if config.PinNumber == -1 {
		return nil
	}
	p := gpio.ByNumber(config.PinNumber)
	if p == nil {
		return fmt.Errorf("pir: failed to find pin %d", config.PinNumber)
	}
	if err := p.In(gpio.Down, gpio.Both); err != nil {
		return errors.Wrapf(err, "pir: failed to pull down %s", p)
	}
	go func() {
		for {
			p.WaitForEdge(-1)
			if p.Read() == gpio.High {
				log.Printf("pir: high")
				// TODO(maruel): sub-second resolution?
				now := time.Now()
				nowStr := []byte(fmt.Sprintf("%d %s", now.Unix(), now))
				err := b.Publish(modules.Message{"pir", nowStr}, modules.MinOnce, false)
				if err != nil {
					log.Printf("pir: failed to publish: %v", err)
				}
				config.Lock()
				if config.Cmd.Topic != "" {
					err = b.Publish(config.Cmd.ToMsg(), modules.ExactlyOnce, false)
				}
				config.Unlock()
				if err != nil {
					log.Printf("pir: failed to publish: %v", err)
				}
			} else {
				log.Printf("pir: low")
			}
		}
	}()
	return nil
}
