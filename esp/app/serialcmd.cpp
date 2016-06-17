// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

#include "user_config.h"
#include <SmingCore/SmingCore.h>

#include "conf.h"
#include "painter.h"
#include "perf.h"
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
  Serial.printf("host.verbose: %d\n", config.host.verbose);
  Serial.printf("display: %d\n", config.has_display);
  Serial.printf("display.enabled: %d\n", config.display.enabled);
  Serial.printf("display.I2Cspeed: %d\n", config.display.I2Cspeed);
  Serial.printf("romURL: \"%s\"\n", config.romURL);
  Serial.printf("mqtt.host: \"%s\"\n", config.mqtt.host);
  Serial.printf("mqtt.port: %d\n", config.mqtt.port);
  Serial.printf("mqtt.username: \"%s\"\n", config.mqtt.username);
  Serial.printf("mqtt.password: \"%s\"\n", config.mqtt.password);
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
  } else if (!strcmp(args, "apa102.SPIspeed")) {
    config.has_apa102 = true;
    config.apa102.has_SPIspeed = true;
    config.apa102.SPIspeed = atoi(value);
  } else if (!strcmp(args, "host.name")) {
    config.has_host = true;
    config.host.has_name = true;
    strcpy(config.host.name, value);
  } else if (!strcmp(args, "host.highSpeed")) {
    config.has_host = true;
    config.host.has_highSpeed = true;
    config.host.highSpeed = atoi(value);
  } else if (!strcmp(args, "host.verbose")) {
    config.has_host = true;
    config.host.has_verbose = true;
    config.host.verbose = atoi(value);
  } else if (!strcmp(args, "display.enabled")) {
    config.has_display = true;
    config.display.has_enabled = true;
    config.display.enabled = atoi(value);
  } else if (!strcmp(args, "display.I2Cspeed")) {
    config.has_display = true;
    config.display.has_I2Cspeed = true;
    config.display.I2Cspeed = atoi(value);
  } else if (!strcmp(args, "romURL")) {
    config.has_romURL = true;
    strcpy(config.romURL, value);
  } else if (!strcmp(args, "mqtt.host")) {
    config.has_mqtt = true;
    config.mqtt.has_host = true;
    strcpy(config.mqtt.host, value);
  } else if (!strcmp(args, "mqtt.port")) {
    config.has_mqtt = true;
    config.mqtt.has_port = true;
    config.mqtt.port = atoi(value);
  } else if (!strcmp(args, "mqtt.username")) {
    config.has_mqtt = true;
    config.mqtt.has_username = true;
    strcpy(config.mqtt.username, value);
  } else if (!strcmp(args, "mqtt.password")) {
    config.has_mqtt = true;
    config.mqtt.has_password = true;
    strcpy(config.mqtt.password, value);
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
  Serial.println("  cat <file>        - show a file in spiffs");
  Serial.println("  clear             - deletes the current configuration");
  Serial.println("  config            - display current config");
  Serial.println("  connect           - connect to wifi");
  Serial.println("  format            - format the SPIFFS partition");
  Serial.println("  help              - display this message");
  Serial.println("  info              - show esp8266 and connectivity info");
  Serial.println("  ls                - list files in spiffs");
  Serial.println("  ota               - perform ota update, switch rom and reboot");
  Serial.println("  restart           - restart the esp8266");
  Serial.println("  set <key> <value> - set a configuration value");
  Serial.println("  switch            - switch to the other rom and reboot");
  Serial.println();
}

