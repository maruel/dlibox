// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// +build !linux

package gpiofs

import (
	"errors"
	"io"
)

// Handle opens the GPIO export only once for the process lifetime and returns
// it.
func Handle() (io.Writer, error) {
	return nil, errors.New("not supported on this OS")
}
