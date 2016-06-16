// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

#include "user_config.h"
#include <SmingCore/SmingCore.h>
#include <SPI.h>

#include "conf.h"
#include "ota.h"
#include "painter.h"
#include "perf.h"
#include "serialcmd.h"
#include "ssd1306.h"
#include "wifi.h"

#define LED_PIN 2 // GPIO2

namespace {

void onReady() {
  digitalWrite(LED_PIN, 1);
  initConfig();
  initSerialCommand();
  initPerf();
  initSSD1306();
  initPainter();
  initWifi();
  digitalWrite(LED_PIN, 0);
}

}  // namespace

void init() {
  system_set_os_print(0);
  pinMode(LED_PIN, OUTPUT);
  // The system is ready a few millisecond later. It is possible the system
  // boots for rboot OTA update (?) so don't do anything stupid before being
  // ready.
  System.onReady(onReady);
}
