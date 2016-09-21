// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package pins

import (
	"testing"

	"github.com/maruel/dlibox/go/pio/host"
	"github.com/maruel/dlibox/go/pio/host/hosttest"
)

func TestAll(t *testing.T) {
	if len(all) != len(All()) {
		t.Fail()
	}
}

func TestByNumber(t *testing.T) {
	if ByNumber(1) != nil {
		t.Fatal("1 exist")
	}
	if ByNumber(2) != gpio2 {
		t.Fatal("2 missing")
	}
}

func TestByName(t *testing.T) {
	if ByName("GPIO0") != nil {
		t.Fatal("GPIO0 doesn't exist")
	}
	if ByName("GPIO2") != gpio2 {
		t.Fatal("GPIO2 should have been found")
	}
}

func TestByFunction(t *testing.T) {
	if ByFunction("SPI1_MOSI") != nil {
		t.Fatal("spi doesn't exist")
	}
	if ByFunction("I2C1_SDA") != gpio2 {
		t.Fatal("I2C1_SDA should have been found")
	}
}

//

var (
	gpio2 = &hosttest.Pin{Name: "GPIO2", Num: 2, Fn: "I2C1_SDA"}
	gpio3 = &hosttest.Pin{Name: "GPIO3", Num: 3, Fn: "I2C1_SCL"}
)

func init() {
	all = []host.PinIO{gpio2, gpio3}
	byFunction = map[string]host.PinIO{
		gpio2.Function(): gpio2,
		gpio3.Function(): gpio3,
	}
}
