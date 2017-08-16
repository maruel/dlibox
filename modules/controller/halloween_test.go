// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package controller

import (
	"testing"

	"github.com/maruel/dlibox/modules"
	"github.com/maruel/dlibox/msgbus"
	"github.com/maruel/ut"
)

func TestHalloween(t *testing.T) {
	b := msgbus.New()
	config := Halloween{}
	config.ResetDefault()
	config.Enabled = true
	config.Modes = map[string]State{"foo/1": Incoming}
	config.Cmds = map[State][]modules.Command{Incoming: {{"bar", "1"}}}
	h, err := initHalloween(msgbus.RebaseSub(msgbus.RebasePub(b, "foo"), "foo"), &config)
	ut.AssertEqual(t, nil, err)
	ut.AssertEqual(t, nil, b.Publish(msgbus.Message{"foo/1", []byte("1")}, msgbus.BestEffort, false))
	// TODO(maruel): Settle wasn't implement in a concurrent safe manner.
	// b.Settle()
	// This test is non-deterministic.
	//ut.AssertEqual(t, Incoming, h.state)
	ut.AssertEqual(t, nil, h.Close())
}
