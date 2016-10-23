// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"sync"

	"github.com/maruel/dlibox/go/modules"
)

// MQTT contains settings for connecting to a MQTT server, or not.
type MQTT struct {
	sync.Mutex
	Host string
	User string
	Pass string
}

func (m *MQTT) ResetDefault() {
	m.Lock()
	defer m.Unlock()
	m.Host = "localhost"
}

func (m *MQTT) Validate() error {
	m.Lock()
	defer m.Unlock()
	return nil
}

func initMQTT(config *MQTT) (modules.Bus, error) {
	bus, err := initMQTTInner(config)
	return modules.Rebase(bus, "dlibox/"+hostName), err
}

func initMQTTInner(config *MQTT) (modules.Bus, error) {
	if len(config.Host) == 0 {
		return &modules.LocalBus{}, nil
	}
	server, err := modules.New(config.Host, hostName, config.User, config.Pass)
	if err != nil {
		// Fallback to a local bus.
		return &modules.LocalBus{}, err
	}
	return server, nil
}
