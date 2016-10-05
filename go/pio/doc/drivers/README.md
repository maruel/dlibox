# pio - Device driver developpers

Documentation for _device driver developers_ who either wants to developper a
device driver in their own code base or want to submit a contribution to extend
the supported hardware.


## Abstract

Go developped a fairly large hardware hacker community in part because the
language and its tooling have the following properties:

* Easy to cross compile to ARM/Linux via `GOOS=linux GOARCH=arm go build .`.
* Significantly faster to execute than python and node.js.
* Significantly lighter in term of memory use than Java or node.js.
* Significantly more productive to code than C/C++.
* Builds reasonably fast on ARM.
* Fairly good OS level support: Debian pre-provided Go package (albeit a tad
  old) makes it easy to apt-get install on arm64, or arm32 users have access to
  package on [golang.org](https://golang.org).

Many Go packages, both generic and specialized, were created to fill the space.
This library came out of the desire to have a _designed_ API (contrary to
growing organically) with strict [code requirements](#requirements) and a
[strong, opiniated philosophy](../../#philosophy) to enable long term
maintenance.


## Goals

* Not more abstract than absolutely needed. Use concrete types whenever
  possible.
* Orthogonality and composability
  * Each component must own an orthogonal part of the platform and each
    components can be composed together.
* Extensible:
  * Users can provide additional drivers that are seamlessly loaded
    with a structured ordering of priority.
* Performance:
  * Execution as performant as possible.
  * Overhead as minimal as possible, i.e. irrelevant driver are not be
    attempted to be loaded, uses memory mapped GPIO registers instead of sysfs
    whenever possible, etc.
* Coverage:
  * Be as OS agnostic as possible. Abstract OS specific concepts like
    [sysfs](https://godoc.org/github.com/maruel/dlibox/go/pio/host/sysfs).
  * Each driver implements and exposes as much of the underlying device
    capability as possible and relevant.
  * [cmd/](../../cmd/) implements useful directly usable tool.
  * [devices/](../../devices/) implements common device drivers.
  * [host/](../../host/) must implement a large base of common platforms that
    _just work_. This is in addition to extensibility.
* Simplicity:
  * Static typing is _thoroughly used_, to reduce the risk of runtime failure.
  * Minimal coding is needed to accomplish a task.
  * Use of the library is defacto portable.
  * Include fakes for buses and device interfaces to simplify the life of
    device driver developers.
* Stability
  * API must be stable without precluding core refactoring.
  * Breakage in the API should happen at a yearly parce at most once the library
    got to a stable state.
* Strong distinction about the driver (as a user of a
  [protocols.Conn](https://godoc.org/github.com/maruel/dlibox/go/pio/protocols#Conn)
  instance) and an application writer (as a user of a device driver). It's the
  _application_ that controls the objects' lifetime.
* Strong distinction between _enablers_ and _devices_. See
  [Background](#background) below.


## Requirements

All the code must fit the following requirements.

**Fear not!** We know the list _is_ daunting but as you create your pull request
to add something at [experimental/](../../experimental/), we'll happily guide
you in the process to help improve the code to meet the expected standard. The
end goal is to write *high quality maintainable code* and use this as a learning
experience.

* The code must be Go idiomatic.
  * Constructor `NewXXX()` returns an object of concrete type.
  * Functions accept interfaces.
  * Leverage standard interfaces like
    [io.Writer](https://golang.org/pkg/io/#Writer) and
    [image.Image](https://golang.org/pkg/image/#Image) where possible.
  * No `interface{}` unless strictly required.
  * Minimal use of factories except for protocol level registries.
  * No `init()` code that accesses peripherals on process startup. These belongs
    to
    [Driver.Init()](https://godoc.org/github.com/maruel/dlibox/go/pio#Driver).
* Exact naming
  * Driver for a chipset must have the name of the chipset or the chipset
    family. Don't use `oleddisplay`, use `ssd1306`.
  * Driver must use the real chip name, not a marketing name by a third party.
    Don't use `dotstar` (as marketed by Adafruit), use `apa102` (as created
    by APA Electronic co. LTD.).
  * A link to the datasheet must be included in the package doc unless NDA'ed
    or inaccessible.
* Testability
  * Code must be testable and tested without a device.
  * When relevant, include a smoke test under [tests/](../../tests/). The smoke
    test tests a real device to confirm the driver physically works for devices.
* Usability
  * Provide a standalone executable in [cmd/](../../cmd/) to expose the
    functionality.  It is acceptable to only expose a small subset of the
    functionality but _the tool must have purpose_.
  * Provide a `func Example()` along your test to describe basic usage of your
    driver. See the official [testing
    package](https://golang.org/pkg/testing/#hdr-Examples) for more details.
* Performance
  * Drivers controling an output device must have a fast path that can be used
    to directly write in the device's native format, e.g.
    [io.Writer](https://golang.org/pkg/io/#Writer).
  * Drivers controling an output device must have a generic path accepting
    higher level interface when found in the stdlib, e.g.
    [image.Image](https://golang.org/pkg/image/#Image).
  * Floating point arithmetic should only be used when absolutely necesary in
    the driver code. Most of the cases can be replaced with fixed point
    arithmetic, for example
    [devices.Milli](https://godoc.org/github.com/maruel/dlibox/go/pio/devices#Milli).
    Floating point arithmetic is acceptable in the unit tests and tools in
    [cmd/](../../cmd/) but should not be abused.
  * Drivers must be implemented with performance in mind. For example I²C
    operations should be batched to minimize overhead.
  * Benchmark must be implemented for non trivial processing running on the host.
* Code must compile on all OSes, with minimal use of OS-specific thunk as
  strictly needed.
* Struct implementing an interface must validate at compile time with `var _
  <Interface> = &<Type>{}`.
* License is Apache v2.0.


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
[experimental/](../../experimental/) directory. The following process must be
followed:
* One or multiple developers have created a driver out of tree.
* The driver is deemed to work.
* The driver meets minimal quality bar under the promise of being improved. See
  [Requirements](#requirements) for the extensive list.
* Follow [CONTRIBUTING.md](CONTRIBUTING.md) demands.
* Create a Pull Request for integration under
  [experimental/](../../experimental/) and respond to the code review.

At this point, it is available for use to everyone but is not loaded defacto by
[host.Init()](https://godoc.org/github.com/maruel/dlibox/go/pio/host#Init).

There is no API compatibility guarantee for drivers under
[experimental/](../../experimental/).


### Stable

A driver in [experimental/](../../experimental/) can be promoted to stable in
either [devices/](../../devices/) or [host/](../../host/) as relevant. The
following process must be followed:
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
[experimental/](../../experimental/) and announced to be deleted after _TO BE
DETERMINED_ amount of time.


### Contributing a new driver

A new proposed driver must be first implemented out of tree and fit all the
items in [Requirements](#requirements) listed above. First propose it as
[Experimental](#experimental), then ask to promote it to [Stable](#stable).


## Background

#### Classes of hardware

This document distinguishes two classes of drivers:

* Enablers: they are what make the interconnects work, so that you can then
  use real stuff. That's buses (I²C, SPI, GPIO, BT, UART). This is what can be
  used as point-to-point protocols. They enable you to do something but are not
  the essence of what you want to do. They can also be MCUs like AVR, ESP8266,
  etc.
* Devices: they are the end goal, to do something functional. There are multiple
  subclasses of devices like sensors, output devices, etc.

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

* Preferred library used by first time Go users and by experts.
* Becomes the defacto HAL library.
* Becomes the central link for hardware support.


## Risks

The risks below are being addressed via a strong commitment to [driver lifetime
management](#driver-lifetime-management) and having a high quality bar via an
explicit list of [requirements](#requirements).


### Users

* The library is rejected by users as being too cryptic or hard to use.
* The device drivers are unreliable or non functional, as observed by users.
* Poor usability of the core interfaces.
* Missing drivers.


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
