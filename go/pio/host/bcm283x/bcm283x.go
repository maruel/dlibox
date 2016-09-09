// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package bcm283x

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/maruel/dlibox/go/pio/host"
	"github.com/maruel/dlibox/go/pio/host/cpu"
	"github.com/maruel/dlibox/go/pio/host/internal/gpiomem"
	"github.com/maruel/dlibox/go/pio/host/ir"
	"github.com/maruel/dlibox/go/pio/host/pins"
)

// Function specifies the active functionality of a pin. The alternative
// function is GPIO pin dependent.
type Function uint8

// Each pin can have one of 7 functions.
const (
	In   Function = 0
	Out  Function = 1
	Alt0 Function = 4
	Alt1 Function = 5
	Alt2 Function = 6
	Alt3 Function = 7
	Alt4 Function = 3
	Alt5 Function = 2
)

const functionName = "InOutAlt5Alt4Alt0Alt1Alt2Alt3"

var functionIndex = [...]uint8{0, 2, 5, 9, 13, 17, 21, 25, 29}

func (i Function) String() string {
	if i >= Function(len(functionIndex)-1) {
		return fmt.Sprintf("Function(%d)", i)
	}
	return functionName[functionIndex[i]:functionIndex[i+1]]
}

// Pin is a GPIO number (GPIOnn) on BCM238(5|6|7).
//
// If you search for pin per their position on the P1 header, look at ../rpi
// package.
//
// Pin implements host.PinIO.
type Pin struct {
	// Immutable.
	number      int
	name        string
	defaultPull host.Pull

	// Mutable
	lock      sync.Mutex
	valueFile *os.File // handle to /sys/class/gpio/gpio*/value.
	edgeFile  *os.File
	event     [1]syscall.EpollEvent
	epollFd   int // Only thing that actually changes when Edges() is disabled.
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

// In setups a pin as an input and implements host.PinIn.
//
// Specifying a value for pull other than host.PullNoChange causes this
// function to be slightly slower (about 1ms).
//
// Will fail if requesting to change a pin that is set to special functionality.
func (p *Pin) In(pull host.Pull) error {
	if gpioMemory32 == nil {
		return globalError
	}
	if !p.setFunction(In) {
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
		time.Sleep(sleep160cycles)
		// 0x98    RW   GPIO Pin Pull-up/down Enable Clock 0 (GPIO0-31)
		// 0x9C    RW   GPIO Pin Pull-up/down Enable Clock 1 (GPIO32-53)
		offset := 38 + p.number/32
		gpioMemory32[offset] = 1 << uint(p.number%32)

		time.Sleep(sleep160cycles)
		gpioMemory32[37] = 0
		gpioMemory32[offset] = 0
	}
	return nil
}

// Read return the current pin level and implements host.PinIn.
//
// This function is very fast. It works even if the pin is set as output.
func (p *Pin) Read() host.Level {
	if gpioMemory32 == nil {
		return host.Low
	}
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
	p.lock.Lock()
	defer p.lock.Unlock()
	if err := p.setEdge(true); err != nil {
		return nil, err
	}
	c := make(chan host.Level)
	go func() {
		defer close(c)
		var b [1]byte
		for {
			if _, err := p.valueFile.Seek(0, 0); err != nil {
				return
			}
			for {
				p.lock.Lock()
				ep := p.epollFd
				p.lock.Unlock()
				if ep == 0 {
					return
				}
				if nr, err := syscall.EpollWait(ep, p.event[:], -1); err != nil {
					return
				} else if nr < 1 {
					continue
				}
				if _, err := p.valueFile.Read(b[:]); err != nil {
					return
				}
				break
			}
			if b[0] == '1' {
				c <- host.High
			} else {
				c <- host.Low
			}
		}
	}()
	return c, nil
}

// Out sets a pin as output and implements host.PinOut. The caller should
// immediately call SetLow() or SetHigh() afterward.
//
// Will fail if requesting to change a pin that is set to special functionality.
func (p *Pin) Out() error {
	if gpioMemory32 == nil {
		return globalError
	}
	// TODO(maruel): Ensure any Edges() loop is canceled.
	if !p.setFunction(In) {
		return errors.New("failed to change pin mode")
	}
	return nil
}

// Set sets a pin already set for output as host.High or host.Low and implements
// host.PinOut.
//
// This function is very fast.
func (p *Pin) Set(l host.Level) {
	if gpioMemory32 != nil {
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
}

// Special functionality.

// Function returns the current GPIO pin function.
func (p *Pin) Function() Function {
	if gpioMemory32 == nil {
		return Alt5
	}
	// 0x00    RW   GPIO Function Select 0 (GPIO0-9)
	// 0x04    RW   GPIO Function Select 1 (GPIO10-19)
	// 0x08    RW   GPIO Function Select 2 (GPIO20-29)
	// 0x0C    RW   GPIO Function Select 3 (GPIO30-39)
	// 0x10    RW   GPIO Function Select 4 (GPIO40-49)
	// 0x14    RW   GPIO Function Select 5 (GPIO50-53)
	return Function((gpioMemory32[p.number/10] >> uint((p.number%10)*3)) & 7)
}

// DefaultPull returns the default pull for the function.
//
// The CPU doesn't return the current pull.
func (p *Pin) DefaultPull() host.Pull {
	return p.defaultPull
}

// Internal code.

// setEdge changes the edge detection setting for the pin.
//
// It is the function that opens the gpio sysfs file handle.
func (p *Pin) setEdge(enable bool) error {
	p.lock.Lock()
	defer p.lock.Unlock()
	if !enable {
		if p.epollFd == 0 {
			// Disabling twice is fine.
			return nil
		}
		_ = syscall.Close(p.epollFd)
		p.epollFd = 0
		// Do not close the handles, just disable the interrupts.
		_, err := p.edgeFile.Write([]byte("none"))
		return err
	}
	if p.epollFd != 0 {
		// Enabling while already enabled is bad.
		return errors.New("already enabled")
	}
	if err := p.open(); err != nil {
		return err
	}
	_, err := p.edgeFile.Write([]byte("both"))
	return err
}

// open opens the gpio sysfs handles. Assumes lock is held.
func (p *Pin) open() error {
	var err error
	if p.valueFile == nil {
		// Assume the pin is exported first. The reason is that exporting a pin that
		// is already exported causes a write failure, which is difficult to
		// differentiate from other errors.
		// On the other hand, accessing /sys/class/gpio/gpio*/value when it is not
		// exported returns a permission denied error. :/
		if p.valueFile, err = os.OpenFile(fmt.Sprintf("/sys/class/gpio/gpio%d/value", p), os.O_RDONLY, 0600); err != nil {
			// Export the pin.
			if err = openExport(); err == nil {
				if _, err = exportHandle.Write([]byte(strconv.Itoa(p.number))); err == nil {
					p.valueFile, err = os.OpenFile(fmt.Sprintf("/sys/class/gpio/gpio%d/value", p), os.O_RDONLY, 0600)
				}
			}
		}
	}
	if p.edgeFile == nil && err == nil {
		// TODO(maruel): Figure out the problem or better use the register instead
		// of the file.
		for i := 0; i < 30 && p.edgeFile == nil; i++ {
			p.edgeFile, err = os.OpenFile(fmt.Sprintf("/sys/class/gpio/gpio%d/edge", p), os.O_WRONLY, 0600)
			// TODO(maruel): Figure out what the hell.
			time.Sleep(time.Millisecond)
		}
	}
	if p.epollFd == 0 && err == nil {
		if p.epollFd, err = syscall.EpollCreate(1); err == nil {
			const EPOLLPRI = 2
			const EPOLL_CTL_ADD = 1
			p.event[0].Events = EPOLLPRI
			p.event[0].Fd = int32(p.valueFile.Fd())
			err = syscall.EpollCtl(p.epollFd, EPOLL_CTL_ADD, int(p.valueFile.Fd()), &p.event[0])
		}
	}
	return err
}

// setFunction changes the GPIO pin function.
//
// Returns false if the pin was in AltN. Only accepts In and Out
func (p *Pin) setFunction(f Function) bool {
	if f != In && f != Out {
		return false
	}
	if actual := p.Function(); actual != In && actual != Out {
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

var globalError error

// Changing pull resistor require a 150 cycles sleep. Use 160 to be safe.
var sleep160cycles time.Duration = time.Second * 160 / time.Duration(cpu.MaxSpeed)

func setIfAlt0(p *Pin, special *host.Pin) {
	if p.Function() == Alt0 {
		if (*special).String() != "INVALID" {
			//fmt.Printf("%s and %s have same functionality\n", p, *special)
		}
		*special = p
	}
}

func setIfAlt(p *Pin, special0 *host.Pin, special1 *host.Pin, special2 *host.Pin, special3 *host.Pin, special4 *host.Pin, special5 *host.Pin) {
	switch p.Function() {
	case Alt0:
		if special0 != nil {
			if (*special0).String() != "INVALID" {
				//fmt.Printf("%s and %s have same functionality\n", p, *special0)
			}
			*special0 = p
		}
	case Alt1:
		if special1 != nil {
			if (*special1).String() != "INVALID" {
				//fmt.Printf("%s and %s have same functionality\n", p, *special1)
			}
			*special1 = p
		}
	case Alt2:
		if special2 != nil {
			if (*special2).String() != "INVALID" {
				//log.Printf("%s and %s have same functionality\n", p, *special2)
			}
			*special2 = p
		}
	case Alt3:
		if special3 != nil {
			if (*special3).String() != "INVALID" {
				//log.Printf("%s and %s have same functionality\n", p, *special3)
			}
			*special3 = p
		}
	case Alt4:
		if special4 != nil {
			if (*special4).String() != "INVALID" {
				//log.Printf("%s and %s have same functionality\n", p, *special4)
			}
			*special4 = p
		}
	case Alt5:
		if special5 != nil {
			if (*special5).String() != "INVALID" {
				//log.Printf("%s and %s have same functionality\n", p, *special5)
			}
			*special5 = p
		}
	}
}

func init() {
	GPIO0 = &Pin{number: 0, name: "GPIO0", defaultPull: host.Up}
	GPIO1 = &Pin{number: 1, name: "GPIO1", defaultPull: host.Up}
	GPIO2 = &Pin{number: 2, name: "GPIO2", defaultPull: host.Up}
	GPIO3 = &Pin{number: 3, name: "GPIO3", defaultPull: host.Up}
	GPIO4 = &Pin{number: 4, name: "GPIO4", defaultPull: host.Up}
	GPIO5 = &Pin{number: 5, name: "GPIO5", defaultPull: host.Up}
	GPIO6 = &Pin{number: 6, name: "GPIO6", defaultPull: host.Up}
	GPIO7 = &Pin{number: 7, name: "GPIO7", defaultPull: host.Up}
	GPIO8 = &Pin{number: 8, name: "GPIO8", defaultPull: host.Up}
	GPIO9 = &Pin{number: 9, name: "GPIO9", defaultPull: host.Down}
	GPIO10 = &Pin{number: 10, name: "GPIO10", defaultPull: host.Down}
	GPIO11 = &Pin{number: 11, name: "GPIO11", defaultPull: host.Down}
	GPIO12 = &Pin{number: 12, name: "GPIO12", defaultPull: host.Down}
	GPIO13 = &Pin{number: 13, name: "GPIO13", defaultPull: host.Down}
	GPIO14 = &Pin{number: 14, name: "GPIO14", defaultPull: host.Down}
	GPIO15 = &Pin{number: 15, name: "GPIO15", defaultPull: host.Down}
	GPIO16 = &Pin{number: 16, name: "GPIO16", defaultPull: host.Down}
	GPIO17 = &Pin{number: 17, name: "GPIO17", defaultPull: host.Down}
	GPIO18 = &Pin{number: 18, name: "GPIO18", defaultPull: host.Down}
	GPIO19 = &Pin{number: 19, name: "GPIO19", defaultPull: host.Down}
	GPIO20 = &Pin{number: 20, name: "GPIO20", defaultPull: host.Down}
	GPIO21 = &Pin{number: 21, name: "GPIO21", defaultPull: host.Down}
	GPIO22 = &Pin{number: 22, name: "GPIO22", defaultPull: host.Down}
	GPIO23 = &Pin{number: 23, name: "GPIO23", defaultPull: host.Down}
	GPIO24 = &Pin{number: 24, name: "GPIO24", defaultPull: host.Down}
	GPIO25 = &Pin{number: 25, name: "GPIO25", defaultPull: host.Down}
	GPIO26 = &Pin{number: 26, name: "GPIO26", defaultPull: host.Down}
	GPIO27 = &Pin{number: 27, name: "GPIO27", defaultPull: host.Down}
	GPIO28 = &Pin{number: 28, name: "GPIO28", defaultPull: host.Float}
	GPIO29 = &Pin{number: 29, name: "GPIO29", defaultPull: host.Float}
	GPIO30 = &Pin{number: 30, name: "GPIO30", defaultPull: host.Down}
	GPIO31 = &Pin{number: 31, name: "GPIO31", defaultPull: host.Down}
	GPIO32 = &Pin{number: 32, name: "GPIO32", defaultPull: host.Down}
	GPIO33 = &Pin{number: 33, name: "GPIO33", defaultPull: host.Down}
	GPIO34 = &Pin{number: 34, name: "GPIO34", defaultPull: host.Up}
	GPIO35 = &Pin{number: 35, name: "GPIO35", defaultPull: host.Up}
	GPIO36 = &Pin{number: 36, name: "GPIO36", defaultPull: host.Up}
	GPIO37 = &Pin{number: 37, name: "GPIO37", defaultPull: host.Down}
	GPIO38 = &Pin{number: 38, name: "GPIO38", defaultPull: host.Down}
	GPIO39 = &Pin{number: 39, name: "GPIO39", defaultPull: host.Down}
	GPIO40 = &Pin{number: 40, name: "GPIO40", defaultPull: host.Down}
	GPIO41 = &Pin{number: 41, name: "GPIO41", defaultPull: host.Down}
	GPIO42 = &Pin{number: 42, name: "GPIO42", defaultPull: host.Down}
	GPIO43 = &Pin{number: 43, name: "GPIO43", defaultPull: host.Down}
	GPIO44 = &Pin{number: 44, name: "GPIO44", defaultPull: host.Float}
	GPIO45 = &Pin{number: 45, name: "GPIO45", defaultPull: host.Float}
	GPIO46 = &Pin{number: 46, name: "GPIO46", defaultPull: host.Up}
	GPIO47 = &Pin{number: 47, name: "GPIO47", defaultPull: host.Up}
	GPIO48 = &Pin{number: 48, name: "GPIO48", defaultPull: host.Up}
	GPIO49 = &Pin{number: 49, name: "GPIO49", defaultPull: host.Up}
	GPIO50 = &Pin{number: 50, name: "GPIO50", defaultPull: host.Up}
	GPIO51 = &Pin{number: 51, name: "GPIO51", defaultPull: host.Up}
	GPIO52 = &Pin{number: 52, name: "GPIO52", defaultPull: host.Up}
	GPIO53 = &Pin{number: 53, name: "GPIO53", defaultPull: host.Up}
	host.AllPins = []host.PinIO{
		GPIO0, GPIO1, GPIO2, GPIO3, GPIO4, GPIO5, GPIO6, GPIO7, GPIO8, GPIO9,
		GPIO10, GPIO11, GPIO12, GPIO13, GPIO14, GPIO15, GPIO16, GPIO17, GPIO18,
		GPIO19, GPIO20, GPIO21, GPIO22, GPIO23, GPIO24, GPIO25, GPIO26, GPIO27,
		GPIO28, GPIO29, GPIO30, GPIO31, GPIO32, GPIO33, GPIO34, GPIO35, GPIO36,
		GPIO37, GPIO38, GPIO39, GPIO40, GPIO41, GPIO42, GPIO43, GPIO44, GPIO45,
		GPIO46, GPIO47, GPIO48, GPIO49, GPIO50, GPIO51, GPIO52, GPIO53,
	}

	GPCLK0 = pins.INVALID
	GPCLK1 = pins.INVALID
	GPCLK2 = pins.INVALID
	I2C_SCL0 = pins.INVALID
	I2C_SDA0 = pins.INVALID
	I2C_SCL1 = pins.INVALID
	I2C_SDA1 = pins.INVALID
	IR_IN = pins.INVALID
	IR_OUT = pins.INVALID
	PCM_CLK = pins.INVALID
	PCM_FS = pins.INVALID
	PCM_DIN = pins.INVALID
	PCM_DOUT = pins.INVALID
	PWM0_OUT = pins.INVALID
	PWM1_OUT = pins.INVALID
	SPI0_CE0 = pins.INVALID
	SPI0_CE1 = pins.INVALID
	SPI0_CLK = pins.INVALID
	SPI0_MISO = pins.INVALID
	SPI0_MOSI = pins.INVALID
	SPI1_CE0 = pins.INVALID
	SPI1_CE1 = pins.INVALID
	SPI1_CE2 = pins.INVALID
	SPI1_CLK = pins.INVALID
	SPI1_MISO = pins.INVALID
	SPI1_MOSI = pins.INVALID
	SPI2_MISO = pins.INVALID
	SPI2_MOSI = pins.INVALID
	SPI2_CLK = pins.INVALID
	SPI2_CE0 = pins.INVALID
	SPI2_CE1 = pins.INVALID
	SPI2_CE2 = pins.INVALID
	UART_RXD0 = pins.INVALID
	UART_CTS0 = pins.INVALID
	UART_CTS1 = pins.INVALID
	UART_RTS0 = pins.INVALID
	UART_RTS1 = pins.INVALID
	UART_TXD0 = pins.INVALID
	UART_RXD1 = pins.INVALID
	UART_TXD1 = pins.INVALID

	var mem *gpiomem.Mem
	mem, globalError = gpiomem.Open()
	if globalError != nil {
		return
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

	in, out := ir.Pins()
	if in != -1 {
		IR_IN = host.GetPinByNumber(in)
	}
	if out != -1 {
		IR_OUT = host.GetPinByNumber(out)
	}
}
