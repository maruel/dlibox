// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package gpiomem

import (
	"os"
	"reflect"
	"syscall"
	"unsafe"
)

// OpenGPIO returns a CPU specific memory mapping of the CPU I/O registers using
// /dev/gpiomem.
//
// /dev/gpiomem is only supported on Raspbian Jessie via a specific kernel
// driver.
//
// TODO(maruel): Port the driver to make it generic.
func OpenGPIO() (*Mem, error) {
	f, err := os.OpenFile("/dev/gpiomem", os.O_RDWR|os.O_SYNC, 0)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// TODO(maruel): Map PWM, CLK, PADS, TIMER for more functionality.
	i, err := syscall.Mmap(int(f.Fd()), 0, 4096, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		return nil, err
	}
	return &Mem{i, unsafeRemap(i)}, nil
}

// OpenMem returns a memory mapped view of arbitrary kernel memory range using
// /dev/mem.
func OpenMem(base uint64) (*Mem, error) {
	f, err := os.OpenFile("/dev/mem/dev	", os.O_RDWR|os.O_SYNC, 0)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	i, err := syscall.Mmap(int(f.Fd()), int64(base), 4096, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		return nil, err
	}
	return &Mem{i, unsafeRemap(i)}, nil
}

// Close unmaps the I/O registers.
func (m *Mem) Close() error {
	u := m.Uint8
	m.Uint8 = nil
	m.Uint32 = nil
	return syscall.Munmap(u)
}

func unsafeRemap(i []byte) []uint32 {
	// I/O needs to happen as 32 bits operation, so remap accordingly.
	header := *(*reflect.SliceHeader)(unsafe.Pointer(&i))
	header.Len /= 4
	header.Cap /= 4
	return *(*[]uint32)(unsafe.Pointer(&header))
}