void cmdinfo() {
  Serial.println();
  Serial.printf("SDK:            v%s\n", system_get_sdk_version());
  Serial.printf("Free Heap:      %d\n", system_get_free_heap_size());
  Serial.printf("CPU Frequency:  %d MHz\n", system_get_cpu_freq());
  Serial.printf("System Chip ID: %x\n", system_get_chip_id());
  Serial.printf("SPI Flash ID:   %x\n", spi_flash_get_id());
  uint32_t size = (1 << ((spi_flash_get_id() >> 16) & 0xff));
  Serial.printf("SPI Flash Size: 0x%x (%d)\n", size, size);
  spiffs_config s = spiffs_get_storage_config();
  Serial.printf("SPIFFS Size:    0x%x (%d)\n", s.phys_size, s.phys_size);
  Serial.printf("SPIFFS Address: 0x%x\n", s.phys_addr);
  Serial.printf("SPIFFS Erase:   0x%x\n", s.phys_erase_block);
  Serial.printf("SPIFFS Block:   0x%x\n", s.log_block_size);
  Serial.printf("SPIFFS Page:    0x%x\n", s.log_page_size);
  Serial.println();
  Serial.printf("Wifi client enabled: %d\n", WifiStation.isEnabled());
  Serial.printf("Wifi client SSID:    %s\n", WifiStation.getSSID().c_str());
  Serial.printf("Wifi client IP:      %s\n", WifiStation.getIP().toString().c_str());
  uint8 hwaddr[6] = {0};
  wifi_get_macaddr(STATION_IF, hwaddr);
  Serial.printf("Wifi client MAC:     " MACSTR "\n", MAC2STR(hwaddr));
  Serial.printf("Wifi client RSSI:    %d dBm\n", WifiStation.getRssi());
  Serial.printf("Wifi client channel: %d\n", WifiStation.getChannel());
  Serial.printf("AccessPoint enabled: %d\n", WifiAccessPoint.isEnabled());
  Serial.printf("AccessPoint IP:      %s\n", WifiAccessPoint.getIP().toString().c_str());
  wifi_get_macaddr(SOFTAP_IF, hwaddr);
  Serial.printf("AccessPoint MAC:     " MACSTR "\n", MAC2STR(hwaddr));
}

void cmdls(const char* args) {
  Vector<String> files = fileList();
  Serial.printf("filecount %d\n", files.count());
  for (uint16_t i = 0; i < files.count(); i++) {
    Serial.println(files[i]);
  }
}

void cmdperf() {
  Serial.println("Ovrhead ms");
  uint32_t s = Perf[LOAD_RENDER].sum();
  uint32_t x = s / 1000;
  uint8_t y = (s-x*1000)/100;
  Serial.printf("Rnd/s%3u.%1u\n", min(x, 1000u), y);
  s = Perf[LOAD_SPI].sum();
  x = s / 1000;
  y = (s-x*1000)/100;
  Serial.printf("SPI/s%3u.%1u\n", min(x, 1000u), y);
  s = Perf[LOAD_I2C].avg();
  x = s / 1000;
  y = (s-x*1000)/100;
  Serial.printf("I2C/f%3u.%1u\n", min(x, 1000u), y);
  Serial.printf("ms/f %5u\n", Perf[FRAMES].avgDelta());
  Serial.println(lastRenderName);
}

void onCommand(const char* cmd, char* args) {
  if (!strcmp(cmd, "cat")) {
    cmdcat(args);
  } else if (!strcmp(cmd, "clear")) {
    clearConfig();
  } else if (!strcmp(cmd, "config")) {
    cmdconfig(args);
  } else if (!strcmp(cmd, "connect")) {
    if (config.wifiClient.ssid[0]) {
      WifiStation.config(config.wifiClient.ssid, config.wifiClient.password);
      WifiStation.enable(true);
    } else {
      Serial.println("wifi client not set, use 'set'");
    }
  } else if (!strcmp(cmd, "format")) {
    spiffs_format();
    Serial.println("SPIFFS formatted.");
  } else if (!strcmp(cmd, "help")) {
    cmdhelp();
  } else if (!strcmp(cmd, "info")) {
    cmdinfo();
  } else if (!strcmp(cmd, "ls")) {
    cmdls(args);
  } else if (!strcmp(cmd, "ota")) {
    OtaUpdate();
  } else if (!strcmp(cmd, "perf")) {
    cmdperf();
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
  // Serial.commandProcessing(false); ?
  Serial.printf("\nCurrently running rom %d.\n", rboot_get_current_rom());
  Serial.println("Type 'help' and press enter for instructions.");
  Serial.println();
  Serial.setCallback(serialCallBack);
}
