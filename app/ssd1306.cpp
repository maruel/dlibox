// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

#include "user_config.h"
#include "ada_ssd1306.h"
#include "ssd1306.h"
#include <SmingCore/SmingCore.h>
#include <Wire.h>

#define OLED_RESET 0  // GPIO0
Adafruit_SSD1306 display(OLED_RESET);

#if (SSD1306_LCDHEIGHT != 48)
#error("Height incorrect, please fix Adafruit_SSD1306.h!");
#endif

namespace ssd1306 {

void init() {
  // Set for the wemos i²c pins.
  Wire.pins(5, 4);
  display.begin();
  display.clearDisplay();
  display.display();
}

}   // namespace ssd1306
