# Host

Host contains everything that relates to the host itself, including its CPU. It
contains all the buses that are connection points where devices can be connected
on.

* [ir](ir) exposes infra red remote support via lircd
* [pine64](pine64) exposes [Pine64](https://www.pine64.org/) specific hardware
  functionality, i.e.  headers pinout.
* [pins](pins) exposes GPIO functionality as found on the CPU driver, if any is
  found. Otherwise fallbacks to gpio sysfs, if available. 
* [rpi](rpi) exposes [Raspberry Pi](https://www.raspberrypi.org/) specific
  hardware functionality, i.e. headers pinout.

Please refer to
[![GoDoc](https://godoc.org/github.com/maruel/dlibox/go/pio/host?status.svg)](https://godoc.org/github.com/maruel/dlibox/go/pio/host).
