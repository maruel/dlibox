// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Package i2c implements a sane I²C sysfs library that works with multiple
// devices.
package i2c

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"syscall"
	"unsafe"

	"github.com/maruel/dlibox/go/pio/buses"
)

// Bus is an open I²C bus.
//
// It can be used to communicate with multiple devices from multiple goroutines.
type Bus struct {
	f    *os.File
	l    sync.Mutex // In theory the kernel probably has an internal lock but not taking any chance.
	fn   functionality
	addr uint16
}

// Make opens an I²C bus via its sysfs interface as described at
// https://www.kernel.org/doc/Documentation/i2c/dev-interface It is not
// Raspberry Pi specific.
//
// The resulting object is safe for concurent use.
func Make(bus int) (*Bus, error) {
	f, err := os.OpenFile(fmt.Sprintf("/dev/i2c-%d", bus), os.O_RDWR, os.ModeExclusive)
	if err != nil {
		// Try to be helpful here. There are generally two cases:
		// - /dev/i2c-X doesn't exist. In this case, /boot/config.txt has to be
		//   edited to enable I²C then the device must be rebooted.
		// - permission denied. In this case, the user has to be added to plugdev.
		if os.IsNotExist(err) {
			return nil, errors.New("I²C is not configured; please follow instructions at https://github.com/maruel/dlibox/tree/master/go/setup")
		}
		return nil, fmt.Errorf("are you member of group 'plugdev'? please follow instructions at https://github.com/maruel/dlibox/tree/master/go/setup. %s", err)
	}
	b := &Bus{f: f, addr: 0xFFFF}

	// TODO(maruel): Changing the speed is currently doing this for all devices.
	// https://github.com/raspberrypi/linux/issues/215
	// Need to access /sys/module/i2c_bcm2708/parameters/baudrate

	// Query to know if 10 bits addresses are supported.
	if err = b.ioctl(ioctlFuncs, uintptr(unsafe.Pointer(&b.fn))); err != nil {
		return nil, err
	}
	return b, nil
}

// Close closes the handle to the I²C driver. It is not a requirement to close
// before process termination.
func (b *Bus) Close() error {
	b.l.Lock()
	defer b.l.Unlock()
	err := b.f.Close()
	b.f = nil
	return err
}

// Tx execute a transaction as a single operation unit.
func (b *Bus) Tx(ios []buses.IOFull) error {
	// Do a quick check first.
	if len(ios) == 0 {
		return nil
	}
	for i := range ios {
		if ios[i].Addr >= 0x400 || (ios[i].Addr >= 0x80 && b.fn&func10BIT_ADDR == 0) {
			return nil
		}
		if len(ios[i].Buf) == 0 {
			return errors.New("buffer is empty")
		}
		if len(ios[i].Buf) > 65535 {
			return errors.New("buffer too large")
		}
	}
	op := ios[len(ios)-1].Op
	if op != buses.WriteStop && op != buses.ReadStop {
		return errors.New("last operation must be Stop")
	}
	return b.txFast(ios)
}

// txSlow sends and receives data as a single transaction by simply using a mutex.
func (b *Bus) txSlow(ios []buses.IOFull) error {
	// A limitation of txSlow() only.
	addr := ios[0].Addr
	for i := 1; i < len(ios); i++ {
		if addr != ios[i].Addr {
			return errors.New("add operations must be on the same address")
		}
	}

	b.l.Lock()
	defer b.l.Unlock()
	if err := b.setAddr(addr); err != nil {
		return err
	}
	// TODO(maruel): Merge multiple operations together.
	for i := 0; i < len(ios); i++ {
		switch ios[i].Op {
		case buses.Write, buses.WriteStop:
			if _, err := b.f.Write(ios[i].Buf); err != nil {
				return err
			}
		case buses.Read, buses.ReadStop:
			if _, err := b.f.Read(ios[i].Buf); err != nil {
				return err
			}
		}
	}
	return nil
}

// setAddr must be called with lock held.
func (b *Bus) setAddr(addr uint16) error {
	if b.addr != addr {
		if err := b.ioctl(ioctlSlave, uintptr(addr)); err != nil {
			return err
		}
		b.addr = addr
	}
	return nil
}

// txFast does a transaction as a single kernel call.
//
// Causes memory allocation but still less costly than doing multiple kernel
// calls.
func (b *Bus) txFast(ios []buses.IOFull) error {
	// Convert the messages to the internal format.
	msgs := make([]i2cMsg, len(ios))
	last := buses.WriteStop
	for i := range ios {
		msgs[i].addr = ios[i].Addr
		switch ios[i].Op {
		case buses.Write, buses.WriteStop:
			if last == buses.Write {
				msgs[i].flags = flagNOSTART
			}
		case buses.Read, buses.ReadStop:
			if last == buses.Read {
				msgs[i].flags = flagRD | flagNOSTART
			} else {
				msgs[i].flags = flagRD
			}
		}
		last = ios[i].Op
		msgs[i].length = uint16(len(ios[i].Buf))
		msgs[i].buf = uintptr(unsafe.Pointer(&ios[i].Buf[0]))
	}
	// Doesn't seem to work, need investigation.
	p := rdwrIoctlData{
		msgs:  uintptr(unsafe.Pointer(&msgs[0])),
		nmsgs: uint32(len(msgs)),
	}
	pp := uintptr(unsafe.Pointer(&p))

	b.l.Lock()
	defer b.l.Unlock()
	return b.ioctl(ioctlRdwr, pp)
}

func (b *Bus) ioctl(op uint, arg uintptr) error {
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, b.f.Fd(), uintptr(op), arg); errno != 0 {
		return fmt.Errorf("i²c ioctl: %s", syscall.Errno(errno))
	}
	return nil
}

