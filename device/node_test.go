// Copyright 2017 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package device

import (
	"testing"

	"github.com/maruel/dlibox/nodes"
)

func TestGetNodeDev(t *testing.T) {
	n := &nodes.Node{
		Name:   "a node",
		Config: &nodes.IR{},
	}
	d, err := genNodeDev(nodes.ID("node1"), n)
	if err != nil {
		t.Fatal(err)
	}
	if err := d.Close(); err != nil {
		t.Fatal(err)
	}
}
