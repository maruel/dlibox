// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package bcm283x

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/maruel/dlibox/go/pio/host"
	"github.com/maruel/dlibox/go/pio/host/internal/gpiomem"
	"github.com/maruel/dlibox/go/pio/host/internal/sysfs"
	"github.com/maruel/dlibox/go/pio/host/ir"
)

// Functional is pins.Functional on this CPU.
var Functional = map[string]host.Pin{
	"GPCLK0":    host.INVALID,
	"GPCLK1":    host.INVALID,
	"GPCLK2":    host.INVALID,
	"I2C_SCL0":  host.INVALID,
	"I2C_SDA0":  host.INVALID,
	"I2C_SCL1":  host.INVALID,
	"I2C_SDA1":  host.INVALID,
	"IR_IN":     host.INVALID,
	"IR_OUT":    host.INVALID,
	"PCM_CLK":   host.INVALID,
	"PCM_FS":    host.INVALID,
	"PCM_DIN":   host.INVALID,
	"PCM_DOUT":  host.INVALID,
	"PWM0_OUT":  host.INVALID,
	"PWM1_OUT":  host.INVALID,
	"SPI0_CE0":  host.INVALID,
	"SPI0_CE1":  host.INVALID,
	"SPI0_CLK":  host.INVALID,
	"SPI0_MISO": host.INVALID,
	"SPI0_MOSI": host.INVALID,
	"SPI1_CE0":  host.INVALID,
	"SPI1_CE1":  host.INVALID,
	"SPI1_CE2":  host.INVALID,
	"SPI1_CLK":  host.INVALID,
	"SPI1_MISO": host.INVALID,
	"SPI1_MOSI": host.INVALID,
	"SPI2_MISO": host.INVALID,
	"SPI2_MOSI": host.INVALID,
	"SPI2_CLK":  host.INVALID,
	"SPI2_CE0":  host.INVALID,
	"SPI2_CE1":  host.INVALID,
	"SPI2_CE2":  host.INVALID,
	"UART_RXD0": host.INVALID,
	"UART_CTS0": host.INVALID,
	"UART_CTS1": host.INVALID,
	"UART_RTS0": host.INVALID,
	"UART_RTS1": host.INVALID,
	"UART_TXD0": host.INVALID,
	"UART_RXD1": host.INVALID,
	"UART_TXD1": host.INVALID,
}

// Pin is a GPIO number (GPIOnn) on BCM238(5|6|7).
//
// If you search for pin per their position on the P1 header, look at ../rpi
// package.
//
// Pin implements host.PinIO.
type Pin struct {
	number      int
	name        string
	defaultPull host.Pull
	edge        *sysfs.Pin // Mutable, set once, then never set back to nil
}

var (
	GPIO0  *Pin // I2C_SDA0
	GPIO1  *Pin // I2C_SCL0
	GPIO2  *Pin // I2C_SDA1
	GPIO3  *Pin // I2C_SCL1
	GPIO4  *Pin // GPCLK0
	GPIO5  *Pin // GPCLK1
	GPIO6  *Pin // GPCLK2
	GPIO7  *Pin // SPI0_CE1
	GPIO8  *Pin // SPI0_CE0
	GPIO9  *Pin // SPI0_MISO
	GPIO10 *Pin // SPI0_MOSI
	GPIO11 *Pin // SPI0_CLK
	GPIO12 *Pin // PWM0_OUT
	GPIO13 *Pin // PWM1_OUT
	GPIO14 *Pin // UART_TXD0, UART_TXD1
	GPIO15 *Pin // UART_RXD0, UART_RXD1
	GPIO16 *Pin // UART_CTS0, SPI1_CE2, UART_CTS1
	GPIO17 *Pin // UART_RTS0, SPI1_CE1, UART_RTS1
	GPIO18 *Pin // PCM_CLK, SPI1_CE0, PWM0_OUT
	GPIO19 *Pin // PCM_FS, SPI1_MISO, PWM1_OUT
	GPIO20 *Pin // PCM_DIN, SPI1_MOSI, GPCLK0
	GPIO21 *Pin // PCM_DOUT, SPI1_CLK, GPCLK1
	GPIO22 *Pin //
	GPIO23 *Pin //
	GPIO24 *Pin //
	GPIO25 *Pin //
	GPIO26 *Pin //
	GPIO27 *Pin //
	GPIO28 *Pin // I2C_SDA0, PCM_CLK
	GPIO29 *Pin // I2C_SCL0, PCM_FS
	GPIO30 *Pin // PCM_DIN, UART_CTS0, UARTS_CTS1
	GPIO31 *Pin // PCM_DOUT, UART_RTS0, UARTS_RTS1
	GPIO32 *Pin // GPCLK0, UART_TXD0, UARTS_TXD1
	GPIO33 *Pin // UART_RXD0, UARTS_RXD1
	GPIO34 *Pin // GPCLK0
	GPIO35 *Pin // SPI0_CE1
	GPIO36 *Pin // SPI0_CE0, UART_TXD0
	GPIO37 *Pin // SPI0_MISO, UART_RXD0
	GPIO38 *Pin // SPI0_MOSI, UART_RTS0
	GPIO39 *Pin // SPI0_CLK, UART_CTS0
	GPIO40 *Pin // PWM0_OUT, SPI2_MISO, UART_TXD1
	GPIO41 *Pin // PWM1_OUT, SPI2_MOSI, UART_RXD1
	GPIO42 *Pin // GPCLK1, SPI2_CLK, UART_RTS1
	GPIO43 *Pin // GPCLK2, SPI2_CE0, UART_CTS1
	GPIO44 *Pin // GPCLK1, I2C_SDA0, I2C_SDA1, SPI2_CE1
	GPIO45 *Pin // PWM1_OUT, I2C_SCL0, I2C_SCL1, SPI2_CE2
	GPIO46 *Pin //
	GPIO47 *Pin // SDCard
	GPIO48 *Pin // SDCard
	GPIO49 *Pin // SDCard
	GPIO50 *Pin // SDCard
	GPIO51 *Pin // SDCard
	GPIO52 *Pin // SDCard
	GPIO53 *Pin // SDCard
)

