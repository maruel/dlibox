// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package sysfs

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"unsafe"

	"github.com/maruel/dlibox/go/pio/host/internal/pins2"
	"github.com/maruel/dlibox/go/pio/protocols/gpio"
	"github.com/maruel/dlibox/go/pio/protocols/pins"
	"github.com/maruel/dlibox/go/pio/protocols/spi"
)

// EnumerateSPI returns the available SPI buses.
//
// The first int is the bus number, the second is the chip select line.
func EnumerateSPI() ([][2]int, error) {
	prefix := "/dev/spidev"
	items, err := filepath.Glob(prefix + "*")
	if err != nil {
		return nil, err
	}
	out := make([][2]int, 0, len(items))
	for _, item := range items {
		parts := strings.Split(item[len(prefix):], ".")
		if len(parts) != 2 {
			continue
		}
		bus, err := strconv.Atoi(parts[0])
		if err != nil {
			continue
		}
		cs, err := strconv.Atoi(parts[1])
		if err != nil {
			continue
		}
		out = append(out, [2]int{bus, cs})
	}
	return out, nil
}

// SPI is an open SPI bus.
type SPI struct {
	f          *os.File
	busNumber  int
	chipSelect int
	clk        gpio.PinOut
	mosi       gpio.PinOut
	miso       gpio.PinIn
	cs         gpio.PinOut
}

func newSPI(busNumber, chipSelect int, speed int64) (*SPI, error) {
	if busNumber < 0 || busNumber > 255 {
		return nil, errors.New("invalid bus")
	}
	if chipSelect < 0 || chipSelect > 255 {
		return nil, errors.New("invalid chip select")
	}
	f, err := os.OpenFile(fmt.Sprintf("/dev/spidev%d.%d", busNumber, chipSelect), os.O_RDWR, os.ModeExclusive)
	if err != nil {
		return nil, err
	}
	s := &SPI{f: f, busNumber: busNumber, chipSelect: chipSelect}
	if speed != 0 {
		if err := s.Speed(speed); err != nil {
			s.Close()
			return nil, err
		}
	}
	if err := s.Configure(spi.Mode3, 8); err != nil {
		s.Close()
		return nil, err
	}
	return s, nil
}

// Close closes the handle to the SPI driver. It is not a requirement to close
// before process termination.
func (s *SPI) Close() error {
	err := s.f.Close()
	s.f = nil
	return err
}

func (s *SPI) Speed(hz int64) error {
	if hz < 1000 {
		return errors.New("invalid speed")
	}
	return s.setFlag(spiIOCMaxSpeedHz, uint64(hz))
}

func (s *SPI) Configure(mode spi.Mode, bits int) error {
	if bits < 1 || bits > 256 {
		return errors.New("invalid bits")
	}
	if err := s.setFlag(spiIOCMode, uint64(mode)); err != nil {
		return err
	}
	return s.setFlag(spiIOCBitsPerWord, uint64(bits))
}

func (s *SPI) Write(b []byte) (int, error) {
	return s.f.Write(b)
}

// Tx sends and receives data simultaneously.
func (s *SPI) Tx(w, r []byte) error {
	p := spiIOCTransfer{
		tx:          uint64(uintptr(unsafe.Pointer(&w[0]))),
		rx:          uint64(uintptr(unsafe.Pointer(&r[0]))),
		length:      uint32(len(w)),
		bitsPerWord: 8,
	}
	return s.ioctl(spiIOCTx|0x40000000, unsafe.Pointer(&p))
}

// CLK implements spi.Conn.
//
// It will fail if host.Init() wasn't called. host.Init() is transparently
// called by host.MakeSPI().
func (s *SPI) CLK() gpio.PinOut {
	s.initPins()
	return s.clk
}

// MISO implements spi.Conn.
//
// It will fail if host.Init() wasn't called. host.Init() is transparently
// called by host.MakeSPI().
func (s *SPI) MISO() gpio.PinIn {
	s.initPins()
	return s.miso
}

// MOSI implements spi.Conn.
//
// It will fail if host.Init() wasn't called. host.Init() is transparently
// called by host.MakeSPI().
func (s *SPI) MOSI() gpio.PinOut {
	s.initPins()
	return s.mosi
}

// CS implements spi.Conn.
//
// It will fail if host.Init() wasn't called. host.Init() is transparently
// called by host.MakeSPI().
func (s *SPI) CS() gpio.PinOut {
	s.initPins()
	return s.cs
}

// Private details.

const (
	cSHigh    spi.Mode = 0x4
	lSBFirst  spi.Mode = 0x8
	threeWire spi.Mode = 0x10
	loop      spi.Mode = 0x20
	noCS      spi.Mode = 0x40
)

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

func (s *SPI) setFlag(op uint, arg uint64) error {
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

func (s *SPI) ioctl(op uint, arg unsafe.Pointer) error {
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, s.f.Fd(), uintptr(op), uintptr(arg)); errno != 0 {
		return fmt.Errorf("spi ioctl: %s", syscall.Errno(errno))
	}
	return nil
}

func (s *SPI) initPins() {
	if s.clk == nil {
		ok := false
		if s.clk, ok = pins2.ByFunction[fmt.Sprintf("SPI%d_CLK", s.busNumber)]; !ok {
			s.clk = pins.INVALID
		}
		if s.miso, ok = pins2.ByFunction[fmt.Sprintf("SPI%d_MISO", s.busNumber)]; !ok {
			s.miso = pins.INVALID
		}
		if s.mosi, ok = pins2.ByFunction[fmt.Sprintf("SPI%d_MISO", s.busNumber)]; !ok {
			s.mosi = pins.INVALID
		}
		if s.cs, ok = pins2.ByFunction[fmt.Sprintf("SPI%d_CS%d", s.busNumber, s.chipSelect)]; !ok {
			s.cs = pins.INVALID
		}
	}
}

var _ spi.Conn = &SPI{}
