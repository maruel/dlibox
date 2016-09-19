// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package internal

import (
	"syscall"
	"time"
)

// Nanosleep sleeps for a short amount of time doing a busy loop.
func Nanosleep(d time.Duration) {
	// time.Sleep() sleeps for really too long, calling it repeatedly with
	// minimal value will give the caller a wake rate of 5KHz or so, depending on
	// the host. This makes it useless for bitbanging protocols.
	//
	// runtime.nanotime() is not exported so it cannot be used to busy loop for
	// very short sleep (10µs or less).
	time := syscall.NsecToTimespec(d.Nanoseconds())
	leftover := syscall.Timespec{}
	for {
		if err := syscall.Nanosleep(&time, &leftover); err != nil {
			time = leftover
			continue
		}
		break
	}
}
