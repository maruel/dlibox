// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"image"

	"github.com/maruel/dlibox/go/pio/conn/i2c"
	"github.com/maruel/dlibox/go/pio/devices"
	"github.com/maruel/dlibox/go/pio/devices/ssd1306"
	"github.com/maruel/dlibox/go/pio/devices/ssd1306/image1bit"
	"github.com/maruel/dlibox/go/psf"
)

func initDisplay(d *Display) (devices.Display, error) {
	i2cBus, err := i2c.New(-1)
	if err != nil {
		return nil, err
	}
	display, err := ssd1306.NewI2C(i2cBus, 128, 64, false)
	if err != nil {
		return nil, err
	}
	f12, err := psf.Load("Terminus12x6")
	if err != nil {
		return nil, err
	}
	f20, err := psf.Load("Terminus20x10")
	if err != nil {
		return nil, err
	}
	img, err := image1bit.New(image.Rect(0, 0, display.W, display.H))
	if err != nil {
		return nil, err
	}
	f20.Draw(img, 0, 0, image1bit.On, nil, "dlibox!")
	f12.Draw(img, 0, display.H-f12.H-1, image1bit.On, nil, "is awesome")
	if _, err = display.Write(img.Buf); err != nil {
		return nil, err
	}
	return display, nil
}
