// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package device

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/maruel/dlibox/modules/nodes/pir"
	"github.com/maruel/dlibox/msgbus"
	"github.com/pkg/errors"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
)

type PIR struct {
	sync.Mutex
	pirs []pir.Dev
}

func (p *PIR) init(b msgbus.Bus) error {
	for _, cfg := range p.pirs {
		pin := gpioreg.ByName(cfg.Pin)
		if pin == nil {
			return fmt.Errorf("pir %s: failed to find pin %s", cfg.Name, cfg.Pin)
		}
		if err := pin.In(gpio.PullDown, gpio.BothEdges); err != nil {
			return errors.Wrapf(err, "pir %s: failed to pull down %s", cfg.Name, pin)
		}
		go runPir(b, cfg, pin)
	}
	return nil
}

func runPir(b msgbus.Bus, cfg pir.Dev, pin gpio.PinIn) {
	for {
		pin.WaitForEdge(-1)
		if pin.Read() == gpio.High {
			log.Printf("pir %s: high", cfg.Name)
			// TODO(maruel): sub-second resolution?
			now := time.Now()
			nowStr := []byte(fmt.Sprintf("%d %s", now.Unix(), now))
			// Fix convention.
			err := b.Publish(msgbus.Message{cfg.Name + "/pir", nowStr}, msgbus.BestEffort, false)
			if err != nil {
				log.Printf("pir %s: failed to publish: %v", cfg.Name, err)
			}
			if err != nil {
				log.Printf("pir %s: failed to publish: %v", cfg.Name, err)
			}
		} else {
			log.Printf("pir %s: low", cfg.Name)
		}
	}
}
