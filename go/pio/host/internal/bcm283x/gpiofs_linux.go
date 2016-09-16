// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// GPIO sysfs handling code, as described at
// https://www.kernel.org/doc/Documentation/gpio/sysfs.txt
// See bcm238x.go for more details on how this code is used.
//
// GPIO sysfs is just one way of accessing the GPIO pins. A fun page is
// http://elinux.org/RPi_GPIO_Code_Samples which lists many ways.
//
// The only reason GPIO sysfs is used is because it's the only way to do edge
// triggered interrupts. Doing this requires cooperation from a driver in the
// kernel.
//
// All other functionality is using /dev/gpiomem since it is infinitely faster,
// and GPIO sysfs doesn't expose pull resistors.

package bcm283x

import (
	"io"
	"os"
)

// exportHandle is the handle to /sys/class/gpio/export
var exportHandle io.WriteCloser

func openExport() error {
	if exportHandle == nil {
		f, err := os.OpenFile("/sys/class/gpio/export", os.O_WRONLY, 0600)
		if err != nil {
			return err
		}
		exportHandle = f
	}
	return nil
}
