// Copyright 2017 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"log"
	_ "net/http/pprof"

	"github.com/maruel/dlibox/go/msgbus"
)

// mainController is the main function when running as the controller.
func mainController(hostname, mqttHost, mqttUser, mqttPass string, port int) error {
	// Config.
	config := ConfigMgr{}
	config.ResetDefault()
	if err := config.Load(); err != nil {
		log.Printf("Loading config failed: %v", err)
	}
	defer config.Close()

	/*
		b, err := json.MarshalIndent(config, "", "  ")
		if err != nil {
			return err
		}
	*/

	// TODO(maruel): Also initializes itself as a node?
	bus := msgbus.New()
	if mqttHost != "" {
		root := "dlibox/" + hostname
		server, err := msgbus.NewMQTT(mqttHost, hostname, mqttUser, mqttPass, msgbus.Message{root + "/$online", []byte("false")})
		if err != nil {
			// TODO(maruel): Have it continuously try to automatically reconnect.
			log.Printf("Failed to connect to server: %v", err)
		} else {
			bus = server
		}
	}
	initState(bus, nil)

	w, err := initWeb(bus, port, &config.Config, nil)
	if err != nil {
		return err
	}
	defer w.Close()

	return watchFile()
}
