# pio - Usage


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

* The interfaces are defined in [protocols/](protocols/).
* The concrete objects _implementing_ the interfaces are in [host/](host/).
* The device drivers _using_ these interfaces are located in
  [devices/](devices/).

A device can be connected on a bus, let's say a strip of LED connected over SPI.
You need to connect the device driver of the LEDs to the SPI bus handle in your
application.


## Initialization

The function to initialize the default registered drivers is
[host.Init()](https://godoc.org/github.com/maruel/dlibox/go/pio/host#Init):

```go
  state, err := host.Init()
```

It returns information about:

* The drivers loaded and active.
* The drivers skipped, because the relevant hardware wasn't found.
* The drivers that failed to load. The app may still run without these drivers.

In addition, it may return an error when there's a structural issue, for example
two drivers with the same name were registered. This is a catastrophic failure.

The package [host](https://godoc.org/github.com/maruel/dlibox/go/pio/host)
**register** all the drivers under [host/](host/).

Calling
[host.NewI2CAuto()](https://godoc.org/github.com/maruel/dlibox/go/pio/host#example-NewI2CAuto)
or
[host.NewSPIAuto()](https://godoc.org/github.com/maruel/dlibox/go/pio/host#NewSPIAuto)
implicitly calls
[host.Init()](https://godoc.org/github.com/maruel/dlibox/go/pio/host#Init) on
your behalf to save you typing.


## Connection

A connection
[protocols.Conn](https://godoc.org/github.com/maruel/dlibox/go/pio/protocols#Conn)
is a point-to-point connection between the host and a device. A `Conn` can be
made of [gpio
pins](https://godoc.org/github.com/maruel/dlibox/go/pio/protocols/gpio#PinIO) or
via a [bus driver](https://godoc.org/github.com/maruel/dlibox/go/pio/host).


### I²C connection

An
[i2c.Conn](https://godoc.org/github.com/maruel/dlibox/go/pio/protocols/i2c#Conn)
is not a
[protocols.Conn](https://godoc.org/github.com/maruel/dlibox/go/pio/protocols#Conn).
This is because an I²C bus is **not** a point-to-point connection but instead is
a real bus where multiple devices can be connected simultaneously, like an USB
bus. To create a virtual connection to a device, the device address is required
via
[i2c.Dev](https://godoc.org/github.com/maruel/dlibox/go/pio/protocols/i2c#Dev):

```go
  bus, _ := host.NewI2CAuto()
  dev := i2c.Dev{bus, 0x76}
  var _ protocols.Bus = dev
```


## Samples

Please look at the device driver documentation for further examples. Tools in
[cmd/](cmd/) can also be used as the basis of your projects.


### IR

Displaying IR remote keys via [lirc](http://www.lirc.org/):

```go
package main

import (
  "fmt"
  "log"

  "github.com/maruel/dlibox/go/pio/devices/lirc"
)

func main() {
  conn, err := lirc.New()
  if err != nil {
    log.Fatal(err)
  }
  for msg := range conn.Channel() {
    if !msg.Repeat {
      fmt.Printf("%12s from %12s\n", msg.Key, msg.RemoteType)
    }
  }
}
```


### OLED 128x64 display

Displaying an animated GIF via a
[ssd1306](https://godoc.org/github.com/maruel/dlibox/go/pio/devices/ssd1306).
The frames in the GIF are resized and centered first to reduce the CPU
overhead.

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

func convertAndResizeAndCenter(w, h int, src image.Image) *image.Gray {
  src = resize.Thumbnail(uint(w), uint(h), src, resize.Bicubic)
  img := image.NewGray(image.Rect(0, 0, w, h))
  r := src.Bounds()
  r = r.Add(image.Point{(w - r.Max.X) / 2, (h - r.Max.Y) / 2})
  draw.Draw(img, r, src, image.Point{}, draw.Src)
  return img
}

func main() {
  bus, err := host.NewI2CAuto()
  if err != nil {
    log.Fatal(err)
  }
  dev, err := ssd1306.NewI2C(bus, 128, 64, false)
  if err != nil {
    log.Fatal(err)
  }
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
  imgs := make([]*image.Gray, len(g.Image))
  for i := range g.Image {
    imgs[i] = convertAndResizeAndCenter(dev.W, dev.H, g.Image[i])
  }
  for i := 0; ; i++ {
    index := i % len(imgs)
    c := time.After(time.Duration(10*g.Delay[index]) * time.Millisecond)
    img := imgs[index]
    dev.Draw(img.Bounds(), img, image.Point{})
    <-c
  }
}
```
