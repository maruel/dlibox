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

	"github.com/maruel/dlibox/go/pio/drivers"
	"github.com/maruel/dlibox/go/pio/protocols/gpio"
	"github.com/maruel/dlibox/go/pio/protocols/pins"
	"github.com/maruel/dlibox/go/pio/protocols/uart"
)

// enumerateUART returns the available serial buses.
func enumerateUART() ([]int, error) {
	// Do not use "/sys/class/tty/ttyS0/" as these are all owned by root.
	prefix := "/dev/ttyS"
	items, err := filepath.Glob(prefix + "*")
	if err != nil {
		return nil, err
	}
	out := make([]int, 0, len(items))
	for _, item := range items {
		i, err := strconv.Atoi(item[len(prefix):])
		if err != nil {
			continue
		}
		out = append(out, i)
	}
	return out, nil
}

// uART is an open serial bus via sysfs.
//
// TODO(maruel): It's not yet implemented so nothing is exported for now.
// Should probably defer to an already working library like
// https://github.com/tarm/serial
type uART struct {
	f         *os.File
	busNumber int
}

func newUART(busNumber int) (*uART, error) {
	// Use the devfs path for now.
	f, err := os.OpenFile(fmt.Sprintf("/dev/ttyS%d", busNumber), os.O_RDWR, os.ModeExclusive)
	if err != nil {
		return nil, err
	}
	u := &uART{f: f, busNumber: busNumber}
	return u, nil
}

func (u *uART) Close() error {
	err := u.f.Close()
	u.f = nil
	return err
}

func (u *uART) Configure(stopBit uart.Stop, parity uart.Parity, bits int) error {
	return errors.New("not implemented")
}

func (u *uART) Write(b []byte) (int, error) {
	return 0, errors.New("not implemented")
}

func (u *uART) Tx(w, r []byte) error {
	return errors.New("not implemented")
}

func (u *uART) Speed(hz int64) error {
	return errors.New("not implemented")
}

func (u *uART) RX() gpio.PinIn {
	return pins.INVALID
}

func (u *uART) TX() gpio.PinOut {
	return pins.INVALID
}

func (u *uART) RTS() gpio.PinIO {
	return pins.INVALID
}

func (u *uART) CTS() gpio.PinIO {
	return pins.INVALID
}

// TODO(maruel): Put again once the implementation is functional.
//var _ uart.Conn = &uART{}

// driverUART implements drivers.Driver.
type driverUART struct {
}

func (d *driverUART) String() string {
	return "sysfs-uart"
}

func (d *driverUART) Type() drivers.Type {
	return drivers.Bus
}

func (d *driverUART) Init() (bool, error) {
	return true, nil
}
