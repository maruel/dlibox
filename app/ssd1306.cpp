// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

#include "user_config.h"
#include <SmingCore/SmingCore.h>
#include <Wire.h>

#include "ada_ssd1306.h"
#include "ssd1306.h"

#define OLED_RESET 0  // GPIO0
Adafruit_SSD1306 display(OLED_RESET);

#if (SSD1306_LCDHEIGHT != 48)
#error("Height incorrect, please fix Adafruit_SSD1306.h!");
#endif

// Font size:
// - 1: 10 characters wide; 6 lines
// - 2: 5 characters wide; 3 lines
void initSSD1306() {
  // Set for the wemos i²c pins.
  Wire.pins(5, 4);
  // TODO(maruel): Change speed according to config.display.I2Cspeed.
  display.begin();
  display.clearDisplay();
  display.setTextSize(1);
  display.setTextColor(WHITE);
  display.setCursor(0,0);
  display.println("dlibox");
  display.display();
}
