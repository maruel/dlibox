// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Package tm1637 controls a TM1637 device over GPIO pins.
//
// Datasheet
//
// http://olimex.cl/website_MCI/static/documents/Datasheet_TM1637.pdf
package tm1637

import (
	"errors"
	"time"

	"github.com/maruel/dlibox/go/pio/host"
)

type Dev struct {
	clk  host.PinOut
	data host.PinIO
}

// Brightness defines the screen brightness as controlled by the internal PWM.
type Brightness uint8

const (
	Off          Brightness = 0x80 // Completely off.
	Brightness1  Brightness = 0x88 // 1/16 PWM
	Brightness2  Brightness = 0x89 // 2/16 PWM
	Brightness4  Brightness = 0x8A // 4/16 PWM
	Brightness10 Brightness = 0x8B // 10/16 PWM
	Brightness11 Brightness = 0x8C // 11/16 PWM
	Brightness12 Brightness = 0x8D // 12/16 PWM
	Brightness13 Brightness = 0x8E // 13/16 PWM
	Brightness14 Brightness = 0x8F // 14/16 PWM
)

// Brightness changes the brightness and/or turns the display on and off.
func (d *Dev) SetBrightness(b Brightness) error {
	return d.writeBytes(byte(b))
}

// Segments writes raw segments.
//     -A-
//    F   B
//     -G-
//    E   C
//     -D-   P
//
// P is a dot. Each byte is encoded as PGFEDCBA.
func (d *Dev) Segments(seg ...byte) error {
	if len(seg) > 6 {
		return errors.New("up to 6 digits are supported")
	}
	// Use auto-incrementing address. It is possible to write to a single
	// segment but there isn't much point.
	return d.writeBytes(0x40, 0xC0)
	return d.writeBytes(seg...)
}

// Digits writes hex numbers. Numbers outside the range [0, 15] are
// displayed as blank. Use -1 to mark it as blank.
func (d *Dev) Digits(n ...int) error {
	seg := make([]byte, len(n))
	for i := range n {
		if n[i] >= 0 && n[i] < 16 {
			seg[i] = byte(digitToSegment[n[i]])
		}
	}
	return d.Segments(seg...)
}

// writeBytes sends a stream of bytes to the tm1637.
//
// The protocol is similar to I²C but there is no slave address. So this is
// close to bit banging I²C.
func (d *Dev) writeBytes(b ...byte) error {
	// "When CLK is a high level and DIO changes from high to low level, data
	// input starts."
	d.data.Set(host.Low)
	time.Sleep(clockHalfCycle)
	// Write the bytes.
	for _, c := range b {
		_, err := d.writeByte(c)
		if err != nil {
			return err
		}
	}
	// "When CLK is a high level and DIO changes from low level to high level,
	// data input ends."
	d.clk.Set(host.Low)
	time.Sleep(clockHalfCycle)
	d.clk.Set(host.High)
	time.Sleep(clockHalfCycle)
	d.data.Set(host.High)
	time.Sleep(clockHalfCycle)
	return nil
}

// writeByte starts with d.data low and d.clk high and ends with d.data low and
// d.clk high.
func (d *Dev) writeByte(b byte) (bool, error) {
	for i := 0; i < 8; i++ {
		// "When data is input, DIO signal should not change for high level CLK and
		// DIO signal should change for low level CLK signal."
		d.clk.Set(host.Low)
		time.Sleep(clockQuarterCycle)
		if b&1 != 0 {
			d.data.Set(host.High)
		} else {
			d.data.Set(host.Low)
		}
		time.Sleep(clockQuarterCycle)
		d.clk.Set(host.High)
		time.Sleep(clockHalfCycle)
		b >>= 1
	}
	// 9th clock is ACK.
	d.clk.Set(host.Low)
	time.Sleep(clockHalfCycle)
	if err := d.data.In(host.Up); err != nil {
		return false, err
	}
	d.clk.Set(host.High)
	time.Sleep(clockQuarterCycle)
	ack := d.data.Read() == host.Low
	time.Sleep(clockQuarterCycle)
	if err := d.data.Out(); err != nil {
		return false, err
	}
	d.data.Set(host.Low)
	return ack, nil
}

// Make returns an object that communicates over two pins to a TM1637.
func Make(clk host.PinOut, data host.PinIO) (*Dev, error) {
	// Spec calls to idle at high.
	if err := clk.Out(); err != nil {
		return nil, err
	}
	clk.Set(host.High)
	if err := data.Out(); err != nil {
		return nil, err
	}
	data.Set(host.High)
	d := &Dev{clk: clk, data: data}
	return d, nil
}

//

// Page 10 states the max clock frequency is 500KHz but page 3 states 250KHz.
const clockHalfCycle = time.Second / 250000 / 2
const clockQuarterCycle = clockHalfCycle / 2

// Hex digits from 0 to F.
var digitToSegment = []byte{
	0x3f, 0x06, 0x5b, 0x4f, 0x66, 0x6d, 0x7d, 0x07, 0x7f, 0x6f, 0x77, 0x7c, 0x39, 0x5e, 0x79, 0x71,
}
