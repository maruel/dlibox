// Copyright 2017 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Package device implements the dlibox device.
//
// When running as a device, it connects to the MQTT server to listen to
// commands from the controller.
package device

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/maruel/dlibox/nodes"
	"github.com/maruel/dlibox/shared"
	"github.com/maruel/interrupt"
	"github.com/maruel/msgbus"
	"periph.io/x/periph/host"
)

// Main is the main function when running as a device (a node).
func Main(server string, bus msgbus.Bus, port int) error {
	log.Printf("device.Main(%s, ..., %d)", server, port)

	// Everything is under the namespace "dlibox/"
	bus = msgbus.RebasePub(msgbus.RebaseSub(bus, "dlibox"), "dlibox")
	root := shared.Hostname()
	dbus := msgbus.RebaseSub(msgbus.RebasePub(bus, root), root)
	retained(dbus, "$online", "initializing")

	// Initialize periph.
	state, err := host.Init()
	if err != nil {
		retained(dbus, "$online", err.Error())
		return err
	}

	if port != 0 {
		if err = webServer(server, port); err != nil {
			return err
		}
	}

	// Poll until the controller is up and running. This ensures that the node
	// are correctly published.
	// TODO(maruel): Polling is not useful, Subscribe($online) will work.
	polled := false
	for !interrupt.IsSet() {
		m, err := bus.Retained("$online")
		if err == nil && len(m) == 1 {
			s := string(m[0].Payload)
			if s == "true" {
				break
			}
			log.Printf("$online: %q", s)
		} else {
			log.Printf("Failed to get retained message $online: %v", err)
		}
		if !polled {
			retained(dbus, "$online", "waiting for "+server)
			polled = true
		}
		time.Sleep(time.Second)
	}
	if polled {
		retained(dbus, "$online", "initializing")
	}

	// TODO(maruel): Uses modified Homie convention. The main modification is
	// that it is the controller that decides which nodes to expose ($nodes), not
	// the device itself. This means no need for configuration on the device
	// itself, assuming all devices run all the same code.
	// https://github.com/marvinroger/homie#device-attributes
	shared.InitState(dbus, state)

	if c, err := dbus.Subscribe("reset", msgbus.ExactlyOnce); err == nil {
		go func() {
			<-c
			// Exiting the process means it'll restart normally and will initialize
			// properly.
			interrupt.Set()
		}()
	} else {
		log.Printf("failed to subscribe to reset: %v", err)
	}

	cfg, err := getConfig(dbus)
	if err != nil {
		pubErr(dbus, "failed to initialize: %v", err)
		return err
	}
	d := dev{nodes: map[nodes.ID]nodeDev{}}
	for id, n := range cfg.Nodes {
		n, err := genNodeDev(id, n)
		if err != nil {
			pubErr(dbus, "failed to initialize: unknown node %q: %v", id, err)
			return fmt.Errorf("unknown node %q: %v", id, err)
		}
		d.nodes[id] = n
	}

	if !interrupt.IsSet() {
		retained(dbus, "$online", "true")
	}
	return shared.WatchFile()
}

func getConfig(b msgbus.Bus) (*nodes.Dev, error) {
	msgs, err := b.Retained("#")
	if err != nil {
		return nil, err
	}

	// Unpack all the node definitions.
	defs := map[nodes.ID]map[string][]byte{}
	nds := nodes.SerializedDev{}
	for _, msg := range msgs {
		if msg.Topic == "$name" {
			nds.Name = string(msg.Payload)
			continue
		}
		if strings.HasPrefix(msg.Topic, "$") {
			// Device description.
			continue
		}
		parts := strings.SplitN(msg.Topic, "/", 2)
		if len(parts) != 2 {
			// Node value.
			continue
		}
		nodeID := nodes.ID(parts[0])
		if err := nodeID.Validate(); err != nil {
			return nil, err
		}
		if _, ok := defs[nodeID]; !ok {
			defs[nodeID] = map[string][]byte{}
		}
		defs[nodeID][parts[1]] = msg.Payload
	}

	// Process each node.
	for nodeID, nodedef := range defs {
		n, err := processNode(nodedef)
		if err != nil {
			return nil, fmt.Errorf("node %q: %v", nodeID, err)
		}
		nds.Nodes[nodeID] = n
	}
	return nds.ToDev()
}

func processNode(nodedef map[string][]byte) (*nodes.SerializedNode, error) {
	n := &nodes.SerializedNode{
		Name: string(nodedef["$name"]),
		Type: nodes.Type(string(nodedef["$type"])),
	}
	if err := n.Type.Validate(); err != nil {
		return nil, fmt.Errorf("node %q: %v", n.Name, err)
	}
	propdefs := map[nodes.ID]map[string][]byte{}
	for topic, payload := range nodedef {
		if strings.HasPrefix(topic, "$") {
			continue
		}
		parts := strings.SplitN(topic, "/", 2)
		if len(parts) != 2 {
			continue
		}
		propID := nodes.ID(parts[0])
		if err := propID.Validate(); err != nil {
			return nil, fmt.Errorf("invalid property %s/%s: %v", propID, err)
		}
		if _, ok := propdefs[propID]; !ok {
			propdefs[propID] = map[string][]byte{}
		}
		propdefs[propID][parts[1]] = payload
	}

	propnames := strings.Split(string(nodedef["$properties"]), ",")
	for _, propname := range propnames {
		propID := nodes.ID(propname)
		if err := propID.Validate(); err != nil {
			return nil, fmt.Errorf("invalid property %s/%s: %v", propID, err)
		}
		propdef := propdefs[propID]
		n.Properties[propID] = nodes.Property{
			Unit:     string(propdef["$unit"]),
			DataType: string(propdef["$datatype"]),
			Format:   string(propdef["$format"]),
			Settable: string(propdef["$settable"]) == "true",
		}
	}
	return n, nil
}

//

func pubErr(b msgbus.Bus, f string, arg ...interface{}) {
	msg := fmt.Sprintf(f, arg)
	log.Print(msg)
	b.Publish(msgbus.Message{Topic: "$error", Payload: []byte(msg)}, msgbus.ExactlyOnce, false)
}

func retained(b msgbus.Bus, topic, payload string) {
	retainedBytes(b, topic, []byte(payload))
}

func retainedBytes(b msgbus.Bus, topic string, payload []byte) {
	if err := b.Publish(msgbus.Message{topic, payload}, msgbus.MinOnce, true); err != nil {
		log.Printf("Failed to publish %s: %v", topic, err)
	}
}
