// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"fmt"

	"github.com/maruel/dlibox/go/modules"
)

type Command modules.Message

func (c *Command) Validate() error {
	switch c.Topic {
	case "painter/setautomated", "painter/setnow", "painter/setuser":
		return Pattern(c.Payload).Validate()
	case "":
		return nil
	default:
		return fmt.Errorf("unsupported command %v", c.Topic)
	}
}
