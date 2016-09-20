// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package hosttest

import (
	"bytes"
	"errors"
	"fmt"
	"sync"

	"github.com/maruel/dlibox/go/pio/host"
)

// I2CIO registers the I/O that happened on either a real or fake I2C bus.
type I2CIO struct {
	Addr  uint16
	Write []byte
	Read  []byte
}

// I2CRecord implements host.I2C that records everything written to it.
//
// This can then be used to feed to I2CPlayback to do "replay" based unit tests.
type I2CRecord struct {
	Bus  host.I2C // Bus can be nil if only writes are being recorded.
	Lock sync.Mutex
	Ops  []I2CIO
}

// Tx implements host.I2C.
func (i *I2CRecord) Tx(addr uint16, w, r []byte) error {
	i.Lock.Lock()
	defer i.Lock.Unlock()
	if i.Bus == nil {
		if len(r) != 0 {
			return errors.New("read unsupported when no bus is connected")
		}
	} else {
		if err := i.Bus.Tx(addr, w, r); err != nil {
			return err
		}
	}
	io := I2CIO{Addr: addr, Write: make([]byte, len(w))}
	if len(r) != 0 {
		io.Read = make([]byte, len(r))
	}
	copy(io.Write, w)
	copy(io.Read, r)
	i.Ops = append(i.Ops, io)
	return nil
}

// I2CPlayblack implements host.I2C and plays back a recorded I/O flow.
//
// While "replay" type of unit tests are of limited value, they still present
// an easy way to do basic code coverage.
//
// BUG(maruel): Have it work as a memory mapped registers, which is how most
// devices being tested work.
type I2CPlayback struct {
	Lock sync.Mutex
	Ops  []I2CIO
}

// Tx implements host.I2C.
func (i *I2CPlayback) Tx(addr uint16, w, r []byte) error {
	i.Lock.Lock()
	defer i.Lock.Unlock()
	if len(i.Ops) == 0 {
		// log.Fatal() ?
		return errors.New("unexpected Tx()")
	}
	if addr != i.Ops[0].Addr {
		return fmt.Errorf("unexpected addr %d != %d", addr, i.Ops[0].Addr)
	}
	if !bytes.Equal(i.Ops[0].Write, w) {
		return fmt.Errorf("unexpected write %#v != %#v", w, i.Ops[0].Write)
	}
	if len(i.Ops[0].Read) != len(r) {
		return fmt.Errorf("unexpected read buffer length %d != %d", len(r), len(i.Ops[0].Read))
	}
	copy(r, i.Ops[0].Read)
	i.Ops = i.Ops[1:]
	return nil
}

var _ host.I2C = &I2CRecord{}
var _ host.I2C = &I2CPlayback{}
