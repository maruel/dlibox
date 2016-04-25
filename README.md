dotsar
======

Drives an APA-102 / Dotstar LED strip via a Raspberry Pi and expose a web server
to control it.


## Features

- Supports emulating the LED strip at the console to test while waiting for the
  LEDs to arrive from your provider.
- Animation can be driven at 400Hz. Includes many stock animations.
- Switching between animations is done with a nice 500ms ease-in-out transition.
- Includes many transitions and stock animations.
- Boots automatically on Raspberry Pi startup within seconds.
- Easy to update to newer version as features are added.
- Writen in Go, easy to hack on.


## Features planned

- Act as an alarm clock configurable via the Web UI.
- PNGs can be uploaded via the Web UI to create custom animations.
- Automatic self-update with the latest code every night.


## Steps

1. Buy [~100$ of hardware](HARDWARE.md).
2. [Set up the Raspberry Pi](setup/).
3. Hook it on the wall.


## References

- The company that makes the LEDs is
  [APA](http://www.neon-world.com/patent_en.html).
  - [One pager 'datasheet' by APA](http://www.neon-world.com/pdf/led.pdf).
  - [Tim's APA-102 datasheet](https://cpldcpu.files.wordpress.com/2014/08/apa-102-super-led-specifications-2013-en.pdf).
  - [Tim's APA-102C datasheet](https://cpldcpu.files.wordpress.com/2014/08/apa-102c-super-led-specifications-2014-en.pdf).
  - Adafruit hosts a repackaged (duh) [PDF of the
    datasheet](https://www.adafruit.com/datasheets/APA102.pdf).
- [Pololu](http://www.neon-world.com/patent_en.html) has great information about
  the APA-102C.
- [Tim](https://github.com/cpldcpu) made two excellent posts about the APA-102C:
  [#1](https://cpldcpu.wordpress.com/2014/08/27/apa102/) and
  [#2](https://cpldcpu.wordpress.com/2014/11/30/understanding-the-apa102-superled/).
- Right-sizing power cables; http://www.powerstream.com/Wire_Size.htm.
