// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package sysfs

import (
	"fmt"
	"log"

	"github.com/maruel/dlibox/go/pio/protocols/gpio"
)

func ExampleLEDByName() {
	// Commented out due to cycle import.
	//if _, err := host.Init(); err != nil {
	//	log.Fatalf("failed to initialize pio: %v", err)
	//}
	for _, led := range LEDs {
		fmt.Printf("- %s: %s\n", led, led.Function())
	}
	led, err := LEDByName("LED0")
	if err != nil {
		log.Fatalf("failed to find LED: %v", err)
	}
	led.Out(gpio.Low)
}
