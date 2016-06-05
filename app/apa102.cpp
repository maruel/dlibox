// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

#include "user_config.h"
#include <SmingCore/SmingCore.h>
#include <SPI.h>

#include "anim1d.h"

// maxOut is the maximum intensity of each channel on a APA102 LED.
const uint16_t maxOut = 0x1EE1;

// Ramp converts input from [0, 0xFF] as intensity to lightness on a scale of
// [0, 0x1EE1] or other desired range [0, max].
//
// It tries to use the same curve independent of the scale used. max can be
// changed to change the color temperature or to limit power dissipation.
//
// It's the reverse of lightness; https://en.wikipedia.org/wiki/Lightness
uint16_t Ramp(uint8_t l, uint16_t max) {
  if (l == 0) {
    // Make sure black is black.
    return 0;
  }
  if (max == 0 || max > maxOut) {
    // If 'max' is not specified or is above maxOut, reset the maximum value.
    max = maxOut;
  } else if (max < 255) {
    max = 255;
  }
  // linearCutOff defines the linear section of the curve. Inputs between
  // [0, linearCutOff] are mapped linearly to the output. It is 1% of maximum
  // output.
  uint32_t linearCutOff = uint32_t((max + 50) / 100);
  uint32_t l32 = uint32(l);
  if (l32 < linearCutOff) {
    return uint16_t(l32);
  }

  // Maps [linearCutOff, 255] to use [linearCutOff*max/255, max] using a x^3
  // ramp.
  // Realign input to [0, 255-linearCutOff]. It now maps to
  // [0, max-linearCutOff*max/255].
  //const inRange = 255
  l32 -= linearCutOff;
  uint32_t inRange = 255 - linearCutOff;
  uint32_t outRange = max - linearCutOff;
  uint32_t offset = inRange >> 1;
  uint32_t y = (l32*l32*l32 + offset) / inRange;
  return uint16_t((y*outRange+(offset*offset))/inRange/inRange + linearCutOff);
}

// ColorToAPA102 converts a color into the 4 bytes needed to control an APA-102
// LED.
//
// The response as seen by the human eye is very non-linear. The APA-102
// provides an overall brightness PWM but it is relatively slower and results
// in human visible flicker. On the other hand the minimal color (1/255) is
// still too intense at full brightness, so for very dark color, it is worth
// using the overall brightness PWM. The goal is to use brightness!=31 as
// little as possible.
//
// Global brightness frequency is 580Hz and color frequency at 19.2kHz.
// https://cpldcpu.wordpress.com/2014/08/27/apa102/
// Both are multiplicative, so brightness@50% and color@50% means an effective
// 25% duty cycle but it is not properly distributed, which is the main problem.
//
// It is unclear to me if brightness is exactly in 1/31 increment as I don't
// have an oscilloscope to confirm. Same for color in 1/255 increment.
//
// Each channel duty cycle ramps from 100% to 1/(31*255) == 1/7905.
//
// Return brighness, blue, green, red.
void ColorToAPA102(const Color &c, uint8_t* dst)  {
  uint16_t r = Ramp(c.R, 0);
  uint16_t g = Ramp(c.G, 0);
  uint16_t b = Ramp(c.B, 0);
  if (r <= 255 && g <= 255 && b <= 255) {
    dst[0] = 0xE0 + 1;
    dst[1] = b;
    dst[2] = g;
    dst[3] = r;
  } else if (r <= 511 && g <= 511 && b <= 511) {
    dst[0] = 0xE0 + 2;
    dst[1] = b>>1;
    dst[2] = g>>1;
    dst[3] = r>>1;
  } else if (r <= 1023 && g <= 1023 && b <= 1023) {
    dst[0] = 0xE0 + 4;
    dst[1] = (b+2)>>2;
    dst[2] = (g+2)>>2;
    dst[3] = (r+2)>>2;
  } else {
    // In this case we need to use a ramp of 255-1 even for lower colors.
    dst[0] = 0xE0 + 31;
    dst[1] = (b+15)/31;
    dst[2] = (g+15)/31;
    dst[3] = (r+15)/31;
  }
}

// Serializes converts a buffer of colors to the APA102 SPI format.
void Raster(const Frame& pixels, uint8_t *buf) {
  // https://cpldcpu.files.wordpress.com/2014/08/apa-102c-super-led-specifications-2014-en.pdf
  int numLights = pixels.len;
  // End frames are needed to be able to push enough SPI clock signals due to
  // internal half-delay of data signal from each individual LED. See
  // https://cpldcpu.wordpress.com/2014/11/30/understanding-the-apa102-superled/
  int l = 4*(numLights+1) + numLights/2/8 + 1;
  /*
  if (len(*buf) != l) {
    *buf = make([]byte, l)
    // It is not necessary to set the end frames to 0xFFFFFFFF.
    // Set end frames right away.
    //s := (*buf)[4+4*numLights:]
    //for i := range s {
    //	s[i] = 0xFF
    //}
  }
  */
  for (int i = 0; i < 4; i++) {
    buf[i] = 0;
  }
  // Start frame is all zeros. Just skip it.
  uint8_t *s = &buf[4];
  for (int i = 0; i < numLights; i++) {
    ColorToAPA102(pixels.c[i], &s[4*i]);
  }
  for (int i = 4+4*numLights; i < l; i++) {
    buf[i] = 0;
  }
}