// Special functions that can be assigned to a GPIO. The values are probed and
// set at runtime. Changing the value of the variables has no effect.
var (
	GPCLK0    host.Pin // GPIO4, GPIO20, GPIO32, GPIO34 (also named GPIO_GCLK)
	GPCLK1    host.Pin // GPIO5, GPIO21, GPIO42, GPIO44
	GPCLK2    host.Pin // GPIO6, GPIO43
	I2C_SCL0  host.Pin // GPIO1, GPIO29, GPIO45
	I2C_SDA0  host.Pin // GPIO0, GPIO28, GPIO44
	I2C_SCL1  host.Pin // GPIO3, GPIO45
	I2C_SDA1  host.Pin // GPIO2, GPIO44
	IR_IN     host.Pin // (any GPIO)
	IR_OUT    host.Pin // (any GPIO)
	PCM_CLK   host.Pin // GPIO18, GPIO28 (I2S)
	PCM_FS    host.Pin // GPIO19, GPIO29
	PCM_DIN   host.Pin // GPIO20, GPIO30
	PCM_DOUT  host.Pin // GPIO21, GPIO31
	PWM0_OUT  host.Pin // GPIO12, GPIO18, GPIO40
	PWM1_OUT  host.Pin // GPIO13, GPIO19, GPIO41, GPIO45
	SPI0_CE0  host.Pin // GPIO8,  GPIO36
	SPI0_CE1  host.Pin // GPIO7,  GPIO35
	SPI0_CLK  host.Pin // GPIO11, GPIO39
	SPI0_MISO host.Pin // GPIO9,  GPIO37
	SPI0_MOSI host.Pin // GPIO10, GPIO38
	SPI1_CE0  host.Pin // GPIO18
	SPI1_CE1  host.Pin // GPIO17
	SPI1_CE2  host.Pin // GPIO16
	SPI1_CLK  host.Pin // GPIO21
	SPI1_MISO host.Pin // GPIO19
	SPI1_MOSI host.Pin // GPIO20
	SPI2_MISO host.Pin // GPIO40
	SPI2_MOSI host.Pin // GPIO41
	SPI2_CLK  host.Pin // GPIO42
	SPI2_CE0  host.Pin // GPIO43
	SPI2_CE1  host.Pin // GPIO44
	SPI2_CE2  host.Pin // GPIO45
	UART_RXD0 host.Pin // GPIO15, GPIO33, GPIO37
	UART_CTS0 host.Pin // GPIO16, GPIO30, GPIO39
	UART_CTS1 host.Pin // GPIO16, GPIO30
	UART_RTS0 host.Pin // GPIO17, GPIO31, GPIO38
	UART_RTS1 host.Pin // GPIO17, GPIO31
	UART_TXD0 host.Pin // GPIO14, GPIO32, GPIO36
	UART_RXD1 host.Pin // GPIO15, GPIO33, GPIO41
	UART_TXD1 host.Pin // GPIO14, GPIO32, GPIO40
)

// PinIO implementation.

// Number implements host.Pin
func (p *Pin) Number() int {
	return p.number
}

// String implements host.Pin
func (p *Pin) String() string {
	return p.name
}

