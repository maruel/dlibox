// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package allwinner_pl

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/maruel/dlibox/go/pio"
	"github.com/maruel/dlibox/go/pio/host/internal"
	"github.com/maruel/dlibox/go/pio/host/internal/gpiomem"
	"github.com/maruel/dlibox/go/pio/host/sysfs"
	"github.com/maruel/dlibox/go/pio/protocols/gpio"
)

var Pins = []Pin{
	{offset: 0, name: "PL0", defaultPull: gpio.Up},
	{offset: 1, name: "PL1", defaultPull: gpio.Up},
	{offset: 2, name: "PL2", defaultPull: gpio.Float},
	{offset: 3, name: "PL3", defaultPull: gpio.Float},
	{offset: 4, name: "PL4", defaultPull: gpio.Float},
	{offset: 5, name: "PL5", defaultPull: gpio.Float},
	{offset: 6, name: "PL6", defaultPull: gpio.Float},
	{offset: 7, name: "PL7", defaultPull: gpio.Float},
	{offset: 8, name: "PL8", defaultPull: gpio.Float},
	{offset: 9, name: "PL9", defaultPull: gpio.Float},
	{offset: 10, name: "PL10", defaultPull: gpio.Float},
	{offset: 11, name: "PL11", defaultPull: gpio.Float},
	{offset: 12, name: "PL12", defaultPull: gpio.Float},
}

type Pin struct {
	offset      uint8      // as per register offset calculation
	name        string     // name as per datasheet
	defaultPull gpio.Pull  // default pull at startup
	edge        *sysfs.Pin // Mutable, set once, then never set back to nil
}

// http://forum.pine64.org/showthread.php?tid=474
// about number calculation.
var (
	PL0  gpio.PinIO = &Pins[0]  // 352
	PL1  gpio.PinIO = &Pins[1]  // 353
	PL2  gpio.PinIO = &Pins[2]  // 357
	PL3  gpio.PinIO = &Pins[3]  // 358
	PL4  gpio.PinIO = &Pins[4]  // 359
	PL5  gpio.PinIO = &Pins[5]  // 360
	PL6  gpio.PinIO = &Pins[6]  // 361
	PL7  gpio.PinIO = &Pins[7]  // 362
	PL8  gpio.PinIO = &Pins[8]  // 363
	PL9  gpio.PinIO = &Pins[9]  // 364
	PL10 gpio.PinIO = &Pins[10] //
	PL11 gpio.PinIO = &Pins[11] //
	PL12 gpio.PinIO = &Pins[12] //
)

// PinIO implementation.

// Number implements gpio.PinIO
//
// It returns the GPIO pin number as represented by gpio sysfs.
func (p *Pin) Number() int {
	return 11*32 + int(p.offset)
}

// String implements gpio.PinIO
func (p *Pin) String() string {
	return fmt.Sprintf("%s(%d)", p.name, p.Number())
}

func (p *Pin) Function() string {
	switch f := p.function(); f {
	case in:
		return "In/" + p.Read().String() + "/" + p.Pull().String()
	case out:
		return "Out/" + p.Read().String()
	case alt1:
		if s := mapping[p.offset][0]; len(s) != 0 {
			return s
		}
		return "<Alt1>"
	case alt2:
		if s := mapping[p.offset][1]; len(s) != 0 {
			return s
		}
		return "<Alt2>"
	case alt3:
		if s := mapping[p.offset][2]; len(s) != 0 {
			return s
		}
		return "<Alt3>"
	case alt4:
		if s := mapping[p.offset][3]; len(s) != 0 {
			return s
		}
		return "<Alt4>"
	case alt5:
		if s := mapping[p.offset][4]; len(s) != 0 {
			return s
		}
		return "<Alt5>"
	case disabled:
		return "<Disabled>"
	default:
		return "<Internal error>"
	}
}

