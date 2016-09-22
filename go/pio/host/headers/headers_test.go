// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package headers

import (
	"testing"

	"github.com/maruel/dlibox/go/pio/protocols/gpio/gpiotest"
	"github.com/maruel/dlibox/go/pio/protocols/pins"
)

func TestAll(t *testing.T) {
	if len(allHeaders) != len(All()) {
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
	gpio2 = &gpiotest.Pin{Name: "GPIO2", Num: 2, Fn: "I2C1_SDA"}
	gpio3 = &gpiotest.Pin{Name: "GPIO3", Num: 3, Fn: "I2C1_SCL"}
)

func init() {
	allHeaders = map[string][][]pins.Pin{
		"P1": {
			{pins.GROUND, pins.V3_3},
			{gpio2, gpio3},
		},
	}
}
