// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

#include <Homie.h>
#include "stuff.h"

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
  DISALLOW_COPY_AND_ASSIGN(Samples);
};

class Ticks {
public:
  Ticks() : value_(0) {}

  void tick() {
    value_++;
  }

  int pop() {
    int value = value_;
    value_ = 0;
    return value;
  }

protected:
  int value_;
};

class Count {
public:
  Count() : value_(0) {}

  void addOne() {
    value_++;
  }
  void add(int i) {
    value_ += i;
  }

  int get() {
    return value_;
  }

protected:
  int value_;
};

// PerfNode is an Homie MQTT node that exposes read-only internal performance
// data about the instance.
// The messages are buffered to reduce chatter.
class PerfNode : public HomieNode {
public:
  PerfNode(const char *name, int delayms)
      : HomieNode(name, "perf"), delayms(delayms) {
    advertise("render");
    advertise("spi");
    advertise("i2c");
    advertise("frames");
    render_.init(10);
    spi_.init(10);
    i2c_.init(10);
  }

  void onRender() {
    render_.add();
  }
  void onSPI() {

  }
  void onI2C() {

  }
protected:
  void flush() {
    setProperty("render").send(String(render_.avgDelta()));
    setProperty("spi").send(String(spi_.avgDelta()));
    setProperty("i2c").send(String(i2c_.avgDelta()));
    setProperty("frames").send(String(frames_));
  }

  Samples render_;
  Samples spi_;
  Samples i2c_;
  Count frames_;
  // Delay between updates.
  int delayms;

  DISALLOW_COPY_AND_ASSIGN(PerfNode);
};

extern PerfNode Perf;
