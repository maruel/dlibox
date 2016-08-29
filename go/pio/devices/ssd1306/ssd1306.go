// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Package ssd1306 controls a 128x64 monochrome OLED display via a ssd1306
// controler.
//
// Datasheet
//
// https://cdn-shop.adafruit.com/datasheets/SSD1306.pdf
package ssd1306

// Some have SPI enabled;
// https://hallard.me/adafruit-oled-display-driver-for-pi/
// https://learn.adafruit.com/ssd1306-oled-displays-with-raspberry-pi-and-beaglebone-black?view=all

import (
	"errors"

	"github.com/maruel/dlibox/go/pio/buses"
)

// FrameRate determines scrolling speed.
type FrameRate byte

const (
	FrameRate2   FrameRate = 7
	FrameRate3   FrameRate = 4
	FrameRate4   FrameRate = 5
	FrameRate5   FrameRate = 0
	FrameRate25  FrameRate = 6
	FrameRate64  FrameRate = 1
	FrameRate128 FrameRate = 2
	FrameRate256 FrameRate = 3
)

// Orientation is used for scrolling.
type Orientation byte

const (
	Left    Orientation = 0x27
	Right   Orientation = 0x26
	UpRight Orientation = 0x29
	UpLeft  Orientation = 0x2A
)

// Dev is an open handle to the display controler.
type Dev struct {
	d buses.Dev
	W int
	H int
}

// Make returns a Dev object that communicates over I²C to SSD1306 display
// controler.
//
// If rotated, turns the display by 180°
func Make(i buses.I2C, w, h int, rotated bool) (*Dev, error) {
	d := &Dev{d: buses.Dev{i, 0x3C}, W: w, H: h}

	contrast := byte(0x7F) // (default value)

	// Set COM output scan direction; C0 means normal; C8 means reversed
	comScan := byte(0xC8)
	// See page 40.
	columnAddr := byte(0xA1)
	if rotated {
		comScan = 0xC0
		columnAddr = byte(0xA0)
	}
	// Initialize the device by fully reseting all values.
	// https://cdn-shop.adafruit.com/datasheets/SSD1306.pdf
	// Page 64 has the full recommended flow.
	// Page 28 lists all the commands.
	init := []byte{
		0xAE,                // Display off
		0xA8, byte(d.H - 1), // Set MUX ratio
		0xD3, 0x00, // Set display offset; 0
		0x40,       // Start display start line; 0
		columnAddr, // Set segment remap; RESET is column 127.
		comScan,
		0xDA, 0x12, // Set COM pins hardware configuration; see page 40
		0x81, contrast, // Set contrast control
		0xA4,       // Set display to use GDDRAM content
		0xA6,       // Set normal display (0xA7 for reversed bitness i.e. bit set is black) (?)
		0xD5, 0x40, // Set osc frequency and divide ratio; power on reset value is 0x3F.
		0x8D, 0x14, // Enable charge pump regulator; page 62

		// Not sure
		0xD9, 0xF1, // Set pre-charge period.
		//0xDB, 0x40, // Set Vcomh deselect level; page 32
		0x20, 0x00, // Set memory addressing mode to horizontal (can be page, horizontal or vertical)
		0x2E,                // Deactivate scroll
		0x00 | 0x00,         // Set column offset (lower nibble)
		0x10 | 0x00,         // Set column offset (higher nibble)
		0xA8, byte(d.H - 1), // Set multiplex ratio (number of lines to display)
		// TODO(maruel): should probably clear the buffer before enabling display, otherwise the previous buffer is shown until refresh.
		0xAF, // Display on
	}
	if _, err := d.d.Write(init); err != nil {
		return nil, err
	}
	return d, nil
}

// Write writes a buffer of pixels to the display.
func (d *Dev) Write(pixels []byte) (int, error) {
	if len(pixels) != d.H*d.W/8 {
		return 0, errors.New("invalid pixel stream")
	}

	// Run everything as one big transaction to reduce downtime on the bus.
	hdr := []byte{
		0xA4,                      // Write data
		0x40 | 0,                  // Start line
		0x21, 0x00, byte(d.W - 1), // Set column address (Width)
		0x22, 0x00, byte(d.H/8 - 1), // Set page address (Pages)
	}

	//*
	d.d.Write(hdr)
	d.d.Write(append([]byte{0x40}, pixels...))
	return 0, nil
	/*/
	// TODO(maruel): Use oscilloscope to figure out why this doesn't work.
	start := []byte{
		0x40, // Pixel data
	}
	ios := []buses.IO{
		{buses.WriteStop, hdr},
		{buses.Write, start},
		{buses.WriteStop, pixels},
	}
	if err := d.d.Tx(ios); err != nil {
		return 0, err
	}
	return len(pixels), nil
	//*/
}

// Scroll scrolls the entire.
func (d *Dev) Scroll(o Orientation, rate FrameRate) error {
	// TODO(maruel): Allow to specify page.
	// TODO(maruel): Allow to specify offset.
	if o == Left || o == Right {
		// page 28
		// STOP, <op>, dummy, <start page>, <rate>,  <end page>, <dummy>, <dummy>, <ENABLE>
		_, err := d.d.Write([]byte{0x2E, byte(o), 0x00, 0x00, byte(rate), 0x07, 0x00, 0xFF, 0x2F})
		return err
	}
	// page 29
	// STOP, <op>, dummy, <start page>, <rate>,  <end page>, <offset>, <ENABLE>
	// page 30: 0xA3 permits to set rows for scroll area.
	_, err := d.d.Write([]byte{0x2E, byte(o), 0x00, 0x00, byte(rate), 0x07, 0x01, 0x2F})
	return err
}

// StopScroll stops any scrolling previously set.
//
// It will only take effect after redrawing the ram.
//
// TODO(maruel): Doesn't work.
func (d *Dev) StopScroll() error {
	_, err := d.d.Write([]byte{0x2E})
	return err
}

// SetContrast changes the screen contrast.
//
// TODO(maruel): Doesn't work.
func (d *Dev) SetContrast(level byte) error {
	_, err := d.d.Write([]byte{0x81, level})
	return err
}

// Enable or disable the display.
//
// TODO(maruel): Doesn't work.
func (d *Dev) Enable(on bool) error {
	b := byte(0xAE)
	if on {
		b = 0xAF
	}
	_, err := d.d.Write([]byte{b})
	return err
}
