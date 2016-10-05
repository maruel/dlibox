# pio - Application developpers

Documentation for _application developers_ who want to write Go applications
leveraging `pio`.

The complete API documentation, including examples, is at
[![GoDoc](https://godoc.org/github.com/maruel/dlibox/go/pio?status.svg)](https://godoc.org/github.com/maruel/dlibox/go/pio).


## Introduction

pio uses a driver registry to efficiently load the relevant drivers on the host
it is running on. It differentiates between drivers that _enable_ functionality
on the host and drivers for devices connected _to_ the host.

Most micro computers expose at least some of the following:
[I²C bus](https://godoc.org/github.com/maruel/dlibox/go/pio/protocols/i2c#Conn),
[SPI bus](https://godoc.org/github.com/maruel/dlibox/go/pio/protocols/spi#Conn),
[gpio
pins](https://godoc.org/github.com/maruel/dlibox/go/pio/protocols/gpio#PinIO),
[analog
pins](https://godoc.org/github.com/maruel/dlibox/go/pio/protocols/analog),
[UART](https://godoc.org/github.com/maruel/dlibox/go/pio/protocols/uart), I2S
and PWM.

* The interfaces are defined in [protocols/](../../protocols/).
* The concrete objects _implementing_ the interfaces are in
  [host/](../../host/).
* The device drivers _using_ these interfaces are located in
  [devices/](../../devices/).

A device can be connected on a bus, let's say a LEDs strip connected over SPI.
You need to connect the device driver of the LEDs to the SPI bus handle in your
application.


## Project state

The library is **not stable** yet and breaking changes continously happen.
Please version the libary using [one of go vendoring
tools](https://github.com/golang/go/wiki/PackageManagementTools) and sync
frequently.


## Initialization

The function to initialize the default registered drivers is
[host.Init()](https://godoc.org/github.com/maruel/dlibox/go/pio/host#Init). It
returns at
[pio.State](https://godoc.org/github.com/maruel/dlibox/go/pio#State):

```go
state, err := host.Init()
```

[pio.State](https://godoc.org/github.com/maruel/dlibox/go/pio#State) contains
information about:

* The drivers loaded and active.
* The drivers skipped, because the relevant hardware wasn't found.
* The drivers that failed to load. The app may still run without these drivers.

In addition,
[host.Init()](https://godoc.org/github.com/maruel/dlibox/go/pio/host#Init) may
return an error when there's a structural issue, for example two drivers with
the same name were registered. This is a catastrophic failure. The package
[host](https://godoc.org/github.com/maruel/dlibox/go/pio/host) registers all the
drivers under [host/](../../host/).


## Connection

A connection
[protocols.Conn](https://godoc.org/github.com/maruel/dlibox/go/pio/protocols#Conn)
is a **point-to-point** connection between the host and a device where you are
the master driving the I/O.

[protocols.Conn](https://godoc.org/github.com/maruel/dlibox/go/pio/protocols#Conn)
implements [io.Writer](https://golang.org/pkg/io/#Writer) for write-only
devices, so you can use functions like
[io.Copy()](https://golang.org/pkg/io/#Copy) to push data over a connection.

A `Conn` can be multiplexed over the underlying bus. For example an I²C bus may
have multiple connections (slaves) to the master, each addressed by the device
address. The same is true on SPI via the `CS` line. On the other hand, UART
connection is always point-to-point. A `Conn` can even be created out of gpio
pins via bit banging.


### I²C connection

An
[i2c.Conn](https://godoc.org/github.com/maruel/dlibox/go/pio/protocols/i2c#Conn)
is **not** a
[protocols.Conn](https://godoc.org/github.com/maruel/dlibox/go/pio/protocols#Conn).
This is because an I²C bus is **not** a point-to-point connection but instead is
a real bus where multiple devices can be connected simultaneously, like an USB
bus. To create a virtual connection to a device, the device address is required
via
[i2c.Dev](https://godoc.org/github.com/maruel/dlibox/go/pio/protocols/i2c#Dev):

```go
// Open the first available I²C bus:
bus, _ := i2c.New(-1)
// Address the device with address 0x76 on the I²C bus:
dev := i2c.Dev{bus, 0x76}
// This is effectively a point-to-point connection:
var _ protocols.Conn = &dev
```

Since many devices have their address hardcoded, it's up to the device driver to
specify the address.


#### exp/io compatibility

To convert a
[i2c.Dev](https://godoc.org/github.com/maruel/dlibox/go/pio/protocols/i2c#Dev)
to a
[exp/io/i2c/driver.Conn](https://godoc.org/golang.org/x/exp/io/i2c/driver#Conn),
use the following:

```go
type adaptor struct {
    protocols.Conn
}

func (a *adaptor) Close() error {
    // It's not to the device to close the bus.
    return nil
}
```

### SPI connection

An
[spi.Conn](https://godoc.org/github.com/maruel/dlibox/go/pio/protocols/spi#Conn)
**is** a
[protocols.Conn](https://godoc.org/github.com/maruel/dlibox/go/pio/protocols#Conn).


#### exp/io compatibility

To convert a
[spi.Conn](https://godoc.org/github.com/maruel/dlibox/go/pio/protocols/spi#Conn)
to a
[exp/io/spi/driver.Conn](https://godoc.org/golang.org/x/exp/io/spi/driver#Conn),
use the following:

```go
type adaptor struct {
    spi.Conn
}

func (a *adaptor) Configure(k, v int) error {
    if k == driver.MaxSpeed {
        return a.Conn.Speed(int64(v))
    }
    // The match is not exact, as spi.Conn.Configure() configures simultaneously
    // mode and bits.
    return errors.New("TODO: implement")
}

func (a *adaptor) Close() error {
    // It's not to the device to close the bus.
    return nil
}
```


### GPIO

[gpio
pins](https://godoc.org/github.com/maruel/dlibox/go/pio/protocols/gpio#PinIO)
can be leveraged for arbitrary use, like buttons, control LEDs, etc. You may
construct an I²C or a SPI bus over raw GPIO pins via
[experimental/bitbang](https://godoc.org/github.com/maruel/dlibox/go/pio/experimental/bitbang).


## Samples

See [SAMPLES.md](SAMPLES.md) for various examples.
