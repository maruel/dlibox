// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package allwinner_pl

import "github.com/maruel/dlibox/go/pio"

func init() {
	pio.MustRegister(&driver{})
}
