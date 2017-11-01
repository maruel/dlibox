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
	"strings"

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

	go hack(dbus)
	return shared.WatchFile()
}

func pubErr(b msgbus.Bus, f string, arg ...interface{}) {
	msg := fmt.Sprintf(f, arg...)
	log.Print(msg)
	if err := b.Publish(msgbus.Message{Topic: "$error", Payload: []byte(msg)}, msgbus.ExactlyOnce); err != nil {
		log.Print(err)
	}
}

func hack(b msgbus.Bus) {
	cPIR, err := b.Subscribe("raspberrypi-e59e/pir/pir", msgbus.ExactlyOnce)
	if err != nil {
		pubErr(b, "Subscription failed: %v", err)
	} else {
		go func() {
			for range cPIR {
				// Trigger stuff.
				if err := b.Publish(msgbus.Message{Topic: "raspberrypi-e59e/sound/sound", Payload: []byte("doorbell")},
					msgbus.ExactlyOnce); err != nil {
					pubErr(b, "Publish failed: %v", err)
				}
				i := trim(animHaut)
				if err := b.Publish(msgbus.Message{Topic: "raspberrypi-e59e/apa/anim1d", Payload: []byte(i), Retained: true},
					msgbus.ExactlyOnce); err != nil {
					pubErr(b, "Publish failed: %v", err)
				}
				j := trim(animCercueil)
				if err := b.Publish(msgbus.Message{Topic: "raspberrypi-73f5/apa/anim1d", Payload: []byte(j), Retained: true},
					msgbus.ExactlyOnce); err != nil {
					pubErr(b, "Publish failed: %v", err)
				}
				/*
					k := trim(animBouilloir)
					if err := b.Publish(msgbus.Message{Topic: "raspberrypi-681e/apa/anim1d", Payload: []byte(k), Retained: true},
						msgbus.ExactlyOnce); err != nil {
						pubErr(b, "Publish failed: %v", err)
					}
				*/
			}
		}()
	}
}

func trim(s string) string {
	return strings.Replace(strings.Replace(strings.Replace(s, " ", "", -1), "\n", "", -1), "\t", "", -1)
}

const animHaut = `
{
	"After": {
		"Left": {
				"Curve": "ease-in-out",
				"Patterns": [
					"#ffa900",
					"#1f0f00"
				],
				"ShowMS": 0,
				"TransitionMS": 1000,
				"_type": "Loop"
		},
		"Offset": "66.66%",
		"Right": {
			"Curve": "steps(1,end)",
			"Patterns": [
				"#ffffff",
				"#000000",
				"#ffffff",
				"#000000",
				"#000000",
				"#ffffff",
				"#000000",
				"#000000",
				"#000000",
				"#000000"
			],
			"ShowMS": 100,
			"TransitionMS": 500,
			"_type": "Loop"
		},
		"_type": "Split"
	},
	"Before": {
		"Patterns": [
			"#000000",
			"#ffffff"
		],
		"ShowMS": 100,
		"TransitionMS": 0,
		"_type": "Loop"
	},
	"OffsetMS": 4000,
	"TransitionMS": 3000,
	"_type": "Transition"
}`

const animCercueil = `
{
	"After": {
		"Patterns": [
			"#000000",
			"#00ff00",
			"#00ff00",
			"#00ff00",
			"#000000",
			"#ffa500",
			"#ffa500",
			"#ffa500"
		],
		"ShowMS": 500,
		"TransitionMS": 1000,
		"_type": "Loop"
	},
	"Before": {
		"Patterns": [
			"#000000",
			"#ffffff"
		],
		"ShowMS": 100,
		"TransitionMS": 0,
		"_type": "Loop"
	},
	"OffsetMS": 4000,
	"TransitionMS": 3000,
	"_type": "Transition"
}`

const animBouilloir = `
{
	"After": {
		"Curve": "ease-out",
		"Patterns": [
			{
				"Patterns": [
					{
						"Child": {
							"Child": {
								"Curve": "direct",
								"Left": "#ffa900",
								"Right": "#000000",
								"_type": "Gradient"
							},
							"Length": 16,
							"Offset": 0,
							"_type": "Subset"
						},
						"MovePerHour": 108000,
						"_type": "Rotate"
					},
					{
						"_type": "Aurore"
					}
				],
				"_type": "Add"
			},
			{
				"Patterns": [
					{
						"_type": "Aurore"
					},
					{
						"C": "#ffffff",
						"_type": "NightStars"
					}
				],
				"_type": "Add"
			}
		],
		"ShowMS": 10000,
		"TransitionMS": 5000,
		"_type": "Loop"
	},
	"Before": {
		"Patterns": [
			"#000000",
			"#ffffff"
		],
		"ShowMS": 100,
		"TransitionMS": 0,
		"_type": "Loop"
	},
	"OffsetMS": 4000,
	"TransitionMS": 3000,
	"_type": "Transition"
}
`
