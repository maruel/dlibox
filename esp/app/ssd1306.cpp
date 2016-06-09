// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

#include "user_config.h"
#include <SmingCore/SmingCore.h>
#include <Wire.h>

#include "ada_ssd1306.h"
#include "conf.h"
#include "images.h"
#include "painter.h"
#include "perf.h"
#include "ssd1306.h"

#define OLED_RESET 0  // GPIO0
Adafruit_SSD1306 display(64, 48, OLED_RESET);

namespace {

const uint8_t* const images[] = {
  coeur,
  dragon,
};

Timer displayTimer;

int index;

}  // namespace


void cycle() {
  uint32_t now = micros();
  display.clearDisplay();
  //if (index == lengthof(images)) {
  if (true) {
    display.setCursor(0, 0);
    display.printf("Ovrhead ms");
    uint32_t s = Perf[LOAD_RENDER].sum();
    uint32_t x = s / 1000;
    uint8_t y = (s-x*1000)/100;
    display.printf("Rnd/s%3u.%1u", min(x, 1000u), y);
    s = Perf[LOAD_SPI].sum();
    x = s / 1000;
    y = (s-x*1000)/100;
    display.printf("SPI/s%3u.%1u", min(x, 1000u), y);
    s = Perf[LOAD_I2C].avg();
    x = s / 1000;
    y = (s-x*1000)/100;
    display.printf("I2C/f%3u.%1u", min(x, 1000u), y);
    display.printf("ms/f %5u", Perf[FRAMES].avgDelta());
    display.println(lastRenderName);
    // This is very slow. Should probably send a separate task for this since we
    // already monopolized for a long time!
    display.display();
  } else {
    display.drawBitmap(0, 0, images[index], display.width(), display.height(), 1);
  }
  display.display();
  // It's very close to 64ms limit!
  Perf[LOAD_I2C].add(micros() - now);
  index = (index + 1) % (lengthof(images)+1);
}

// Font size:
// - 1: 10 characters wide; 6 lines
// - 2: 5 characters wide; 3 lines
void initSSD1306() {
  if (config.display.enabled) {
    // Set for the wemos i²c pins.
    Wire.pins(5, 4);
    // TODO(maruel): Change speed according to config.display.I2Cspeed.
    display.begin();
    display.clearDisplay();
    display.setTextSize(1);
    display.setTextColor(WHITE);
    display.setCursor(0,0);
    display.println("dlibox");
    display.display();
    displayTimer.initializeMs(2000, cycle).start();
  }
}
