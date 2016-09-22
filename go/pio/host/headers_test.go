// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package host

import (
	"testing"

	"github.com/maruel/dlibox/go/pio/host/hosttest"
	"github.com/maruel/dlibox/go/pio/protocols/gpio"
	"github.com/maruel/dlibox/go/pio/protocols/pins"
)

func TestPinsHeaders(t *testing.T) {
	if len(allHeaders) != len(PinsHeaders()) {
		t.Fail()
	}
}

func TestPinIsConnected(t *testing.T) {
	if !PinIsConnected(pins.V3_3) {
		t.Fail()
	}
	if PinIsConnected(pins.V5) {
		t.Fail()
	}
	if !PinIsConnected(gpio2) {
		t.Fail()
	}
}

//

var (
	gpio2 = &hosttest.Pin{Name: "GPIO2", Num: 2, Fn: "I2C1_SDA"}
	gpio3 = &hosttest.Pin{Name: "GPIO3", Num: 3, Fn: "I2C1_SCL"}
)

func init() {
	allHeaders = map[string][][]gpio.PinIO{
		"P1": {
			{pins.GROUND, pins.V3_3},
			{gpio2, gpio3},
		},
	}
}
