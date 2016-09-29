# pio

pio is a peripheral I/O library in Go. The documentation, including examples, is at:
[![GoDoc](https://godoc.org/github.com/maruel/dlibox/go/pio?status.svg)](https://godoc.org/github.com/maruel/dlibox/go/pio)
and usage can be found at [USAGE.md](USAGE.md).


## Installation

pio includes many ready-to-use tools!

```bash
go get github.com/maruel/dlibox/go/pio/cmd/...
```

To cross-compile and send an executable to your micro computer:

```bash
cd $GOPATH/src/github.com/maruel/dlibox/go/pio/cmd/bme280
GOOS=linux GOARCH=arm go build .
scp bme280 raspberrypi:.
```


## State

The library is **not stable** yet and breaking changes continously happen.
Please version the libary using [one of go vendoring
tools](https://github.com/golang/go/wiki/PackageManagementTools) and sync
frequently.


## Usage

Here's a complete example to get the current temperature, barometric pressure
and relative humidity using a bme280:

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
  bus, err := host.NewI2CAuto()
  if err != nil {
    log.Fatal(err)
  }
  defer bus.Close()
  dev, err := bme280.NewI2C(bus, bme280.O2x, bme280.O2x, bme280.O2x, bme280.S500ms, bme280.FOff)
  if err != nil {
    log.Fatal(err)
  }
  defer dev.Close()
  var env devices.Environment
  if err = dev.Read(&env); err != nil {
    log.Fatal(err)
  }
  fmt.Printf("%8s %10s %9s\n", env.Temperature, env.Pressure, env.Humidity)
}
```

See more at [USAGE.md](USAGE.md).


## Authors

The main author is [Marc-Antoine Ruel](https://github.com/maruel). The full list
is in [AUTHORS](AUTHORS) and [CONTRIBUTORS](CONTRIBUTORS).


## Design

See [DESIGN.md](DESIGN.md) for the goals, requirements and driver lifetime
management.


## Contributions

We gladly accept contributions via GitHub pull requests, as long as the author
has signed the Google Contributor License. Please see
[CONTRIBUTING.md](CONTRIBUTING.md) for more details.


### Disclaimer

This is not an official Google product (experimental or otherwise), it
is just code that happens to be owned by Google.
