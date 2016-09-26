// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// +build !linux

package internal

import (
	"runtime"
	"time"
)

// Nanosleep sleeps for a short amount of time doing a busy loop.
func Nanosleep(d time.Duration) {
	// TODO(maruel): That's not optimal.
	runtime.LockOSThread()
	for start := time.Now(); time.Since(start) < d; {
	}
	runtime.UnlockOSThread()
}
