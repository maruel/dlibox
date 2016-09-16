// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// +build !linux

package sysfs

import "errors"

var errNotImpl = errors.New("not implemented on non-linux OSes")

// Init returns an error no non-Linux OS.
func Init() error {
	return errNotImpl
}
