// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package controller

import (
	"testing"

	"github.com/maruel/dlibox/controller/rules"
	"github.com/maruel/msgbus"
	"github.com/maruel/ut"
)

func TestHalloween(t *testing.T) {
	b := msgbus.New()
	c := halloweenRule{}
	c.ResetDefault()
	c.Enabled = true
	c.Modes = map[string]halloweenState{"foo/1": incoming}
	c.Cmds = map[halloweenState][]rules.Command{incoming: {{Topic: "bar", Payload: "1"}}}
	h, err := initHalloween(msgbus.RebaseSub(msgbus.RebasePub(b, "foo"), "foo"), &c)
	ut.AssertEqual(t, nil, err)
	ut.AssertEqual(t, nil, b.Publish(msgbus.Message{"foo/1", []byte("1")}, msgbus.BestEffort, false))
	// TODO(maruel): Settle wasn't implement in a concurrent safe manner.
	// b.Settle()
	// This test is non-deterministic.
	//ut.AssertEqual(t, incoming, h.state)
	ut.AssertEqual(t, nil, h.Close())
}
