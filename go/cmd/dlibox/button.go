// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"errors"
	"log"
	"sort"
	"sync"
	"time"

	"github.com/maruel/dlibox/go/donotuse/conn/gpio"
	"github.com/maruel/dlibox/go/modules"
	"github.com/maruel/interrupt"
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

func initButton(b modules.Bus, r map[string][]byte, config *Button) error {
	if len(r) == 0 {
		// TODO(maruel): Temporary hack to disable this code.
		return nil
	}
	if config.PinNumber == -1 {
		return nil
	}
	pin := gpio.ByNumber(config.PinNumber)
	if pin == nil {
		return errors.New("pin not found")
	}
	if err := pin.In(gpio.Up, gpio.Both); err != nil {
		return err
	}

	names := make([]string, 0, len(r))
	for n := range r {
		names = append(names, n)
	}
	sort.Strings(names)
	go func() {
		index := 0
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
			pin.WaitForEdge(-1)
			if state := pin.Read(); state != last {
				last = state
				if state == gpio.Low {
					index = (index + 1) % len(names)
					if err := b.Publish(modules.Message{"painter/setuser", r[names[index]]}, modules.ExactlyOnce, false); err != nil {
						log.Printf("button: failed to publish: %v", err)
					}
				}
				if err := b.Publish(modules.Message{"button", []byte(state.String())}, modules.ExactlyOnce, false); err != nil {
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
