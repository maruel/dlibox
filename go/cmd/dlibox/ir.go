// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import "github.com/maruel/dlibox/go/donotuse/devices/lirc"

func initIR(painter *painter, config *IR) error {
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
						painter.SetPattern(string(pat))
					}
				}
			}
		}
	}()
	return nil
}
