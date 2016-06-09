// This file was retrieved from:
// https://github.com/mcauser/Adafruit_SSD1306/blob/esp8266-64x48/Adafruit_SSD1306.h
// which includes changes to support a 64x48 display

/*********************************************************************
This is a library for our Monochrome OLEDs based on SSD1306 drivers

  Pick one up today in the adafruit shop!
  ------> http://www.adafruit.com/category/63_98

These displays use SPI to communicate, 4 or 5 pins are required to
interface

Adafruit invests time and resources providing this open source code,
please support Adafruit and open-source hardware by purchasing
products from Adafruit!

Written by Limor Fried/Ladyada  for Adafruit Industries.
BSD license, check license.txt for more information
All text above, and the splash screen must be included in any redistribution
*********************************************************************/
#ifndef _Adafruit_SSD1306_H_
#define _Adafruit_SSD1306_H_

#if defined(__SAM3X8E__)
 typedef volatile RwReg PortReg;
 typedef uint32_t PortMask;
 #define HAVE_PORTREG
#elif defined(ARDUINO_ARCH_SAMD)
// not supported
#elif defined(ESP8266) || defined(ARDUINO_STM32_FEATHER)
  typedef volatile uint32_t PortReg;
  typedef uint32_t PortMask;
#else
  typedef volatile uint8_t PortReg;
  typedef uint8_t PortMask;
 #define HAVE_PORTREG
#endif

#include <SPI.h>
#include <Libraries/Adafruit_GFX/Adafruit_GFX.h>

#define BLACK 0
#define WHITE 1
#define INVERSE 2

#define SSD1306_I2C_ADDRESS   0x3C  // 011110+SA0+RW - 0x3C or 0x3D
// Address for 128x32 is 0x3C
// Address for 128x64 is 0x3D (default) or 0x3C (if SA0 is grounded)

#define SSD1306_SETCONTRAST 0x81
#define SSD1306_DISPLAYALLON_RESUME 0xA4
#define SSD1306_DISPLAYALLON 0xA5
#define SSD1306_NORMALDISPLAY 0xA6
#define SSD1306_INVERTDISPLAY 0xA7
#define SSD1306_DISPLAYOFF 0xAE
#define SSD1306_DISPLAYON 0xAF

#define SSD1306_SETDISPLAYOFFSET 0xD3
#define SSD1306_SETCOMPINS 0xDA

#define SSD1306_SETVCOMDETECT 0xDB

#define SSD1306_SETDISPLAYCLOCKDIV 0xD5
#define SSD1306_SETPRECHARGE 0xD9

#define SSD1306_SETMULTIPLEX 0xA8

#define SSD1306_SETLOWCOLUMN 0x00
#define SSD1306_SETHIGHCOLUMN 0x10

#define SSD1306_SETSTARTLINE 0x40

#define SSD1306_MEMORYMODE 0x20
#define SSD1306_COLUMNADDR 0x21
#define SSD1306_PAGEADDR   0x22

#define SSD1306_COMSCANINC 0xC0
#define SSD1306_COMSCANDEC 0xC8

#define SSD1306_SEGREMAP 0xA0

#define SSD1306_CHARGEPUMP 0x8D

#define SSD1306_EXTERNALVCC 0x1
#define SSD1306_SWITCHCAPVCC 0x2

// Scrolling #defines
#define SSD1306_ACTIVATE_SCROLL 0x2F
#define SSD1306_DEACTIVATE_SCROLL 0x2E
#define SSD1306_SET_VERTICAL_SCROLL_AREA 0xA3
#define SSD1306_RIGHT_HORIZONTAL_SCROLL 0x26
#define SSD1306_LEFT_HORIZONTAL_SCROLL 0x27
#define SSD1306_VERTICAL_AND_RIGHT_HORIZONTAL_SCROLL 0x29
#define SSD1306_VERTICAL_AND_LEFT_HORIZONTAL_SCROLL 0x2A

class Adafruit_SSD1306 : public Adafruit_GFX {
 public:
  // SPI - we indicate DataCommand, ChipSelect, Reset.
  Adafruit_SSD1306(uint16_t w, uint16_t h, int8_t DC, int8_t RST, int8_t CS)
      : Adafruit_GFX(w, h), dc(DC), rst(RST), cs(CS), buffer(new uint8_t[w*h/8]) {
  }
  // I²C - only need the reset pin
  Adafruit_SSD1306(uint16_t w, uint16_t h, int8_t reset = -1)
      : Adafruit_GFX(w, h), dc(-1), cs(-1), rst(reset), buffer(new uint8_t[w*h/8]) {
  }

  void begin(uint8_t switchvcc = SSD1306_SWITCHCAPVCC, uint8_t i2caddr = SSD1306_I2C_ADDRESS, bool reset=true);
  void ssd1306_command(uint8_t c);
  // TODO(maruel): const.
  void ssd1306_commands(uint8_t *c, uint16_t len);

  // Clear the display black.
  void clearDisplay();
  void invertDisplay(bool i);
  // Write the display buffer to the controller.
  void display();

  // Activate a left|right handed scroll for rows start through stop
  // Hint, the display is 16 rows tall. To scroll the whole display, run:
  // display.scrollright(0x00, 0x0F)
  void startScrollHor(bool left, uint8_t start, uint8_t stop);

  // Activate a diagonal scroll for rows start through stop
  // Hint, the display is 16 rows tall. To scroll the whole display, run:
  // display.scrollright(0x00, 0x0F)
  void startScrollDiag(bool left, uint8_t start, uint8_t stop);
  void stopscroll();

  // Dim the display.
  void dim(bool dim);

  // Most basic function, draw a pixel.
  void drawPixel(int16_t x, int16_t y, uint16_t color);

  virtual void drawFastVLine(int16_t x, int16_t y, int16_t h, uint16_t color);
  virtual void drawFastHLine(int16_t x, int16_t y, int16_t w, uint16_t color);

 protected:
  uint8_t calcContrast();
  void setupSPI();

  int8_t _i2caddr, _vccstate, dc, rst, cs;
  uint8_t *const buffer;
#ifdef HAVE_PORTREG
  PortReg *mosiport, *clkport, *csport, *dcport;
  PortMask mosipinmask, clkpinmask, cspinmask, dcpinmask;
#endif

  inline void drawFastVLineInternal(int16_t x, int16_t y, int16_t h, uint16_t color) __attribute__((always_inline));
  inline void drawFastHLineInternal(int16_t x, int16_t y, int16_t w, uint16_t color) __attribute__((always_inline));
};

#endif /* _Adafruit_SSD1306_H_ */
