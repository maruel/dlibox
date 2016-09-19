// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// +build !arm

package bcm283x

import "errors"

// Init initializes the Broadcom bcm283x CPU GPIO registers if relevant.
func Init() error {
	return errors.New("not supported on this platform")
}
