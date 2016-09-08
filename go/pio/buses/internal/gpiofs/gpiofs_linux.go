// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package gpiofs

import (
	"io"
	"os"
)

var err error
var exportHandle io.Writer

// Handle opens the GPIO export only once for the process lifetime and returns
// it.
func Handle() (io.Writer, error) {
	if exportHandle == nil && err == nil {
		var f *os.File
		f, err = os.OpenFile("/sys/class/gpio/export", os.O_WRONLY, 0600)
		exportHandle = f
	}
	return exportHandle, err
}
