// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package headers

import (
	"testing"

	"github.com/maruel/dlibox/go/pio/host"
	"github.com/maruel/dlibox/go/pio/host/hosttest"
	"github.com/maruel/dlibox/go/pio/host/pins"
)

func TestAll(t *testing.T) {
	if len(all) != len(All()) {
		t.Fail()
	}
}

func TestIsConnected(t *testing.T) {
	if !IsConnected(pins.V3_3) {
		t.Fail()
	}
	if IsConnected(pins.V5) {
		t.Fail()
	}
	if !IsConnected(gpio2) {
		t.Fail()
	}
}

//

var (
	gpio2 = &hosttest.Pin{Name: "GPIO2", Num: 2, Fn: "I2C1_SDA"}
	gpio3 = &hosttest.Pin{Name: "GPIO3", Num: 3, Fn: "I2C1_SCL"}
)

func init() {
	all = map[string][][]host.PinIO{
		"P1": {
			{pins.GROUND, pins.V3_3},
			{gpio2, gpio3},
		},
	}
}
