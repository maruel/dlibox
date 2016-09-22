# Devices

[![GoDoc](https://godoc.org/github.com/maruel/dlibox/go/pio/devices?status.svg)](https://godoc.org/github.com/maruel/dlibox/go/pio/devices)


## Drivers

Contains driver for the following devices:

* [apa102](apa102) is APA102 LEDs strip driver

![apa102](https://raw.githubusercontent.com/wiki/maruel/dlibox/apa102.jpg)

* [bme280](bme280) is a very precise environment sensor

![bme280](https://raw.githubusercontent.com/wiki/maruel/dlibox/bme280.jpg)

* [ir](ir) exposes infra red remote support via lircd.

\<insert image of a remote here\>

* [ssd1306](ssd1306) drives a small OLED display of 128x64 or less

![ssd1306](https://raw.githubusercontent.com/wiki/maruel/dlibox/ssd1306.jpg)

* [tm1637](ssd1306) drives a small segment based numerical display up to 6
  digits

![tm1637](https://raw.githubusercontent.com/wiki/maruel/dlibox/tm1637.jpg)


## Other

In addition, contains packages for:

* [devicetest](devicetest) fakes for unit testing
* [piblaster](piblaster) wrapper to generate PWM via DMA
