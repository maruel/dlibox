// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package device

import (
	"log"

	"github.com/maruel/dlibox/nodes"
	"github.com/maruel/msgbus"
	"periph.io/x/periph/devices/lirc"
)

type irDev struct {
	NodeBase
	Cfg *nodes.IR
}

func (i *irDev) init(b msgbus.Bus) error {
	bus, err := lirc.New()
	if err != nil {
		return err
	}
	go i.run(b, bus)
	return nil
}

func (i *irDev) run(b msgbus.Bus, bus *lirc.Conn) {
	c := bus.Channel()
	for {
		select {
		case msg, ok := <-c:
			if !ok {
				break
			}
			if !msg.Repeat {
				if err := b.Publish(msgbus.Message{"ir", []byte(msg.Key)}, msgbus.BestEffort, false); err != nil {
					log.Printf("%s: failed to publish: %v", i, err)
				}
			}
		}
	}
}