// In implemented gpio.PinIn.
//
// This requires opening a gpio sysfs file handle. The pin will be exported at
// /sys/class/gpio/gpio*/. Note that the pin will not be unexported at
// shutdown.
//
// Not all pins support edge detection Allwinner processors!
func (p *Pin) In(pull gpio.Pull, edge gpio.Edge) error {
	if gpioMemory == nil {
		return errors.New("subsystem not initialized")
	}
	if !p.setFunction(in) {
		return fmt.Errorf("failed to set pin %s as input", p.name)
	}
	if pull != gpio.PullNoChange {
		off := 7 + p.offset/16
		shift := 2 * (p.offset % 16)
		// Do it in a way that is concurrent safe.
		gpioMemory[off] &^= 3 << shift
		switch pull {
		case gpio.Down:
			gpioMemory[off] = 2 << shift
		case gpio.Up:
			gpioMemory[off] = 1 << shift
		default:
		}
	}
	if edge != gpio.None {
		// This is a race condition but this is fine; at worst PinByNumber() is
		// called twice but it is guaranteed to return the same value. p.edge is
		// never set to nil.
		if p.edge == nil {
			var err error
			if p.edge, err = sysfs.PinByNumber(p.Number()); err != nil {
				return err
			}
		}
		if err := p.edge.In(gpio.PullNoChange, edge); err != nil {
			return err
		}
	} else if p.edge != nil {
		if err := p.edge.In(gpio.PullNoChange, edge); err != nil {
			return err
		}
	}
	return nil
}

func (p *Pin) Read() gpio.Level {
	// Pn_DAT  n*0x24+0x10  Port n Data Register (n from 1(B) to 7(H))
	return gpio.Level(gpioMemory[4]&(1<<p.offset) != 0)
}

// WaitForEdge does edge detection and implements gpio.PinIn.
func (p *Pin) WaitForEdge(timeout time.Duration) bool {
	if p.edge != nil {
		return p.edge.WaitForEdge(timeout)
	}
	return false
}

func (p *Pin) Pull() gpio.Pull {
	if gpioMemory == nil {
		return gpio.PullNoChange
	}
	off := 7 + p.offset/16
	var v uint32
	// Pn_PULL  n*0x24+0x1C Port n Pull Register (n from 1(B) to 7(H))
	v = gpioMemory[off]
	switch (v >> (2 * (p.offset % 16))) & 3 {
	case 0:
		return gpio.Float
	case 1:
		return gpio.Up
	case 2:
		return gpio.Down
	default:
		// Confused.
		return gpio.PullNoChange
	}
}

func (p *Pin) Out(l gpio.Level) error {
	if gpioMemory == nil {
		return errors.New("subsystem not initialized")
	}
	if !p.setFunction(out) {
		return fmt.Errorf("failed to set pin %s as output", p.name)
	}
	// TODO(maruel): Set the value *before* changing the pin to be an output, so
	// there is no glitch.
	bit := uint32(1 << p.offset)
	if l {
		gpioMemory[4] |= bit
	} else {
		gpioMemory[4] &^= bit
	}
	return nil
}

func (p *Pin) PWM(duty int) error {
	return errors.New("pwm is not supported")
}

//

// function returns the current GPIO pin function.
func (p *Pin) function() function {
	if gpioMemory == nil {
		return disabled
	}
	off := p.offset / 8
	shift := 4 * (p.offset % 8)
	// Pn_CFGx n*0x24+0x0x  Port n Configure Register x (n from 1(B) to 7(H))
	return function((gpioMemory[off] >> shift) & 7)
}

// setFunction changes the GPIO pin function.
//
// Returns false if the pin was in AltN. Only accepts in and out
func (p *Pin) setFunction(f function) bool {
	if f != in && f != out {
		return false
	}
	// Interrupt based edge triggering is Alt5 but this is only supported on some
	// pins.
	// TODO(maruel): This check should use a whitelist of pins.
	if actual := p.function(); actual != in && actual != out && actual != disabled && actual != alt5 {
		// Pin is in special mode.
		return false
	}
	off := p.offset / 8
	shift := 4 * (p.offset % 8)
	mask := uint32(disabled) << shift
	v := (uint32(f) << shift) ^ mask
	// First disable, then setup. This is concurrent safe.
	gpioMemory[off] |= mask
	gpioMemory[off] &^= v
	if p.function() != f {
		panic(f)
	}
	return true
}

