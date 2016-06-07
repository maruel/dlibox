// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

#include "user_config.h"
#include "ota.h"
#include <SmingCore/SmingCore.h>
#include <SPI.h>

#include "anim1d.h"
#include "apa102.h"

#define LED_PIN 2 // GPIO2

Timer procTimer;
bool state = true;

const uint8_t Rainbow_[] = {
  0x20, 0x00, 0x25, 0x26, 0x00, 0x30, 0x2b, 0x00, 0x3c, 0x31, 0x00, 0x47, 0x37, 0x00, 0x53,
  0x3c, 0x00, 0x5e, 0x42, 0x00, 0x6a, 0x48, 0x00, 0x75, 0x4d, 0x00, 0x81, 0x53, 0x00, 0x8c,
  0x59, 0x00, 0x97, 0x5e, 0x00, 0xa3, 0x64, 0x00, 0xae, 0x6a, 0x00, 0xba, 0x6f, 0x00, 0xc5,
  0x78, 0x00, 0xd6, 0x7e, 0x00, 0xe2, 0x83, 0x00, 0xed, 0x89, 0x00, 0xf9, 0x82, 0x00, 0xff,
  0x74, 0x00, 0xff, 0x66, 0x00, 0xff, 0x52, 0x00, 0xff, 0x44, 0x00, 0xff, 0x36, 0x00, 0xff,
  0x29, 0x00, 0xff, 0x1b, 0x00, 0xff, 0x06, 0x00, 0xff, 0x00, 0x05, 0xff, 0x00, 0x0f, 0xff,
  0x00, 0x19, 0xff, 0x00, 0x28, 0xff, 0x00, 0x33, 0xff, 0x00, 0x3d, 0xff, 0x00, 0x47, 0xff,
  0x00, 0x56, 0xff, 0x00, 0x60, 0xff, 0x00, 0x6b, 0xff, 0x00, 0x7a, 0xff, 0x00, 0x84, 0xff,
  0x00, 0x8e, 0xff, 0x00, 0x9e, 0xff, 0x00, 0xa8, 0xff, 0x00, 0xb2, 0xff, 0x00, 0xc1, 0xff,
  0x00, 0xcc, 0xff, 0x00, 0xdb, 0xff, 0x00, 0xe5, 0xff, 0x00, 0xef, 0xff, 0x00, 0xff, 0xff,
  0x00, 0xff, 0xe5, 0x00, 0xff, 0xbf, 0x00, 0xff, 0xa5, 0x00, 0xff, 0x7f, 0x00, 0xff, 0x66,
  0x00, 0xff, 0x3f, 0x00, 0xff, 0x26, 0x00, 0xff, 0x00, 0x07, 0xff, 0x00, 0x12, 0xff, 0x00,
  0x19, 0xff, 0x00, 0x24, 0xff, 0x00, 0x2b, 0xff, 0x00, 0x36, 0xff, 0x00, 0x3d, 0xff, 0x00,
  0x48, 0xff, 0x00, 0x53, 0xff, 0x00, 0x5b, 0xff, 0x00, 0x66, 0xff, 0x00, 0x70, 0xff, 0x00,
  0x78, 0xff, 0x00, 0x83, 0xff, 0x00, 0x8e, 0xff, 0x00, 0x95, 0xff, 0x00, 0xa0, 0xff, 0x00,
  0xab, 0xff, 0x00, 0xb2, 0xff, 0x00, 0xbd, 0xff, 0x00, 0xc8, 0xff, 0x00, 0xd3, 0xff, 0x00,
  0xde, 0xff, 0x00, 0xe5, 0xff, 0x00, 0xf0, 0xff, 0x00, 0xfb, 0xff, 0x00, 0xff, 0xf7, 0x00,
  0xff, 0xeb, 0x00, 0xff, 0xdf, 0x00, 0xff, 0xd7, 0x00, 0xff, 0xcc, 0x00, 0xff, 0xc0, 0x00,
  0xff, 0xb4, 0x00, 0xff, 0xa8, 0x00, 0xff, 0x9c, 0x00, 0xff, 0x91, 0x00, 0xff, 0x85, 0x00,
  0xff, 0x79, 0x00, 0xff, 0x6d, 0x00, 0xff, 0x62, 0x00, 0xff, 0x56, 0x00, 0xff, 0x4a, 0x00,
  0xff, 0x3e, 0x00, 0xff, 0x33, 0x00, 0xff, 0x23, 0x00, 0xff, 0x17, 0x00, 0xff, 0x0b, 0x00,
  0xff, 0x00, 0x00, 0xff, 0x00, 0x00, 0xff, 0x00, 0x00, 0xff, 0x00, 0x00, 0xff, 0x00, 0x00,
  0xff, 0x00, 0x00, 0xff, 0x00, 0x00, 0xff, 0x00, 0x00, 0xff, 0x00, 0x00, 0xff, 0x00, 0x00,
  0xff, 0x00, 0x00, 0xff, 0x00, 0x00, 0xff, 0x00, 0x00, 0xff, 0x00, 0x00, 0xff, 0x00, 0x00,
  0xff, 0x00, 0x00, 0xff, 0x00, 0x00, 0xf9, 0x00, 0x00, 0xf0, 0x00, 0x00, 0xe5, 0x00, 0x00,
  0xdc, 0x00, 0x00, 0xd1, 0x00, 0x00, 0xc5, 0x00, 0x00, 0xbd, 0x00, 0x00, 0xb1, 0x00, 0x00,
  0xa9, 0x00, 0x00, 0x9d, 0x00, 0x00, 0x92, 0x00, 0x00, 0x86, 0x00, 0x00, 0x7e, 0x00, 0x00,
  0x72, 0x00, 0x00, 0x67, 0x00, 0x00, 0x5b, 0x00, 0x00, 0x50, 0x00, 0x00, 0x44, 0x00, 0x00,
  0x39, 0x00, 0x00, 0x2e, 0x00, 0x00, 0x25, 0x00, 0x00, 0x00, 0x00, 0x00,
};

const Frame Rainbow = Frame{(Color *)(Rainbow_), sizeof(Rainbow_)/3};

const uint8_t Gray_[] = {
  0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f,
  0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f,
  0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f,
  0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f,
  0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f,
  0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f,
  0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f,
  0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f,
  0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f,
  0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f,
  0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f,
  0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f,
  0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f,
  0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f,
  0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f,
  0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f,
  0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f,
  0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f,
  0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f,
  0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f,
  0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f,
  0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f,
  0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f,
  0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f,
  0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f,
  0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f,
  0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f,
  0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f,
  0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f,
};

const Frame Gray = Frame{(Color *)(Gray_), sizeof(Gray_)/3};

#define numLights 144
uint8_t t[4*(numLights+1) + numLights/2/8 + 1];

void blink() {
  digitalWrite(LED_PIN, state);
  state = !state;
  if (state) {
    Raster(Rainbow, t);
  } else {
    Raster(Gray, t);
  }
  // TODO(maruel): Use an asynchronous version.
  // TODO(maruel): Use a writeBytes() that doesn't overwrite the buffer.
  // TODO(maruel): Start rendering the next buffer; use double-buffering.
  SPI.transfer(t, sizeof(t));
}

void init() {
  pinMode(LED_PIN, OUTPUT);
  //system_update_cpu_freq(SYS_CPU_160MHZ);
  //wifi_set_sleep_type(NONE_SLEEP_T);
  spiffs_mount();
  initSerialCommand();
  WifiAccessPoint.enable(false);
  SPI.begin();

  // Run at ~60Hz to see how far we can push it.
  procTimer.initializeMs(33, blink).start();
}
