// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package gpiomem

// OpenGPIO returns a CPU specific memory mapping of the CPU I/O registers using
// /dev/gpiomem.
//
// /dev/gpiomem is only supported on Raspbian Jessie via a specific kernel
// driver.
func OpenGPIO() (*Mem, error) {
	return openGPIO()
}

// OpenMem returns a memory mapped view of arbitrary kernel memory range using
// /dev/mem.
func OpenMem(base uint64) (*Mem, error) {
	return openMem(base)
}
