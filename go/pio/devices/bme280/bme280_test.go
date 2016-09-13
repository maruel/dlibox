// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package bme280

import "testing"

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
