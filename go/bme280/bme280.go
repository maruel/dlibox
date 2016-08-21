// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Package bme280 controls a Bosch BME280 device over I²C.
//
// Implemented as per datasheet at
// https://cdn-shop.adafruit.com/datasheets/BST-BME280_DS001-10.pdf
//
// Mostly implemented it to prove myself I can still read a datasheet at my age.
package bme280

import (
	"errors"

	"github.com/maruel/dlibox/go/rpi"
)

// Oversampling affects how much time is taken to measure each of temperature,
// pressure and humidity.
//
// Temperature must be measured for pressure and humidity to be measured. The
// duration is approximatively:
//     duration_in_ms = 1 + 2*temp + 2*press+0.5 + 2*humidy+0.5
//
// Using high oversampling and low standby results in highest power
// consumption, but this is still below 1mA so we generally don't care.
type Oversampling uint8

const (
	No   Oversampling = 0
	O1x  Oversampling = 1
	O2x  Oversampling = 2
	O4x  Oversampling = 3
	O8x  Oversampling = 4
	O16x Oversampling = 5
)

// Standby is the time the BME280 waits idle between measurements. This reduces
// power consumption when the host won't read the values as fast as the
// measurements are done.
type Standby uint8

const (
	S500us Standby = 0
	S10ms  Standby = 6
	S20ms  Standby = 7
	S62ms  Standby = 1
	S125ms Standby = 2
	S250ms Standby = 3
	S500ms Standby = 4
	S1s    Standby = 5
)

// Filter specifies the internal IIR filter to get steady measurements without
// using oversampling. This is mainly used to reduce power consumption.
type Filter uint8

const (
	FOff Filter = 0
	F2   Filter = 1
	F4   Filter = 2
	F8   Filter = 3
	F16  Filter = 4
)

type BME280 struct {
	i *rpi.I2C
	c calibration
}

// Read returns measurements as C, kPa and % of relative humidity.
func (b *BME280) Read() (float32, float32, float32, error) {
	// Pressure: 0xF7~0xF9
	// Temperature: 0xFA~0xFC
	// Humidity: 0xFD~0xFE
	buf := [0xFF - 0xF7]byte{}
	if err := b.i.ReadReg(0xF7, buf[:]); err != nil {
		return 0, 0, 0, err
	}
	pRaw := int32(buf[0])<<12 | int32(buf[1])<<4 | int32(buf[2])>>4
	tRaw := int32(buf[3])<<12 | int32(buf[4])<<4 | int32(buf[5])>>4
	hRaw := int32(buf[6])<<8 | int32(buf[7])
	t, tFine := b.c.compensateTempFloat(tRaw)
	p := b.c.compensatePressureFloat(pRaw, tFine)
	h := b.c.compensateHumidityFloat(hRaw, tFine)
	return t, p, h, nil
}

// MakeBME280 returns a strip that communicates over I²C to BME280
// environmental sensor.
//
// Recommended values are O8x for oversampling, S20ms for standby and FOff for
// filter.
func MakeBME280(i *rpi.I2C, temperature, pressure, humidity Oversampling, standby Standby, filter Filter) (*BME280, error) {
	b := &BME280{i: i}
	b.i.Address(bme280Address)

	// The device starts in 2ms as per datasheet. No need to wait for boot to be
	// finished.

	// Read the chipd ID right away. It should be 0x60.
	buf := [0xA2 - 0x88]byte{}
	if err := b.i.ReadReg(0xD0, buf[:1]); err != nil {
		return nil, err
	}
	if buf[0] != 0x60 {
		return nil, errors.New("unexpected chip id; is this a BME280?")
	}

	// Read t1~3, p1~9, 8bits padding, h1.
	if err := b.i.ReadReg(0x88, buf[:]); err != nil {
		return nil, err
	}
	b.c.t1 = uint16(buf[0]) | uint16(buf[1])<<8
	b.c.t2 = int16(buf[2]) | int16(buf[3])<<8
	b.c.t3 = int16(buf[4]) | int16(buf[5])<<8
	b.c.p1 = uint16(buf[6]) | uint16(buf[7])<<8
	b.c.p2 = int16(buf[8]) | int16(buf[9])<<8
	b.c.p3 = int16(buf[10]) | int16(buf[11])<<8
	b.c.p4 = int16(buf[12]) | int16(buf[13])<<8
	b.c.p5 = int16(buf[14]) | int16(buf[15])<<8
	b.c.p6 = int16(buf[16]) | int16(buf[17])<<8
	b.c.p7 = int16(buf[18]) | int16(buf[19])<<8
	b.c.p8 = int16(buf[20]) | int16(buf[21])<<8
	b.c.p9 = int16(buf[22]) | int16(buf[23])<<8
	b.c.h1 = uint8(buf[25])

	// Read h2~6
	if err := b.i.ReadReg(0xE1, buf[:0xE8-0xE1]); err != nil {
		return nil, err
	}
	b.c.h2 = int16(buf[0]) | int16(buf[1])<<8
	b.c.h3 = uint8(buf[2])
	b.c.h4 = int16(buf[3])<<4 | int16(buf[4])&0xF
	b.c.h5 = int16(buf[4])&0xF0 | int16(buf[5])<<4
	b.c.h6 = int8(buf[6])

	// Write config and start it.
	if _, err := b.i.Write([]byte{
		// ctrl_meas; put it to sleep otherwise the config update may be ignored.
		0xF4, byte(temperature)<<5 | byte(pressure)<<2 | byte(sleep),
		// ctrl_hum
		0xF2, byte(humidity),
		// config
		0xF5, byte(standby)<<5 | byte(filter)<<2,
		// ctrl_meas
		0xF4, byte(temperature)<<5 | byte(pressure)<<2 | byte(normal),
	}); err != nil {
		return nil, err
	}
	return b, nil
}

