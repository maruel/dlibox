# pio

pio is a peripheral I/O library in Go. The documentation, including examples, is at
[![GoDoc](https://godoc.org/github.com/maruel/dlibox/go/pio?status.svg)](https://godoc.org/github.com/maruel/dlibox/go/pio).
Usage and HowTos can be found at [USAGE.md](USAGE.md).


## Installation

_For end users_:

pio includes many ready-to-use tools!

```bash
go get github.com/maruel/dlibox/go/pio/cmd/...
```

To cross-compile and send an executable to your ARM based micro computer (e.g.
Raspberry Pi):

```bash
cd $GOPATH/src/github.com/maruel/dlibox/go/pio/cmd/bme280
GOOS=linux GOARCH=arm go build .
scp bme280 raspberrypi:.
```

The pio project doesn't release binaries, you are expected to build from
sources.


## Usage

_For application developpers_:

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
    // Open a handle to the first available I²C bus:
    bus, err := host.NewI2CAuto()
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
    if err = dev.Read(&env); err != nil {
        log.Fatal(err)
    }
    fmt.Printf("%8s %10s %9s\n", env.Temperature, env.Pressure, env.Humidity)
}
```

See more examples at [USAGE.md](USAGE.md#samples)!


### State

The library is **not stable** yet and breaking changes continously happen.
Please version the libary using [one of go vendoring
tools](https://github.com/golang/go/wiki/PackageManagementTools) and sync
frequently.


## Design

_For device drivers developpers_:

See [DESIGN.md](DESIGN.md) for the goals, requirements and driver lifetime
management. It is a required reading (it's okay to skim a bit but don't tell
anyone, shhh!) before contribution.


### Authors

The main author is [Marc-Antoine Ruel](https://github.com/maruel). The full list
is in [AUTHORS](AUTHORS) and [CONTRIBUTORS](CONTRIBUTORS).


### Contributions

We gladly accept contributions via GitHub pull requests, as long as the author
has signed the Google Contributor License. Please see
[CONTRIBUTING.md](CONTRIBUTING.md) for more details.


### Disclaimer

This is not an official Google product (experimental or otherwise), it
is just code that happens to be owned by Google.
