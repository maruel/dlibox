// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

#include "user_config.h"
#include <SmingCore/SmingCore.h>

#include "anim1d.h"

namespace {

// waveLengthToRGB returns a color over a rainbow.
//
// This code was inspired by public domain code on the internet.
void waveLength2RGB(uint16_t w, Color* c) {
  if (w < 380) {
    c->R = 0;
    c->G = 0;
    c->B = 0;
  } else if (w < 420) {
    // Red peaks at 1/3 at 420.
    c->R = uint8(196 - (170*(440-w))/(440-380));
    c->G = 0;
    c->B = uint8(26 + (229*(w-380))/(420-380));
  } else if (w < 440) {
    c->R = uint8((0x89 * (440 - w)) / (440 - 420));
    c->G = 0;
    c->B = 255;
  } else if (w < 490) {
    c->R = 0;
    c->G = uint8((255 * (w - 440)) / (490 - 440));
    c->B = 255;
  } else if (w < 510) {
    c->R = 0;
    c->G = 255;
    c->B = uint8((255 * (510 - w)) / (510 - 490));
  } else if (w < 580) {
    c->R = uint8((255 * (w - 510)) / (580 - 510));
    c->G = 255;
    c->B = 0;
  } else if (w < 645) {
    c->R = 255;
    c->G = uint8((255 * (645 - w)) / (645 - 580));
    c->B = 0;
  } else if (w < 700) {
    c->R = 255;
    c->G = 0;
    c->B = 0;
  } else if (w < 781) {
    c->R = uint8(26 + (229*(780-w))/(780-700));
    c->G = 0;
    c->B = 0;
  } else {
    c->R = 0;
    c->G = 0;
    c->B = 0;
  }
}

}  // namespace


void Rainbow::NextFrame(Frame& f, int timeMS) {
  const uint16_t start = 380;
  const uint16_t end = 781;
  const uint16_t delta = end - start;
  // TODO(maruel): Change the scale to be nicer.
  //scale := logn(2)
  //step := 1. / float32(len(pixels))
  for (uint16_t i = 0; i < f.len; i++) {
    //j := log1p(float32(len(pixels)-i-1)*step) / scale
    uint16_t x = (delta*i+1)/f.len;
    waveLength2RGB(start + x, &f.pixels[i]);
  }
}

