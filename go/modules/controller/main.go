// Copyright 2017 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Packages the static files in a .go file.
//go:generate go run package.go -out static_files_gen.go ../../../web

package controller

import (
	"fmt"
	"log"
	_ "net/http/pprof"

	"github.com/maruel/dlibox/go/modules/shared"
	"github.com/maruel/dlibox/go/msgbus"
)

// Main is the main function when running as the controller.
func Main(bus msgbus.Bus, port int) error {
	log.Printf("controller.Main(..., %d)", port)
	d := DBMgr{}
	if err := d.Load(); err != nil {
		log.Printf("Loading DB failed: %v", err)
	}
	defer d.Close()

	// Note: <devID>/$online will not function properly for the controller, use
	// $online.
	shared.InitState(msgbus.RebasePub(bus, shared.Hostname()), nil)

	w, err := initWeb(bus, port, &d.DB, nil)
	if err != nil {
		return err
	}
	defer w.Close()

	// Publish all the nodes.
	for devID, dev := range d.DB.Config.Devices {
		b := msgbus.RebasePub(bus, string(devID))
		retained(b, "$name", dev.Name)
		for nodeID, def := range dev.ToNodes() {
			bn := msgbus.RebasePub(b, string(nodeID))
			retained(bn, "$name", def.Name)
			retained(bn, "$type", string(def.Type))
			for pID, p := range def.Properties {
				bp := msgbus.RebasePub(bn, string(pID))
				retained(bp, "$unit", p.Unit)
				retained(bp, "$datatype", p.DataType)
				retained(bp, "$format", p.Format)
				retained(bp, "$settable", fmt.Sprintf("%t", p.Settable))
			}
			retainedBytes(bn, "$config", def.Config)
		}
	}
	retained(bus, "$online", "true")
	return shared.WatchFile()
}

func retained(b msgbus.Bus, topic, payload string) {
	retainedBytes(b, topic, []byte(payload))
}

func retainedBytes(b msgbus.Bus, topic string, payload []byte) {
	if err := b.Publish(msgbus.Message{topic, payload}, msgbus.MinOnce, true); err != nil {
		log.Printf("Failed to publish %s: %v", topic, err)
	}
}