//

const bme280Address = 0x76

// mode is stored in config
type mode byte

const (
	sleep  mode = 0 // no operation, all registers accessible, lowest power, selected after startup
	forced mode = 1 // perform one measurement, store results and return to sleep mode
	normal mode = 3 // perpetual cycling of measurements and inactive periods
)

type status byte

const (
	measuring status = 8 // set when conversion is running
	im_update status = 1 // set when NVM data are being copied to image registers
)

// Register table:
// 0x00..0x87  --
// 0x88..0xA1  Calibration data
// 0xA2..0xCF  --
// 0xD0        Chip id; reads as 0x60
// 0xD1..0xDF  --
// 0xE0        Reset by writing 0xB6 to it
// 0xE1..0xF0  Calibration data
// 0xF1        --
// 0xF2        ctrl_hum; ctrl_meas must be writen to after for change to this register to take effect
// 0xF3        status
// 0xF4        ctrl_meas
// 0xF5        config
// 0xF6        --
// 0xF7        press_msb
// 0xF8        press_lsb
// 0xF9        press_xlsb
// 0xFA        temp_msb
// 0xFB        temp_lsb
// 0xFC        temp_xlsb
// 0xFD        hum_msb
// 0xFE        hum_lsb

// https://cdn-shop.adafruit.com/datasheets/BST-BME280_DS001-10.pdf
// Page 23

type calibration struct {
	t1                             uint16
	t2, t3                         int16
	p1                             uint16
	p2, p3, p4, p5, p6, p7, p8, p9 int16
	h2                             int16 // Reordered for packing
	h1, h3                         uint8
	h4, h5                         int16
	h6                             int8
}

// compensateTempInt returns temperature in DegC, resolution is 0.01 DegC.
// Output value of 5123 equals 51.23 C.
func (c *calibration) compensateTempInt(raw int32) (int32, int32) {
	x := ((raw>>3 - int32(c.t1)<<1) * int32(c.t2)) >> 1
	var2 := ((((raw>>4 - int32(c.t1)) * (raw>>4 - int32(c.t1))) >> 2) * int32(c.t3)) >> 14
	tFine := x + var2
	return (tFine*5 + 128) >> 8, tFine
}

// compensatePressureInt returns pressure in Pa in Q24.8 format (24 integer bits
// and 8 fractional bits). Output value of 24674867 represents 24674867/256 =
// 96386.2 Pa = 963.862 hPa.
func (c *calibration) compensatePressureInt(raw, tFine int32) uint32 {
	x := int64(tFine) - 128000
	y := x * x * int64(c.p6)
	y += x * int64(c.p5) << 17
	y += int64(c.p4) << 35
	x = (x*x*int64(c.p3))>>8 + (x*int64(c.p2))<<12
	x += 1 << 47
	x *= int64(c.p1) >> 33
	if x == 0 {
		return 0
	}
	p := ((int64(1048576-raw)<<31 - y) * 3125) / x
	x = (int64(c.p9) * int64(p) >> 13 * int64(raw) >> 13) >> 25
	y = (int64(c.p8) * p) >> 19
	return uint32((p+x+y)>>8 + int64(c.p7)<<4)
}

// compensateHumidityInt returns humidity in %RH in Q22.10 format (22 integer
// and 10 fractional bits). Output value of 47445 represents 47445/1024 =
// 46.333%
func (c *calibration) compensateHumidityInt(raw, tFine int32) uint32 {
	x := tFine - 76800
	x = ((((raw<<14 - int32(c.h4)<<20 - int32(c.h5)*x) + 16384) >> 15) * ((((((x*int32(c.h6))>>10*(((x*int32(c.h3))>>11)+32768))>>10)+2097152)*int32(c.h2) + 8192) >> 14))
	x = (((x >> 15 * x >> 15) >> 7) * int32(c.h1)) >> 4
	if x < 0 {
		return 0
	}
	if x > 419430400 {
		return 419430400 >> 12
	}
	return uint32(x >> 12)
}

// Page 49

func (c *calibration) compensateTempFloat(raw int32) (float32, int32) {
	x := (float64(raw)/16384. - float64(c.t1)/1024.) * float64(c.t2)
	y := (float64(raw)/131072. - float64(c.t1)/8192.) * float64(c.t3)
	tFine := int32(x + y)
	return float32((x + y) / 5120.), tFine
}

func (c *calibration) compensatePressureFloat(raw, tFine int32) float32 {
	x := float64(tFine)*0.5 - 64000.
	y := x * x * float64(c.p6) / 32768.
	y += x * float64(c.p5) * 2.
	y = y*0.25 + float64(c.p4)*65536.
	x = (float64(c.p3)*x*x/524288. + float64(c.p2)*x) / 524288.
	x = (1. + x/32768.) * float64(c.p1)
	if x <= 0 {
		return 0
	}
	p := float64(1048576 - raw)
	p = (p - y/4096.) * 6250. / x
	x = float64(c.p9) * p * p / 2147483648.
	y = p * float64(c.p8) / 32768.
	return float32(p+(x+y+float64(c.p7))/16.) / 1000.
}

func (c *calibration) compensateHumidityFloat(raw, tFine int32) float32 {
	h := float64(tFine - 76800)
	h = (float64(raw) - float64(c.h4)*64. + float64(c.h5)/16384.*h) * float64(c.h2) / 65536. * (1. + float64(c.h6)/67108864.*h*(1.+float64(c.h3)/67108864.*h))
	h *= 1. - float64(c.h1)*h/524288.
	if h > 100. {
		return 100.
	}
	if h < 0. {
		return 0.
	}
	return float32(h)
}
