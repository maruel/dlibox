// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package host

import (
	"testing"

	"github.com/maruel/dlibox/go/pio/host/internal/pins2"
	"github.com/maruel/dlibox/go/pio/protocols/gpio"
)

func TestPins(t *testing.T) {
	if len(pins2.All) != len(Pins()) {
		t.Fail()
	}
}

func TestPinByNumber(t *testing.T) {
	if PinByNumber(1) != nil {
		t.Fatal("1 exist")
	}
	if PinByNumber(2) != gpio2 {
		t.Fatal("2 missing")
	}
}

func TestPinByName(t *testing.T) {
	if PinByName("GPIO0") != nil {
		t.Fatal("GPIO0 doesn't exist")
	}
	if PinByName("GPIO2") != gpio2 {
		t.Fatal("GPIO2 should have been found")
	}
}

func TestPinByFunction(t *testing.T) {
	if PinByFunction("SPI1_MOSI") != nil {
		t.Fatal("spi doesn't exist")
	}
	if PinByFunction("I2C1_SDA") != gpio2 {
		t.Fatal("I2C1_SDA should have been found")
	}
}

//

func init() {
	pins2.All = []gpio.PinIO{gpio2, gpio3}
	pins2.ByFunction = map[string]gpio.PinIO{
		gpio2.Function(): gpio2,
		gpio3.Function(): gpio3,
	}
	setIR()
}
