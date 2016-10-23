// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package anim1d

import "math"

func abs(x float32) float32 {
	return float32(math.Abs(float64(x)))
}

func ceil(x float32) float32 {
	return float32(math.Ceil(float64(x)))
}

func hypot(x, y float32) float32 {
	return float32(math.Hypot(float64(x), float64(y)))
}

func logn(x float32) float32 {
	return float32(math.Log(float64(x)))
}

func log1p(x float32) float32 {
	return float32(math.Log1p(float64(x)))
}

func sin(x float32) float32 {
	return float32(math.Sin(float64(x)))
}

func roundF(x float32) float32 {
	if x < 0 {
		return ceil(x - 0.5)
	}
	return ceil(x + 0.5)
}
