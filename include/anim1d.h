// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

#ifndef __ANIM1D_H__
#define __ANIM1D_H__

struct Color {
  uint8_t R;
  uint8_t G;
  uint8_t B;
};

struct Frame {
  Color *c;
  uint16_t len;
};

#endif
