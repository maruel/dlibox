// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

#include "user_config.h"
#include <SmingCore/SmingCore.h>

#include "conf.h"
#include "perf.h"

Samples Perf[LAST_PERF];

uint32_t frameCount;

void initPerf() {
  Perf[FRAMES].init(config.apa102.frameRate*2);
  Perf[LOAD_RENDER].init(config.apa102.frameRate);
  Perf[LOAD_SPI].init(config.apa102.frameRate);
  Perf[LOAD_I2C].init(5);
  if (config.host.highSpeed) {
    System.setCpuFrequency(eCF_160MHz);
    //system_update_cpu_freq(SYS_CPU_160MHZ);
  }
}