// Function implements host.Pin
func (p *Pin) Function() string {
	switch f := p.function(); f {
	case in:
		return "In/" + p.Read().String()
	case out:
		return "Out/" + p.Read().String()
	case alt0:
		if s := mapping[p.number][0]; len(s) != 0 {
			return s
		}
		return "<Alt0>"
	case alt1:
		if s := mapping[p.number][1]; len(s) != 0 {
			return s
		}
		return "<Alt1>"
	case alt2:
		if s := mapping[p.number][2]; len(s) != 0 {
			return s
		}
		return "<Alt2>"
	case alt3:
		if s := mapping[p.number][3]; len(s) != 0 {
			return s
		}
		return "<Alt3>"
	case alt4:
		if s := mapping[p.number][4]; len(s) != 0 {
			return s
		}
		return "<Alt4>"
	case alt5:
		if s := mapping[p.number][5]; len(s) != 0 {
			return s
		}
		return "<Alt5>"
	default:
		return "<Unknown>"
	}
}

// In setups a pin as an input and implements host.PinIn.
//
// Specifying a value for pull other than host.PullNoChange causes this
// function to be slightly slower (about 1ms).
//
// Will fail if requesting to change a pin that is set to special functionality.
func (p *Pin) In(pull host.Pull) error {
	if gpioMemory32 == nil {
		return errors.New("subsystem not initialized")
	}
	if !p.setFunction(in) {
		return errors.New("failed to change pin mode")
	}
	if pull != host.PullNoChange {
		// Changing pull resistor requires a specific dance as described at
		// https://www.raspberrypi.org/wp-content/uploads/2012/02/BCM2835-ARM-Peripherals.pdf
		// page 101.

		// Set Pull
		// 0x94    RW   GPIO Pin Pull-up/down Enable (00=Float, 01=Down, 10=Up)
		gpioMemory32[37] = uint32(pull)

		// Datasheet states caller needs to sleep 150 cycles.
		sleep150cycles()
		// 0x98    RW   GPIO Pin Pull-up/down Enable Clock 0 (GPIO0-31)
		// 0x9C    RW   GPIO Pin Pull-up/down Enable Clock 1 (GPIO32-53)
		offset := 38 + p.number/32
		gpioMemory32[offset] = 1 << uint(p.number%32)

		sleep150cycles()
		gpioMemory32[37] = 0
		gpioMemory32[offset] = 0
	}
	return nil
}

// Read return the current pin level and implements host.PinIn.
//
// This function is very fast. It works even if the pin is set as output.
func (p *Pin) Read() host.Level {
	// 0x34    R    GPIO Pin Level 0 (GPIO0-31)
	// 0x38    R    GPIO Pin Level 1 (GPIO32-53)
	return host.Level((gpioMemory32[13+p.number/32] & (1 << uint(p.number&31))) != 0)
}

// Edges creates a edge detection loop and implements host.PinIn.
//
// This requires opening a gpio sysfs file handle. Make sure the user is member
// of group 'gpio'. The pin will be exported at /sys/class/gpio/gpio*/. Note
// that the pin will not be unexported at shutdown.
//
// For edge detection, the processor samples the input at its CPU clock rate
// and looks for '011' to rising and '100' for falling detection to avoid
// glitches. Because gpio sysfs is used, the latency is unpredictable.
func (p *Pin) Edges() (chan host.Level, error) {
	// This is a race condition but this is fine; at worst GetPin() is called
	// twice but it is guaranteed to return the same value. p.edge is never set
	// to nil.
	if p.edge == nil {
		var err error
		if p.edge, err = sysfs.GetPin(p.Number()); err != nil {
			return nil, err
		}
	}
	return p.edge.Edges()
}

// Out sets a pin as output and implements host.PinOut. The caller should
// immediately call SetLow() or SetHigh() afterward.
//
// Will fail if requesting to change a pin that is set to special functionality.
func (p *Pin) Out() error {
	if gpioMemory32 == nil {
		return errors.New("subsystem not initialized")
	}
	if !p.setFunction(out) {
		return errors.New("failed to change pin mode")
	}
	return nil
}

// Set sets a pin already set for output as host.High or host.Low and implements
// host.PinOut.
//
// This function is very fast.
func (p *Pin) Set(l host.Level) {
	// 0x1C    W    GPIO Pin Output Set 0 (GPIO0-31)
	// 0x20    W    GPIO Pin Output Set 1 (GPIO32-53)
	base := 7 + p.number/32
	if l == host.Low {
		// 0x28    W    GPIO Pin Output Clear 0 (GPIO0-31)
		// 0x2C    W    GPIO Pin Output Clear 1 (GPIO32-53)
		base += 3
	}
	gpioMemory32[base] = 1 << uint(p.number&31)
}

// Special functionality.

// DefaultPull returns the default pull for the function.
//
// The CPU doesn't return the current pull.
func (p *Pin) DefaultPull() host.Pull {
	return p.defaultPull
}

// Internal code.

