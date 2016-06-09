// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

#include "user_config.h"
#include <SmingCore/SmingCore.h>
#include <pb_encode.h>
#include <pb_decode.h>

#include "conf.h"

namespace {

const char *const CONFIG_FILE = "config";

}  // namespace

char chipID[9];
char hostName[7+sizeof(chipID)];
Config config;

void initConfig() {
  spiffs_mount();
  sprintf(chipID, "%08x", system_get_chip_id());
  sprintf(hostName, "dlibox-%s", chipID);

  if (fileExist(CONFIG_FILE)) {
    // TODO(maruel): It'd be nice to stream from spiffs to nanopb.
    // See nanopb/examples/network_server/common.c as an example.
    String raw = fileGetContent(CONFIG_FILE);
    pb_istream_t stream = pb_istream_from_buffer((const pb_byte_t*)raw.c_str(), raw.length());
    // pb_decode_noinit()
    if (!pb_decode(&stream, Config_fields, &config)) {
      memset(&config, 0, sizeof(config));
    }
  } else {
    /*
    strcpy(config.wifiClient.ssid, "AA");
    strcpy(config.wifiClient.password, "BB");
    config.wifiClient.has_ssid = true;
    config.wifiClient.has_password = true;
    config.has_wifiClient = true;
    */
  }
}

void saveConfig() {
  pb_byte_t buffer[Config_size];
  pb_ostream_t stream = pb_ostream_from_buffer(buffer, sizeof(buffer));
  if (pb_encode(&stream, Config_fields, &config)) {
    file_t file = fileOpen(CONFIG_FILE, eFO_CreateNewAlways | eFO_WriteOnly);
    fileWrite(file, buffer, stream.bytes_written);
    fileClose(file);
  }
}
