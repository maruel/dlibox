// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package device

import (
	"fmt"
	"log"
	"time"

	"github.com/maruel/dlibox/nodes"
	"github.com/maruel/msgbus"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
)

type pirDev struct {
	NodeBase
	Cfg *nodes.PIR
}

func (p *pirDev) init(b msgbus.Bus) error {
	pin := gpioreg.ByName(p.Cfg.Pin)
	if pin == nil {
		return fmt.Errorf("%s: failed to find pin %s", p, p.Cfg.Pin)
	}
	if err := pin.In(gpio.PullDown, gpio.BothEdges); err != nil {
		return fmt.Errorf("%s: failed to pull down %s: %v", p, pin, err)
	}
	go p.run(b, pin)
	return nil
}

func (p *pirDev) run(b msgbus.Bus, pin gpio.PinIn) {
	for {
		pin.WaitForEdge(-1)
		if pin.Read() == gpio.High {
			log.Printf("%s: high", p)
			// TODO(maruel): sub-second resolution?
			now := time.Now()
			nowStr := []byte(fmt.Sprintf("%d %s", now.Unix(), now))
			// Fix convention.
			err := b.Publish(msgbus.Message{"pir", nowStr}, msgbus.BestEffort, false)
			if err != nil {
				log.Printf("%s: failed to publish: %v", p, err)
			}
			if err != nil {
				log.Printf("%s: failed to publish: %v", p, err)
			}
		} else {
			log.Printf("%s: low", p)
		}
	}
}
