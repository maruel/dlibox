// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package rules

import (
	"errors"
	"fmt"

	"github.com/maruel/msgbus"
)

// Command is an MQTT message to send to take action.
//
// Most commands shall respect the Homie convention.
// https://github.com/marvinroger/homie
type Command struct {
	Topic   string
	Payload string
}

// ToMsg converts the command to a MQTT Message.
func (c *Command) ToMsg() msgbus.Message {
	return msgbus.Message{Topic: c.Topic, Payload: []byte(c.Payload)}
}

// Validate ensures the command is valid.
func (c *Command) Validate() error {
	switch c.Topic {
	case "leds/temperature", "leds/intensity":
		// TODO(maruel): Validate number?
		return nil
	case "painter/setautomated", "painter/setnow", "painter/setuser":
		// TODO(maruel): Add back once migrated out of ../cmd/dlibox/config.go.
		//return painter1d.Pattern(c.Payload).Validate()
		return nil
	case "":
		if len(c.Payload) != 0 {
			return errors.New("empty topic requires empty payload")
		}
		return nil
	default:
		return fmt.Errorf("unsupported command %v", c.Topic)
	}
}
