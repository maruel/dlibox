// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Package i2c implements a sane I²C sysfs library that works with multiple
// devices.
package i2c

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"syscall"
	"unsafe"
)

// Bus is an open I²C bus.
//
// It can be used to communicate with multiple devices from multiple goroutines.
type Bus struct {
	f *os.File
	l sync.Mutex // In theory the kernel probably has an internal lock but not taking any chance.
}

// Make opens an I²C bus via its sysfs interface as described at
// https://www.kernel.org/doc/Documentation/i2c/dev-interface It is not
// Raspberry Pi specific.
//
// `bus` should normally be 1, unless I2C0 was manually enabled.
//
// Spec: http://cache.nxp.com/documents/user_manual/UM10204.pdf
func Make(bus int) (*Bus, error) {
	f, err := os.OpenFile(fmt.Sprintf("/dev/i2c-%d", bus), os.O_RDWR, os.ModeExclusive)
	if err != nil {
		// Try to be helpful here. There are generally two cases:
		// - /dev/i2c-X doesn't exist. In this case, /boot/config.txt has to be
		//   edited to enable I²C then the device must be rebooted.
		// - permission denied. In this case, the user has to be added to plugdev.
		if os.IsNotExist(err) {
			return nil, errors.New("I²C is not configured; please follow instructions at https://github.com/maruel/dlibox/tree/master/go/setup")
		}
		return nil, fmt.Errorf("are you member of group 'plugdev'? please follow instructions at https://github.com/maruel/dlibox/tree/master/go/setup. %s", err)
	}
	i := &Bus{f: f}

	// TODO(maruel): Changing the speed is currently doing this for all devices.
	// https://github.com/raspberrypi/linux/issues/215
	// Need to access /sys/module/i2c_bcm2708/parameters/baudrate

	return i, nil
}

// Close closes the handle to the I²C driver. It is not a requirement to close
// before process termination.
func (i *Bus) Close() error {
	i.l.Lock()
	defer i.l.Unlock()
	err := i.f.Close()
	i.f = nil
	return err
}

// Device selects a device on a I²C bus.
func (i *Bus) Device(addr uint16) *Dev {
	return &Dev{i, addr}
}

// tx sends and receives data simultaneously.
func (i *Bus) tx(msgs []i2cMsg) error {
	p := i2cRdwrIoctlData{
		msgs:  uintptr(unsafe.Pointer(&msgs[0])),
		nmsgs: uint32(len(msgs)),
	}
	i.l.Lock()
	defer i.l.Unlock()
	return i.ioctl(i2cRdwr, uintptr(unsafe.Pointer(&p)))
}

// Dev is a device on a I²C bus.
type Dev struct {
	i    *Bus
	addr uint16
}

// Write writes to the I²C bus without reading.
//
// It's a wrapper for Tx()
func (d *Dev) Write(b []byte) (int, error) {
	if err := d.Tx([]Cmd{{true, b}}); err != nil {
		return 0, err
	}
	return len(b), nil
}

// ReadReg writes the register number to the I²C bus, then reads data.
func (d *Dev) ReadReg(reg byte, b []byte) error {
	return d.Tx([]Cmd{{Buf: []byte{reg}}, {true, b}})
}

// Cmd is one command to do inside a transaction.
type Cmd struct {
	Write bool
	Buf   []byte
}

// Tx does a transaction.
func (d *Dev) Tx(cmds []Cmd) error {
	// Convert the messages to the internal format.
	msgs := make([]i2cMsg, len(cmds))
	for i := range cmds {
		if len(cmds[i].Buf) > 65535 {
			return errors.New("buffer too large")
		}
		msgs[i].addr = d.addr
		if cmds[i].Write {
			msgs[i].flags = 1
		}
		msgs[i].length = uint16(len(cmds[i].Buf))
		msgs[i].buf = uintptr(unsafe.Pointer(&cmds[i].Buf[0]))
	}
	return d.i.tx(msgs)
}

// Private details.

// i2cdev driver IOCTL control codes.
//
// Constants and structure definition can be found at
// /usr/include/linux/i2c-dev.h and /usr/include/linux/i2c.h.
const (
	i2cSlave = 0x703
	i2cRdwr  = 0x707
)

type i2cRdwrIoctlData struct {
	msgs  uintptr // Pointer to i2cMsg
	nmsgs uint32
}

type i2cMsg struct {
	addr   uint16 // Address to communicate with
	flags  uint16 // 1 for read, see i2c.h for more details
	length uint16
	buf    uintptr
}

func (i *Bus) ioctl(op uint, arg uintptr) error {
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, i.f.Fd(), uintptr(op), arg); errno != 0 {
		return fmt.Errorf("i²c ioctl: %s", syscall.Errno(errno))
	}
	return nil
}
