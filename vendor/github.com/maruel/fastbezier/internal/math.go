// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Math functions

package internal

import "math"

// FloatToUint16 converts a floating point value in range [0, 65535] to a
// uint16.
//
// Doesn't return valid values for x < 0 or > 65535.
func FloatToUint16(x float32) uint16 {
	return uint16(math.Floor(float64(x + 0.5)))
}

const reverse = 1. / 65535.

func CubicBezier16(x0, y0, x1, y1 float32, x uint16) uint16 {
	return FloatToUint16(CubicBezier(x0, y0, x1, y1, float32(x)*reverse) * 65535.)
}
