// Copyright 2017 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package device

import (
	"fmt"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"strings"
	"time"

	"github.com/maruel/dlibox/go/modules/nodes"
	"github.com/maruel/dlibox/go/modules/shared"
	"github.com/maruel/dlibox/go/msgbus"
	"github.com/maruel/interrupt"
	"github.com/maruel/serve-dir/loghttp"
	"periph.io/x/periph/host"
)

// Dev is the device (nodes).
//
// The device doesn't store it, it's stored on the MQTT server.
type Dev struct {
	Buttons  []*Button
	Displays []*Display
	LEDs     []*LEDs
	IRs      []IR
	PIRs     []*PIR
	Sound    []*Sound
}

func (d *Dev) Close() error {
	return nil
}

// Main is the main function when running as a device (a node).
func Main(server string, bus msgbus.Bus, port int) error {
	log.Printf("device.Main(%s, ..., %d)", server, port)
	// Initialize periph.
	state, err := host.Init()
	if err != nil {
		return err
	}

	if port != 0 {
		if err = webServer(server, port); err != nil {
			return err
		}
	}

	// Poll until the controller is up and running. This ensures that the node
	// are correctly published.
	for !interrupt.IsSet() {
		m, err := bus.Retained("$online")
		if err != nil || len(m) != 1 {
			log.Printf("Failed to get retained message $online: %v", err)
			time.Sleep(time.Second)
			continue
		}
		if s := string(m[0].Payload); s != "true" {
			log.Printf("$online: %q", s)
			time.Sleep(time.Second)
			continue
		}
		break
	}

	// TODO(maruel): Uses modified Homie convention. The main modification is
	// that it is the controller that decides which nodes to expose ($nodes), not
	// the device itself. This means no need for configuration on the device
	// itself, assuming all devices run all the same code.
	// https://github.com/marvinroger/homie#device-attributes
	root := shared.Hostname()
	rebased := msgbus.RebaseSub(msgbus.RebasePub(bus, root), root)
	shared.InitState(rebased, state)

	if c, err := rebased.Subscribe("reset", msgbus.MinOnce); err == nil {
		go func() {
			<-c
			// Exiting the process means it'll restart normally and will initialize
			// properly.
			interrupt.Set()
		}()
	} else {
		log.Printf("failed to subscribe to reset")
	}

	cfg := getConfig(rebased)
	if cfg == nil {
		return nil
	}
	dev := Dev{}
	/*
		dev.Buttons.init(cfg.Buttons)
		dev.Displays.init(cfg.Displays)
		dev.IRs.init(cfg.IRs)
		dev.PIRs.init(cfg.PIRs)
		dev.Sound.init(cfg.Sound)
	*/
	defer dev.Close()
	rebased.Publish(msgbus.Message{"$online", []byte("true")}, msgbus.MinOnce, true)
	return shared.WatchFile()
}

// webServer is the device's web server. It is quite simple.
func webServer(server string, port int) error {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	if server != "" {
		http.DefaultServeMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "http://"+server, 302)
		})
	}
	s := http.Server{
		Addr:           ln.Addr().String(),
		Handler:        &loghttp.Handler{Handler: http.DefaultServeMux},
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 16,
	}
	go s.Serve(ln)
	log.Printf("Visit: http://%s:%d/debug/pprof for debugging", shared.Hostname(), port)
	return nil
}

func getConfig(b msgbus.Bus) *nodes.Dev {
	msgs, err := b.Retained("#")
	if err != nil {
		// If error, retry.
		log.Printf("Failed to get retained nodes: %v", err)
		return nil
	}

	// Unpack all the node definitions.
	defs := map[nodes.ID]map[string][]byte{}
	name := ""
	for _, msg := range msgs {
		if msg.Topic == "$name" {
			name = string(msg.Payload)
			continue
		}
		if strings.HasPrefix(msg.Topic, "$") {
			// Device description
			continue
		}
		parts := strings.SplitN(msg.Topic, "/", 2)
		if len(parts) != 2 {
			// Node value
			continue
		}
		nodeID := nodes.ID(parts[0])
		if !nodeID.IsValid() {
			pubErr(b, "invalid node %q", nodeID)
			continue
		}
		if _, ok := defs[nodeID]; !ok {
			defs[nodeID] = map[string][]byte{}
		}
		defs[nodeID][parts[1]] = msg.Payload
	}

	nds := nodes.Nodes{}
	// Process each node.
	for nodeID, nodedef := range defs {
		n := nodes.Node{
			Name: string(nodedef["$name"]),
			Type: nodes.Type(string(nodedef["$type"])),
		}
		if !n.Type.IsValid() {
			pubErr(b, "invalid node %q type %s", n.Name, n.Type)
			continue
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
			if !propID.IsValid() {
				pubErr(b, "invalid property %s/%s", nodeID, propID)
				continue
			}
			if _, ok := propdefs[propID]; !ok {
				propdefs[propID] = map[string][]byte{}
			}
			propdefs[propID][parts[1]] = payload
		}

		propnames := strings.Split(string(nodedef["$properties"]), ",")
		for _, propname := range propnames {
			propID := nodes.ID(propname)
			if !propID.IsValid() {
				pubErr(b, "invalid property %s/%s", nodeID, propID)
				continue
			}
			propdef := propdefs[propID]
			n.Properties[propID] = nodes.Property{
				Unit:     string(propdef["$unit"]),
				DataType: string(propdef["$datatype"]),
				Format:   string(propdef["$format"]),
				Settable: string(propdef["$settable"]) == "true",
			}
		}
		nds[nodeID] = n
	}
	d := nds.ToDev()
	d.Name = name
	return d
}

func pubErr(b msgbus.Bus, f string, arg ...interface{}) {
	msg := fmt.Sprintf(f, arg)
	log.Print(msg)
	b.Publish(msgbus.Message{Topic: "error", Payload: []byte(msg)}, msgbus.ExactlyOnce, false)
}
