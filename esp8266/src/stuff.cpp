// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

#include "nodes.h"

int isBool(const String& value) {
  if (value == "true" || value == "1") {
    return 1;
  }
  if (value == "false" || value == "0") {
    return 0;
  }
  return -1;
}

int toInt(const String& value, int min, int max) {
  int v = atoi(value.c_str());
  if (v < min) {
    return min;
  }
  if (v > max) {
    return max;
  }
  return v;
}
