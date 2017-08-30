// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package device

import (
	"fmt"
	"log"
	"time"

	"github.com/maruel/dlibox/nodes"
	"github.com/maruel/interrupt"
	"github.com/maruel/msgbus"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
)

type buttonDev struct {
	NodeBase
	Cfg *nodes.Button
}

func (b *buttonDev) init(bus msgbus.Bus) error {
	pin := gpioreg.ByName(b.Cfg.Pin)
	if pin == nil {
		return fmt.Errorf("%s: failed to find pin %s", b, b.Cfg.Pin)
	}
	if err := pin.In(gpio.PullDown, gpio.BothEdges); err != nil {
		return fmt.Errorf("%s: failed to pull down %s: %v", b, pin, err)
	}
	go b.run(bus, pin)
	return nil
}

func (b *buttonDev) run(bus msgbus.Bus, pin gpio.PinIn) {
	//index := 0
	last := gpio.High
	for {
		// Types of press:
		// - Short press (<2s)
		// - 2s press
		// - 4s press
		// - double-click (incompatible with repeated short press)
		pin.WaitForEdge(-1)
		if state := pin.Read(); state != last {
			last = state
			log.Printf("%s: %s", b, state)
			// TODO(maruel): sub-second resolution?
			now := time.Now()
			nowStr := []byte(fmt.Sprintf("%d %s", now.Unix(), now))
			// Fix convention.
			err := bus.Publish(msgbus.Message{Topic: "button", Payload: nowStr}, msgbus.BestEffort)
			if err != nil {
				log.Printf("%s: failed to publish: %v", b, err)
			}
			if err != nil {
				log.Printf("%s: failed to publish: %v", b, err)
			}
		} else {
			log.Printf("%s: %s", b, state)
		}
		select {
		case <-interrupt.Channel:
			return
		case <-time.After(time.Millisecond):
		}
	}
}
