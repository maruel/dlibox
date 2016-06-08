// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

#include "user_config.h"
#include "conf.h"
#include <SmingCore/SmingCore.h>
#include <pb_encode.h>
#include <pb_decode.h>

namespace {

const char *const CONFIG_FILE = "config";

}  // namespace

Config config;

void initConfig() {
  spiffs_mount();

  if (fileExist(CONFIG_FILE)) {
    // TODO(maruel): It'd be nice to stream from spiffs to nanopb.
    // See nanopb/examples/network_server/common.c as an example.
    String raw = fileGetContent(CONFIG_FILE);
    pb_istream_t stream = pb_istream_from_buffer((const pb_byte_t*)raw.c_str(), raw.length());
    // pb_decode_noinit()
    if (!pb_decode(&stream, Config_fields, &config)) {
      memset(&config, 0, sizeof(config));
    }
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
