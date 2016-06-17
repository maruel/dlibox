# dlibox

The box for funny people. It's main purpose is to drive APA-102 LEDs.

There is two versions, one in [Go](go/) to run on a Raspberry Pi and one in
[C++](esp/) to run on a ESP8266.

The [Raspberry Pi version](go/) can do more, but the [ESP8266 version](esp/)
cost much less. Both can communicate together via MQTT and discovery is done
through mDNS.


## APA-102 References

- The company that makes the LEDs is
  [APA](http://www.neon-world.com/patent_en.html).
  - [One pager 'datasheet' by APA](http://www.neon-world.com/pdf/led.pdf).
  - [Tim's APA-102 datasheet](https://cpldcpu.files.wordpress.com/2014/08/apa-102-super-led-specifications-2013-en.pdf).
  - [Tim's APA-102C datasheet](https://cpldcpu.files.wordpress.com/2014/08/apa-102c-super-led-specifications-2014-en.pdf).
  - Adafruit hosts a repackaged [PDF of the
    datasheet](https://www.adafruit.com/datasheets/APA102.pdf).
- [Pololu](http://www.neon-world.com/patent_en.html) has great information about
  the APA-102C.
- [Tim](https://github.com/cpldcpu) made two excellent posts about the APA-102C:
  [#1](https://cpldcpu.wordpress.com/2014/08/27/apa102/) and
  [#2](https://cpldcpu.wordpress.com/2014/11/30/understanding-the-apa102-superled/).
- Right-sizing power cables; http://www.powerstream.com/Wire_Size.htm.
