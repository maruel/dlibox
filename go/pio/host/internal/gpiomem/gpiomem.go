// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package gpiomem

// Mem is the memory mapped CPU I/O registers.
type Mem struct {
	Uint8  []uint8
	Uint32 []uint32
}
