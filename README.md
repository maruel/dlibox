# dlibox

The box for funny people. It's main purpose is to drive APA-102 LEDs and
integrate into a home automation system by communicating via MQTT. It also have
its own interface to create complex animations and can work standalone. It
integrates sensor support (temperature) and output (small displays).

There is two versions, one in [Go](go/) to run on a Raspberry Pi/Orange
Pi/Pine64/etc and one incomplete in [C++](esp/) to run on a ESP8266.

The [Raspberry Pi version](go/) can do more, but the [ESP8266 version](esp/)
cost much less. Both can communicate together via MQTT and discovery is done
through mDNS.

Look at [HARDWARE.md](HARDWARE.md) for more information about what to buy.


## Related projects

### Rule engine

In this summary, I'm only looking at open sources projects that can run in
standalone mode without the need of internet connectivity. Otherwise you can use
https://ifttt.com, https://firebase.google.com/, the trash from Apple, etc.

- http://www.openhab.org/
  - Summary: by far the most popular and well supported rule engine. At the
    moment of writing, openHAB 2 is still on beta but we'll assume this version.
    It's only drawbacks are: server sluggishness (!!!) and hard requirement Java
    8, which is tricky to install on Armbian (arm64).
  - Doc: http://docs.openhab.org/
  - Server version 2 in Java 8 (Java 7 for openHAB1)
    - Server is extremely sluggish to start on a Raspberry Pi (several minutes!)
  - Web frontend in Polymer
  - Native Android and iOS apps.
  - Rules uses complex [Xtend expression but syntax close to
    Java](https://github.com/openhab/openhab/wiki/Scripts)
    - Native tool to edit the rules, supported by Eclipse foundation
  - Has a foundation to support the project long term
  - Has broadest hardware support:
    [Nest](https://github.com/openhab/openhab/wiki/Nest-Binding-Example),
    Insteon, [Sonos](https://github.com/openhab/openhab/wiki/Sonos-Binding),
    Philips Hue, Z-Wave,
    [Asterix](https://github.com/openhab/openhab/wiki/Asterisk-Binding)(!), etc.
  - Supports MQTT but requires a separate broker (e.g. mosquitto).
  - Everything about this project sounds heavy weight.
- https://home-assistant.io/
  - Summary: a new lighter weight entrant that runs in a docker image.
  - Samples: https://home-assistant.io/cookbook/
  - Server in python 3.
  - Web frontend in Polymer
  - Rule language in yaml
  - Supports MQTT but requires a separate broker (e.g. mosquitto).
  - Supports Chromecast, Philips Hue,
    [Z-Wave](https://home-assistant.io/getting-started/z-wave/).
- https://git.io/homieiot
  - Summary: lightweight esp8266 specific automation framework, also runs in a
    docker image.
  - Doc: https://homie-esp8266.readme.io/
  - Server in nodeJS
  - Rules are in JSON
    - Rules are edited with [Node-RED ](http://nodered.org/), written by IBM,
      also in nodeJS.
  - Supports MQTT but requires a separate broker (e.g. mosquitto).
    - Has a [nice
      schema](https://github.com/marvinroger/homie/tree/master#device-properties)
      for devices.


## Tools for video surveillance integration

- [ZoneMinder](https://www.zoneminder.com/) is a complete solution
  - It finally (!) added an API in 2016
  - Android and iOS app: http://pliablepixels.github.io/
  - Someone made a [docker](https://github.com/QuantumObject/docker-zoneminder)
- [Restreamer](https://datarhei.github.io/restreamer/) reencodes video on the
  fly for web viewing; for simpler solution
- [Motion](http://www.lavrsen.dk/foswiki/bin/view/Motion/WebHome) is a bit
  anticated
  - Someone made a [docker](https://github.com/kfei/dockmotion)
- RAW ffpmeg; not to be ignored, [it's always an
  option](https://docs.google.com/presentation/d/1EvaSzUjQc4zUNJxDPMzsDFTwGt1HssSBUX1-jF_HsQc)