// function returns the current GPIO pin function.
func (p *Pin) function() function {
	if gpioMemory32 == nil {
		return alt5
	}
	// 0x00    RW   GPIO Function Select 0 (GPIO0-9)
	// 0x04    RW   GPIO Function Select 1 (GPIO10-19)
	// 0x08    RW   GPIO Function Select 2 (GPIO20-29)
	// 0x0C    RW   GPIO Function Select 3 (GPIO30-39)
	// 0x10    RW   GPIO Function Select 4 (GPIO40-49)
	// 0x14    RW   GPIO Function Select 5 (GPIO50-53)
	return function((gpioMemory32[p.number/10] >> uint((p.number%10)*3)) & 7)
}

// setFunction changes the GPIO pin function.
//
// Returns false if the pin was in AltN. Only accepts in and out
func (p *Pin) setFunction(f function) bool {
	if f != in && f != out {
		return false
	}
	if f == in && p.edge != nil {
		p.edge.DisableEdge()
	}
	if actual := p.function(); actual != in && actual != out {
		return false
	}
	// 0x00    RW   GPIO Function Select 0 (GPIO0-9)
	// 0x04    RW   GPIO Function Select 1 (GPIO10-19)
	// 0x08    RW   GPIO Function Select 2 (GPIO20-29)
	// 0x0C    RW   GPIO Function Select 3 (GPIO30-39)
	// 0x10    RW   GPIO Function Select 4 (GPIO40-49)
	// 0x14    RW   GPIO Function Select 5 (GPIO50-53)
	off := p.number / 10
	shift := uint(p.number%10) * 3
	gpioMemory32[off] = (gpioMemory32[off] &^ (7 << shift)) | (uint32(f) << shift)
	return true
}

//

// function specifies the active functionality of a pin. The alternative
// function is GPIO pin dependent.
type function uint8

// Each pin can have one of 7 functions.
const (
	in   function = 0
	out  function = 1
	alt0 function = 4
	alt1 function = 5
	alt2 function = 6
	alt3 function = 7
	alt4 function = 3
	alt5 function = 2
)

// Mapping as
// https://www.raspberrypi.org/wp-content/uploads/2012/02/BCM2835-ARM-Peripherals.pdf
// pages 90-91.
// Offset  Mode Value
// 0x00    RW   GPIO Function Select 0 (GPIO0-9)
// 0x04    RW   GPIO Function Select 1 (GPIO10-19)
// 0x08    RW   GPIO Function Select 2 (GPIO20-29)
// 0x0C    RW   GPIO Function Select 3 (GPIO30-39)
// 0x10    RW   GPIO Function Select 4 (GPIO40-49)
// 0x14    RW   GPIO Function Select 5 (GPIO50-53)
// 0x18    -    Reserved
// 0x1C    W    GPIO Pin Output Set 0 (GPIO0-31)
// 0x20    W    GPIO Pin Output Set 1 (GPIO32-53)
// 0x24    -    Reserved
// 0x28    W    GPIO Pin Output Clear 0 (GPIO0-31)
// 0x2C    W    GPIO Pin Output Clear 1 (GPIO32-53)
// 0x30    -    Reserved
// 0x34    R    GPIO Pin Level 0 (GPIO0-31)
// 0x38    R    GPIO Pin Level 1 (GPIO32-53)
// 0x3C    -    Reserved
// 0x40    RW   GPIO Pin Event Detect Status 0 (GPIO0-31)
// 0x44    RW   GPIO Pin Event Detect Status 1 (GPIO32-53)
// 0x48    -    Reserved
// 0x4C    RW   GPIO Pin Rising Edge Detect Enable 0 (GPIO0-31)
// 0x50    RW   GPIO Pin Rising Edge Detect Enable 1 (GPIO32-53)
// 0x54    -    Reserved
// 0x58    RW   GPIO Pin Falling Edge Detect Enable 0 (GPIO0-31)
// 0x5C    RW   GPIO Pin Falling Edge Detect Enable 1 (GPIO32-53)
// 0x60    -    Reserved
// 0x64    RW   GPIO Pin High Detect Enable 0 (GPIO0-31)
// 0x68    RW   GPIO Pin High Detect Enable 1 (GPIO32-53)
// 0x6C    -    Reserved
// 0x70    RW   GPIO Pin Low Detect Enable 0 (GPIO0-31)
// 0x74    RW   GPIO Pin Low Detect Enable 1 (GPIO32-53)
// 0x78    -    Reserved
// 0x7C    RW   GPIO Pin Async Rising Edge Detect 0 (GPIO0-31)
// 0x80    RW   GPIO Pin Async Rising Edge Detect 1 (GPIO32-53)
// 0x84    -    Reserved
// 0x88    RW   GPIO Pin Async Falling Edge Detect 0 (GPIO0-31)
// 0x8C    RW   GPIO Pin Async Falling Edge Detect 1 (GPIO32-53)
// 0x90    -    Reserved
// 0x94    RW   GPIO Pin Pull-up/down Enable (00=Float, 01=Down, 10=Up)
// 0x98    RW   GPIO Pin Pull-up/down Enable Clock 0 (GPIO0-31)
// 0x9C    RW   GPIO Pin Pull-up/down Enable Clock 1 (GPIO32-53)
// 0xA0    -    Reserved
// 0xB0    -    Test (byte)
var gpioMemory32 []uint32

