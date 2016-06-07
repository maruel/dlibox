// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

#include "user_config.h"
#include "config.h"
#include <SmingCore/SmingCore.h>

namespace config {

const char *const CONFIG_FILE = "config";

String config;

void init() {
  if (fileExist(CONFIG_FILE)) {
    config = fileGetContent(CONFIG_FILE);
  } else {
    config = "{}";
  }
}

// setValue sets a new value to one key.
void setValue(const char *key, const char *value) {
  // TODO(maruel): Slow as hell.
  DynamicJsonBuffer jsonBuffer;
  JsonObject& root = jsonBuffer.parseObject(config);
  root[key] = value;
  root.printTo(config);
  fileSetContent(CONFIG_FILE, config);
}

const char * getValue(const char *key) {
  // TODO(maruel): Slow as hell.
  DynamicJsonBuffer jsonBuffer;
  JsonObject& root = jsonBuffer.parseObject(config);
  return (const char*)root[key];
}

}  // namespace config
