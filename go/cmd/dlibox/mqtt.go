// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"sync"

	"github.com/maruel/dlibox/go/msgbus"
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
	m.Host = "tcp://localhost:1883"
}

func (m *MQTT) Validate() error {
	m.Lock()
	defer m.Unlock()
	return nil
}

func initMQTT(config *MQTT) (msgbus.Bus, error) {
	bus, err := initMQTTInner(config)
	root := "dlibox/" + hostName
	return msgbus.RebaseSub(msgbus.RebasePub(bus, root), root), err
}

func initMQTTInner(config *MQTT) (msgbus.Bus, error) {
	if len(config.Host) == 0 {
		return msgbus.New(), nil
	}
	server, err := msgbus.NewMQTT(config.Host, hostName, config.User, config.Pass, msgbus.Message{})
	if err != nil {
		// Fallback to a local bus.
		return msgbus.New(), err
	}
	return server, nil
}
