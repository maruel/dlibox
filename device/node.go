// Copyright 2017 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package device

import (
	"fmt"
	"io"
	"reflect"

	"github.com/maruel/dlibox/nodes"
)

// NodeBase is the base type for all kind of supported nodes.
type NodeBase struct {
	id   nodes.ID
	name string
	typ  nodes.Type
}

func (n *NodeBase) String() string {
	return fmt.Sprintf("%s(%s/%s)", n.typ, n.id, n.name)
}

// Close is the default implementation that does nothing.
func (n *NodeBase) Close() error {
	return nil
}

func (n *NodeBase) wrap(err error) error {
	return fmt.Errorf("%s: %v", n, err)
}

type nodeDev interface {
	fmt.Stringer
	io.Closer
}

// dev is the device (nodes).
//
// The device doesn't store it, it's stored on the MQTT server.
type dev struct {
	nodes map[nodes.ID]nodeDev
}

var knownTypes = map[interface{}]interface{}{
	&nodes.Anim1D{}:  &anim1DDev{},
	&nodes.Button{}:  &buttonDev{},
	&nodes.Display{}: &displayDev{},
	&nodes.IR{}:      &irDev{},
	&nodes.PIR{}:     &pirDev{},
	&nodes.Sound{}:   &soundDev{},
}

var typesMap = map[reflect.Type]reflect.Type{}

func init() {
	for k, v := range knownTypes {
		typesMap[reflect.TypeOf(k).Elem()] = reflect.TypeOf(v).Elem()
	}

	// Verification.
	for t, r := range nodes.TypesMap {
		v, ok := typesMap[r]
		if !ok {
			panic(fmt.Sprintf("missing type %s", t))
		}
		if v.Field(0).Name != "NodeBase" {
			panic(fmt.Sprintf("NodeBase must be the first element in %s", v.Name()))
		}
		if v.Field(1).Name != "Cfg" {
			panic(fmt.Sprintf("Cfg must be the second element in %s", v.Name()))
		}
	}
}

// genNodeDev returns a nodeDev for a given.
func genNodeDev(id nodes.ID, n *nodes.Node) (nodeDev, error) {
	r, ok := typesMap[reflect.TypeOf(n.Config).Elem()]
	if !ok {
		return nil, fmt.Errorf("unknown type for %T", n.Config)
	}
	v := reflect.New(r)
	e := v.Elem()
	e.Field(0).Set(reflect.ValueOf(NodeBase{id: id, name: n.Name, typ: n.Type()}))
	e.Field(1).Set(reflect.ValueOf(n.Config))
	/*
		switch v := n.Config.(type) {
		case *nodes.Anim1D:
			d.nodes[id] = &anim1DDev{NodeBase: b, cfg: v}
		case *nodes.Button:
			d.nodes[id] = &buttonDev{nodeBase: b, cfg: v}
		case *nodes.Display:
			d.nodes[id] = &displayDev{nodeBase: b, cfg: v}
		case *nodes.IR:
			d.nodes[id] = &irDev{nodeBase: b, cfg: v}
		case *nodes.PIR:
			d.nodes[id] = &pirDev{nodeBase: b, cfg: v}
		case *nodes.Sound:
			d.nodes[id] = &soundDev{nodeBase: b, cfg: v}
		default:
			pubErr(dbus, "failed to initialize: unknown node %q: %T", id, n)
			return fmt.Errorf("unknown node %q: %T", id, n)
		}
	*/
	return v.Interface().(nodeDev), nil
}
