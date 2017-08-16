// Copyright 2017 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	_ "net/http/pprof"
)

// mainDevice is the main function when running as a device (a node).
func mainDevice() error {
	return watchFile()
}