// Changing pull resistor require a 150 cycles sleep.
//
// Do not inline so the temporary value is not optimized out.
//
//go:noinline
func sleep150cycles() uint32 {
	// Do not call into any kernel function, since this causes a high chance of
	// being preempted.
	// Abuse the fact that gpioMemory32 is uncached memory.
	// TODO(maruel): No idea if this is too much or enough.
	var out uint32
	for i := 0; i < 150; i++ {
		out += gpioMemory32[i]
	}
	return out
}

func setIfAlt0(p *Pin, special *host.Pin) {
	if p.function() == alt0 {
		if (*special).String() != "INVALID" {
			//fmt.Printf("%s and %s have same functionality\n", p, *special)
		}
		*special = p
	}
}

func setIfAlt(p *Pin, special0 *host.Pin, special1 *host.Pin, special2 *host.Pin, special3 *host.Pin, special4 *host.Pin, special5 *host.Pin) {
	switch p.function() {
	case alt0:
		if special0 != nil {
			if (*special0).String() != "INVALID" {
				//fmt.Printf("%s and %s have same functionality\n", p, *special0)
			}
			*special0 = p
		}
	case alt1:
		if special1 != nil {
			if (*special1).String() != "INVALID" {
				//fmt.Printf("%s and %s have same functionality\n", p, *special1)
			}
			*special1 = p
		}
	case alt2:
		if special2 != nil {
			if (*special2).String() != "INVALID" {
				//log.Printf("%s and %s have same functionality\n", p, *special2)
			}
			*special2 = p
		}
	case alt3:
		if special3 != nil {
			if (*special3).String() != "INVALID" {
				//log.Printf("%s and %s have same functionality\n", p, *special3)
			}
			*special3 = p
		}
	case alt4:
		if special4 != nil {
			if (*special4).String() != "INVALID" {
				//log.Printf("%s and %s have same functionality\n", p, *special4)
			}
			*special4 = p
		}
	case alt5:
		if special5 != nil {
			if (*special5).String() != "INVALID" {
				//log.Printf("%s and %s have same functionality\n", p, *special5)
			}
			*special5 = p
		}
	}
}

// Pins is all the supported pins. The bcm 283x exports continuously numbered
// pins.
var Pins = [54]Pin{
	{number: 0, name: "GPIO0", defaultPull: host.Up},
	{number: 1, name: "GPIO1", defaultPull: host.Up},
	{number: 2, name: "GPIO2", defaultPull: host.Up},
	{number: 3, name: "GPIO3", defaultPull: host.Up},
	{number: 4, name: "GPIO4", defaultPull: host.Up},
	{number: 5, name: "GPIO5", defaultPull: host.Up},
	{number: 6, name: "GPIO6", defaultPull: host.Up},
	{number: 7, name: "GPIO7", defaultPull: host.Up},
	{number: 8, name: "GPIO8", defaultPull: host.Up},
	{number: 9, name: "GPIO9", defaultPull: host.Down},
	{number: 10, name: "GPIO10", defaultPull: host.Down},
	{number: 11, name: "GPIO11", defaultPull: host.Down},
	{number: 12, name: "GPIO12", defaultPull: host.Down},
	{number: 13, name: "GPIO13", defaultPull: host.Down},
	{number: 14, name: "GPIO14", defaultPull: host.Down},
	{number: 15, name: "GPIO15", defaultPull: host.Down},
	{number: 16, name: "GPIO16", defaultPull: host.Down},
	{number: 17, name: "GPIO17", defaultPull: host.Down},
	{number: 18, name: "GPIO18", defaultPull: host.Down},
	{number: 19, name: "GPIO19", defaultPull: host.Down},
	{number: 20, name: "GPIO20", defaultPull: host.Down},
	{number: 21, name: "GPIO21", defaultPull: host.Down},
	{number: 22, name: "GPIO22", defaultPull: host.Down},
	{number: 23, name: "GPIO23", defaultPull: host.Down},
	{number: 24, name: "GPIO24", defaultPull: host.Down},
	{number: 25, name: "GPIO25", defaultPull: host.Down},
	{number: 26, name: "GPIO26", defaultPull: host.Down},
	{number: 27, name: "GPIO27", defaultPull: host.Down},
	{number: 28, name: "GPIO28", defaultPull: host.Float},
	{number: 29, name: "GPIO29", defaultPull: host.Float},
	{number: 30, name: "GPIO30", defaultPull: host.Down},
	{number: 31, name: "GPIO31", defaultPull: host.Down},
	{number: 32, name: "GPIO32", defaultPull: host.Down},
	{number: 33, name: "GPIO33", defaultPull: host.Down},
	{number: 34, name: "GPIO34", defaultPull: host.Up},
	{number: 35, name: "GPIO35", defaultPull: host.Up},
	{number: 36, name: "GPIO36", defaultPull: host.Up},
	{number: 37, name: "GPIO37", defaultPull: host.Down},
	{number: 38, name: "GPIO38", defaultPull: host.Down},
	{number: 39, name: "GPIO39", defaultPull: host.Down},
	{number: 40, name: "GPIO40", defaultPull: host.Down},
	{number: 41, name: "GPIO41", defaultPull: host.Down},
	{number: 42, name: "GPIO42", defaultPull: host.Down},
	{number: 43, name: "GPIO43", defaultPull: host.Down},
	{number: 44, name: "GPIO44", defaultPull: host.Float},
	{number: 45, name: "GPIO45", defaultPull: host.Float},
	{number: 46, name: "GPIO46", defaultPull: host.Up},
	{number: 47, name: "GPIO47", defaultPull: host.Up},
	{number: 48, name: "GPIO48", defaultPull: host.Up},
	{number: 49, name: "GPIO49", defaultPull: host.Up},
	{number: 50, name: "GPIO50", defaultPull: host.Up},
	{number: 51, name: "GPIO51", defaultPull: host.Up},
	{number: 52, name: "GPIO52", defaultPull: host.Up},
	{number: 53, name: "GPIO53", defaultPull: host.Up},
}

