// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package bme280

import (
	"testing"

	"github.com/maruel/dlibox/go/pio/devices"
	"github.com/maruel/dlibox/go/pio/host/hosttest"
)

// Real data extracted from a device.
var calib = calibration{
	t1: 28176,
	t2: 26220,
	t3: 350,
	p1: 38237,
	p2: -10824,
	p3: 3024,
	p4: 7799,
	p5: -99,
	p6: -7,
	p7: 9900,
	p8: -10230,
	p9: 4285,
	h2: 366, // Note they are inversed for bit packing.
	h1: 75,
	h3: 0,
	h4: 309,
	h5: 0,
	h6: 30,
}

func TestRead(t *testing.T) {
	// This data was generated with "bme280 -r"
	bus := hosttest.I2CPlayback{
		Ops: []hosttest.I2CIO{
			// Chipd ID detection.
			{Addr: 0x76, Write: []byte{0xd0}, Read: []byte{0x60}},
			// Calibration data.
			{
				Addr:  0x76,
				Write: []byte{0x88},
				Read:  []byte{0x10, 0x6e, 0x6c, 0x66, 0x32, 0x0, 0x5d, 0x95, 0xb8, 0xd5, 0xd0, 0xb, 0x77, 0x1e, 0x9d, 0xff, 0xf9, 0xff, 0xac, 0x26, 0xa, 0xd8, 0xbd, 0x10, 0x0, 0x4b},
			},
			// Calibration data.
			{Addr: 0x76, Write: []byte{0xe1}, Read: []byte{0x6e, 0x1, 0x0, 0x13, 0x5, 0x0, 0x1e}},
			// Configuration.
			{Addr: 0x76, Write: []byte{0xf4, 0x6c, 0xf2, 0x3, 0xf5, 0xe0, 0xf4, 0x6f}, Read: nil},
			// Read.
			{Addr: 0x76, Write: []byte{0xf7}, Read: []byte{0x4a, 0x52, 0xc0, 0x80, 0x96, 0xc0, 0x7a, 0x76}},
		},
	}
	dev, err := NewI2C(&bus, O4x, O4x, O4x, S20ms, FOff)
	if err != nil {
		t.Fatal(err)
	}
	env := devices.Environment{}
	if err := dev.Read(&env); err != nil {
		t.Fatalf("Read(): %v", err)
	}
	if env.MilliCelcius != 23720 {
		t.Fatalf("temp %d", env.MilliCelcius)
	}
	if env.Pascal != 100943 {
		t.Fatalf("pressure %d", env.Pascal)
	}
	if env.Humidity != 6531 {
		t.Fatalf("humidity %d", env.Humidity)
	}
}

func TestCalibrationFloat(t *testing.T) {
	// Real data extracted from measurements from this device.
	tRaw := int32(524112)
	pRaw := int32(309104)
	hRaw := int32(30987)

	// Compare the values with the 3 algorithms.
	temp, tFine := calib.compensateTempFloat(tRaw)
	pres := calib.compensatePressureFloat(pRaw, tFine)
	humi := calib.compensateHumidityFloat(hRaw, tFine)
	if tFine != 117494 {
		t.Fatalf("tFine %d", tFine)
	}
	if !floatEqual(temp, 22.948120) {
		// 22.95°C
		t.Fatalf("temp %f", temp)
	}
	if !floatEqual(pres, 100.046074) {
		// 100.046kPa
		t.Fatalf("pressure %f", pres)
	}
	if !floatEqual(humi, 63.167889) {
		// 63.17%
		t.Fatalf("humidity %f", humi)
	}
}

func TestCalibrationInt(t *testing.T) {
	// Real data extracted from measurements from this device.
	tRaw := int32(524112)
	pRaw := int32(309104)
	hRaw := int32(30987)

	temp, tFine := calib.compensateTempInt(tRaw)
	pres64 := calib.compensatePressureInt64(pRaw, tFine)
	pres32 := calib.compensatePressureInt32(pRaw, tFine)
	humi := calib.compensateHumidityInt(hRaw, tFine)
	if tFine != 117407 {
		t.Fatalf("tFine %d", tFine)
	}
	if temp != 2293 {
		// 2293/100 = 22.93°C
		// Delta is <0.02°C which is pretty good.
		t.Fatalf("temp %d", temp)
	}
	if pres64 != 25611063 {
		// 25611063/256/1000 = 100.043214844
		// Delta is 3Pa which is ok.
		t.Fatalf("pressure64 %d", pres64)
	}
	if pres32 != 100045 {
		// 100045/1000 = 100.045kPa
		// Delta is 1Pa which is pretty good.
		t.Fatalf("pressure32 %d", pres32)
	}
	if humi != 64686 {
		// 64686/1024 = 63.17%
		// Delta is <0.01% which is pretty good.
		t.Fatalf("humidity %d", humi)
	}
}

var epsilon float32 = 0.00000001

func floatEqual(a, b float32) bool {
	return (a-b) < epsilon && (b-a) < epsilon
}
