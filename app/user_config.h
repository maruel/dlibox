// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

#ifndef __USER_CONFIG_H__
#define __USER_CONFIG_H__

#ifdef __cplusplus
extern "C" {
#endif

  // UART config
  // Running at high speed seems to make the connection flaky.
  //#define SERIAL_BAUD_RATE 921600
  #define SERIAL_BAUD_RATE 115200

  // ESP SDK config
  #define LWIP_OPEN_SRC
  #define USE_US_TIMER

  // Default types
  #define __CORRECT_ISO_CPP_STDLIB_H_PROTO
  #include <limits.h>
  #include <stdint.h>

  // Override c_types.h include and remove buggy espconn
  #define _C_TYPES_H_
  #define _NO_ESPCON_

  // Updated, compatible version of c_types.h
  // Just removed types declared in <stdint.h>
  #include <espinc/c_types_compatible.h>

  // System API declarations
  #include <esp_systemapi.h>

  // C++ Support
  #include <esp_cplusplus.h>
  // Extended string conversion for compatibility
  #include <stringconversion.h>
  // Network base API
  #include <espinc/lwip_includes.h>

  // Beta boards
  #define BOARD_ESP01

#ifdef __cplusplus
}
#endif

#endif
