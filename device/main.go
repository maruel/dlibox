// Copyright 2017 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Package device implements the dlibox device.
//
// When running as a device, it connects to the MQTT server to listen to
// commands from the controller.
package device

import (
	"errors"
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
	shared.RetainedStr(dbus, "$online", "initializing")

	// Initialize periph.
	state, err := host.Init()
	if err != nil {
		shared.RetainedStr(dbus, "$online", err.Error())
		return err
	}

	if port != 0 {
		if err = webServer(server, port); err != nil {
			return err
		}
	}

	// Wait until the controller is up and running. This ensures that the node
	// are correctly published.
	c, err := bus.Subscribe("$online", msgbus.ExactlyOnce)
	if err != nil {
		return err
	}
	for {
		select {
		case msg, ok := <-c:
			if !ok {
				return errors.New("MQTT server died")
			}
			s := string(msg.Payload)
			if s == "true" {
				goto done
			}
			log.Printf("$online: %q", s)
		case <-interrupt.Channel:
			break
		}
	}
done:
	if !interrupt.IsSet() {
		shared.RetainedStr(dbus, "$online", "initializing")
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
		shared.RetainedStr(dbus, "$online", "true")
	}
	return shared.WatchFile()
}

func getConfig(b msgbus.Bus) (*nodes.Dev, error) {
	msgs, err := msgbus.Retained(b, 10*time.Second, "$name", "$nodes")
	if err != nil {
		return nil, err
	}
	nds := nodes.SerializedDev{Name: string(msgs["$name"])}
	nodesID := string(msgs["nodes"])
	if len(nodesID) != 0 {
		for _, id := range strings.Split(nodesID, ",") {
			nodeID := nodes.ID(id)
			if err := nodeID.Validate(); err != nil {
				return nil, err
			}
			// TODO(maruel): Query all nodes concurrently to reduce the effect of round
			// trip latency.
			n, err := processNode(b, id)
			if err != nil {
				return nil, fmt.Errorf("node %q: %v", nodeID, err)
			}
			nds.Nodes[nodeID] = n
		}
	}
	return nds.ToDev()
}

func processNode(b msgbus.Bus, nodeID string) (*nodes.SerializedNode, error) {
	msgs, err := msgbus.Retained(b, 10*time.Second, "$name", "$properties", "$types")
	if err != nil {
		return nil, err
	}

	n := &nodes.SerializedNode{
		Name: string(msgs["$name"]),
		Type: nodes.Type(string(msgs["$type"])),
	}
	if err := n.Type.Validate(); err != nil {
		return nil, fmt.Errorf("node %q: %v", n.Name, err)
	}

	// TODO(maruel): Query concurrently.
	propnames := strings.Split(string(msgs["$properties"]), ",")
	for _, propname := range propnames {
		propID := nodes.ID(propname)
		if err := propID.Validate(); err != nil {
			return nil, fmt.Errorf("invalid property %s/%s: %v", nodeID, propID, err)
		}
		pm, err := msgbus.Retained(b, 10*time.Second, "$datatype", "$format", "$settable", "$unit")
		if err != nil {
			return nil, err
		}
		n.Properties[propID] = nodes.Property{
			Unit:     string(pm["$unit"]),
			DataType: string(pm["$datatype"]),
			Format:   string(pm["$format"]),
			Settable: string(pm["$settable"]) == "true",
		}
	}
	return n, nil
}

//

func pubErr(b msgbus.Bus, f string, arg ...interface{}) {
	msg := fmt.Sprintf(f, arg)
	log.Print(msg)
	if err := b.Publish(msgbus.Message{Topic: "$error", Payload: []byte(msg)}, msgbus.ExactlyOnce); err != nil {
		log.Print(err)
	}
}
