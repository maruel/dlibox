# dlibox for Go

Drives an APA-102 LED strip via a Raspberry Pi and expose a web server to
control it. The package includes many other utilities, i2c, spi, GPOI edge
triggering, bme280, ssd1306, PSF font support, 1bit image, 1D stateless
animation.

[![GoDoc](https://godoc.org/github.com/maruel/dlibox/go?status.svg)](https://godoc.org/github.com/maruel/dlibox/go)


## Features

- Supports emulating the LED strip at the console to test while waiting for the
  LEDs to arrive from your provider.
- Includes many transitions and stock animations.
- Boots automatically on Raspberry Pi startup, within seconds.
- Easy to update to newer version as features are added.
- Writen in Go, easy to hack on.


## Features planned

- Act as an alarm clock configurable via the Web UI.
- PNGs can be uploaded via the Web UI to create custom animations.
- Automatic self-update with the latest code every night.
- mDNS based discovery with multi-device synchronization.


## Steps

1. Buy [~100$ of hardware](HARDWARE.md).
2. [Set up the Raspberry Pi](setup/).
3. Hook it on the wall.
