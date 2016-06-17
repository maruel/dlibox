// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package anim1d

// cubicBezier returns [0, 1] for input `t` based on the cubic bezier curve
// (x0,y0), (x1, y1).
// Extracted from https://github.com/golang/mobile/blob/master/exp/sprite/clock/tween.go
func cubicBezier(x0, y0, x1, y1, x float32) float32 {
	t := x
	for i := 0; i < 5; i++ {
		t2 := t * t
		t3 := t2 * t
		d := 1 - t
		d2 := d * d

		nx := 3*d2*t*x0 + 3*d*t2*x1 + t3
		dxdt := 3*d2*x0 + 6*d*t*(x1-x0) + 3*t2*(1-x1)
		if dxdt == 0 {
			break
		}

		t -= (nx - x) / dxdt
		if t <= 0 || t >= 1 {
			break
		}
	}
	if t < 0 {
		t = 0
	}
	if t > 1 {
		t = 1
	}

	// Solve for y using t.
	t2 := t * t
	t3 := t2 * t
	d := 1 - t
	d2 := d * d
	y := 3*d2*t*y0 + 3*d*t2*y1 + t3

	return y
}
