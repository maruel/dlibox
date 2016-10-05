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

A device can be connected on a bus, let's say a strip of LED connected over SPI.
You need to connect the device driver of the LEDs to the SPI bus handle in your
application.


## State

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
the same name were registered. This is a catastrophic failure.

The package [host](https://godoc.org/github.com/maruel/dlibox/go/pio/host)
registers all the drivers under [host/](../../host/).


## Connection

A connection
[protocols.Conn](https://godoc.org/github.com/maruel/dlibox/go/pio/protocols#Conn)
is a **point-to-point** connection between the host and a device.

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
bus, _ := i2c.New(-1)
dev := i2c.Dev{bus, 0x76}
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

Please look at the device driver documentation for further examples.

You are encouraged to look at tools in [cmd/](cmd/). These can be used as the
basis of your projects.


### IR (infra red remote)

_Purpose:_ display IR remote keys.

This sample uses lirc (http://www.lirc.org/). This assumes you installed lirc
and configured it. See
[devices/lirc](https://godoc.org/github.com/maruel/dlibox/go/pio/devices/lirc)
for more information.

```go
package main

import (
    "fmt"
    "log"

    "github.com/maruel/dlibox/go/pio/devices/lirc"
)

func main() {
    // Open a handle to lircd:
    conn, err := lirc.New()
    if err != nil {
        log.Fatal(err)
    }

    // Open a channel to receive IR messages and print them out as they are
    // received, skipping repeated messages:
    for msg := range conn.Channel() {
        if !msg.Repeat {
            fmt.Printf("%12s from %12s\n", msg.Key, msg.RemoteType)
        }
    }
}
```


### OLED 128x64 display

_Purpose:_ display an animated GIF.

This sample uses a
[ssd1306](https://godoc.org/github.com/maruel/dlibox/go/pio/devices/ssd1306).
The frames in the GIF are resized and centered first to reduce the CPU overhead.

```go
package main

import (
    "image"
    "image/draw"
    "image/gif"
    "log"
    "os"
    "time"

    "github.com/maruel/dlibox/go/pio/devices/ssd1306"
    "github.com/maruel/dlibox/go/pio/host"
    "github.com/nfnt/resize"
)

// convertAndResizeAndCenter takes an image, resizes and centers it on a
// image.Gray of size w*h.
func convertAndResizeAndCenter(w, h int, src image.Image) *image.Gray {
    src = resize.Thumbnail(uint(w), uint(h), src, resize.Bicubic)
    img := image.NewGray(image.Rect(0, 0, w, h))
    r := src.Bounds()
    r = r.Add(image.Point{(w - r.Max.X) / 2, (h - r.Max.Y) / 2})
    draw.Draw(img, r, src, image.Point{}, draw.Src)
    return img
}

func main() {
    // Load all the drivers:
    if _, err := host.Init(); err != nil {
        log.Fatal(err)
    }

    // Open a handle to the first available I²C bus:
    bus, err := i2c.New(-1)
    if err != nil {
        log.Fatal(err)
    }

    // Open a handle to a ssd1306 connected on the I²C bus:
    dev, err := ssd1306.NewI2C(bus, 128, 64, false)
    if err != nil {
        log.Fatal(err)
    }

    // Decodes an animated GIF as specified on the command line:
    if len(os.Args) != 2 {
        log.Fatal("please provide the path to an animated GIF")
    }
    f, err := os.Open(os.Args[1])
    if err != nil {
        log.Fatal(err)
    }
    g, err := gif.DecodeAll(f)
    f.Close()
    if err != nil {
        log.Fatal(err)
    }

    // Converts every frame to image.Gray and resize them:
    imgs := make([]*image.Gray, len(g.Image))
    for i := range g.Image {
        imgs[i] = convertAndResizeAndCenter(dev.W, dev.H, g.Image[i])
    }

    // Display the frames in a loop:
    for i := 0; ; i++ {
        index := i % len(imgs)
        c := time.After(time.Duration(10*g.Delay[index]) * time.Millisecond)
        img := imgs[index]
        dev.Draw(img.Bounds(), img, image.Point{})
        <-c
    }
}
```

## GPIO Edge detection

_Purpose:_ Signals when a button was pressed or a motion detector detected a
movement.

The
[gpio.PinIn.Edge()](https://godoc.org/github.com/maruel/dlibox/go/pio/protocols/gpio#PinIn)
function permits a edge detection without a busy loop. This is useful for **motion
detectors**, **buttons** and other kinds of inputs where a busy loop would burn
CPU for no reason.

```go
package main

import (
    "fmt"
    "log"

    "github.com/maruel/dlibox/go/pio/host"
    "github.com/maruel/dlibox/go/pio/protocols/gpio"
)

func main() {
    // Load all the drivers:
    if _, err := host.Init(); err != nil {
        log.Fatal(err)
    }

    // Lookup a pin by its number:
    p, err := gpio.ByNumber(16)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("%s: %s\n", p, p.Function())

    // Set it as input, with an internal pull down resistor:
    if err = p.In(gpio.Down, gpio.Both); err != nil {
        log.Fatal(err)
    }

    // Wait for edges as detected by the hardware, and print the value read:
    for {
        p.WaitForEdge()
        fmt.Printf("-> %s\n", p.Read())
    }
}
```


## Measuring weather

_Purpose:_ gather temperature, pressure and relative humidity.

This sample uses a
[bme280](https://godoc.org/github.com/maruel/dlibox/go/pio/devices/bme280).

```go
package main

import (
    "fmt"
    "log"

    "github.com/maruel/dlibox/go/pio/devices"
    "github.com/maruel/dlibox/go/pio/devices/bme280"
    "github.com/maruel/dlibox/go/pio/host"
)

func main() {
    // Load all the drivers:
    if _, err := host.Init(); err != nil {
        log.Fatal(err)
    }

    // Open a handle to the first available I²C bus:
    bus, err := i2c.New(-1)
    if err != nil {
        log.Fatal(err)
    }
    defer bus.Close()

    // Open a handle to a bme280 connected on the I²C bus:
    dev, err := bme280.NewI2C(bus, bme280.O2x, bme280.O2x, bme280.O2x, bme280.S500ms, bme280.FOff)
    if err != nil {
        log.Fatal(err)
    }
    defer dev.Close()

    // Read temperature from the sensor:
    var env devices.Environment
    if err = dev.Sense(&env); err != nil {
        log.Fatal(err)
    }
    fmt.Printf("%8s %10s %9s\n", env.Temperature, env.Pressure, env.Humidity)
}
```
