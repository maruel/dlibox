// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package sysfs

import (
	"os"
	"syscall"

	"github.com/maruel/dlibox/go/pio/drivers"
)

type event [1]syscall.EpollEvent

func (e event) wait(ep int) (int, error) {
	return syscall.EpollWait(ep, e[:], -1)
}

func (e event) makeEvent(f *os.File) (int, error) {
	epollFd, err := syscall.EpollCreate(1)
	if err != nil {
		return 0, err
	}
	const EPOLLPRI = 2
	const EPOLL_CTL_ADD = 1
	fd := f.Fd()
	e[0].Events = EPOLLPRI
	e[0].Fd = int32(fd)
	return epollFd, syscall.EpollCtl(epollFd, EPOLL_CTL_ADD, int(fd), &e[0])
}

func init() {
	drivers.Register(&driverGPIO{})
}
