// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

#ifndef __STUFF_H_
#define __STUFF_H_
#pragma once

#include <Arduino.h>

template<class T, size_t N>
inline size_t lengthof(T (&)[N]) { return N; }

#define DISALLOW_COPY_AND_ASSIGN(TypeName)                                     \
  TypeName(const TypeName &) = delete;                                         \
  void operator=(const TypeName &) = delete

// Converts a string to true/false value. Returns -1 for indeterminate.
int isBool(const String &v);
// Converts a string to integer, 
int toInt(const String &v, int min, int max);

#endif
