// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Feel free to copy-paste this file along the license when you need a quick
// and dirty SPI client.

package rpi

import (
	"fmt"
	"io"
	"os"
	"sync"
	"sync/atomic"
	"syscall"
	"unsafe"
)

// spi is a thread-safe SPI writer.
type spi struct {
	closed int32
	path   string
	speed  int
	lock   sync.Mutex
	f      *os.File
}

func makeSPI(path string, speed int) (*spi, error) {
	if path == "" {
		path = "/dev/spidev0.0"
	}
	f, err := os.OpenFile(path, os.O_RDWR, os.ModeExclusive)
	if err != nil {
		return nil, err
	}
	s := &spi{path: path, speed: speed, f: f}
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

// MakeSPI is to be used when testing directly to the bus bypassing the DotStar
// controller.
func MakeSPI(path string, speed int) (io.WriteCloser, error) {
	return makeSPI(path, speed)
}

func (s *spi) Close() error {
	if !atomic.CompareAndSwapInt32(&s.closed, 0, 1) {
		return io.ErrClosedPipe
	}
	s.lock.Lock()
	defer s.lock.Unlock()
	var err error
	if s.f != nil {
		err = s.f.Close()
		s.f = nil
	}
	return err
}

// Write pushes a buffer as-is.
func (s *spi) Write(b []byte) (int, error) {
	if atomic.LoadInt32(&s.closed) != 0 {
		return 0, io.ErrClosedPipe
	}
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.f.Write(b)
}

// Read returns a buffer as-is.
func (s *spi) Read(b []byte) (int, error) {
	if atomic.LoadInt32(&s.closed) != 0 {
		return 0, io.ErrClosedPipe
	}
	s.lock.Lock()
	defer s.lock.Unlock()
	n, err := s.f.Read(b)
	if err == nil && n != len(b) {
		err = io.ErrShortBuffer
	}
	return n, err
}

// Private details.

// spidev driver IOCTL control codes.
const (
	spiIOCMode        = 0x16B01
	spiIOCBitsPerWord = 0x16B03
	spiIOCMaxSpeedHz  = 0x46B04
)

func (s *spi) setFlag(op uint, arg uint64) error {
	if atomic.LoadInt32(&s.closed) != 0 {
		return io.ErrClosedPipe
	}
	s.lock.Lock()
	defer s.lock.Unlock()
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

func (s *spi) ioctl(op uint, arg unsafe.Pointer) error {
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, s.f.Fd(), uintptr(op), uintptr(arg)); errno != 0 {
		return fmt.Errorf("spi ioctl: %s", syscall.Errno(errno))
	}
	return nil
}
