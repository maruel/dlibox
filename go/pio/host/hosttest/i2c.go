// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package hosttest

import (
	"errors"
	"sync"

	"github.com/maruel/dlibox/go/pio/host"
)

// I2CIO registers the I/O that happened on fake I2C.
type I2CIO struct {
	Addr  uint16
	Write []byte
}

// I2C implements host.I2C. It registers everything written to it.
//
// BUG(maruel): I2C does not support reading yet.
type I2C struct {
	sync.Mutex
	Writes []I2CIO
}

// Tx currently only support writes.
func (i *I2C) Tx(addr uint16, w, r []byte) error {
	if len(r) != 0 {
		return errors.New("not implemented")
	}
	io := I2CIO{addr, make([]byte, len(w))}
	copy(io.Write, w)
	i.Lock()
	defer i.Unlock()
	i.Writes = append(i.Writes, io)
	return nil
}

var _ host.I2C = &I2C{}
