// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

#include "user_config.h"
#include <SmingCore/SmingCore.h>
#include <SPI.h>

#include "conf.h"
#include "ota.h"
#include "painter.h"
#include "serialcmd.h"
#include "ssd1306.h"
#include "wifi.h"

#define LED_PIN 2 // GPIO2

//  digitalWrite(LED_PIN, state);

void init() {
  system_set_os_print(0);
  pinMode(LED_PIN, OUTPUT);
  initConfig();
  if (config.highSpeed) {
    // System.setCpuFrequency(eCF_160MHz); ?
    system_update_cpu_freq(SYS_CPU_160MHZ);
  }
  initSSD1306();
  initSerialCommand();
  resetWifi();
  initPainter();
}
