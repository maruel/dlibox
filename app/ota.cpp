// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

#include "user_config.h"
#include "ota.h"
#include <SmingCore/SmingCore.h>

// TODO(maruel): Find out how to leverage mDNS for this?
// TODO(maruel): Add minimal authentication. We're not talking about any serious
// security but just basic verification.
// This file is located at out/firmware/. Use:
//   go get github.com/maruel/serve-dir && cd out/firmware && serve-dir -port 8010
const char * const ROM_URL = "http://host:8010/rom0.bin";

namespace {

rBootHttpUpdate* otaUpdater = NULL;

void OtaUpdate_CallBack(bool result) {
  // TODO(maruel): Print SHA1 of firmware.
  if (result == true) {
    uint8 slot = !rboot_get_current_rom();
    Serial.printf("Firmware updated, rebooting to rom %d...\r\n", slot);
    rboot_set_current_rom(slot);
    System.restart();
  } else {
    Serial.println("Firmware update failed!");
    // TODO(maruel): Memory leak?
    //delete otaUpdater;
    //otaUpdater = NULL;
  }
}

void serialCallBack(Stream& stream, char arrivedChar, unsigned short availableCharsCount) {
  if (arrivedChar == '\n') {
    char str[availableCharsCount];
    for (int i = 0; i < availableCharsCount; i++) {
      str[i] = stream.read();
      if (str[i] == '\r' || str[i] == '\n') {
        str[i] = '\0';
      }
    }
    if (!strcmp(str, "cat")) {
      Vector<String> files = fileList();
      if (files.count() > 0) {
        Serial.printf("dumping file %s:\r\n", files[0].c_str());
        Serial.println(fileGetContent(files[0]));
      } else {
        Serial.println("Empty spiffs!");
      }
    } else if (!strcmp(str, "connect")) {
      // TODO(maruel): Get from config.
      const char * const WIFI_SSID = "XXX";
      const char * const WIFI_PWD = "YYY";
      WifiStation.config(WIFI_SSID, WIFI_PWD);
      WifiStation.enable(true);
    } else if (!strcmp(str, "help")) {
      Serial.println();
      Serial.println("available commands:");
      Serial.println("  cat     - show first file in spiffs");
      // config - display current config.
      // set <key> <value> - set a value in the config.
      Serial.println("  connect - connect to wifi");
      Serial.println("  help    - display this message");
      Serial.println("  info    - show esp8266 info");
      Serial.println("  ip      - show current ip address");
      Serial.println("  ls      - list files in spiffs");
      Serial.println("  ota     - perform ota update, switch rom and reboot");
      Serial.println("  restart - restart the esp8266");
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
      Serial.printf("ip: %s mac: %s\r\n", WifiStation.getIP().toString().c_str(), WifiStation.getMAC().c_str());
    } else if (!strcmp(str, "ls")) {
      Vector<String> files = fileList();
      Serial.printf("filecount %d\r\n", files.count());
      for (unsigned int i = 0; i < files.count(); i++) {
        Serial.println(files[i]);
      }
    } else if (!strcmp(str, "ota")) {
      OtaUpdate(ROM_URL);
    } else if (!strcmp(str, "restart")) {
      System.restart();
    } else if (!strcmp(str, "switch")) {
      SwitchROM();
    } else {
      Serial.println("unknown command");
    }
  }
}

}

void OtaUpdate(const char * rom_url) {
  // TODO(maruel): User configurable host then construct the URLs.
  Serial.println("Updating...");
  if (otaUpdater) {
    delete otaUpdater;
  }
  otaUpdater = new rBootHttpUpdate();
  rboot_config bootconf = rboot_get_config();
  otaUpdater->addItem(bootconf.roms[!bootconf.current_rom], rom_url);
  otaUpdater->setCallback(OtaUpdate_CallBack);
  otaUpdater->start();
}

void SwitchROM() {
  uint8 before = rboot_get_current_rom();
  uint8 after = !before;
  Serial.printf("Swapping from rom %d to rom %d.\r\n", before, after);
  rboot_set_current_rom(after);
  Serial.println("Restarting...\r\n");
  System.restart();
}

void initSerialCommand() {
  Serial.begin(SERIAL_BAUD_RATE);  // TODO(maruel): Needed?
  Serial.systemDebugOutput(true);
  Serial.printf("\r\nCurrently running rom %d.\r\n", rboot_get_current_rom());
  Serial.println("Type 'help' and press enter for instructions.");
  Serial.println();
  Serial.setCallback(serialCallBack);
}
