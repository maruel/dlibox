# dlibox

Yet-another-home-automation project

[![GoDoc](https://godoc.org/github.com/maruel/dlibox?status.svg)](https://godoc.org/github.com/maruel/dlibox) [![Go Report Card](https://goreportcard.com/badge/github.com/maruel/dlibox)](https://goreportcard.com/report/github.com/maruel/dlibox)

Why another one?

- **Performant**
  - When there's a power outage, it boots within 11 seconds on a RPi3.
    - I don't want to wait 3 minutes for Java (OpenHAB) or node.js (node-red) to
      startup. Go executables start instantaneously.
  - It is designed to run extremely well on single core systems like the
    [C.H.I.P.](https://getchip.com/) or the Raspberry Pi Zero.
- **Maintainable**
  - Devices can be deployed via
    [github.com/periph/bootstrap](https://github.com/periph/bootstrap).
  - dlibox self-updates every night.
  - The controller and the device (node) are the same Go executable. It can be
    simply scp'ed if desired.
  - The device has **no** local configuration beside the MQTT server name, which
    defaults to the host `dlibox` so if you setup your controller hostname to
    `dlibox`, you litterally have no configuration to do on the devices.
  - Uses a derivative of the [Homie
    convention](https://github.com/marvinroger/homie) which is well designed.
    The tweak is that it's the *controller* that tells the device what nodes it
    shall present. This simplifies management.
  - Communicates over MQTT, which is a stable protocol and a stable
    implementation.
- **Secure**
  - No internet connectivity is needed nor used. Everything is local on the LAN.
    What is in your house stays in your house.
  - Devices deployed via bootstrap apply Debian security updates automatically
    every night.
- **Usable**
  - Web App served directly from the controller. Everything is accessed via this
    Web App. It is installable on mobile phones to use it like a App. It makes
    it trivial to make dashboards with old tablets.
- **Featureful**
  - Can drive multiple strips of LEDs like the APA-102 in a **fully synchronous
    manner**, thanks to
    [github.com/maruel/anim1d](https://github.com/maruel/anim1d). anim1d permits
    to create complex animations that are synchronized across multiple nodes.
    This permits very long runs of LEDs strips that are fully synchronized by
    using multiple computers, one per few hundred LED.
  - Supports general 'home automation' like sensors and displays.
  - Leverages [periph.io](https://periph.io) for all hardware access.

Look at [HARDWARE.md](HARDWARE.md) for more information about what to buy.

There's an incomplete device implemented in [C++](esp/) to run on a ESP8266 that
will act as an Homie node.
