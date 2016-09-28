// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package sysfs

import "log"

func ExampleEnumerateI2C() {
	numbers, err := EnumerateI2C()
	if err != nil {
		log.Fatalf("failed to enumerate I²C buses: %v", err)
	}
	if len(numbers) == 0 {
		log.Fatalf("no I²C bus found")
	}
	bus, err := NewI2C(numbers[0])
	if err != nil {
		log.Fatalf("failed to open I²C: %v", err)
	}
	defer bus.Close()

	// Use bus.
}