// This excludes the functions in and out.
var mapping = [54][6]string{
	{"I2C_SDA0"}, // 0
	{"I2C_SCL0"},
	{"I2C_SDA1"},
	{"I2C_SCL1"},
	{"GPCLK0"},
	{"GPCLK1"},
	{"GPCLK2"},
	{"SPI0_CE1"},
	{"SPI0_CE0"},
	{"SPI0_MISO"},
	{"SPI0_MOSI"}, // 10
	{"SPI0_CLK"},
	{"PWM0_OUT"},
	{"PWM1_OUT"},
	{"UART_TXD0", "", "", "", "", "UART_TXD1"},
	{"UART_RXD0", "", "", "", "", "UART_RXD1"},
	{"", "", "", "UART_CTS0", "SPI1_CE2", "UART_CTS1"},
	{"", "", "", "UART_RTS0", "SPI1_CE1", "UART_RTS1"},
	{"PCM_CLK", "", "", "", "SPI1_CE0", "PWM0_OUT"},
	{"PCM_FS", "", "", "", "SPI1_MISO", "PWM1_OUT"},
	{"PCM_DIN", "", "", "", "SPI1_MOSI", "GPCLK0"}, // 20
	{"PCM_DOUT", "", "", "", "SPI1_CLK", "GPCLK1"},
	{},
	{},
	{},
	{},
	{},
	{},
	{"I2C_SDA0", "", "PCM_CLK", "", "", ""},
	{"I2C_SCL0", "", "PCM_FS", "", "", ""},
	{"", "", "PCM_DIN", "UART_CTS0", "", "UART_CTS1"}, // 30
	{"", "", "PCM_DOUT", "UART_RTS0", "", "UART_RTS"},
	{"GPCLK0", "", "", "UART_TXD0", "", "UART_TXD1"},
	{"", "", "", "UART_RXD0", "", "UART_RXD1"},
	{"GPCLK0"},
	{"SPI0_CE1"},
	{"SPI0_CE0", "", "UART_TXD0", "", "", ""},
	{"SPI0_MISO", "", "UART_RXD0", "", "", ""},
	{"SPI0_MOSI", "", "UART_RTS0", "", "", ""},
	{"SPI0_CLK", "", "UART_CTS0", "", "", ""},
	{"PWM0_OUT", "", "", "", "SPI2_MISO", "UART_TXD1"}, // 40
	{"PWM1_OUT", "", "", "", "SPI2_MOSI", "UART_RXD1"},
	{"GPCLK1", "", "", "", "SPI2_CLK", "UART_RTS1"},
	{"GPCLK2", "", "", "", "SPI2_CE0", "UART_CTS1"},
	{"GPCLK1", "I2C_SDA0", "I2C_SDA1", "", "SPI2_CE1", ""},
	{"PWM1_OUT", "I2C_SCL0", "I2C_SCL1", "", "SPI2_CE2", ""},
}

