// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

#include "user_config.h"
#include <SmingCore/SmingCore.h>
#include <SPI.h>

#include "apa102.h"
#include "anim1d.h"
#include "conf.h"
#include "painter.h"
#include "perf.h"
#include "ssd1306.h"

namespace {

Timer paintTimer;

Frame buf;

uint32_t start;

// TODO(maruel): Make sure of as much const memory as possible? Since it's
// Harvard, it's probably a waste of time.
// TODO(maruel): Lots of static initialized in there, need to fix this.
Color red{0xFF, 0, 0};
Color white{0xFF, 0xFF, 0xFF};
PColor gray({0x7F, 0x7F, 0x7F});
Rainbow rainbow;
Rotate rainbowRotated(&rainbow, 60);
Color candyChunk[] = {white, white, white, white, white, red, red, red, red, red};
Frame candyBar(candyChunk, lengthof(candyChunk));
Repeated candyBarRepeated(candyBar);
Rotate candyBarRotated(&candyBarRepeated, 60);
IPattern* frames[] = {
  &rainbow,
  &rainbowRotated,
  &gray,
  &candyBarRepeated,
  &candyBarRotated,
};
Cycle cycle(frames, lengthof(frames), 3000);
IPattern* p = &cycle;
IPattern* pNew = NULL;

void painterLoop() {
  // We need both as we need 49 days roll over for NextFrame() call but we need
  // the precision for Perf.
  uint32_t nowMS = millis();
  uint32_t nowUS = micros();
  // TODO(maruel): Atomic.
  if (pNew != NULL) {
    delete p;
    p = pNew;
    pNew = NULL;
    start = nowMS;
  }
  Perf[FRAMES].add(nowMS);
  // It is not guaranteed that the IPattern draws on every pixel. Make sure
  // that pixels not drawn on are black.
  memset(buf.pixels, 0, sizeof(Color) * buf.len);
  // TODO(maruel): Memory fragmentation.
  lastName = p->NextFrame(buf, nowMS-start);
  uint32_t render = Write(buf, maxAPA102Out / 4);
  // Time taken to render.
  // Max that can be calculated is 64ms.
  Perf[LOAD_RENDER].add(render-nowUS);
}

}  // namespace

String lastName;

void initPainter() {
  if (config.apa102.frameRate && config.apa102.numLights) {
    initAPA102();
    start = millis();
    buf.reset(config.apa102.numLights);
    paintTimer.initializeMs(1000/config.apa102.frameRate, painterLoop).start();
  }
}