// Private details.

// i2cdev driver IOCTL control codes.
//
// Constants and structure definition can be found at
// /usr/include/linux/i2c-dev.h and /usr/include/linux/i2c.h.
const (
	ioctlRetries = 0x701 // TODO(maruel): Expose this
	ioctlTimeout = 0x702 // TODO(maruel): Expose this; in units of 10ms
	ioctlSlave   = 0x703
	ioctlTenBits = 0x704 // TODO(maruel): Expose this but the header says it's broken (!?)
	ioctlFuncs   = 0x705
	ioctlRdwr    = 0x707
)

// flags
const (
	flagTEN          = 0x0010 // this is a ten bit chip address
	flagRD           = 0x0001 // read data, from slave to master
	flagSTOP         = 0x8000 // if I2C_FUNC_PROTOCOL_MANGLING
	flagNOSTART      = 0x4000 // if I2C_FUNC_NOSTART
	flagREV_DIR_ADDR = 0x2000 // if I2C_FUNC_PROTOCOL_MANGLING
	flagIGNORE_NAK   = 0x1000 // if I2C_FUNC_PROTOCOL_MANGLING
	flagNO_RD_ACK    = 0x0800 // if I2C_FUNC_PROTOCOL_MANGLING
	flagRECV_LEN     = 0x0400 // length will be first received byte

)

type functionality uint64

const (
	funcI2C                    = 0x00000001
	func10BIT_ADDR             = 0x00000002
	funcPROTOCOL_MANGLING      = 0x00000004 // I2C_M_IGNORE_NAK etc.
	funcSMBUS_PEC              = 0x00000008
	funcNOSTART                = 0x00000010 // I2C_M_NOSTART
	funcSMBUS_BLOCK_PROC_CALL  = 0x00008000 // SMBus 2.0
	funcSMBUS_QUICK            = 0x00010000
	funcSMBUS_READ_BYTE        = 0x00020000
	funcSMBUS_WRITE_BYTE       = 0x00040000
	funcSMBUS_READ_BYTE_DATA   = 0x00080000
	funcSMBUS_WRITE_BYTE_DATA  = 0x00100000
	funcSMBUS_READ_WORD_DATA   = 0x00200000
	funcSMBUS_WRITE_WORD_DATA  = 0x00400000
	funcSMBUS_PROC_CALL        = 0x00800000
	funcSMBUS_READ_BLOCK_DATA  = 0x01000000
	funcSMBUS_WRITE_BLOCK_DATA = 0x02000000
	funcSMBUS_READ_I2C_BLOCK   = 0x04000000 // I2C-like block xfer
	funcSMBUS_WRITE_I2C_BLOCK  = 0x08000000 // w/ 1-byte reg. addr.
)

func (f functionality) String() string {
	var out []string
	if f&funcI2C != 0 {
		out = append(out, "I2C")
	}
	if f&func10BIT_ADDR != 0 {
		out = append(out, "10BIT_ADDR")
	}
	if f&funcPROTOCOL_MANGLING != 0 {
		out = append(out, "PROTOCOL_MANGLING")
	}
	if f&funcSMBUS_PEC != 0 {
		out = append(out, "SMBUS_PEC")
	}
	if f&funcNOSTART != 0 {
		out = append(out, "NOSTART")
	}
	if f&funcSMBUS_BLOCK_PROC_CALL != 0 {
		out = append(out, "SMBUS_BLOCK_PROC_CALL")
	}
	if f&funcSMBUS_QUICK != 0 {
		out = append(out, "SMBUS_QUICK")
	}
	if f&funcSMBUS_READ_BYTE != 0 {
		out = append(out, "SMBUS_READ_BYTE")
	}
	if f&funcSMBUS_WRITE_BYTE != 0 {
		out = append(out, "SMBUS_WRITE_BYTE")
	}
	if f&funcSMBUS_READ_BYTE_DATA != 0 {
		out = append(out, "SMBUS_READ_BYTE_DATA")
	}
	if f&funcSMBUS_WRITE_BYTE_DATA != 0 {
		out = append(out, "SMBUS_WRITE_BYTE_DATA")
	}
	if f&funcSMBUS_READ_WORD_DATA != 0 {
		out = append(out, "SMBUS_READ_WORD_DATA")
	}
	if f&funcSMBUS_WRITE_WORD_DATA != 0 {
		out = append(out, "SMBUS_WRITE_WORD_DATA")
	}
	if f&funcSMBUS_PROC_CALL != 0 {
		out = append(out, "SMBUS_PROC_CALL")
	}
	if f&funcSMBUS_READ_BLOCK_DATA != 0 {
		out = append(out, "SMBUS_READ_BLOCK_DATA")
	}
	if f&funcSMBUS_WRITE_BLOCK_DATA != 0 {
		out = append(out, "SMBUS_WRITE_BLOCK_DATA")
	}
	if f&funcSMBUS_READ_I2C_BLOCK != 0 {
		out = append(out, "SMBUS_READ_I2C_BLOCK")
	}
	if f&funcSMBUS_WRITE_I2C_BLOCK != 0 {
		out = append(out, "SMBUS_WRITE_I2C_BLOCK")
	}
	return strings.Join(out, "|")
}

type rdwrIoctlData struct {
	msgs  uintptr // Pointer to i2cMsg
	nmsgs uint32
}

type i2cMsg struct {
	addr   uint16 // Address to communicate with
	flags  uint16 // 1 for read, see i2c.h for more details
	length uint16
	buf    uintptr
}
