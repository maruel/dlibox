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
	"github.com/maruel/interrupt"
	"github.com/pkg/errors"
)

// Button contains settings for controlling the lights through a button.
type Button struct {
	sync.Mutex
	PinNumber int
}

func (b *Button) ResetDefault() {
	b.Lock()
	defer b.Unlock()
	b.PinNumber = -1
}

func (b *Button) Validate() error {
	b.Lock()
	defer b.Unlock()
	return nil
}

func initButton(b modules.Bus, config *Button) error {
	if config.PinNumber == -1 {
		return nil
	}
	p := gpio.ByNumber(config.PinNumber)
	if p == nil {
		return fmt.Errorf("button: failed to find pin %d", config.PinNumber)
	}
	if err := p.In(gpio.Up, gpio.Both); err != nil {
		return errors.Wrapf(err, "button: failed to pull up %s", p)
	}

	/*
		names := make([]string, 0, len(r))
		for n := range r {
			names = append(names, n)
		}
		sort.Strings(names)
	*/
	go func() {
		//index := 0
		last := gpio.High
		for {
			// Types of press:
			// - Short press (<2s)
			// - 2s press
			// - 4s press
			// - double-click (incompatible with repeated short press)
			//
			// Functions:
			// - Bonne nuit
			// - Next / Prev
			// - Éteindre (longer press après bonne nuit?)
			p.WaitForEdge(-1)
			if state := p.Read(); state != last {
				last = state
				/*
					if state == gpio.Low {
						index = (index + 1) % len(names)
						if err := b.Publish(modules.Message{"painter/setuser", r[names[index]]}, modules.BestEffort, false); err != nil {
							log.Printf("button: failed to publish: %v", err)
						}
					}
				*/
				if err := b.Publish(modules.Message{"button", []byte(state.String())}, modules.BestEffort, false); err != nil {
					log.Printf("button: failed to publish: %v", err)
				}
			}
			select {
			case <-interrupt.Channel:
				return
			case <-time.After(time.Millisecond):
			}
		}
	}()
	return nil
}
