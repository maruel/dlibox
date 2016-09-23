// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// +build !linux

package sysfs

import (
	"errors"
	"os"
)

type event struct{}

func (e event) wait(ep int) (int, error) {
	return 0, errors.New("unreachable code")
}

func (e event) makeEvent(f *os.File) (int, error) {
	return 0, errors.New("unreachable code")
}
