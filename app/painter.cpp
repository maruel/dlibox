// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

#include "user_config.h"
#include <SmingCore/SmingCore.h>
#include <SPI.h>

#include "apa102.h"
#include "anim1d.h"
#include "ssd1306.h"

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
Rotate rainbowRotated(&rainbow, 60);
Color candyChunk[] = {white, white, white, white, white, red, red, red, red, red};
Frame candyBar(candyChunk, sizeof(candyChunk)/sizeof(*candyChunk));
Repeated candyBarRepeated(candyBar);
Rotate candyBarRotated(&candyBarRepeated, 60);
IPattern *frames[] = {
  &rainbow,
  &rainbowRotated,
  &gray,
  &candyBarRepeated,
  &candyBarRotated,
};
Cycle cycle(frames, sizeof(frames)/sizeof(*frames), 3000);
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
  // It is not guaranteed that the IPattern draws on every pixel. Make sure that
  // pixels not drawn on are black.
  memset(buf.pixels, 0, sizeof(Color) * buf.len);
  String name = p->NextFrame(buf, now-start);
  Write(buf);
  display.setCursor(0, 0);
  display.clearDisplay();
  display.println("dlibox");
  display.println(name);
  display.display();
}

}  // namespace

void initPainter() {
  initAPA102();
  start = millis();
  // Run at ~15Hz to see how far we can push it.
  paintTimer.initializeMs(66, painterLoop).start();
}
