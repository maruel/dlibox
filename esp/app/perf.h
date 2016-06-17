// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

#ifndef __PERF_H__
#define __PERF_H__

/*
template<int N_>
struct SamplesT {
  static const uint16_t N = N_;
  uint16_t samples[N];
  uint16_t index;

  void add(uint16_t t) {
    if (N) {
      samples[index] = t;
      index = (index+1) % N;
    }
  }

  uint32_t sum() {
    uint32_t s = 0;
    for (uint16_t i = 0; i < N; i++) {
      s += uint32_t(samples[i]);
    }
    return s;
  }

  uint16_t avg() {
    if (N) {
      return uint16_t(sum() / uint32_t(N));
    }
    return 0;
  }

  // Return value should be divided by N-1.
  uint16_t sumDelta() {
    uint16_t s = 0;
    for (uint16_t i = 0; i < N; i++) {
      if (i != index) {
        uint16_t j = (i + N - 1) % N;
        s += (samples[i] - samples[j]);
      }
    }
    return s;
  }

  uint16_t avgDelta() {
    if (N>1) {
      return sumDelta() / (N-1);
    }
    return 0;
  }
};
*/

struct Samples {
  uint16_t N;
  uint16_t index;
  uint16_t* samples;

  Samples() : N(0), index(0), samples(NULL) {}
  ~Samples() { delete samples; }

  void init(uint16_t n) {
    delete samples;
    N = n;
    index = 0;
    samples = new uint16_t[n]();
  }

  void add(uint16_t t) {
    if (N) {
      samples[index] = t;
      index = (index+1) % N;
    }
  }

  uint32_t sum() {
    uint32_t s = 0;
    for (uint16_t i = 0; i < N; i++) {
      s += uint32_t(samples[i]);
    }
    return s;
  }

  uint16_t avg() {
    if (N) {
      return uint16_t(sum() / uint32_t(N));
    }
    return 0;
  }

  // Return value should be divided by N-1.
  uint16_t sumDelta() {
    uint16_t s = 0;
    for (uint16_t i = 0; i < N; i++) {
      if (i != index) {
        uint16_t j = (i + N - 1) % N;
        s += (samples[i] - samples[j]);
      }
    }
    return s;
  }

  uint16_t avgDelta() {
    if (N>1) {
      return sumDelta() / (N-1);
    }
    return 0;
  }

private:
  Samples(const Samples&);
};

enum PerfMeasurement : uint8_t {
  FRAMES = 0,
  LOAD_RENDER,
  LOAD_SPI,
  LOAD_I2C,
  LAST_PERF,
};

extern Samples Perf[LAST_PERF];
extern uint32_t frameCount;

// Must be called after initializing config and should be called on reconfig.
void initPerf();

#endif
