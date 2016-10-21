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
	Mapping map[ir.Key]Command
}

func (i *IR) ResetDefault() {
	i.Lock()
	defer i.Unlock()
	i.Mapping = map[ir.Key]Command{
		ir.KEY_NUMERIC_0: {"painter/setuser", []byte("\"#000000\"")},
		ir.KEY_100PLUS:   {"painter/setuser", []byte("\"#ffffff\"")},
	}
}

func (i *IR) Validate() error {
	i.Lock()
	defer i.Unlock()
	for k, v := range i.Mapping {
		if err := v.Validate(); err != nil {
			return errors.Wrap(err, fmt.Sprintf("can't load command for key %s", k))
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
					config.Lock()
					cmd := config.Mapping[msg.Key]
					config.Unlock()
					if len(cmd.Topic) != 0 {
						b.Publish(modules.Message(cmd), modules.ExactlyOnce, false)
					}
				}
			}
		}
	}()
	return nil
}