//

// function specifies the active functionality of a pin. The alternative
// function is GPIO pin dependent.
type function uint8

// Page 23~24
// Each pin can have one of 7 functions.
const (
	in       function = 0
	out      function = 1
	alt1     function = 2
	alt2     function = 3
	alt3     function = 4
	alt4     function = 5
	alt5     function = 6
	disabled function = 7
)

// Page 410 GPIO PL.
var gpioMemory []uint32

var _ gpio.PinIO = &Pin{}

// See ../allwinner/allwinner.go for details.
// TODO(maruel): Figure out what the S_ prefix means.
var mapping = [13][5]string{
	{"S_RSB_SCK", "S_I2C_SCL", "", "", "S_PL_EINT0"}, // PL0
	{"S_RSB_SDA", "S_I2C_SDA", "", "", "S_PL_EINT1"}, // PL1
	{"S_UART_TX", "", "", "", "S_PL_EINT2"},          // PL2
	{"S_UART_RX", "", "", "", "S_PL_EINT3"},          // PL3
	{"S_JTAG_MS", "", "", "", "S_PL_EINT4"},          // PL4
	{"S_JTAG_CK", "", "", "", "S_PL_EINT5"},          // PL5
	{"S_JTAG_DO", "", "", "", "S_PL_EINT6"},          // PL6
	{"S_JTAG_DI", "", "", "", "S_PL_EINT7"},          // PL7
	{"S_I2C_CSK", "", "", "", "S_PL_EINT8"},          // PL8
	{"S_I2C_SDA", "", "", "", "S_PL_EINT9"},          // PL9
	{"S_PWM", "", "", "", "S_PL_EINT10"},             // PL10
	{"S_CIR_RX", "", "", "", "S_PL_EINT11"},          // PL11
	{"", "", "", "", "S_PL_EINT12"},                  // PL12
}

// getBaseAddress queries the virtual file system to retrieve the base address
// of the GPIO registers for GPIO pins in group PL.
//
// Defaults to 0x01F02C00 as per datasheet if could query the file system.
func getBaseAddress() uint64 {
	base := uint64(0x01F02C00)
	link, err := os.Readlink("/sys/bus/platform/drivers/sun50i-r-pinctrl/driver")
	if err != nil {
		return base
	}
	parts := strings.SplitN(path.Base(link), ".", 2)
	if len(parts) != 2 {
		return base
	}
	base2, err := strconv.ParseUint(parts[0], 16, 64)
	if err != nil {
		return base
	}
	return base2
}

// driver implements pio.Driver.
type driver struct {
}

func (d *driver) String() string {
	return "allwinner_pl"
}

func (d *driver) Type() pio.Type {
	return pio.Processor
}

func (d *driver) Prerequisites() []string {
	return []string{"allwinner"}
}

func (d *driver) Init() (bool, error) {
	if !internal.IsAllwinner() {
		// BUG(maruel): Fix detection, need to specifically look for A64!
		return false, errors.New("A64 CPU not detected")
	}
	mem, err := gpiomem.OpenMem(getBaseAddress())
	if err != nil {
		return true, err
	}
	gpioMemory = mem.Uint32
	for i := range Pins {
		p := &Pins[i]
		if err := gpio.Register(p); err != nil {
			return true, err
		}
		if f := p.Function(); f[:2] != "In" && f[:3] != "Out" {
			gpio.MapFunction(f, p)
		}
	}
	return true, nil
}

var _ pio.Driver = &driver{}
var _ gpio.PinIn = &Pin{}
var _ gpio.PinOut = &Pin{}
var _ gpio.PinIO = &Pin{}
