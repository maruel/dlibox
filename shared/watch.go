// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

//+build !linux

package shared

import "github.com/maruel/interrupt"

// WatchFile returns when the process' executable is modified or interrupt is
// set.
func WatchFile() error {
	<-interrupt.Channel
	return nil
}
