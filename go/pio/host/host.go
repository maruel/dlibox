// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package host

import "github.com/maruel/dlibox/go/pio/drivers"

// Init calls drivers.Init() and returns it as-is.
//
// The only difference is that by calling host.Init(), you are guaranteed to
// have all drivers implemented to be implicitly loaded.
func Init() ([]drivers.Driver, []error) {
	return drivers.Init()
}
