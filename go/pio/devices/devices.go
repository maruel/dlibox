// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package devices

import (
	"image"
	"image/color"
	"io"
)

// Display represents a pixel output device. It is a write-only interface.
//
// What Display represents can be as varied as a 1 bit OLED display or a strip
// of LED lights.
type Display interface {
	// Writer can be used when the native display pixel format is known. Each
	// write must cover exactly the whole screen as a single packed stream of
	// pixels.
	io.Writer
	// ColorModel returns the device native color model. It is generally
	// color.NRGBA for color display and color.Palette for black and white
	// display.
	ColorModel() color.Model
	// Bounds returns the size of the output device. Generally displays should
	// have Min at {0, 0} but this is not guaranteed in multiple displays setup
	// or when an instance of this interface represents a section of a larger
	// logical display.
	Bounds() image.Rectangle
	// Draw updates the display with this image starting at 'sp' offset into the
	// display into 'r'. The code will likely be faster if the image is in the
	// display's native color format.
	Draw(r image.Rectangle, src image.Image, sp image.Point)
}

// Environment represents measurements from an environmental sensor.
type Environment struct {
	MilliCelcius int32 // 0.001°C; range [-273150, >1000000]
	Pascal       int32 // 1Pa; range [0, >1000000]
	Humidity     int32 // 0.01%rH or 0.1milli-rH; range [0, 10000]
}

// Environmental represents an environmental sensor.
type Environmental interface {
	// Read returns the value read from the sensor. Unsupported metrics are not
	// modified.
	Read(env *Environment) error
}
