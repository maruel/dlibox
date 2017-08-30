// Copyright 2017 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Package nodes describes the shared nodes definition language between the
// controller and the devices.
package nodes

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

// Validator ensures a configuration is valid.
type Validator interface {
	Validate() error
}

// ID is a valid homie ID.
type ID string

// Validate implements Validator.
func (i ID) Validate() error {
	if !re.MatchString(string(i)) {
		return fmt.Errorf("invalid id: %q", i)
	}
	return nil
}

// Dev is the configuration for all the nodes on a device.
//
// The Dev is the micro-computer (a Rasberry Pi, C.H.I.P., ESP8266) that
// exposes nodes.
//
// The controller stores this data on the MQTT server upon startup and
// upon configuration update.
//
// The device fetches this data from the MQTT server upon startup. So when a
// device configuration occurs, the device should restart. This is really just
// restarting the Go process so this is nearly instant.
type Dev struct {
	// Name is the display name of this nodes collection: the device.
	Name  string
	Nodes map[ID]*Node
}

// ToSerialized is used by the controller to serialize the device description.
//
// Serialization should never fail.
func (d *Dev) ToSerialized() *SerializedDev {
	nds := &SerializedDev{Name: d.Name}
	for id, n := range d.Nodes {
		c, err := json.Marshal(n.Config)
		if err != nil {
			panic(err)
		}
		nds.Nodes[id] = &SerializedNode{
			Name:       n.Name,
			Type:       n.Type(),
			Properties: n.Config.toProperties(),
			Config:     c,
		}
	}
	return nds
}

// Validate implements Validator.
func (d *Dev) Validate() error {
	if len(d.Name) == 0 {
		return errors.New("dev: missing Name")
	}
	// TODO(maruel): Maybe do it in order to have deterministic result?
	for id, node := range d.Nodes {
		if err := id.Validate(); err != nil {
			return fmt.Errorf("dev %s: %v", d.Name, err)
		}
		if err := node.Validate(); err != nil {
			return fmt.Errorf("dev %s: node %s: %v", d.Name, id, err)
		}
	}
	return nil
}

// Node is the descriptor for a configured node on a device.
type Node struct {
	// Name is the display name of this node.
	Name   string
	Config NodeCfg
}

// Type returns the node's Type.
func (n *Node) Type() Type {
	return Type(strings.ToLower(reflect.TypeOf(n.Config).Elem().Name()))
}

// Validate implements Validator.
func (n *Node) Validate() error {
	if len(n.Name) == 0 {
		return errors.New("node: missing Name")
	}
	if err := n.Type().Validate(); err != nil {
		//return err
		return fmt.Errorf("node %s: unknown Config %T", n.Name, n.Config)
	}
	return n.Config.Validate()
}

// Serialized form.

// SerializedDev is the serialized form of Dev as stored on the MQTT server.
type SerializedDev struct {
	Name  string `json:"$name"`
	Nodes map[ID]*SerializedNode
}

// ToDev is used by the device to deserialize the configuration.
func (s *SerializedDev) ToDev() (*Dev, error) {
	// Parse every nodes, return a processed config.
	d := &Dev{Name: s.Name, Nodes: map[ID]*Node{}}
	for id, n := range s.Nodes {
		nd, err := n.toNode(id)
		if err != nil {
			return nil, err
		}
		d.Nodes[id] = nd
	}
	return d, nil
}

// SerializedNode is loosely based on
// https://github.com/marvinroger/homie#node-attributes
//
// It is the serialized form of Node.
type SerializedNode struct {
	Name       string          `json:"$name"`
	Type       Type            `json:"$type"`
	Properties map[ID]Property `json:"$properties"`
	Config     []byte          `json:"$config"`
}

func (s *SerializedNode) toNode(id ID) (*Node, error) {
	r := TypesMap[s.Type]
	if r == nil {
		return nil, fmt.Errorf("node %s: unknown type %s", id, s.Type)
	}
	v := reflect.New(r).Interface().(NodeCfg)
	if err := json.Unmarshal(s.Config, v); err != nil {
		return nil, fmt.Errorf("node %s: failed to unmarshal config: %v", id, err)
	}
	return &Node{Name: s.Name, Config: v}, nil
}

// Property defines one property of a node.
//
// A Node can have multiple properties. For example a buzzer could have separate
// knob for frequency and intensity.
type Property struct {
	Unit     string `json:"$unit"`
	DataType string `json:"$datatype"`
	Format   string `json:"$format"`
	Settable bool   `json:"$settable"`
}

//

var re = regexp.MustCompile(`^[a-z0-9][a-z0-9\-]+$`)
