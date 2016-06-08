// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

#include "user_config.h"
#include "conf.h"
#include "ota.h"
#include <SmingCore/SmingCore.h>

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

}  // namespace

// TODO(maruel): Find out how to leverage mDNS + MQTT for this?
// TODO(maruel): Add minimal authentication. We're not talking about any serious
// security but just basic verification.
void OtaUpdate() {
  // This file is located at out/firmware/. Use:
  //   go get github.com/maruel/serve-dir && serve-dir -root out/firmware -port 8010
  Serial.printf("Updating from %s...", config.romURL);
  if (otaUpdater) {
    delete otaUpdater;
  }
  otaUpdater = new rBootHttpUpdate();
  rboot_config bootconf = rboot_get_config();
  otaUpdater->addItem(bootconf.roms[!bootconf.current_rom], config.romURL);
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
