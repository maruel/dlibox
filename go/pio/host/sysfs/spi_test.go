// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package sysfs

import "log"

func ExampleEnumerateSPI() {
	numbers, err := EnumerateSPI()
	if err != nil {
		log.Fatalf("failed to enumerate SPI buses: %v", err)
	}
	if len(numbers) == 0 {
		log.Fatalf("no SPI bus found")
	}
	bus, err := NewSPI(numbers[0][0], numbers[0][1], 0)
	if err != nil {
		log.Fatalf("failed to open SPI: %v", err)
	}
	defer bus.Close()

	// Use bus.
}
