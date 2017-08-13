// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

#ifndef __ANIM1D_H__
#define __ANIM1D_H__

#include "anim1d_msg.pb.h"

struct Frame;

struct IPattern {
  virtual ~IPattern() {}
  // 49.71 days is enough for everyone! After, it will reset to 0.
  virtual String  Render(Frame& f, uint32_t timeMS) = 0;
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
    // TODO(maruel): Can't, this object is often copied.
    //delete pixels;
  }
  virtual String Render(Frame& f, uint32_t timeMS) {
    memcpy(f.pixels, pixels, sizeof(Color)*min(f.len, len));
    return "Frame";
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
      pixels = new Color[l]();
    }
    len = l;
  }
};

struct PColor : public IPattern {
  Color c;

  PColor(const Color& c) : c(c) {};
  virtual ~PColor() {}
  virtual String Render(Frame& f, uint32_t timeMS) {
    for (uint16_t i = 0; i < f.len; i++) {
      f.pixels[i] = c;
    }
    return "Color";
  }

  void from(const MColor& m) {
    c.from(m);
  }
};

struct Rainbow : public IPattern {
  // TODO(maruel): Keep a local buffer for performance.

  virtual ~Rainbow() {}
  virtual String Render(Frame& f, uint32_t timeMS);
};

struct Repeated : public IPattern {
  Frame frame;

  // TODO(maruel): The pixels are aliased.
  Repeated(const Frame& f) : frame(f) {}
  virtual ~Repeated() {}
  virtual String Render(Frame& f, uint32_t timeMS) {
    for (uint16_t i = 0; i < f.len; i += frame.len) {
      memcpy(f.pixels+i, frame.pixels, sizeof(Color)*min(f.len-i, frame.len));
    }
    return "Repeated";
  }
};

// Cycle cycles between multiple patterns. It can be used as an animatable
// looping frame.
//
// TODO(maruel): Blend between frames with Curve, defaults to step.
// TODO(maruel): Merge with Loop.
struct Cycle : public IPattern {
  IPattern **children;
  uint16_t nb_children;
  uint16_t durationMS;

  // Takes ownership of the patterns.
  Cycle(IPattern **c, uint16_t n, uint16_t d) : children(c), nb_children(n), durationMS(d) {
  }
  virtual ~Cycle() {
    for (uint16_t i = 0; i < nb_children; i++) {
      delete children[i];
    }
    delete children;
  }
  virtual String Render(Frame& f, uint32_t timeMS) {
    if (nb_children != 0) {
      uint32_t x = timeMS/uint32_t(durationMS);
      return children[x%nb_children]->Render(f, timeMS);
    }
    return "Cycle";
  }
};

// Rotate rotates a pattern that can also cycle either way.
//
// Use negative to go left. Can be used for 'candy bar'.
//
// Similar to PingPong{} except that it doesn't bounce.
//
// Use 5x oversampling with Scale{} to create smoother animation
struct Rotate : public IPattern {
  IPattern *child;
  uint16_t moveMS; // Expressed in duration of each frame.
  Frame buf;

  // Takes ownership of the pattern.
  Rotate(IPattern *c, uint16_t m) : child(c), moveMS(m) {
  }
  virtual ~Rotate() {
    delete child;
  }
  virtual const char * Name() {
    return "Rotate";
  }
  virtual String Render(Frame& f, uint32_t timeMS) {
    if (f.len == 0 || child == NULL) {
      return "Rotate{}";
    }
    buf.reset(f.len);
    String name("Rotate{");
    name += child->Render(buf, timeMS);
    uint16_t offset = (timeMS/uint32_t(moveMS)) % f.len;
    if (offset < 0) {
      offset = f.len + offset;
    }
    memcpy(&f.pixels[offset], buf.pixels, sizeof(Color)*(f.len-offset));
    memcpy(f.pixels, &buf.pixels[f.len-offset], sizeof(Color)*offset);
    name += "}";
    return name;
  }
};

#endif
