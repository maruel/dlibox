// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Package spi implements a sane SPI sysfs library that is generic.
package spi

import (
	"errors"
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

// Bus is an open SPI bus.
type Bus struct {
	f *os.File
}

// Make opens a *Bus via its sysfs interface as described at
// https://www.kernel.org/doc/Documentation/spi/spidev. It is not Raspberry Pi
// specific.
//
// `bus` should normally be 0, unless SPI1 was manually enabled.
//
// `chipSelect` should normally be 0, unless CE lines were manually enabled.
//
// `speed` must be specified and should be in the high Khz or low Mhz range,
// it's a good idea to start at 4000000 (4Mhz) and go upward as long as the
// signal is good.
func Make(bus, chipSelect int, speed int64) (*Bus, error) {
	if speed < 1000 {
		return nil, errors.New("invalid speed")
	}
	f, err := os.OpenFile(fmt.Sprintf("/dev/spidev%d.%d", bus, chipSelect), os.O_RDWR, os.ModeExclusive)
	if err != nil {
		// Try to be helpful here. There are generally two cases:
		// - /dev/spidevX.Y doesn't exist. In this case, /boot/config.txt has to be
		//   edited to enable SPI then the device must be rebooted.
		// - permission denied. In this case, the user has to be added to plugdev.
		if os.IsNotExist(err) {
			return nil, errors.New("SPI is not configured; please follow instructions at https://github.com/maruel/dlibox/tree/master/go/setup")
		}
		return nil, fmt.Errorf("are you member of group 'plugdev'? please follow instructions at https://github.com/maruel/dlibox/tree/master/go/setup. %s", err)
	}
	s := &Bus{f: f}
	if err := s.setFlag(spiIOCMode, 3); err != nil {
		s.Close()
		return nil, err
	}
	if err := s.setFlag(spiIOCBitsPerWord, 8); err != nil {
		s.Close()
		return nil, err
	}
	if err := s.setFlag(spiIOCMaxSpeedHz, uint64(speed)); err != nil {
		s.Close()
		return nil, err
	}
	return s, nil
}

// Close closes the handle to the SPI driver. It is not a requirement to close
// before process termination.
func (s *Bus) Close() error {
	err := s.f.Close()
	s.f = nil
	return err
}

// Write writes to the SPI bus without reading.
func (s *Bus) Write(b []byte) (int, error) {
	return s.f.Write(b)
}

// Tx sends and receives data simultaneously.
func (s *Bus) Tx(w, r []byte) error {
	p := spiIOCTransfer{
		tx:          uint64(uintptr(unsafe.Pointer(&w[0]))),
		rx:          uint64(uintptr(unsafe.Pointer(&r[0]))),
		length:      uint32(len(w)),
		bitsPerWord: 8,
	}
	return s.ioctl(spiIOCTx|0x40000000, unsafe.Pointer(&p))
}

// Private details.

// spidev driver IOCTL control codes.
//
// Constants and structure definition can be found at
// /usr/include/linux/spi/spidev.h.
const (
	spiIOCMode        = 0x16B01
	spiIOCBitsPerWord = 0x16B03
	spiIOCMaxSpeedHz  = 0x46B04
	spiIOCTx          = 0x206B00
)

type spiIOCTransfer struct {
	tx          uint64 // Pointer to byte slice
	rx          uint64 // Pointer to byte slice
	length      uint32
	speedHz     uint32
	delayUsecs  uint16
	bitsPerWord uint8
	csChange    uint8
	txNBits     uint8
	rxNBits     uint8
	pad         uint16
}

func (s *Bus) setFlag(op uint, arg uint64) error {
	if err := s.ioctl(op|0x40000000, unsafe.Pointer(&arg)); err != nil {
		return err
	}
	actual := uint64(0)
	// getFlag() equivalent.
	if err := s.ioctl(op|0x80000000, unsafe.Pointer(&actual)); err != nil {
		return err
	}
	if actual != arg {
		return fmt.Errorf("spi op 0x%x: set 0x%x, read 0x%x", op, arg, actual)
	}
	return nil
}

func (s *Bus) ioctl(op uint, arg unsafe.Pointer) error {
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, s.f.Fd(), uintptr(op), uintptr(arg)); errno != 0 {
		return fmt.Errorf("spi ioctl: %s", syscall.Errno(errno))
	}
	return nil
}
