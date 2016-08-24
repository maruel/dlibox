// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Package pio is a peripheral I/O library. It contains buses, devices, and
// fakes.
//
// (Temporary name, hopefully upstreamable)
//
//   - pio/buses contains implementations of ports to connect devices to, i.e.
//     i²c, spi, IR, etc. 'buses' contains the interfaces and subpackages
//     contain contain concrete types.
//   - pio/devices contains devices that are connected to a bus, i.e. ssd1306
//     (display controller), bm280 (environmental sensor), etc. 'devices'
//     contains the interfaces and subpackages
//     contain contain concrete types.
//   - pio/fakes contains non-hardware fakes, like a fake SPI bus, a fake
//     APA102 LEDs strip, etc.
package pio
