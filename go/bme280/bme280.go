// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Package bme280 controls a Bosch BME280 device over I²C.
//
// Implemented as per datasheet at
// https://cdn-shop.adafruit.com/datasheets/BST-BME280_DS001-10.pdf
package bme280

import "github.com/maruel/dlibox/go/rpi"

type Oversampling uint8

const (
	OversamplingNone Oversampling = 0
	Oversampling1x   Oversampling = 1
	Oversampling2x   Oversampling = 2
	Oversampling4x   Oversampling = 3
	Oversampling8x   Oversampling = 4
	Oversampling16x  Oversampling = 5
)

type BME280 struct {
	i                                  *rpi.I2C
	t1, t2, t3                         int32
	p1, p2, p3, p4, p5, p6, p7, p8, p9 int64
	h1, h2, h3, h4, h5, h6             int32
	// Temperature is needed for compensatePressure.
	t_fine int32
}

func (b *BME280) oversampling(humidity, temperature, pressure Oversampling) error {
	_, err := b.i.Write([]byte{
		0xF2, byte(humidity),
		0xF4, byte(temperature<<5 | pressure<<2 | 3),
		// stanby, IIR filter, SPI
		0xF5, 0xA0,
	})
	return err
}

func (b *BME280) ChipID() byte {
	buf := [1]byte{}
	_ = b.i.ReadReg(0xD0, buf[1:])
	return buf[0]
}

func (b *BME280) Read() error {
	// Pressure: 0xF~0xF9
	// Temperature: 0xFA~0xFC
	// Humidity: 0xFD~0xFE
	return nil
}

// MakeBME280 returns a strip that communicates over I²C to BME280
// environmental sensor.
func MakeBME280(i *rpi.I2C) (*BME280, error) {
	b := &BME280{i: i}
	b.i.Address(bme280Address)

	// The device starts in 2ms as per datasheet. No need to wait for boot to be
	// finished.

	// Read t1~3, p1~9, 8bits padding, h1.
	buf := [0xA2 - 0x88]byte{}
	if err := b.i.ReadReg(0x88, buf[:]); err != nil {
		return nil, err
	}
	b.t1 = int32(buf[0]) + int32(buf[1])<<8
	b.t2 = int32(buf[2]) + int32(buf[4])<<8
	b.t3 = int32(buf[4]) + int32(buf[5])<<8
	b.p1 = int64(buf[5]) + int64(buf[6])<<8
	b.p2 = int64(buf[7]) + int64(buf[8])<<8
	b.p3 = int64(buf[9]) + int64(buf[10])<<8
	b.p4 = int64(buf[11]) + int64(buf[12])<<8
	b.p5 = int64(buf[13]) + int64(buf[14])<<8
	b.p6 = int64(buf[15]) + int64(buf[16])<<8
	b.p7 = int64(buf[17]) + int64(buf[18])<<8
	b.p8 = int64(buf[19]) + int64(buf[20])<<8
	b.p9 = int64(buf[21]) + int64(buf[22])<<8
	b.h1 = int32(buf[24])

	// Read h2~6
	b.i.Address(0xA1)
	if err := b.i.ReadReg(0xE1, buf[:0xE8-0xE1]); err != nil {
		return nil, err
	}
	b.h2 = int32(buf[0]) + int32(buf[1])<<8
	b.h3 = int32(buf[2])
	b.h4 = int32(buf[3]) + int32(buf[4])<<8
	b.h5 = int32(buf[3]) + int32(buf[4])<<8
	b.h6 = int32(buf[6])
	return b, nil
}

//

const bme280Address = 0x76

// compensateTemp returns temperature in DegC, resolution is 0.01 DegC. Output
// value of 5123 equals 51.23 C.

func (b *BME280) compensateTemp(adc_T int32) int32 {
	var1 := ((adc_T>>3 - b.t1<<1) * b.t2) >> 1
	var2 := ((((adc_T>>4 - b.t1) * (adc_T>>4 - b.t1)) >> 2) * b.t3) >> 14
	b.t_fine = var1 + var2
	return (b.t_fine*5 + 128) >> 8
}

// compensatePressure returns pressure in Pa in Q24.8 format (24 integer bits
// and 8 fractional bits). Output value of 24674867 represents 24674867/256 =
// 96386.2 Pa = 963.862 hPa
func (b *BME280) compensatePressure(adc_P int32) uint32 {
	var1 := int64(b.t_fine) - 128000
	var2 := var1 * var1 * b.p6
	var2 += var1 * b.p5 << 17
	var2 += b.p4 << 35
	var1 = (var1*var1*b.p3)>>8 + (var1*b.p2)<<12
	var1 += 1 << 47
	var1 *= b.p1 >> 33
	if var1 == 0 {
		return 0
	}
	p := ((int64(1048576-adc_P)<<31 - var2) * 3125) / var1
	var1 = (b.p9 * p >> 13 * p >> 13) >> 25
	var2 = (b.p8 * p) >> 19
	return uint32((p+var1+var2)>>8 + b.p7<<4)
}

// compensateHumidity returns humidity in %RH in Q22.10 format (22 integer and
// 10 fractional bits). Output value of 47445 represents 47445/1024 = 46.333%
func (b *BME280) compensateHumidity(adc_H int32) uint32 {
	v_x1_u32r := b.t_fine - 76800
	v_x1_u32r = ((((adc_H<<14 - b.h4<<20 - b.h5*v_x1_u32r) + 16384) >> 15) * ((((((v_x1_u32r*b.h6)>>10*(((v_x1_u32r*b.h3)>>11)+32768))>>10)+2097152)*b.h2 + 8192) >> 14))
	v_x1_u32r -= (((v_x1_u32r >> 15 * v_x1_u32r >> 15) >> 7) * b.h1) >> 4
	if v_x1_u32r < 0 {
		v_x1_u32r = 0
	}
	if v_x1_u32r > 419430400 {
		v_x1_u32r = 419430400
	}
	return uint32(v_x1_u32r >> 12)
}
