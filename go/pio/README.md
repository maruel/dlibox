# pio

pio is a peripheral I/O library in Go.

The documentation is split into 3 sections:
* [doc/users/](doc/users/) for users who need ready-to-use tools.
* [doc/apps/](doc/apps/) for application writers to want to use `pio` as a
  library.
  * The complete API documentation, including examples, is at
    [![GoDoc](https://godoc.org/github.com/maruel/dlibox/go/pio?status.svg)](https://godoc.org/github.com/maruel/dlibox/go/pio).
* [doc/drivers/](doc/drivers/) for device driver writers who want to expand
  the list of supported hardware and hopefully contribute to the project.


## Philosophy

1. Optimize for usability.
2. At usability expense, the user can chose to optimize for performance.
3. Use a divide and conquer approach. Each component has exactly one
   responsibility.
4. The driver's writer pleasure is dead last.


## Users

pio includes many ready-to-use tools!

```bash
go get github.com/maruel/dlibox/go/pio/cmd/...
```

See [doc/users/](doc/users/) for more info on:

* Configuring the host
* Using the included tools


## Application developpers

For [application developpers](doc/apps/), to get a quick feel, here's a
complete example to get the current temperature, barometric pressure and
relative humidity using a bme280:

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
    if err = dev.Read(&env); err != nil {
        log.Fatal(err)
    }
    fmt.Printf("%8s %10s %9s\n", env.Temperature, env.Pressure, env.Humidity)
}
```

See more examples at [doc/apps/](doc/apps/#samples)!


## Device drivers developpers

For [device drivers developpers](doc/drivers/), `pio` provides an extensible
driver registry and common bus interfaces. See this page for requirements to
submit contributions.


## Authors

The main author is [Marc-Antoine Ruel](https://github.com/maruel). The full list
is in [AUTHORS](AUTHORS) and [CONTRIBUTORS](CONTRIBUTORS).


## Contributions

We gladly accept contributions via GitHub pull requests, as long as the author
has signed the Google Contributor License. Please see
[doc/drivers/CONTRIBUTING.md](doc/drivers/CONTRIBUTING.md) for more details.


## Disclaimer

This is not an official Google product (experimental or otherwise), it
is just code that happens to be owned by Google.
