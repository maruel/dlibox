// Copyright 2014 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Package interrupt is a single global way to handle process interruption.
//
// It is useful for both long lived process to implement controlled shutdown
// and for CLI tool to handle early termination.
//
// The interrupt signal can be set exactly once in the process lifetime and
// cannot be unset. The signal can optionally be set automatically on
// Ctrl-C/os.Interrupt. When set, it is expected the process to abort any
// on-going execution early.
//
// The signal can be read via two ways:
//
//   select {
//   case <- interrupt.Channel:
//     // Handle abort.
//   case ...
//     ...
//   default:
//   }
//
// or
//
//   if interrupt.IsSet() {
//     // Handle abort.
//   }
package interrupt

import (
	"errors"
	"os"
	"os/signal"
	"sync/atomic"
)

// ErrInterrupted can be used as an error to signal that a process was
// interrupted but didn't fail in any other way. This permits disambiguating
// from any other genuine error.
var ErrInterrupted = errors.New("interrupted")

// Set sets the interrupt signal. It is to be used when the process must exit
// as soon as possible.
func Set() {
	atomic.StoreInt32(&interrupted, 1)
	go func() {
		for {
			interruptedChannel <- true
		}
	}()
}

// IsSet returns true once the interrupt signal was set. It is meant to be used
// when polling for status instead of using channel selection.
func IsSet() bool {
	return atomic.LoadInt32(&interrupted) != 0
}

// Channel continuously sends true once the interrupt signal was set. It can be
// used in select section to handle interrupted process.
var Channel <-chan bool

// HandleCtrlC initializes an handler to handle SIGINT, which is normally sent
// on Ctrl-C.
//
// This function is provided for convenience. To handle other situations, use
// Set() directly.
func HandleCtrlC() {
	chanSignal := make(chan os.Signal)

	go func() {
		<-chanSignal
		Set()
	}()

	signal.Notify(chanSignal, os.Interrupt)
}

var interrupted int32

var interruptedChannel chan<- bool

func init() {
	c := make(chan bool)
	interruptedChannel = c
	Channel = c
}
