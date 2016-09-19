// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package sysfs

import (
	"os"
	"syscall"
)

// Init initializes GPIO sysfs handling code.
//
// Uses gpio sysfs as described at
// https://www.kernel.org/doc/Documentation/gpio/sysfs.txt
//
// GPIO sysfs is often the only way to do edge triggered interrupts. Doing this
// requires cooperation from a driver in the kernel.
//
// The main drawback of GPIO sysfs is that it doesn't expose internal pull
// resistor and it is much slower than using memory mapped hardware registers.
//
// Init returns an error no non-Linux OS.
func Init() error {
	return initLinux()
}

//

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
