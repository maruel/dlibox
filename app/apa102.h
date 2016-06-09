// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

#ifndef __APA102_H__
#define __APA102_H__

#include "anim1d.h"

void initAPA102();

uint16_t Ramp(uint8_t l, uint16_t max);
void ColorToAPA102(const Color &c, uint8_t* dst);
void Raster(const Frame& pixels, uint8_t *buf);

#endif
