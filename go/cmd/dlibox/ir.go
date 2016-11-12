// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
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
		ir.KEY_CHANNELDOWN: {"leds/temperature", "-500"},
		ir.KEY_CHANNEL:     {"leds/temperature", "5000"},
		ir.KEY_CHANNELUP:   {"leds/temperature", "+500"},
		ir.KEY_PREVIOUS:    {"leds/temperature", "3000"},
		ir.KEY_NEXT:        {"leds/temperature", "5000"},
		ir.KEY_PLAYPAUSE:   {"leds/temperature", "6500"},
		ir.KEY_VOLUMEDOWN:  {"leds/intensity", "-15"},
		ir.KEY_VOLUMEUP:    {"leds/intensity", "+15"},
		ir.KEY_EQ:          {"leds/intensity", "128"},
		ir.KEY_NUMERIC_0:   {"leds/intensity", "0"},
		ir.KEY_100PLUS:     {"painter/setuser", "\"#ffffff\""},
		ir.KEY_200PLUS:     {"leds/intensity", "255"},
		ir.KEY_NUMERIC_1:   {"painter/setuser", "\"Rainbow\""},
		ir.KEY_NUMERIC_2:   {"painter/setuser", "{\"Child\":\"Rainbow\",\"MovePerHour\":108000,\"_type\":\"Rotate\"}"},
		ir.KEY_NUMERIC_3:   {"painter/setuser", "{\"Child\":{\"Frame\":\"Lff0000ff0000ff0000ff0000ff0000ffffffffffffffffffffffffffffff\",\"_type\":\"Repeated\"},\"MovePerHour\":21600,\"_type\":\"Rotate\"}"},
		ir.KEY_NUMERIC_4:   {"painter/setuser", "{\"Child\":\"L0100010f0000000f0000000f\",\"_type\":\"Chronometer\"}"},
		ir.KEY_NUMERIC_5:   {"painter/setuser", "{\"Child\":\"Lff0000ff0000ee0000dd0000cc0000bb0000aa0000990000880000770000660000550000440000330000220000110000\",\"MovePerHour\":108000,\"_type\":\"PingPong\"}"},
		ir.KEY_NUMERIC_6:   {"painter/setuser", "{\"C\":\"#ff9000\",\"_type\":\"NightStars\"}"},
		ir.KEY_NUMERIC_7:   {"painter/setuser", "{\"Curve\":\"ease-out\",\"Patterns\":[{\"Patterns\":[{\"Child\":\"Lff0000ff0000ee0000dd0000cc0000bb0000aa0000990000880000770000660000550000440000330000220000110000\",\"MovePerHour\":108000,\"_type\":\"Rotate\"},{\"_type\":\"Aurore\"}],\"_type\":\"Add\"},{\"Patterns\":[{\"_type\":\"Aurore\"},{\"C\":\"#ffffff\",\"_type\":\"NightStars\"}],\"_type\":\"Add\"}],\"ShowMS\":10000,\"TransitionMS\":5000,\"_type\":\"Loop\"}"},
		ir.KEY_NUMERIC_8:   {"painter/setuser", "{\"Left\":{\"Curve\":\"ease-out\",\"Patterns\":[\"#000f00\",\"#00ff00\",\"#1f0f00\",\"#ffa900\"],\"ShowMS\":100,\"TransitionMS\":700,\"_type\":\"Loop\"},\"Offset\":\"50%\",\"Right\":{\"Curve\":\"ease-out\",\"Patterns\":[\"#1f0f00\",\"#ffa900\",\"#000f00\",\"#00ff00\"],\"ShowMS\":100,\"TransitionMS\":700,\"_type\":\"Loop\"},\"_type\":\"Split\"}"},
		ir.KEY_NUMERIC_9:   {"painter/setuser", "{\"Child\":\"Lffffff\",\"MovePerHour\":108000,\"_type\":\"PingPong\"}"},
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
						if err := b.Publish(cmd.ToMsg(), modules.BestEffort, false); err != nil {
							log.Printf("ir: failed to publish: %v", err)
						}
					}
					if err = b.Publish(modules.Message{"ir", []byte(msg.Key)}, modules.BestEffort, false); err != nil {
						log.Printf("ir: failed to publish: %v", err)
					}
				}
			}
		}
	}()
	return nil
}
