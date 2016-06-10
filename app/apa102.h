// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

#ifndef __APA102_H__
#define __APA102_H__

#include "anim1d.h"

// maxAPA102Out is the maximum intensity of each channel on a APA102 LED.
const uint16_t maxAPA102Out = 0x1EE1;

void initAPA102();

// Ramp the color intensity l [0, 255] over maxIntensity on a n³ curve.
// maxIntensity should be 0 for the default, which defaults to maxAPA102Out or
// between [255, maxAPA102Out].
uint16_t Ramp(uint8_t l, uint16_t maxIntensity);
void ColorToAPA102(const Color &c, uint8_t* dst, uint16_t maxIntensity);
void Raster(const Frame& pixels, uint8_t *buf, uint16_t maxIntensity);
uint32_t Write(const Frame& pixels, uint16_t maxIntensity);

#endif
