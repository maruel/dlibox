// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// FIXME: This is an example where this driver is only relevant on ARM systems.
// In this example, the code will be compiled on all platforms but the driver
// will only be registered when the application is compiled for ARM.
//
// FIXME: See https://golang.org/pkg/go/build/#hdr-Build_Constraints for more
// information. For example, rename to _linux.go for code only relevant when
// running on linux, _windows.go for Windows, _darwin.go for OSX, etc.
//
// FIXME: Don't forget to remove all the FIXME comments before sending your PR!
// Otherwise the PR will me immediately refused.

package driver_skeleton

import "github.com/maruel/dlibox/go/pio"

func init() {
	pio.MustRegister(&driver{})
}
