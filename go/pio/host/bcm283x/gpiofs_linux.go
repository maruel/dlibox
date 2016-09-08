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
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"syscall"
	"time"

	"github.com/maruel/dlibox/go/pio/host"
)

// Edges creates a edge detection loop and implements host.PinIn.
//
// This requires opening a gpio sysfs file handle. Make sure the user is member
// of group 'gpio'. The pin will be "exported" at /sys/class/gpio/gpio*/. Note
// that the pin will not be unexported at shutdown.
//
// For edge detection, the processor samples the input at its CPU clock rate
// and looks for '011' to rising and '100' for falling detection to avoid
// glitches. Because gpio sysfs is used, the latency is unpredictable.
func (p Pin) Edges() (chan host.Level, error) {
	if err := p.setEdge(true); err != nil {
		return nil, err
	}
	c := make(chan host.Level)
	go func() {
		defer p.setEdge(false)
		var b [1]byte
		for {
			l := host.Low
			if err := gpios[p].readPoll(b[:]); err != nil {
				// In case of error or unknown value, returns low. The file handle was
				// already opened so the chance of this happening is low.
				log.Printf("%s: error reading edge: %v", p, err)
			} else if b[0] == '1' {
				l = host.High
			}
			c <- l
		}
	}()
	return c, nil
}

// setEdge changes the edge detection setting for the pin.
//
// It is the function that opens the gpio sysfs file handle.
func (p Pin) setEdge(enable bool) error {
	if !enable {
		// Do not close the handles.
		return gpios[p].writeEdge([]byte("none"))
		gpios[p].usingEdge = false
		return nil
	}
	gpios[p].usingEdge = true
	if err := gpios[p].open(p); err != nil {
		return err
	}
	return gpios[p].writeEdge([]byte("both"))
}

// gpio is used for interrupt based edge detection.
type gpio struct {
	value     *os.File // handle to /sys/class/gpio/gpio*/value.
	edge      *os.File
	usingEdge bool
	epollFd   int
	event     [1]syscall.EpollEvent
}

// exportHandle is the handle to /sys/class/gpio/export
var exportHandle io.WriteCloser
var gpios [54]gpio

func (g *gpio) open(p Pin) error {
	var err error
	if g.value == nil {
		// Assume the pin is exported first. The reason is that exporting a pin that
		// is already exported causes a write failure, which is difficult to
		// differentiate from other errors.
		// On the other hand, accessing /sys/class/gpio/gpio*/value when it is not
		// exported returns a permission denied error. :/
		if g.value, err = os.OpenFile(fmt.Sprintf("/sys/class/gpio/gpio%d/value", p), os.O_RDONLY, 0600); err != nil {
			// Export the pin.
			if err = openExport(); err == nil {
				if _, err = exportHandle.Write([]byte(strconv.Itoa(int(p)))); err == nil {
					g.value, err = os.OpenFile(fmt.Sprintf("/sys/class/gpio/gpio%d/value", p), os.O_RDONLY, 0600)
				}
			}
		}
	}
	if g.edge == nil && err == nil {
		// TODO(maruel): Figure out the problem or better use the register instead
		// of the file.
		for i := 0; i < 30 && g.edge == nil; i++ {
			g.edge, err = os.OpenFile(fmt.Sprintf("/sys/class/gpio/gpio%d/edge", p), os.O_WRONLY, 0600)
			time.Sleep(time.Millisecond)
		}
	}
	if g.epollFd == 0 && err == nil {
		if g.epollFd, err = syscall.EpollCreate(1); err == nil {
			const EPOLLPRI = 2
			const EPOLL_CTL_ADD = 1
			g.event[0].Events = EPOLLPRI
			g.event[0].Fd = int32(g.value.Fd())
			err = syscall.EpollCtl(g.epollFd, EPOLL_CTL_ADD, int(g.value.Fd()), &g.event[0])
		}
	}
	if err != nil {
		_ = g.Close()
	}
	return err
}

func (g *gpio) Close() error {
	if g.value != nil {
		_ = g.value.Close()
		g.value = nil
	}
	if g.edge != nil {
		_ = g.edge.Close()
		g.edge = nil
	}
	if g.epollFd != 0 {
		_ = syscall.Close(g.epollFd)
		g.epollFd = 0
	}
	return nil
}

func (g *gpio) writeEdge(b []byte) error {
	_, err := g.edge.Write(b)
	return err
}

func (g *gpio) readPoll(b []byte) error {
	for {
		if nr, err := syscall.EpollWait(g.epollFd, g.event[:], -1); err != nil {
			return err
		} else if nr < 1 {
			continue
		}
		if _, err := g.value.Seek(0, 0); err != nil {
			return err
		}
		_, err := g.value.Read(b)
		return err
	}
}

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

func closeExport() error {
	// TODO(maruel): Unexport pins if ever desired. This is not really a problem
	// in practice.
	w := exportHandle
	exportHandle = nil
	for p := range gpios {
		gpios[p].Close()
	}
	return w.Close()
}
