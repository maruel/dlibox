// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package host

// Make sure CPU and board drivers are registered.
import (
	_ "github.com/maruel/dlibox/go/pio/host/allwinner"
	_ "github.com/maruel/dlibox/go/pio/host/bcm283x"
	_ "github.com/maruel/dlibox/go/pio/host/pine64"
	_ "github.com/maruel/dlibox/go/pio/host/rpi"
)
