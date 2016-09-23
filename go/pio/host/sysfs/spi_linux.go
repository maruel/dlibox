// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package sysfs

import "github.com/maruel/dlibox/go/pio/drivers"

// NewSPI opens a SPI bus via its sysfs interface as described at
// https://www.kernel.org/doc/Documentation/spi/spidev.
//
// busNumber is the bus number as exported by sysfs. For example if the path is
// /dev/spidev0.1, busNumber should be 0 and chipSelect should be 1.
//
// speed can either be 0 for the default speed or should be in the high Khz or
// low Mhz range, it's a good idea to start at 4000000 (4Mhz) and go upward as
// long as the signal is good.
//
// Default configuration is Mode3 and 8 bits.
func NewSPI(busNumber, chipSelect int, speed int64) (*SPI, error) {
	return newSPI(busNumber, chipSelect, speed)
}

func init() {
	drivers.Register(&driverSPI{})
}
