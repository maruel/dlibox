// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package controller

/* This deadlocks.
func TestHalloween(t *testing.T) {
	b := msgbus.New()
	c := halloweenRule{}
	c.ResetDefault()
	c.Enabled = true
	c.Modes = map[string]halloweenState{"foo/1": incoming}
	c.Cmds = map[halloweenState][]rules.Command{incoming: {{Topic: "bar", Payload: "1"}}}
	h, err := initHalloween(msgbus.RebaseSub(msgbus.RebasePub(b, "foo"), "foo"), &c)
	if err != nil {
		t.Fatal(err)
	}
	if err := b.Publish(msgbus.Message{Topic: "foo/1", Payload: []byte("1")}, msgbus.ExactlyOnce); err != nil {
		t.Fatal(err)
	}
	// TODO(maruel): Settle wasn't implement in a concurrent safe manner.
	// b.Settle()
	// This test is non-deterministic.
	//if incoming != h.state { t.Fatal(...) }
	if err := h.Close(); err != nil {
		t.Fatal(err)
	}
}
*/
