// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

#ifndef __APA102_H__
#define __APA102_H__

#include <Homie.h>
#include "anim1d.h"

// maxAPA102Out is the maximum intensity of each channel on a APA102 LED.
const uint16_t maxAPA102Out = 0x1EE1;

void initAPA102();

uint32_t Write(const Frame& pixels);

#endif
