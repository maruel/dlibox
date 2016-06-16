// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

#include "user_config.h"
#include <SmingCore/SmingCore.h>

#include "conf.h"
#include "ota.h"
#include "serialcmd.h"

namespace {

void cmdcat(const char* args) {
  if (args == NULL) {
    Serial.println("specify a file; list them with ls");
    return;
  }
  Serial.println(fileGetContent(args));
}

void cmdconfig(const char* args) {
  Serial.printf("wifiClient: %d\n", config.has_wifiClient);
  Serial.printf("wifiClient.ssid: \"%s\"\n", config.wifiClient.ssid);
  Serial.printf("wifiClient.password: \"%s\"\n", config.wifiClient.password);
  Serial.printf("wifiAP: %d\n", config.has_wifiAP);
  Serial.printf("wifiAP.ssid: \"%s\"\n", config.wifiAP.ssid);
  Serial.printf("wifiAP.password: \"%s\"\n", config.wifiAP.password);
  Serial.printf("apa102: %d\n", config.has_apa102);
  Serial.printf("apa102.frameRate: %d\n", config.apa102.frameRate);
  Serial.printf("apa102.numLights: %d\n", config.apa102.numLights);
  Serial.printf("apa102.SPIspeed: %d\n", config.apa102.SPIspeed);
  Serial.printf("host: %d\n", config.has_host);
  Serial.printf("host.name: \"%s\"\n", config.host.name);
  Serial.printf("host.highSpeed: %d\n", config.host.highSpeed);
  Serial.printf("display: %d\n", config.has_display);
  Serial.printf("display.enabled: %d\n", config.display.enabled);
  Serial.printf("display.I2Cspeed: %d\n", config.display.I2Cspeed);
  Serial.printf("romURL: \"%s\"\n", config.romURL);
}

void cmdset(char* args) {
  char* value = args;
  for (;; value++) {
    if (*value == 0) {
      Serial.println("please provide a value");
      return;
    }
    if (*value == ' ') {
      break;
    }
  }
  *value++ = 0;
  if (!strcmp(args, "wifiClient.ssid")) {
    config.has_wifiClient = true;
    config.wifiClient.has_ssid = true;
    strcpy(config.wifiClient.ssid, value);
  } else if (!strcmp(args, "wifiClient.password")) {
    config.has_wifiClient = true;
    config.wifiClient.has_password = true;
    strcpy(config.wifiClient.password, value);
  } else if (!strcmp(args, "wifiAP.ssid")) {
    config.has_wifiAP = true;
    config.wifiAP.has_ssid = true;
    strcpy(config.wifiAP.ssid, value);
  } else if (!strcmp(args, "wifiAP.password")) {
    config.has_wifiAP = true;
    config.wifiAP.has_password = true;
    strcpy(config.wifiAP.password, value);
  } else if (!strcmp(args, "apa102.frameRate")) {
    config.has_apa102 = true;
    config.apa102.has_frameRate = true;
    config.apa102.frameRate = atoi(value);
  } else if (!strcmp(args, "apa102.numLights")) {
    config.has_apa102 = true;
    config.apa102.has_numLights = true;
    config.apa102.numLights = atoi(value);
  } else {
    Serial.printf("invalid key \"%s\"\n", args);
    return;
  }
  saveConfig();
  Serial.println("Don't forget to restart for settings to take effect!");
}

void cmdhelp() {
  Serial.println();
  Serial.println("available commands:");
  Serial.println("  cat     - show a file in spiffs");
  Serial.println("  clear   - deletes the current configuration");
  Serial.println("  config  - display current config");
  Serial.println("  connect - connect to wifi");
  Serial.println("  help    - display this message");
  Serial.println("  info    - show esp8266 info");
  Serial.println("  ip      - show current ip address");
  Serial.println("  ls      - list files in spiffs");
  Serial.println("  ota     - perform ota update, switch rom and reboot");
  Serial.println("  restart - restart the esp8266");
  Serial.println("  set     - set a configuration value");
  Serial.println("  switch  - switch to the other rom and reboot");
  Serial.println();
}

void cmdinfo() {
  Serial.println();
  Serial.printf("SDK: v%s\n", system_get_sdk_version());
  Serial.printf("Free Heap: %d\n", system_get_free_heap_size());
  Serial.printf("CPU Frequency: %d MHz\n", system_get_cpu_freq());
  Serial.printf("System Chip ID: %x\n", system_get_chip_id());
  Serial.printf("SPI Flash ID: %x\n", spi_flash_get_id());
  Serial.printf("SPI Flash Size: %d\n", (1 << ((spi_flash_get_id() >> 16) & 0xff)));
}

void cmdls(const char* args) {
  Vector<String> files = fileList();
  Serial.printf("filecount %d\n", files.count());
  for (uint16_t i = 0; i < files.count(); i++) {
    Serial.println(files[i]);
  }
}

void onCommand(const char* cmd, char* args) {
  if (!strcmp(cmd, "cat")) {
    cmdcat(args);
  } else if (!strcmp(cmd, "clear")) {
    clearConfig();
  } else if (!strcmp(cmd, "config")) {
    cmdconfig(args);
  } else if (!strcmp(cmd, "connect")) {
    WifiStation.config(config.wifiClient.ssid, config.wifiClient.password);
    WifiStation.enable(true);
  } else if (!strcmp(cmd, "help")) {
    cmdhelp();
  } else if (!strcmp(cmd, "info")) {
    cmdinfo();
  } else if (!strcmp(cmd, "ip")) {
    Serial.printf("client ip: %s mac: %s\n", WifiStation.getIP().toString().c_str(), WifiStation.getMAC().c_str());
    Serial.printf("ap     ip: %s mac: %s\n", WifiAccessPoint.getIP().toString().c_str(), WifiAccessPoint.getMAC().c_str());
  } else if (!strcmp(cmd, "ls")) {
    cmdls(args);
  } else if (!strcmp(cmd, "ota")) {
    OtaUpdate();
  } else if (!strcmp(cmd, "restart")) {
    // TODO(maruel): Hangs instead...
    System.restart();
  } else if (!strcmp(cmd, "set")) {
    cmdset(args);
  } else if (!strcmp(cmd, "switch")) {
    SwitchROM();
  } else {
    Serial.printf("unknown command \"%s\"\n", cmd);
  }
}

void serialCallBack(Stream& stream, char arrivedChar, unsigned short availableCharsCount) {
  // Echo back.
  Serial.print(arrivedChar);

  if (arrivedChar == '\n') {
    char str[availableCharsCount];
    char* args = NULL;
    for (uint16_t i = 0; i < availableCharsCount; i++) {
      char c = stream.read();
      if (c == '\r' || c == '\n') {
        str[i] = '\0';
      } else if (c == ' ' && args == NULL) {
        str[i] = 0;
        args = &str[i+1];
      } else {
        str[i] = c;
      }
    }
    onCommand(str, args);
  }
}

}  // namespace

void initSerialCommand() {
  Serial.begin(SERIAL_BAUD_RATE);  // TODO(maruel): Needed?
  Serial.systemDebugOutput(false);
  system_set_os_print(0);  // may break stuff.
  // Serial.commandProcessing(false); ?
  Serial.printf("\nCurrently running rom %d.\n", rboot_get_current_rom());
  Serial.println("Type 'help' and press enter for instructions.");
  Serial.println();
  Serial.setCallback(serialCallBack);
}
