// Copyright 2017 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Package controller implements the dlibox controller.
//
// When acting as the controller, it connects to the MQTT server and instructs
// devices what events to send back.
package controller

import (
	"fmt"
	"log"

	"github.com/maruel/dlibox/shared"
	"github.com/maruel/interrupt"
	"github.com/maruel/msgbus"
)

// Main is the main function when running as the controller.
func Main(bus msgbus.Bus, port int) error {
	log.Printf("controller.Main(..., %d)", port)
	d := dbMgr{}
	if err := d.Load(); err != nil {
		log.Printf("Loading DB failed: %v", err)
	}
	defer d.Close()

	// Everything is under the namespace "dlibox/"
	dbus := msgbus.RebasePub(msgbus.RebaseSub(bus, "dlibox"), "dlibox")
	// TODO(maruel): Also listen to homie/ for simple nodes.

	// Note: <devID>/$online will not function properly for the controller, use
	// $online.
	shared.InitState(msgbus.RebasePub(dbus, shared.Hostname()), nil)

	w, err := newWebServer(fmt.Sprintf("0.0.0.0:%d", port), true, dbus, &d.db, nil)
	if err != nil {
		return err
	}
	defer w.Close()

	// Publish all the devices.
	for devID, dev := range d.db.Config.Devices {
		b := msgbus.RebasePub(dbus, string(devID))
		shared.RetainedStr(b, "$name", dev.Name)

		// Publish all the device's nodes.
		for nodeID, def := range dev.ToSerialized().Nodes {
			bn := msgbus.RebasePub(b, string(nodeID))
			shared.RetainedStr(bn, "$name", def.Name)
			shared.RetainedStr(bn, "$type", string(def.Type))
			for pID, p := range def.Properties {
				bp := msgbus.RebasePub(bn, string(pID))
				shared.RetainedStr(bp, "$unit", p.Unit)
				shared.RetainedStr(bp, "$datatype", p.DataType)
				shared.RetainedStr(bp, "$format", p.Format)
				shared.RetainedStr(bp, "$settable", fmt.Sprintf("%t", p.Settable))
			}
			shared.Retained(bn, "$config", def.Config)
		}
	}
	if !interrupt.IsSet() {
		shared.RetainedStr(dbus, "$online", "true")
	}
	return shared.WatchFile()
}

func pubErr(b msgbus.Bus, f string, arg ...interface{}) {
	msg := fmt.Sprintf(f, arg...)
	log.Print(msg)
	if err := b.Publish(msgbus.Message{Topic: "$error", Payload: []byte(msg)}, msgbus.ExactlyOnce); err != nil {
		log.Print(err)
	}
}
