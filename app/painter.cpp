// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

#include "user_config.h"
#include <SmingCore/SmingCore.h>
#include <SPI.h>

#include "apa102.h"
#include "anim1d.h"
#include "painter.h"
#include "ssd1306.h"

namespace {

template<int N_>
struct Samples {
  static const uint16_t N = N_;
  uint16_t samples[N];
  uint16_t index;

  void add(uint16_t t) {
    samples[index] = t;
    index = (index+1) % N;
  }

  uint32_t sum() {
    uint32_t s = 0;
    for (uint16_t i = 0; i < N; i++) {
      s += uint32_t(samples[i]);
    }
    return s;
  }

  uint16_t avg() {
    return uint16_t(sum() / uint32_t(N));
  }

  // Return value should be divided by N-1.
  uint16_t sumDelta() {
    uint16_t s = 0;
    for (uint16_t i = 0; i < N; i++) {
      if (i != index) {
        uint16_t j = (i + N - 1) % N;
        s += (samples[i] - samples[j]);
      }
    }
    return s;
  }

  uint16_t avgDelta() {
    return sumDelta() / (N-1);
  }
};

const int FrameRate = 60;
Samples<FrameRate> loadRender;
Samples<FrameRate> loadSPI;
Samples<FrameRate*2> timestamps;

String lastName;

Timer paintTimer;

#define numLights 144
Frame buf{colors: new Color[numLights], len: numLights};

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
IPattern *frames[] = {
  &rainbow,
  &rainbowRotated,
  &gray,
  &candyBarRepeated,
  &candyBarRotated,
};
Cycle cycle(frames, lengthof(frames), 3000);
IPattern *p = &cycle;
IPattern *pNew = NULL;

void painterLoop() {
  // TODO(maruel): Atomic.
  if (pNew != NULL) {
    delete p;
    p = pNew;
    pNew = NULL;
  }
  uint32_t now = millis();
  // It is not guaranteed that the IPattern draws on every pixel. Make sure that
  // pixels not drawn on are black.
  memset(buf.pixels, 0, sizeof(Color) * buf.len);
  String name = p->NextFrame(buf, now-start);
  int32_t render = millis();
  Write(buf, maxAPA102Out / 4);
  uint32_t spi = millis();
  timestamps.add(now);
  loadRender.add(render-now);
  loadSPI.add(spi-render);
  if (name != lastName) {
    lastName = name;
    display.setCursor(0, 0);
    display.clearDisplay();
    display.printf("dlibox  ms");
    display.printf("Render%4u", loadRender.sum());
    display.printf("SPI  %5u", loadSPI.sum());
    display.printf("%ums/frame", timestamps.avgDelta());
    display.println(name);
    // This is very slow. Should probably send a separate task for this since we
    // already monopolized for a long time!
    display.display();
  }
}

}  // namespace

void initPainter() {
  initAPA102();
  start = millis();
  paintTimer.initializeMs(1000/FrameRate, painterLoop).start();
}
