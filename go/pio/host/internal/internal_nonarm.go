// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// +build !arm

package internal

// OS

// IsArmbian returns true if running on a Armbian distribution.
//
// http://www.armbian.com/
func IsArmbian() bool {
	return false
}

// IsRaspbian returns true if running on a Raspbian distribution.
//
// https://raspbian.org/
func IsRaspbian() bool {
	return false
}

// CPU

// IsBCM283x returns true if running on a Broadcom bcm283x based CPU.
func IsBCM283x() bool {
	return false
}

// IsAllwinner returns true if running on an Allwinner based CPU.
//
// https://en.wikipedia.org/wiki/Allwinner_Technology
func IsAllwinner() bool {
	return false
}

// Board

// IsRaspberryPi returns true if running on a raspberry pi board.
//
// https://www.raspberrypi.org/
func IsRaspberryPi() bool {
	return false
}

// IsPine64 returns true if running on a pine64 board.
//
// https://www.pine64.org/
func IsPine64() bool {
	return false
}
