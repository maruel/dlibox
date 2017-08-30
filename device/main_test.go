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
	b := msgbus.Log(msgbus.New())
	shared.RetainedStr(b, "dlibox/$online", "true")
	d := msgbus.RebasePub(b, "dlibox/"+shared.Hostname())
	shared.RetainedStr(d, "reset", "true")
	shared.RetainedStr(d, "$name", "foo")
	shared.RetainedStr(d, "$nodes", "node1")
	shared.RetainedStr(d, "$ignored", "really")
	shared.RetainedStr(d, "node1", "true")
	shared.RetainedStr(d, "node1/$name", "The node")
	shared.RetainedStr(d, "node1/$type", "button")
	shared.RetainedStr(d, "node1/$unit", ".")
	shared.RetainedStr(d, "node1/$datatype", "bool")
	shared.RetainedStr(d, "node1/$format", ".")
	shared.RetainedStr(d, "node1/$settable", "false")
	//retained(d, "node1/$", "button")
	interrupt.Set()
	Main("", b, 0)
}