// getBaseAddress queries the virtual file system to retrieve the base address
// of the GPIO registers.
//
// Defaults to 0x3F200000 as per datasheet if could query the file system.
func getBaseAddress() uint64 {
	items, _ := ioutil.ReadDir("/sys/bus/platform/drivers/pinctrl-bcm2835/")
	for _, item := range items {
		if item.Mode()&os.ModeSymlink != 0 {
			parts := strings.SplitN(path.Base(item.Name()), ".", 2)
			if len(parts) != 2 {
				continue
			}
			base, err := strconv.ParseUint(parts[0], 16, 64)
			if err != nil {
				continue
			}
			return base
		}
	}
	return 0x3F200000
}

func init() {
	GPIO0 = &Pins[0]
	GPIO1 = &Pins[1]
	GPIO2 = &Pins[2]
	GPIO3 = &Pins[3]
	GPIO4 = &Pins[4]
	GPIO5 = &Pins[5]
	GPIO6 = &Pins[6]
	GPIO7 = &Pins[7]
	GPIO8 = &Pins[8]
	GPIO9 = &Pins[9]
	GPIO10 = &Pins[10]
	GPIO11 = &Pins[11]
	GPIO12 = &Pins[12]
	GPIO13 = &Pins[13]
	GPIO14 = &Pins[14]
	GPIO15 = &Pins[15]
	GPIO16 = &Pins[16]
	GPIO17 = &Pins[17]
	GPIO18 = &Pins[18]
	GPIO19 = &Pins[19]
	GPIO20 = &Pins[20]
	GPIO21 = &Pins[21]
	GPIO22 = &Pins[22]
	GPIO23 = &Pins[23]
	GPIO24 = &Pins[24]
	GPIO25 = &Pins[25]
	GPIO26 = &Pins[26]
	GPIO27 = &Pins[27]
	GPIO28 = &Pins[28]
	GPIO29 = &Pins[29]
	GPIO30 = &Pins[30]
	GPIO31 = &Pins[31]
	GPIO32 = &Pins[32]
	GPIO33 = &Pins[33]
	GPIO34 = &Pins[34]
	GPIO35 = &Pins[35]
	GPIO36 = &Pins[36]
	GPIO37 = &Pins[37]
	GPIO38 = &Pins[38]
	GPIO39 = &Pins[39]
	GPIO40 = &Pins[40]
	GPIO41 = &Pins[41]
	GPIO42 = &Pins[42]
	GPIO43 = &Pins[43]
	GPIO44 = &Pins[44]
	GPIO45 = &Pins[45]
	GPIO46 = &Pins[46]
	GPIO47 = &Pins[47]
	GPIO48 = &Pins[48]
	GPIO49 = &Pins[49]
	GPIO50 = &Pins[50]
	GPIO51 = &Pins[51]
	GPIO52 = &Pins[52]
	GPIO53 = &Pins[53]

	GPCLK0 = host.INVALID
	GPCLK1 = host.INVALID
	GPCLK2 = host.INVALID
	I2C_SCL0 = host.INVALID
	I2C_SDA0 = host.INVALID
	I2C_SCL1 = host.INVALID
	I2C_SDA1 = host.INVALID
	IR_IN = host.INVALID
	IR_OUT = host.INVALID
	PCM_CLK = host.INVALID
	PCM_FS = host.INVALID
	PCM_DIN = host.INVALID
	PCM_DOUT = host.INVALID
	PWM0_OUT = host.INVALID
	PWM1_OUT = host.INVALID
	SPI0_CE0 = host.INVALID
	SPI0_CE1 = host.INVALID
	SPI0_CLK = host.INVALID
	SPI0_MISO = host.INVALID
	SPI0_MOSI = host.INVALID
	SPI1_CE0 = host.INVALID
	SPI1_CE1 = host.INVALID
	SPI1_CE2 = host.INVALID
	SPI1_CLK = host.INVALID
	SPI1_MISO = host.INVALID
	SPI1_MOSI = host.INVALID
	SPI2_MISO = host.INVALID
	SPI2_MOSI = host.INVALID
	SPI2_CLK = host.INVALID
	SPI2_CE0 = host.INVALID
	SPI2_CE1 = host.INVALID
	SPI2_CE2 = host.INVALID
	UART_RXD0 = host.INVALID
	UART_CTS0 = host.INVALID
	UART_CTS1 = host.INVALID
	UART_RTS0 = host.INVALID
	UART_RTS1 = host.INVALID
	UART_TXD0 = host.INVALID
	UART_RXD1 = host.INVALID
	UART_TXD1 = host.INVALID
}

