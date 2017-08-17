// Copyright 2017 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Package nodes describes the shared nodes definition language between the
// controller and the devices.
package nodes

import (
	"encoding/json"
	"log"
	_ "net/http/pprof"
	"regexp"

	"github.com/maruel/dlibox/nodes/button"
	"github.com/maruel/dlibox/nodes/display"
	"github.com/maruel/dlibox/nodes/ir"
	"github.com/maruel/dlibox/nodes/leds"
	"github.com/maruel/dlibox/nodes/pir"
	"github.com/maruel/dlibox/nodes/sound"
)

// ID is a valid homie ID.
type ID string

// IsValid returns true if this is a valid id.
func (i ID) IsValid() bool {
	return re.MatchString(string(i))
}

// Dev is the configuration for all the nodes on a device.
//
// The device doesn't store it, it's stored on the MQTT server.
type Dev struct {
	// Name is the display name of this nodes collection, the device.
	Name     string // $name
	Buttons  map[ID]button.Dev
	Displays map[ID]display.Dev
	IRs      map[ID]ir.Dev
	LEDs     map[ID]leds.Dev
	PIRs     map[ID]pir.Dev
	Sound    map[ID]sound.Dev
}

// ToNodes is used by the controller to serialize the device description.
func (d *Dev) ToNodes() Nodes {
	n := Nodes{}
	for id, v := range d.Buttons {
		c, _ := json.Marshal(v)
		n[id] = Node{
			Name: v.Name,
			Type: Button,
			Properties: map[ID]Property{
				// TODO(maruel): Support double-click.
				"button": {DataType: "boolean"},
			},
			Config: c,
		}
	}
	for id, v := range d.Displays {
		c, _ := json.Marshal(v)
		n[id] = Node{
			Name: v.Name,
			Type: SSD1306,
			Properties: map[ID]Property{
				"markee": {
					DataType: "string",
					Settable: true,
				},
				"content": {
					DataType: "string",
					Settable: true,
				},
			},
			Config: c,
		}
	}
	for id, v := range d.Displays {
		c, _ := json.Marshal(v)
		n[id] = Node{
			Name: v.Name,
			Type: IR,
			Properties: map[ID]Property{
				"ir": {DataType: "string"},
			},
			Config: c,
		}
	}
	for id, v := range d.LEDs {
		c, _ := json.Marshal(v)
		n[id] = Node{
			Name: v.Name,
			Type: Anim1D,
			Properties: map[ID]Property{
				// A pattern.
				"anim1d": {
					DataType: "string",
					Settable: true,
				},
			},
			Config: c,
		}
	}
	for id, v := range d.PIRs {
		c, _ := json.Marshal(v)
		n[id] = Node{
			Name: v.Name,
			Type: PIR,
			Properties: map[ID]Property{
				"pir": {DataType: "boolean"},
			},
			Config: c,
		}
	}
	for id, v := range d.Sound {
		c, _ := json.Marshal(v)
		n[id] = Node{
			Name: v.Name,
			Type: Sound,
			Properties: map[ID]Property{
				// An URL.
				"speakers": {
					DataType: "string",
					Settable: true,
				},
			},
			Config: c,
		}
	}
	return n
}

// Nodes is the serialized form of Dev as stored on the MQTT server.
type Nodes map[ID]Node

// ToDev is used by the device to deserialize the configuration.
func (n Nodes) ToDev() *Dev {
	// Parse every nodes, return a processed config.
	d := &Dev{}
	for id, node := range n {
		switch node.Type {
		case Anim1D:
			v := leds.Dev{}
			if err := json.Unmarshal(node.Config, &v); err != nil {
				log.Printf("failed to unmarshal %s", id)
				continue
			}
			d.LEDs[id] = v

		case Button:
			v := button.Dev{}
			if err := json.Unmarshal(node.Config, &v); err != nil {
				log.Printf("failed to unmarshal %s", id)
				continue
			}
			d.Buttons[id] = v

		case SSD1306:
			v := display.Dev{}
			if err := json.Unmarshal(node.Config, &v); err != nil {
				log.Printf("failed to unmarshal %s", id)
				continue
			}
			d.Displays[id] = v

		case IR:
			v := ir.Dev{}
			if err := json.Unmarshal(node.Config, &v); err != nil {
				log.Printf("failed to unmarshal %s", id)
				continue
			}
			d.IRs[id] = v

		case PIR:
			v := pir.Dev{}
			if err := json.Unmarshal(node.Config, &v); err != nil {
				log.Printf("failed to unmarshal %s", id)
				continue
			}
			d.PIRs[id] = v

		case Sound:
			v := sound.Dev{}
			if err := json.Unmarshal(node.Config, &v); err != nil {
				log.Printf("failed to unmarshal %s", id)
				continue
			}
			d.Sound[id] = v

		default:
			log.Printf("id %s unknown type %s", id, node.Type)
		}
	}
	return d
}

// Type is all the known node types.
type Type string

const (
	Anim1D  Type = "anim1d"
	Button  Type = "button"
	SSD1306 Type = "ssd1306"
	IR      Type = "ir"
	PIR     Type = "pir"
	Sound   Type = "sound"
)

// IsValid returns true if the node type is known.
func (t Type) IsValid() bool {
	switch t {
	case Anim1D, Button, SSD1306, IR, PIR, Sound:
		return true
	default:
		return false
	}
}

// Node is loosely based on
// https://github.com/marvinroger/homie#node-attributes
type Node struct {
	Name       string          `json:"$name"`
	Type       Type            `json:"$type"`
	Properties map[ID]Property `json:"$properties`
	Config     []byte          `json:"$config`
}

// Property defines one property of a node.
type Property struct {
	Unit     string `json:"$unit`
	DataType string `json:"$datatype`
	Format   string `json:"$format`
	Settable bool   `json:"$settable`
}

//

var re = regexp.MustCompile(`^[a-z0-9][a-z0-9\-]+$`)
