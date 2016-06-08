// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

#ifndef __CONF_H__
#define __CONF_H__

#include "config.pb.h"

// init initializes SPIFFS and the config variable.
// It should be called first in the setup() function.
void initConfig();

// save saves the config to SPIFFS. Call it after modifying config.
void saveConfig();

extern Config config;

#endif
