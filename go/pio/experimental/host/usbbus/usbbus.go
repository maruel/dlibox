// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package usbbus

import (
	"sort"
	"sync"

	gousb "github.com/kylelemons/gousb/usb"
	"github.com/maruel/dlibox/go/pio"
	"github.com/maruel/dlibox/go/pio/experimental/protocols/usb"
)

// Descriptor is a small subset of the USB descriptor.
type Descriptor struct {
	Bus   uint8
	Addr  uint8
	VenID uint16
	DevID uint16
}

// All returns all the USB devices detected.
func All() []Descriptor {
	lock.Lock()
	defer lock.Unlock()
	out := make([]Descriptor, len(all))
	copy(out, all)
	return out
}

//

var (
	lock sync.Mutex
	all  descriptors
)

type descriptors []Descriptor

func (d descriptors) Len() int      { return len(d) }
func (d descriptors) Swap(i, j int) { d[i], d[j] = d[j], d[i] }
func (d descriptors) Less(i, j int) bool {
	if d[i].Bus < d[j].Bus {
		return true
	}
	if d[i].Bus > d[j].Bus {
		return false
	}
	return d[i].Addr < d[j].Addr
}

// Options:
// - https://github.com/kylelemons/gousb (which was forked multiple times)
//   - https://github.com/truveris/gousb
// - https://github.com/gotmc/libusb
// The only one which does not require libusb but only works on linux:
// - https://github.com/swetland/go-usb/tree/master/src/usb

// dev is an open handle to an USB device.
//
// The device can disappear at any moment.
type dev struct {
	Descriptor
	name string
	d    *gousb.Device
	e    gousb.Endpoint
}

func (d *dev) String() string {
	return d.name
}

func (d *dev) Close() error {
	return d.d.Close()
}

func (d *dev) Write(b []byte) (int, error) {
	return d.e.Write(b)
}

func (d *dev) Tx(w, r []byte) error {
	if _, err := d.e.Write(w); err != nil {
		return err
	}
	if len(r) == 0 {
		return nil
	}
	_, err := d.e.Read(r)
	return err
}

// driver implements pio.Driver.
type driver struct {
}

func (d *driver) String() string {
	return "usb"
}

func (d *driver) Type() pio.Type {
	return pio.Bus
}

func (d *driver) Prerequisites() []string {
	return nil
}

func (d *driver) Init() (bool, error) {
	// I'd much prefer something that just talks to the OS instead of using
	// libusb. Especially we only require a small API surface.
	lock.Lock()
	defer lock.Unlock()
	option2()

	// TODO(maruel): Start an event loop when new devices are plugged in without
	// polling.
	// go func() { for { WaitForDevice(); usb.OnDevice(...) } }()
	return true, nil
}

// Getting go error:
// could not determine kind of name for C.LIBUSB_TRANSFER_TYPE_BULK_STREAM
/*
func option1() error {
	ctx, err := libusb.Init()
	if err != nil {
		return err
	}
	defer ctx.Close()
	devs, err := ctx.GetDeviceList()
	if err != nil {
		// TODO(maruel): This shouldn't be handled this way. Failures happen all
		// the time on USB, this doesn't mean the driver is faulty.
		return err
	}
	for _, dev := range devs {
		desc, err := dev.GetDeviceDescriptor()
		if err != nil {
			continue
		}
		if usb.OnDevice(d.VendorID, d.ProductID, nil) {
			h, err := dev.Open()
			if err != nil {
				continue
			}
			//usb.OnDevice(d.VendorID, d.ProductID, &dev{})
			h.Close()
		}
	}
	return err
}
*/

func option2() error {
	ctx := gousb.NewContext()
	defer ctx.Close()
	all = nil
	devs, err := ctx.ListDevices(func(d *gousb.Descriptor) bool {
		// Return true to keep the device open.
		desc := Descriptor{d.Bus, d.Address, uint16(d.Vendor), uint16(d.Product)}
		all = append(all, desc)
		return usb.OnDevice(uint16(d.Vendor), uint16(d.Product), nil)
	})
	sort.Sort(all)
	if err != nil {
		// TODO(maruel): This shouldn't be handled this way. Failures happen all
		// the time on USB, this doesn't mean the driver is faulty.
		return err
	}
	for _, d := range devs {
		name, err := d.GetStringDescriptor(1)
		if err != nil {
			d.Close()
			continue
		}
		// Control, isochronous or bulk?
		e, err := d.OpenEndpoint(1, 0, 0, 1|uint8(gousb.ENDPOINT_DIR_IN))
		if err != nil {
			d.Close()
			continue
		}
		desc := Descriptor{d.Bus, d.Address, uint16(d.Vendor), uint16(d.Product)}
		usb.OnDevice(uint16(d.Vendor), uint16(d.Product), &dev{desc, name, d, e})
	}
	return nil
}

func init() {
	pio.MustRegister(&driver{})
}

var _ pio.Driver = &driver{}
var _ usb.ConnCloser = &dev{}
