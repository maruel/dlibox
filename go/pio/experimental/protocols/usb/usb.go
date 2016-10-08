// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Package usb implements an USB device registry.
package usb

import (
	"fmt"
	"io"
	"sync"

	"github.com/maruel/dlibox/go/pio/protocols"
)

// ConnCloser hides the hell that libusb is.
type ConnCloser interface {
	io.Closer
	protocols.Conn
}

// Opener takes control of an already opened USB device.
type Opener func(dev ConnCloser) error

// Register registers a driver for an USB device.
//
// When this device is found, the factory will be called with a device handle.
func Register(name string, venid, devid uint16, opener Opener) error {
	lock.Lock()
	defer lock.Unlock()
	if _, ok := byName[name]; ok {
		return fmt.Errorf("registering the same USB %s %d/%d twice", name, venid, devid)
	}
	if _, ok := byNumber[venid]; !ok {
		byNumber[venid] = map[uint16]Opener{}
	}
	if _, ok := byNumber[venid][devid]; ok {
		return fmt.Errorf("registering the same USB %s %d/%d twice", name, venid, devid)
	}

	byName[name] = opener
	byNumber[venid][devid] = opener
	return nil
}

// OnDevice is called when a device is detected on an USB bus.
//
// When called with dev == nil, it still returns true or false to signal if it
// is a device that is registered.
func OnDevice(venid, devid uint16, dev ConnCloser) bool {
	lock.Lock()
	defer lock.Unlock()
	if m := byNumber[venid]; m != nil {
		if opener := m[devid]; opener != nil {
			if dev == nil {
				return true
			}
			if err := opener(dev); err != nil {
				dev.Close()
				return false
			}
			return true
		}
	}
	return false
}

//

var (
	lock     sync.Mutex
	byName   = map[string]Opener{}
	byNumber = map[uint16]map[uint16]Opener{}
)
