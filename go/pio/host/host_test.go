// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package host

import (
	"fmt"
	"log"
	"testing"

	"github.com/maruel/dlibox/go/pio/devices"
	"github.com/maruel/dlibox/go/pio/devices/bme280"
	"github.com/maruel/dlibox/go/pio/protocols/i2c"
	"github.com/maruel/dlibox/go/pio/protocols/spi"
)

func ExampleInit() {
	state, err := Init()
	if err != nil {
		log.Fatalf("failed to initialize pio: %v", err)
	}
	fmt.Printf("Using drivers:\n")
	for _, driver := range state.Loaded {
		fmt.Printf("  - %s\n", driver.String())
	}
	fmt.Printf("Drivers skipped:\n")
	for _, driver := range state.Skipped {
		fmt.Printf("  - %s\n", driver.String())
	}
	if len(state.Failed) != 0 {
		// Having drivers failing to load may not require process termination.
		fmt.Printf("Drivers failed to load:\n")
		for _, f := range state.Failed {
			fmt.Printf("  - %s: %v\n", f.D, f.Err)
		}
	}

	// Use pins, buses, devices, etc.
}

func ExampleNewI2CAuto() {
	if _, err := Init(); err != nil {
		log.Fatalf("failed to initialize pio: %v", err)
	}
	bus, err := NewI2CAuto()
	if err != nil {
		log.Fatalf("failed to initialize I²C: %v", err)
	}
	defer bus.Close()
	if p, ok := bus.(i2c.Pins); ok {
		log.Printf("Using pins SCL: %s  SDA: %s", p.SCL(), p.SDA())
	}

	// Reads off a weather sensor.
	dev, err := bme280.NewI2C(bus, bme280.O2x, bme280.O2x, bme280.O2x, bme280.S20ms, bme280.FOff)
	if err != nil {
		log.Fatalf("failed to initialize bme280: %v", err)
	}
	env := devices.Environment{}
	dev.Read(&env)
	fmt.Printf("%8s %10s %9s\n", env.Temperature, env.Pressure, env.Humidity)
}

func ExampleNewSPIAuto() {
	if _, err := Init(); err != nil {
		log.Fatalf("failed to initialize pio: %v", err)
	}
	bus, err := NewSPIAuto()
	if err != nil {
		log.Fatalf("failed to initialize SPI: %v", err)
	}
	defer bus.Close()
	if p, ok := bus.(spi.Pins); ok {
		log.Printf("Using pins CLK: %s  MOSI: %s  MISO: %s  CS: %s", p.CLK(), p.MOSI(), p.MISO(), p.CS())
	}

	// Reads off a weather sensor.
	dev, err := bme280.NewSPI(bus, bme280.O2x, bme280.O2x, bme280.O2x, bme280.S20ms, bme280.FOff)
	if err != nil {
		log.Fatalf("failed to initialize bme280: %v", err)
	}
	env := devices.Environment{}
	dev.Read(&env)
	fmt.Printf("%8s %10s %9s\n", env.Temperature, env.Pressure, env.Humidity)
}

func TestInit(t *testing.T) {
	if _, err := Init(); err != nil {
		t.Fatalf("failed to initialize pio: %v", err)
	}
}
