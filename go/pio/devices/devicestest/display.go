// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package devicestest

import (
	"errors"
	"image"
	"image/color"
	"image/draw"

	"github.com/maruel/dlibox/go/pio/devices"
)

// Display is a fake devices.Display
type Display struct {
	Img *image.NRGBA
}

func (d *Display) Write(pixels []byte) (int, error) {
	if len(pixels)%3 != 0 {
		return 0, errors.New("invalid RGB stream length")
	}
	copy(d.Img.Pix, pixels)
	return len(pixels), nil
}

func (d *Display) ColorModel() color.Model {
	return d.Img.ColorModel()
}

func (d *Display) Bounds() image.Rectangle {
	return d.Img.Bounds()
}

func (d *Display) Draw(r image.Rectangle, src image.Image, sp image.Point) {
	draw.Draw(d.Img, r, src, sp, draw.Src)
}

var _ devices.Display = &Display{}
