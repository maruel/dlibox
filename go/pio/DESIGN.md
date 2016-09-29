# pio - design


pio is a peripheral I/O library in Go. The documentation, including examples, is at:
[![GoDoc](https://godoc.org/github.com/maruel/dlibox/go/pio?status.svg)](https://godoc.org/github.com/maruel/dlibox/go/pio)

It is recommended to look at the stand alone executables in [cmd/](cmd/) for use
cases.


## Abstract

Go developped a fairly large hardware hacker community in because the language
and its tooling have the following properties:

* Easy to cross compile to Arm/Linux via `GOOS=linux GOARCH=arm go build .`.
* Significantly faster to execute than python.
* Significantly lighter in term of memory use than Java or node.js.
* Builds reasonably fast on Arm.
* Fairly good OS level support: Debian pre-provided Go package (albeit a tad
  old) makes it easy to apt-get install on arm64, or arm32 users have access to
  package on [golang.org](https://golang.org).

Many packages, both generic like [embd](https://github.com/kidoman/embd),
[gobot](https://github.com/hybridgroup/gobot) and specialized (various one-off
drivers), were created to fill the space but there isn’t one clear winner or a
cohesive design pattern that scales to multiple platforms. Many have either
grown organically or have incomplete implementation. Most have a primitive
driver loading mechanism but is generally not flexible enough. A effort is in
progress to create a generic set of interface at
[exp/io](https://golang.org/x/exp/io) but this doesn't span the actual
implementations.

This document exposes a design to create a cohesive and definitive common
library that can be maintained on the long term.


## Goals

* Not more abstract than absolutely needed.
* Extensible:
  * Users can provide additional drivers that are seamlessly loaded
    with a structured ordering of priority.
* Performance:
  * Execution must be as performant as possible.
  * Overhead must be as minimal as possible, i.e. irrelevant driver must not be
    attempted to be loaded.
* Coverage:
  * Each driver must implement and expose as much of the underlying device
    capability as possible and relevant.
  * [cmd/](cmd/) implements useful directly usable tool.
  * [devices/](devices/) implements common device drivers.
  * [host/](host/) must implement a large base of common platforms that _just
    work_. This is in addition to extensibility.
  * Interfacing for common OS provided functionality (i.e. [sysfs](host/sysfs))
    and emulated ones (i.e. [bitbang](host/bitbang)).
* Simplicity:
  * Static typing is *thoroughly used*, to reduce the risk of runtime failure.
  * Minimal coding is needed to accomplish a task.
  * Use of the library is defacto portable.
  * Include fakes for buses and device interfaces to simplify the life of
    device driver developers.
* Stability
  * The library must be stable without precluding core refactoring.
  * Breakage in the API should happen at a yearly parce at most once the library
    got to a stable state.
* Strong distinction about the driver (as a user of a `Conn` instance) and an
  application writer (as a user of a device driver). It's the application that
  controls the object's lifetime.
* Strong distinction between _enablers_ and _devices_. See
  [Background](#background) below.


## Requirements

All the code must fit these requirements:

* The code must be Go idiomatic.
  * Constructor `NewXXX()` returns an object of concrete type.
  * Functions accept interfaces.
  * Leverage standard interfaces like
    [io.Writer](https://golang.org/pkg/io/#Writer) and
    [image.Image](https://golang.org/pkg/image/#Image) where possible.
  * No `interface{}` unless strictly required.
  * Minimal use of factories.
  * No `init()` code that accesses peripherals on process startup.
* Exact naming
  * Driver for a chipset must have the name of the chipset or the chipset
    family. Don't use `oleddisplay`, use `ssd1306`.
  * Driver must use the real chip name, not a marketing name by a third party.
    Don't use `dotstar` (as marketed by Adafruit), use `apa102` (as published
    by APA Electronic co. LTD.).
  * A link to the datasheet should be included in the package doc.
* Testability
  * Code must be testable and tested without a driver.
  * Include smoke-test (working with a real device) to confirm the library
    physically works for devices other than write-only devices.
* Usability
  * Provide a standalone executable in [cmd/](cmd/) to expose the functionality.
    It is acceptable to only expose a small subset of the functionality but the
    tool must have purpose.
* Performance
  * Drivers controling an output device must have a fast path that can be used
    to directly write in the device's native format.
  * Drivers controling an output device must have a generic path accepting
    higher level interface when found in the stdlib, i.e.
    [image.Image](https://golang.org/pkg/image/#Image)
  * Floating point arithmetic should only be used when absolutely necesary in
    the driver code. Most of the cases can be replaced with fixed point
    arithmetic, for example
    [devices.Milli](https://godoc.org/github.com/maruel/dlibox/go/pio/devices#Milli).
    Floating point arithmetic is acceptable in the unit tests and tools in
    [cmd/](cmd/) but should not be abused.
  * Drivers must be implemented with performance in mind.
  * Benchmark must be implemented for non trivial processing.
* Code must compile on all OSes, with minimal use of OS-specific thunk as
  strictly as needed.
* Struct implementing an interface must validate at compile time with `var _
  <Interface> = &<Type>{}`.
* No code under the GPL, LGPL or APL license will be accepted.
  * Users are free to use the library in commercial projects.


## Driver lifetime management

Proper driver lifetime management is key to the success of this project. There
must be clear expectations to add, update and remove drivers for the core
project. As described in [Risks](#Risk) below, poor drivers or high churn rate
will destroy the value proposition.

This is critical as drivers can be silently broken by seemingly innocuous
changes. Because the testing story of hardware is significantly harder than
software-only projects, there’s an inherent faith in the quality of the code
that must be asserted.


### Experimental

Any driver can be requested to be added to the library under
[experimental/](experimental/) directory. The following process must be
followed:
* One or multiple developers have created a driver out of tree.
* The driver is deemed to work.
* The driver meets minimal quality bar under the promise of being improved.
* Follows [CONTRIBUTING.md](CONTRIBUTING.md) demand.
* Create a Pull Request for integration under [experimental/](experimental/) and
  respond to the code review.

At this point, it is available for use to everyone but is not loaded defacto.
There is no API compatibility guarantee.


### Stable

A driver in experimental can be promoted to stable. The following process must
be followed:
* Declare at least one (or multiple) owners that are responsive to reply to
  feature requests and bug reports.
  * There could be a threshold, > _TO BE DETERMINED_ lines, where more than one
    owner is required.
  * Contributors commit to support the driver for the foreseeable future and
    **promptly** do code reviews to keep the driver quality to the expected
    standard.
* There are multiple reports that the driver is functioning as expected.
* If another driver exists for an intersecting class of devices, the other
  driver must enter deprecation phase.
* At this point the driver must maintain its API compatibility promise.


### Deprecation

A driver can be subsumed by a newer driver with a better core implementation or
a new breaking API. The previous driver must be deprecated, moved back to
[experimental/](experimental/) and announced to be deleted after _TO BE
DETERMINED_ amount of time.


### Contributing a new driver

A new proposed driver must be first implemented out of tree and fit all the
items in [Requirements](#requirements) listed above. First propose it as
Experimental, then promote it to Stable.


## Background

#### Classes of hardware

This document distinguishes two classes of drivers:

* Enablers. They are what make the interconnects work, so that you can then
  use real stuff. That's buses (I²C, SPI, GPIO, BT, UART). This is what can be
  used as point-to-point protocols. They enable you to do something but are not
  the essence of what you want to do.
* Devices. They are the end goal, to do something functional. There are multiple
  subclasses of devices like sensors, write-only devices, etc.

The enablers is what will break or make this project. Nobody want to do them
but they are needed. You need a large base of enablers so people can use
anything yet they are hard to get right. You want them all in the same repo so
that when someone builds an app, it supports everything transparently. It just
works.

The device drivers do not need to all be in the same repo, that scales since
people know what is physically connected, but enablers are what needs to be in
the base repository. People do not care that a Pine64 has a different processor
than a Rasberry Pi; both have the same 40 pins header and that's what they care
about. So enablers need to be a great HAL -> the right hardware abstraction
layer (not too deep, not too light) is the core here.

Devices need common interfaces to help with application developers (like
[devices.Display](https://godoc.org/github.com/maruel/dlibox/go/pio/devices#Display)
and
[devices.Environmental](https://godoc.org/github.com/maruel/dlibox/go/pio/devices#Environmental))
but the lack of core repository and coherency is less dramatic.


## Success criteria

* Preferred library used by first time Go users and by experts
* Becomes the defacto HAL library.
* Becomes the central link for hardware support.


## Risks

### Users

* The library is rejected by users as being too cryptic or hard to use.
* The device drivers are unreliable or non functional, as observed by users.
* Poor usability of the core interfaces.


### Contributors

* Lack of API stability; high churn rate.
* Poor fitting of the core interfaces.
* No uptake in external contribution.
* Poor quality of contribution.
* Duplicate ways to accomplish the same thing, without a clear way to define the
  right way.


## Detailed design

### Driver registry

The core of extensibility is implemented as an in-process driver registry. The
things that make it work are:
* Clear priority classes via
  [pio.Type](https://godoc.org/github.com/maruel/dlibox/go/pio#Type).
  Each category is loaded one after the other so a driver of a type can assume
  that all relevant drivers of lower level types were fully loaded.
* Native way to skip a driver on unrelated platform.
  * At compile time via conditional compilation.
  * At runtime via early `Init()` exit.
* Native way to return the state of all registered drivers. The ones loaded, the
  ones skipped and the ones that failed.
* Native way to declare inter-driver dependency. A specialized
  [Processor](https://godoc.org/github.com/maruel/dlibox/go/pio#Type)
  driver may dependent on generic
  [Processor](https://godoc.org/github.com/maruel/dlibox/go/pio#Type)
  driver and the drivers will be loaded sequentially.
* In another other case, the drivers are loaded in parallel for minimum total
  latency.
