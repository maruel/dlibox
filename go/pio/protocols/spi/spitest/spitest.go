// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Package spitest is meant to be used to test drivers over a fake SPI bus.
package spitest

import (
	"errors"
	"io"
	"sync"

	"github.com/maruel/dlibox/go/pio/protocols/gpio"
	"github.com/maruel/dlibox/go/pio/protocols/protocolstest"
	"github.com/maruel/dlibox/go/pio/protocols/spi"
)

// RecordRaw implements spi.Conn. It sends everything written to it to W.
type RecordRaw struct {
	protocolstest.RecordRaw
}

// NewRecordRaw is a shortcut to create a RecordRaw
func NewRecordRaw(w io.Writer) *RecordRaw {
	return &RecordRaw{protocolstest.RecordRaw{W: w}}
}

// Close is a no-op.
func (r *RecordRaw) Close() error {
	r.Lock.Lock()
	defer r.Lock.Unlock()
	return nil
}

// Speed is a no-op.
func (r *RecordRaw) Speed(hz int64) error {
	return nil
}

// Configure is a no-op.
func (r *RecordRaw) Configure(mode spi.Mode, bits int) error {
	return nil
}

// Record implements spi.Conn that records everything written to it.
//
// This can then be used to feed to Playback to do "replay" based unit tests.
type Record struct {
	Conn spi.Conn // Conn can be nil if only writes are being recorded.
	Lock sync.Mutex
	Ops  []protocolstest.IO
}

func (r *Record) String() string {
	return "record"
}

func (r *Record) Write(d []byte) (int, error) {
	if err := r.Tx(d, nil); err != nil {
		return 0, err
	}
	return len(d), nil
}

func (r *Record) Tx(w, read []byte) error {
	r.Lock.Lock()
	defer r.Lock.Unlock()
	if r.Conn == nil {
		if len(read) != 0 {
			return errors.New("read unsupported when no bus is connected")
		}
	} else {
		if err := r.Conn.Tx(w, read); err != nil {
			return err
		}
	}
	io := protocolstest.IO{Write: make([]byte, len(w))}
	if len(read) != 0 {
		io.Read = make([]byte, len(read))
	}
	copy(io.Write, w)
	copy(io.Read, read)
	r.Ops = append(r.Ops, io)
	return nil
}

func (r *Record) Speed(hz int64) error {
	if r.Conn != nil {
		return r.Conn.Speed(hz)
	}
	return nil
}

func (r *Record) Configure(mode spi.Mode, bits int) error {
	if r.Conn != nil {
		return r.Conn.Configure(mode, bits)
	}
	return nil
}

func (r *Record) CLK() gpio.PinOut {
	if p, ok := r.Conn.(spi.Pins); ok {
		return p.CLK()
	}
	return gpio.INVALID
}

func (r *Record) MOSI() gpio.PinOut {
	if p, ok := r.Conn.(spi.Pins); ok {
		return p.MOSI()
	}
	return gpio.INVALID
}

func (r *Record) MISO() gpio.PinIn {
	if p, ok := r.Conn.(spi.Pins); ok {
		return p.MISO()
	}
	return gpio.INVALID
}

func (r *Record) CS() gpio.PinOut {
	if p, ok := r.Conn.(spi.Pins); ok {
		return p.CS()
	}
	return gpio.INVALID
}

// Playblack implements spi.Conn and plays back a recorded I/O flow.
//
// While "replay" type of unit tests are of limited value, they still present
// an easy way to do basic code coverage.
type Playback struct {
	protocolstest.Playback
}

func (p *Playback) Speed(hz int64) error {
	return nil
}

func (p *Playback) Configure(mode spi.Mode, bits int) error {
	return nil
}

var _ spi.Conn = &RecordRaw{}
var _ spi.Conn = &Record{}
var _ spi.Pins = &Record{}
var _ spi.Conn = &Playback{}
