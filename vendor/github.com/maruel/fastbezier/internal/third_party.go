// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Copyright (c) 2009 The Go Authors. All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

   * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
   * Neither the name of Google Inc. nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

package internal

// cubicBezier returns [0, 1] for input `t` based on the cubic bezier curve
// (0,0), (x0,y0), (x1, y1), (1, 1).
//
// Extracted from
// https://github.com/golang/mobile/blob/master/exp/sprite/clock/tween.go
//
// It was adapted to use float32 instead of float64. On x86 based platform,
// there is no performacen difference between these two types but on other
// platforms like ARM, there can be significant difference. Since float32 has 7
// digits of precision and uint16 has 5, it is reasonably precise enough to not
// degrade the end results.
func CubicBezier(x0, y0, x1, y1, x float32) float32 {
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
