// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"testing"

	"github.com/maruel/dlibox/go/modules"
	"github.com/maruel/ut"
)

func TestHalloween(t *testing.T) {
	b := modules.LocalBus{}
	config := Halloween{}
	config.ResetDefault()
	config.Enabled = true
	config.Modes = map[string]State{"foo/1": Incoming}
	config.Cmds = map[State][]Command{Incoming: {{"bar", "1"}}}
	h, err := initHalloween(modules.Rebase(&b, "foo"), &config)
	ut.AssertEqual(t, nil, err)
	ut.AssertEqual(t, nil, b.Publish(modules.Message{"foo/1", []byte("1")}, modules.ExactlyOnce, false))
	b.Settle()
	// This test is non-deterministic.
	//ut.AssertEqual(t, Incoming, h.state)
	ut.AssertEqual(t, nil, h.Close())
}
