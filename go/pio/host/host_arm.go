// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package host

import (
	// Make sure CPU and board drivers are registered.
	_ "github.com/maruel/dlibox/go/pio/host/allwinner"
	_ "github.com/maruel/dlibox/go/pio/host/allwinner_pl"
	_ "github.com/maruel/dlibox/go/pio/host/bcm283x"
	_ "github.com/maruel/dlibox/go/pio/host/pine64"
	_ "github.com/maruel/dlibox/go/pio/host/rpi"
)
