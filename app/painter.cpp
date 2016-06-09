// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

#include "user_config.h"
#include <SmingCore/SmingCore.h>
#include <SPI.h>

#include "apa102.h"
#include "anim1d.h"

namespace {

Timer paintTimer;

#define numLights 144
Frame buf{colors: new Color[numLights], len: numLights};

int start;

PColor gray(Color{0x7F, 0x7F, 0x7F});
Rainbow rainbow;
IPattern *p = &rainbow;
IPattern *pNew = NULL;

void painterLoop() {
  int now = millis();
  // TODO(maruel): Atomic.
  if (pNew != NULL) {
    delete p;
    p = pNew;
    pNew = NULL;
  }
  p->NextFrame(buf, now-start);
  Write(buf);
}

}  // namespace

void initPainter() {
  initAPA102();
  start = millis();
  // Run at ~15Hz to see how far we can push it.
  paintTimer.initializeMs(66, painterLoop).start();
}
