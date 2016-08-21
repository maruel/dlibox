// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package bme280

import "testing"

func TestCalibration(t *testing.T) {
	c := calibration{
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
	tRaw := int32(524112)
	pRaw := int32(309104)
	hRaw := int32(30987)
	temp, tFine := c.compensateTempFloat(tRaw)
	pres := c.compensatePressureFloat(pRaw, tFine)
	humi := c.compensateHumidityFloat(hRaw, tFine)
	if !floatEqual(temp, 22.948120) {
		t.Fatalf("temp %f", temp)
	}
	if !floatEqual(pres, 100.046074) {
		t.Fatalf("pressure %f", pres)
	}
	if !floatEqual(humi, 63.167889) {
		t.Fatalf("humidity %f", humi)
	}
}

var epsilon float32 = 0.00000001

func floatEqual(a, b float32) bool {
	return (a-b) < epsilon && (b-a) < epsilon
}
