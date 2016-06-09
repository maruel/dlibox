// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

#ifndef __ANIM1D_H__
#define __ANIM1D_H__

#include "anim1d_msg.pb.h"

struct Frame;

struct IPattern {
  virtual ~IPattern() {};
  virtual void NextFrame(Frame& f, int timeMS) = 0;
};

// For compactness, differentiate bettwen the Pattern Color and a pixel color.
struct Color {
  uint8_t R;
  uint8_t G;
  uint8_t B;

  void from(const MColor& m) {
    B = m.color;
    G = m.color >> 8;
    R = m.color >> 16;
  }
};

struct Frame : public IPattern {
  Color *pixels;
  uint16_t len;

  Frame() : pixels(NULL), len(0) {
  }

  Frame(Color* c, uint16_t l) : pixels(c), len(l) {
  }

  virtual ~Frame() {
    delete pixels;
  }

  virtual void NextFrame(Frame& f, int timeMS) {
    memcpy(f.pixels, pixels, sizeof(Color)*f.len);
  }

  void from(const MFrame& m) {
    // Hackish, assume a very specific memory layout.
    uint16_t l = *(uint16_t*)m.colors.arg;
    reset(l);
    memcpy(pixels, (const char*)m.colors.arg+2, sizeof(Color)*len);
  }

  void reset(uint16_t l) {
    if (l == len) {
      return;
    }
    delete pixels;
    pixels = NULL;
    if (l != 0) {
      pixels = new Color[l];
    }
    len = l;
  }
};

struct PColor : public IPattern {
  Color c;

  PColor(const Color& c) : c(c) {};
  virtual ~PColor() {}

  virtual void NextFrame(Frame& f, int timeMS) {
    for (int i = 0; i < f.len; i++) {
      f.pixels[i] = c;
    }
  }

  void from(const MColor& m) {
    c.from(m);
  }
};


struct Rainbow : public IPattern {
  // TODO(maruel): Keep a local buffer for performance.

  virtual ~Rainbow() {}
  virtual void NextFrame(Frame& f, int timeMS);
};

#endif
