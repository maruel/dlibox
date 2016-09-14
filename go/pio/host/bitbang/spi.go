// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Specification
//
// Motorola never published a proper specification.
// http://electronics.stackexchange.com/questions/30096/spi-specifications
// http://www.nxp.com/files/microcontrollers/doc/data_sheet/M68HC11E.pdf page 120
// http://www.st.com/content/ccc/resource/technical/document/technical_note/58/17/ad/50/fa/c9/48/07/DM00054618.pdf/files/DM00054618.pdf/jcr:content/translations/en.DM00054618.pdf

package bitbang

import (
	"errors"
	"log"
	"syscall"
	"time"

	"github.com/maruel/dlibox/go/pio/host"
)

// SPI represents a SPI master implemented as bit-banging on 3 or 4 GPIO pins.
type SPI struct {
	sck       host.PinOut // Clock
	sdi       host.PinIn  // MISO
	sdo       host.PinOut // MOSI
	csn       host.PinOut // CS
	mode      host.Mode
	bits      int
	halfCycle time.Duration
}

// Configure implements host.SPI.
func (s *SPI) Configure(mode host.Mode, bits int) error {
	if mode != host.Mode3 {
		return errors.New("not implemented")
	}
	s.mode = mode
	s.bits = bits
	return nil
}

// Tx implements host.SPI.
//
// BUG(maruel): Implement mode.
// BUG(maruel): Implement bits.
// BUG(maruel): Test if read works.
func (s *SPI) Tx(w, r []byte) error {
	if len(r) != 0 && len(w) != len(r) {
		return errors.New("write and read buffers must be the same length")
	}
	if s.csn != nil {
		s.csn.Set(host.Low)
		s.sleepHalfCycle()
	}
	for i := uint(0); i < uint(len(w)*8); i++ {
		s.sdo.Set(w[i/8]&(1<<(i%8)) != 0)
		s.sleepHalfCycle()
		s.sck.Set(host.Low)
		s.sleepHalfCycle()
		if len(r) != 0 {
			if s.sdi.Read() == host.High {
				r[i/8] |= 1 << (i % 8)
			}
		}
		s.sck.Set(host.Low)
	}
	if s.csn != nil {
		s.csn.Set(host.High)
	}
	return nil
}

// Write implements host.SPI.
func (s *SPI) Write(d []byte) (int, error) {
	if err := s.Tx(d, nil); err != nil {
		return 0, err
	}
	return len(d), nil
}

// MakeSPI returns an object that communicates SPI over 3 or 4 pins.
//
// BUG(maruel): Completely untested.
//
// cs can be nil.
func MakeSPI(clk, mosi host.PinOut, miso host.PinIn, cs host.PinOut, speedHz int) (*SPI, error) {
	if err := clk.Out(); err != nil {
		return nil, err
	}
	clk.Set(host.High)
	if err := mosi.Out(); err != nil {
		return nil, err
	}
	mosi.Set(host.High)
	if err := miso.In(host.Up); err != nil {
		return nil, err
	}
	if cs != nil {
		if err := cs.Out(); err != nil {
			return nil, err
		}
		// Low to select.
		cs.Set(host.High)
	}
	s := &SPI{
		sck:       clk,
		sdi:       miso,
		sdo:       mosi,
		csn:       cs,
		mode:      host.Mode3,
		bits:      8,
		halfCycle: time.Second / time.Duration(speedHz) / time.Duration(2),
	}
	return s, nil
}

//

// sleep does a busy loop to act as fast as possible.
func (s *SPI) sleepHalfCycle() {
	// If time.Sleep(s.halfCycle) is used, we can expect roughly 5kHz or so. When
	// getting in the 1MHz range, the sleep is 500ns. Another option is
	// syscall.Nanosleep() or runtime.nanotime but the later is not exported. :(
	//for start := time.Now(); time.Since(start) < s.halfCycle; {
	//}
	time := syscall.NsecToTimespec(s.halfCycle.Nanoseconds())
	leftover := syscall.Timespec{}
	for {
		if err := syscall.Nanosleep(&time, &leftover); err != nil {
			time = leftover
			log.Printf("Nanosleep() -> %v: %v", leftover, err)
			continue
		}
		break
	}
}

var _ host.SPI = &SPI{}
