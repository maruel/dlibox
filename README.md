dotsar
======

Drives a [DotStar](https://www.adafruit.com/datasheets/APA102.pdf) /
[APA102](https://cpldcpu.files.wordpress.com/2014/08/apa-102-super-led-specifications-2013-en.pdf)
/
[APA102-C](https://cpldcpu.files.wordpress.com/2014/08/apa-102c-super-led-specifications-2014-en.pdf)
LED strip via a Raspberry Pi. Note that the datasheet from Adafruits is a
weirdly repackaged PDF.

See [excellent](https://cpldcpu.wordpress.com/2014/08/27/apa102/)
[posts](https://cpldcpu.wordpress.com/2014/11/30/understanding-the-apa102-superled/)
from "Tim".

Supports emulating the LED strip at the console to test while waiting for
the LEDs to arrive from your provider.

## Hardware

- One or many DotStar / APA102-C LED strip.
  - One of the nice thing from these is that the chip is rated for -40°C so
    buying a weatherproof strip + long cables is worth it.
- One 74AHCT125N to safely upgrade the SPI signal from 3.3V to 5V.
- One large (>=10A) 5V power supply per ~400 LEDs. Specs are unclear about max
  rating. TODO(maruel): Calculate effective mA per LED.
- Wires.
- Small breadboard.
- LEDs + Switches (TODO).
- Raspberry Pi (1 or 2) + its own power supply + wifi adaptor / ethernet
  connected for the web server.

You can buy one from Adafruits or a cheaper place like aliexpress.com if you are
not in a hurry and/or outside the US.


## Software

You need to enable the SPI port on the Raspberry Pi. We assume you use Raspbian
Jessie.

etc...
