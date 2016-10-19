// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"github.com/maruel/dlibox/go/anim1d"
	"github.com/maruel/dlibox/go/donotuse/conn/gpio"
)

func initPIR(painter *anim1d.Painter, config *PIR) error {
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
				// TODO(maruel): Locking.
				painter.SetPattern(string(config.Pattern))
			}
		}
	}()
	return nil
}
