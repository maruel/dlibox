// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package device

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/maruel/dlibox/go/modules/nodes/button"
	"github.com/maruel/dlibox/go/msgbus"
	"github.com/maruel/interrupt"
	"github.com/pkg/errors"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
)

type Button struct {
	sync.Mutex
	buttons []button.Dev
}

func (b *Button) init(bus msgbus.Bus) error {
	for _, cfg := range b.buttons {
		pin := gpioreg.ByName(cfg.Pin)
		if pin == nil {
			return fmt.Errorf("button %s: failed to find pin %s", cfg.Name, cfg.Pin)
		}
		if err := pin.In(gpio.PullDown, gpio.BothEdges); err != nil {
			return errors.Wrapf(err, "button %s: failed to pull down %s", cfg.Name, pin)
		}
		go runButton(bus, cfg, pin)
	}
	return nil
}

func runButton(bus msgbus.Bus, cfg button.Dev, pin gpio.PinIn) {
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
			log.Printf("button %s: %s", cfg.Name, state)
			// TODO(maruel): sub-second resolution?
			now := time.Now()
			nowStr := []byte(fmt.Sprintf("%d %s", now.Unix(), now))
			// Fix convention.
			err := bus.Publish(msgbus.Message{cfg.Name + "/button", nowStr}, msgbus.BestEffort, false)
			if err != nil {
				log.Printf("button %s: failed to publish: %v", cfg.Name, err)
			}
			if err != nil {
				log.Printf("button %s: failed to publish: %v", cfg.Name, err)
			}
		} else {
			log.Printf("button %s: %s", cfg.Name, state)
		}
		select {
		case <-interrupt.Channel:
			return
		case <-time.After(time.Millisecond):
		}
	}
}
