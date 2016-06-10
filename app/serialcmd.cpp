// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

#include "user_config.h"
#include <SmingCore/SmingCore.h>

#include "conf.h"
#include "ota.h"
#include "serialcmd.h"

namespace {

void onCommand(const char * str) {
  if (!strcmp(str, "cat")) {
    Vector<String> files = fileList();
    if (files.count() > 0) {
      Serial.printf("dumping file %s:\r\n", files[0].c_str());
      Serial.println(fileGetContent(files[0]));
    } else {
      Serial.println("Empty spiffs!");
    }
  } else if (!strcmp(str, "config")) {
    Serial.printf("wifiClient: %d\r\n", config.has_wifiClient);
    Serial.printf("wifiAP: %d\r\n", config.has_wifiAP);
  } else if (!strcmp(str, "connect")) {
    WifiStation.config(config.wifiClient.ssid, config.wifiClient.password);
    WifiStation.enable(true);
  } else if (!strcmp(str, "help")) {
    Serial.println();
    Serial.println("available commands:");
    Serial.println("  cat     - show first file in spiffs");
    Serial.println("  config    - display current config");
    Serial.println("  connect - connect to wifi");
    // get <key> - get config value.
    Serial.println("  help    - display this message");
    Serial.println("  info    - show esp8266 info");
    Serial.println("  ip      - show current ip address");
    Serial.println("  ls      - list files in spiffs");
    Serial.println("  ota     - perform ota update, switch rom and reboot");
    Serial.println("  restart - restart the esp8266");
    // set <key> <value> - set a value in the config.
    Serial.println("  switch  - switch to the other rom and reboot");
    Serial.println();
  } else if (!strcmp(str, "info")) {
    Serial.printf("\r\nSDK: v%s\r\n", system_get_sdk_version());
    Serial.printf("Free Heap: %d\r\n", system_get_free_heap_size());
    Serial.printf("CPU Frequency: %d MHz\r\n", system_get_cpu_freq());
    Serial.printf("System Chip ID: %x\r\n", system_get_chip_id());
    Serial.printf("SPI Flash ID: %x\r\n", spi_flash_get_id());
    //Serial.printf("SPI Flash Size: %d\r\n", (1 << ((spi_flash_get_id() >> 16) & 0xff)));
  } else if (!strcmp(str, "ip")) {
    Serial.printf("client ip: %s mac: %s\r\n", WifiStation.getIP().toString().c_str(), WifiStation.getMAC().c_str());
    Serial.printf("ap     ip: %s mac: %s\r\n", WifiAccessPoint.getIP().toString().c_str(), WifiAccessPoint.getMAC().c_str());
  } else if (!strcmp(str, "ls")) {
    Vector<String> files = fileList();
    Serial.printf("filecount %d\r\n", files.count());
    for (uint16_t i = 0; i < files.count(); i++) {
      Serial.println(files[i]);
    }
  } else if (!strcmp(str, "ota")) {
    OtaUpdate();
  } else if (!strcmp(str, "restart")) {
    System.restart();
  } else if (!strcmp(str, "switch")) {
    SwitchROM();
  } else {
    Serial.printf("unknown command \"%s\"\r\n", str);
  }
}

void serialCallBack(Stream& stream, char arrivedChar, unsigned short availableCharsCount) {
  if (arrivedChar == '\n') {
    char str[availableCharsCount];
    for (uint16_t i = 0; i < availableCharsCount; i++) {
      // TODO(maruel): Ugh.
      str[i] = stream.read();
      if (str[i] == '\r' || str[i] == '\n') {
        str[i] = '\0';
      }
    }
    onCommand(str);
  }
}

}  // namespace

void initSerialCommand() {
  Serial.begin(SERIAL_BAUD_RATE);  // TODO(maruel): Needed?
  Serial.systemDebugOutput(false);
  system_set_os_print(0);  // may break stuff.
  // Serial.commandProcessing(false); ?
  Serial.printf("\r\nCurrently running rom %d.\r\n", rboot_get_current_rom());
  Serial.println("Type 'help' and press enter for instructions.");
  Serial.println();
  Serial.setCallback(serialCallBack);
}
