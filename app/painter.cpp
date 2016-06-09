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

// TODO(maruel): Make sure of as much const memory as possible? Since it's
// Harvard, it's probably a waste of time.
// TODO(maruel): Lots of static initialized in there, need to fix this.
Color red{0xFF, 0, 0};
Color white{0xFF, 0xFF, 0xFF};
PColor gray({0x7F, 0x7F, 0x7F});
Rainbow rainbow;
Color candyChunk[] = {white, white, white, white, white, red, red, red, red, red};
Repeated candyPartR(Frame(candyChunk, sizeof(candyChunk)));
Rotate candyBar(&candyPartR, 60);
IPattern *frames[] = {&rainbow, &gray, &candyBar};
Cycle cycle(frames, sizeof(frames), 1000);
IPattern *p = &cycle;
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
