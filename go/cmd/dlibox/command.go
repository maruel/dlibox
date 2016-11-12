// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"errors"
	"fmt"

	"github.com/maruel/dlibox/go/modules"
)

type Command struct {
	Topic   string
	Payload string
}

func (c *Command) ToMsg() modules.Message {
	return modules.Message{c.Topic, []byte(c.Payload)}
}

func (c *Command) Validate() error {
	switch c.Topic {
	case "leds/temperature", "leds/intensity":
		// TODO(maruel): Validate number?
		return nil
	case "painter/setautomated", "painter/setnow", "painter/setuser":
		return Pattern(c.Payload).Validate()
	case "":
		if len(c.Payload) != 0 {
			return errors.New("empty topic requires empty payload")
		}
		return nil
	default:
		return fmt.Errorf("unsupported command %v", c.Topic)
	}
}
