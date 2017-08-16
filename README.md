# dlibox

Yet-another-home-automation project

Why another one?

- It is FAST.
  - When there's a power outage, boots within 11 seconds on a RPi3.
    - I don't want to wait 3 minutes for Java (OpenHAB) or node.js (node-red) to
      startup. Go executables starts instantaneously.
  - It is designed to run extremely well on single core systems like the
    [C.H.I.P.](https://getchip.com/) or the Raspberry Pi Zero.
- Optimized for maintainability:
  - Devices are expected to be deployed via
    [github.com/periph/bootstrap](https://github.com/periph/bootstrap). It
    applies Debian security updates automatically every night.
  - dlibox self updates every night.
  - The controller and the device (node) are the same Go executable.
  - The device has no local configuration beside the MQTT server and default to
    the host 'dlibox'.
  - Uses a derivative of the [Homie
    convention](https://github.com/marvinroger/homie) which is well designed.
    The tweak is that it's the *controller* that tells the device what nodes it
    shall present. This simplifies management.
- No internet connectivity is needed. Everything is local on the LAN.
- Web App served directly from the controller. Everything is accessed via this
  Web App. It is installable on mobile phones to use it like a App. It makes it
  trivial to make dashboards with old tablets.
- Can drive multiple strips of LEDs like the APA-102 in a **fully synchronous
  manner**, thanks to
  [github.com/maruel/anim1d][(https://github.com/maruel/anim1d). anim1d permits
  to create complex animations that are synchronized across multiple nodes.
  This permits very long runs of LEDs strips that are fully synchronized by
  using multiple computers, one per few hundred LED.
- Communicates over MQTT.
- Supports general 'home automation' like sensors and displays.
  - Leverages [periph.io](https://periph.io) for all hardware access.

Look at [HARDWARE.md](HARDWARE.md) for more information about what to buy.

There's an incomplete device implemented in [C++](esp/) to run on a ESP8266.

[![GoDoc](https://godoc.org/github.com/maruel/dlibox?status.svg)](https://godoc.org/github.com/maruel/dlibox)
