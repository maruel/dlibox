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

String lastName;

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
  // TODO(maruel): Atomic.
  uint32_t now = millis();
  if (pNew != NULL) {
    delete p;
    p = pNew;
    pNew = NULL;
    start = now;
  }
  Perf[FRAMES].add(now);
  if (config.apa102.numLights != buf.len) {
    buf.reset(config.apa102.numLights);
  } else {
    // It is not guaranteed that the IPattern draws on every pixel. Make sure
    // that pixels not drawn on are black.
    memset(buf.pixels, 0, sizeof(Color) * buf.len);
  }
  if (config.apa102.numLights != 0) {
    String name = p->NextFrame(buf, now-start);
    uint32_t render = Write(buf, maxAPA102Out / 4);
    // Time taken to render.
    Perf[LOAD_RENDER].add(render-now);
    if (name != lastName) {
      lastName = name;
      display.setCursor(0, 0);
      display.clearDisplay();
      display.printf("Ovrhead ms");
      display.printf("Rndr/s%4u", min(Perf[LOAD_RENDER].sum(), 1000u));
      display.printf("SPI/s %4u", min(Perf[LOAD_SPI].sum(), 1000u));
      display.printf("I2C/f %4u", min(Perf[LOAD_I2C].avg(), 1000u));
      display.printf("%ums/f\n", Perf[FRAMES].avgDelta());
      display.println(name);
      // This is very slow. Should probably send a separate task for this since we
      // already monopolized for a long time!
      now = millis();
      display.display();
      Perf[LOAD_I2C].add(millis() - now);
    }
  }
}

}  // namespace

void initPainter() {
  initAPA102();
  start = millis();
  paintTimer.initializeMs(1000/config.apa102.frameRate, painterLoop).start();
}
