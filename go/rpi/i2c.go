// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Feel free to copy-paste this file along the license when you need a quick
// and dirty I²C client.

package rpi

import (
	"errors"
	"fmt"
	"io"
	"os"
	"syscall"
)

// I2C is an open I²C bus.
type I2C struct {
	f    *os.File
	addr uint16
}

// MakeI2C opens an I²C bus via its sysfs interface as described at
// https://www.kernel.org/doc/Documentation/i2c/dev-interface It is not
// Raspberry Pi specific.
//
// `bus` should normally be 1, unless I2C0 was manually enabled.
//
// Spec: http://cache.nxp.com/documents/user_manual/UM10204.pdf
func MakeI2C(bus int) (*I2C, error) {
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
	i := &I2C{f: f}

	// TODO(maruel): Changing the speed is currently doing this for all devices.
	// https://github.com/raspberrypi/linux/issues/215
	// Need to access /sys/module/i2c_bcm2708/parameters/baudrate

	return i, nil
}

// Close closes the handle to the I²C driver. It is not a requirement to close
// before process termination.
func (i *I2C) Close() error {
	err := i.f.Close()
	i.f = nil
	return err
}

// Address changes the address of the device to communicate with.
func (i *I2C) Address(addr uint16) error {
	if i.addr != addr {
		// Addresses are at maximum 10 bits.
		if addr >= 1<<9 {
			return errors.New("invalid address")
		}
		// TODO(maruel): Add support for i2cTenbit when addr >= 0x400.
		if err := i.ioctl(i2cSlave, uintptr(addr)); err != nil {
			return err
		}
		i.addr = addr
	}
	return nil
}

/*
// Write writes one byte to the I²C bus without reading.
func (i *I2C) WriteByte(c byte) error {
	var b [1]byte
	b[0] = c
	_, err := i.Write(b[:])
	return err
}
*/

// Write writes to the I²C bus without reading.
func (i *I2C) Write(b []byte) (int, error) {
	return i.f.Write(b)
}

/*
// Read reads from the I²C bus.
//
// Returns io.ErrShortBuffer if the buffer was not filled.
func (i *I2C) Read(b []byte) (int, error) {
	n, err := i.f.Read(b)
	if err == nil && n != len(b) {
		err = io.ErrShortBuffer
	}
	return n, err
}

func (i *I2C) ReadByte() (byte, error) {
	var b [1]byte
	_, err := i.Read(b[:])
	return b[0], err
}
*/

// Write writes to the I²C bus, first a register, then data.
func (i *I2C) WriteReg(reg byte, data byte) error {
	_, err := i.f.Write(append([]byte{reg, data}))
	return err
}

// Read writes the register number to the I²C bus, then reads data.
func (i *I2C) ReadReg(reg byte, b []byte) error {
	// TODO(maruel): Use Tx() to do a single kernel call.
	if _, err := i.f.Write([]byte{reg}); err != nil {
		return err
	}
	n, err := i.f.Read(b)
	if err == nil && n != len(b) {
		err = io.ErrShortBuffer
	}
	return err
}

/*
// I2CMessage is one operation to do in a batched operation.
type I2CMessage struct {
	Addr  uint16
	Write bool
	Buf   []byte
}

// Tx sends and receives data simultaneously.
func (i *I2C) Tx(msgs ...I2CMessage) error {
	// Convert the messages to the internal format.
	packets := make([]i2cMsg, len(msgs))
	for i := range msgs {
		packets[i].addr = msgs[i].addr
		if msgs[i].write {
			packets[i].flags = 1
		}
		if len(msgs[i].buf) > 65535 {
			return errors.New("buffer too large")
		}
		packets[i].length = uint16(len(msgs[i].buf))
		packets[i].buf = uintptr(unsafe.Pointer(&msgs[i].buf[0]))
	}
	p := i2cRdwrIoctlData{
		msgs:  uintptr(unsafe.Pointer(&packets[0])),
		nmsgs: uint32(len(packets)),
	}
	return i.ioctl(i2cRdwr, uintptr(unsafe.Pointer(&p)))
}
*/

// Private details.

// i2cdev driver IOCTL control codes.
//
// Constants and structure definition can be found at
// /usr/include/linux/i2c-dev.h and /usr/include/linux/i2c.h.
const (
	i2cSlave = 0x703
	i2cRdwr  = 0x707
)

/*
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
*/

func (i *I2C) ioctl(op uint, arg uintptr) error {
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, i.f.Fd(), uintptr(op), arg); errno != 0 {
		return fmt.Errorf("i²c ioctl: %s", syscall.Errno(errno))
	}
	return nil
}