func Init() error {
	mem, err := gpiomem.OpenGPIO()
	if err != nil {
		// Try without /dev/gpiomem.
		mem, err = gpiomem.OpenMem(getBaseAddress())
		if err != nil {
			return err
		}
	}
	gpioMemory32 = mem.Uint32

	// https://www.raspberrypi.org/wp-content/uploads/2012/02/BCM2835-ARM-Peripherals.pdf
	// Page 102.
	setIfAlt0(GPIO0, &I2C_SDA0)
	setIfAlt0(GPIO1, &I2C_SCL0)
	setIfAlt0(GPIO2, &I2C_SDA1)
	setIfAlt0(GPIO3, &I2C_SCL1)
	setIfAlt0(GPIO4, &GPCLK0)
	setIfAlt0(GPIO5, &GPCLK1)
	setIfAlt0(GPIO6, &GPCLK2)
	setIfAlt0(GPIO7, &SPI0_CE1)
	setIfAlt0(GPIO8, &SPI0_CE0)
	setIfAlt0(GPIO9, &SPI0_MISO)
	setIfAlt0(GPIO10, &SPI0_MOSI)
	setIfAlt0(GPIO11, &SPI0_CLK)
	setIfAlt0(GPIO12, &PWM0_OUT)
	setIfAlt0(GPIO13, &PWM1_OUT)
	setIfAlt(GPIO14, &UART_TXD0, nil, nil, nil, nil, &UART_TXD1)
	setIfAlt(GPIO15, &UART_RXD0, nil, nil, nil, nil, &UART_RXD1)
	setIfAlt(GPIO16, nil, nil, nil, &UART_CTS0, &SPI1_CE2, &UART_CTS1)
	setIfAlt(GPIO17, nil, nil, nil, &UART_RTS0, &SPI1_CE1, &UART_RTS1)
	setIfAlt(GPIO18, &PCM_CLK, nil, nil, nil, &SPI1_CE0, &PWM0_OUT)
	setIfAlt(GPIO19, &PCM_FS, nil, nil, nil, &SPI1_MISO, &PWM1_OUT)
	setIfAlt(GPIO20, &PCM_DIN, nil, nil, nil, &SPI1_MOSI, &GPCLK0)
	setIfAlt(GPIO21, &PCM_DOUT, nil, nil, nil, &SPI1_CLK, &GPCLK1)
	// GPIO22-GPIO27 do not have interesting alternate function.
	setIfAlt(GPIO28, &I2C_SDA0, nil, &PCM_CLK, nil, nil, nil)
	setIfAlt(GPIO29, &I2C_SCL0, nil, &PCM_FS, nil, nil, nil)
	setIfAlt(GPIO30, nil, nil, &PCM_DIN, &UART_CTS0, nil, &UART_CTS1)
	setIfAlt(GPIO31, nil, nil, &PCM_DOUT, &UART_RTS0, nil, &UART_RTS1)
	setIfAlt(GPIO32, &GPCLK0, nil, nil, &UART_TXD0, nil, &UART_TXD1)
	setIfAlt(GPIO33, nil, nil, nil, &UART_RXD0, nil, &UART_RXD1)
	setIfAlt0(GPIO34, &GPCLK0)
	setIfAlt0(GPIO35, &SPI0_CE1)
	setIfAlt(GPIO36, &SPI0_CE0, nil, &UART_TXD0, nil, nil, nil)
	setIfAlt(GPIO37, &SPI0_MISO, nil, &UART_RXD0, nil, nil, nil)
	setIfAlt(GPIO38, &SPI0_MOSI, nil, &UART_RTS0, nil, nil, nil)
	setIfAlt(GPIO39, &SPI0_CLK, nil, &UART_CTS0, nil, nil, nil)
	setIfAlt(GPIO40, &PWM0_OUT, nil, nil, nil, &SPI2_MISO, &UART_TXD1)
	setIfAlt(GPIO41, &PWM1_OUT, nil, nil, nil, &SPI2_MOSI, &UART_RXD1)
	setIfAlt(GPIO42, &GPCLK1, nil, nil, nil, &SPI2_CLK, &UART_RTS1)
	setIfAlt(GPIO43, &GPCLK2, nil, nil, nil, &SPI2_CE0, &UART_CTS1)
	setIfAlt(GPIO44, &GPCLK1, &I2C_SDA0, &I2C_SDA1, nil, &SPI2_CE1, nil)
	setIfAlt(GPIO45, &PWM1_OUT, &I2C_SCL0, &I2C_SCL1, nil, &SPI2_CE2, nil)
	// GPIO46 doesn't have interesting alternate function.
	// GPIO47-GPIO53 are connected to the SDCard.

	// TODO(maruel): Remove all the functional variables.
	for i := range Pins {
		if i == 45 {
			break
		}
		if f := Pins[i].Function(); len(f) < 3 || (f[:2] != "In" && f[:3] != "Out") {
			Functional[f] = &Pins[i]
		}
	}

	in, out := ir.Pins()
	if in != -1 {
		IR_IN = &Pins[in]
	}
	if out != -1 {
		IR_OUT = &Pins[out]
	}
	return nil
}
