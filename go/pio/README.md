# pio

Peripheral I/O library.

* [cmd](cmd) contains executables to access right away the devices and buses as
  implemented in this library.
* [devices](devices) contains the device drivers for multiple devices you can
  connect to your host. This package implements interfaces for common use cases,
  like `Display` and `Environmental`.
* [host](host) contains host-specific code, relating either to the CPU, the GPIO
  pins, the buses (I²C, SPI, IR).

For more details,
[![GoDoc](https://godoc.org/github.com/maruel/dlibox/go/pio?status.svg)](https://godoc.org/github.com/maruel/dlibox/go/pio)
