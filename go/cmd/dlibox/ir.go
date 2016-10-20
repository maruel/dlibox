// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"fmt"
	"sync"

	"github.com/maruel/dlibox/go/donotuse/conn/ir"
	"github.com/maruel/dlibox/go/donotuse/devices/lirc"
	"github.com/maruel/dlibox/go/modules"
	"github.com/pkg/errors"
)

// IR contains InfraRed remote information.
type IR struct {
	sync.Mutex
	Mapping map[ir.Key]Pattern // TODO(maruel): We may actually do something more complex than just set a pattern.
}

func (i *IR) ResetDefault() {
	i.Lock()
	defer i.Unlock()
	i.Mapping = map[ir.Key]Pattern{
		ir.KEY_NUMERIC_0: "\"#000000\"",
		ir.KEY_100PLUS:   "\"#ffffff\"",
	}
}

func (i *IR) Validate() error {
	i.Lock()
	defer i.Unlock()
	for k, v := range i.Mapping {
		if err := v.Validate(); err != nil {
			return errors.Wrap(err, fmt.Sprintf("can't load pattern for key %s", k))
		}
	}
	return nil
}

func initIR(b modules.Bus, config *IR) error {
	bus, err := lirc.New()
	if err != nil {
		return err
	}
	go func() {
		c := bus.Channel()
		for {
			select {
			case msg, ok := <-c:
				if !ok {
					break
				}
				if !msg.Repeat {
					// TODO(maruel): Locking.
					if pat := config.Mapping[msg.Key]; len(pat) != 0 {
						b.Publish(modules.Message{"painter/setuser", []byte(pat)}, modules.ExactlyOnce, false)
					}
				}
			}
		}
	}()
	return nil
}
