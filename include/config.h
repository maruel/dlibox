// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

#ifndef __CONFIG_H__
#define __CONFIG_H__

namespace config {

// init initializes the config file. Make sure SPIFFS is initialized
// *before*.
void init();

// setValue sets a new value to one key.
void setValue(const char *key, const char *value);

// getValue returns a value or NULL. The pointer is temporary and will get
// recycled.
const char * getValue(const char *key);

}  // namespace config

#endif
