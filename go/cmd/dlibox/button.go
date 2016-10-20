// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"errors"
	"sort"
	"time"

	"github.com/maruel/dlibox/go/donotuse/conn/gpio"
	"github.com/maruel/interrupt"
)

func initButton(p *painter, r map[string]string, config *Button) error {
	if len(r) == 0 {
		// TODO(maruel): Temporary hack to disable this code.
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
					p.SetPattern(r[names[index]])
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
