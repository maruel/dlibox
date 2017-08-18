// Copyright 2017 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package device

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/maruel/dlibox/shared"
	"github.com/maruel/interrupt"
	"github.com/maruel/msgbus"
)

func TestMain(t *testing.T) {
	if !testing.Verbose() {
		log.SetOutput(ioutil.Discard)
		defer log.SetOutput(os.Stderr)
	}
	b := msgbus.New()
	retained(b, "dlibox/$online", "true")
	d := msgbus.RebasePub(b, "dlibox/"+shared.Hostname())
	retained(d, "reset", "true")
	retained(d, "$name", "foo")
	retained(d, "$ignored", "really")
	retained(d, "node1", "true")
	retained(d, "node1/$name", "The node")
	retained(d, "node1/$type", "button")
	//retained(d, "node1/$", "button")
	interrupt.Set()
	Main("", b, 0)
}
