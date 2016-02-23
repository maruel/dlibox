// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"image/color"
	"testing"

	"github.com/maruel/ut"
)

func TestRGBToANSI(t *testing.T) {
	ut.AssertEqual(t, 0, rgbToANSI(color.NRGBA{}))
	ut.AssertEqual(t, 0, rgbToANSI(color.NRGBA{1, 1, 1, 0}))
	ut.AssertEqual(t, 15, rgbToANSI(color.NRGBA{255, 255, 255, 0}))
	ut.AssertEqual(t, 15, rgbToANSI(color.NRGBA{254, 254, 254, 0}))
	ut.AssertEqual(t, 255, rgbToANSI(color.NRGBA{0xE5, 0xEE, 0xF0, 0}))
}
