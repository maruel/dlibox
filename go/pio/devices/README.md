# Devices

[![GoDoc](https://godoc.org/github.com/maruel/dlibox/go/pio/devices?status.svg)](https://godoc.org/github.com/maruel/dlibox/go/pio/devices)

Contains driver for the following devices:

* [apa102](apa102) is APA102 LEDs strip driver
  ![apa102](https://raw.githubusercontent.com/wiki/maruel/dlibox/apa102.jpg)
* [bme280](bme280) is a very precise environment sensor
  ![bme280](https://raw.githubusercontent.com/wiki/maruel/dlibox/bme280.jpg)
* [ssd1306](ssd1306) drives a small OLED display of 128x64 or less
  ![ssd1306](https://raw.githubusercontent.com/wiki/maruel/dlibox/ssd1306.jpg)
* [tm1637](ssd1306) drives a small OLED display of 128x64 or less
  ![tm1637](https://raw.githubusercontent.com/wiki/maruel/dlibox/tm1637.jpg)

In addition, contains logic for:

* [devicetest](devicetest) fakes for unit testing
* [i2cdev](i2cdev) adapter to convert a host.I2C into a host.Bus
* [piblaster](piblaster) wrapper to generate PWM via DMA
