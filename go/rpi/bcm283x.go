// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// BCM283x GPIO GPIO handling.

//go:generate stringer -type Edge
//go:generate stringer -type Function
//go:generate stringer -type Pin
//go:generate stringer -type Pull

// From go get -u golang.org/x/tools/cmd/stringer

package rpi

import (
	"os"
	"reflect"
	"syscall"
	"time"
	"unsafe"
)

// Function specifies the active functionality of a pin. The alternative
// function is GPIO pin dependent.
//
// Values are specified according to Broadcom spec:
// https://www.raspberrypi.org/wp-content/uploads/2012/02/BCM2835-ARM-Peripherals.pdf
// page 92.
type Function uint8

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

// Level is the level of the pin: Low or High.
type Level bool

const (
	Low  Level = false
	High Level = true
)

func (l Level) String() string {
	if l == Low {
		return "Low"
	}
	return "High"
}

// Edge specifies interrupt based triggering for a pin set as input.
type Edge uint8

const (
	EdgeNone Edge = 0
	Rising   Edge = 1
	Falling  Edge = 2
	EdgeBoth Edge = 3
)

// Pull specifies the internal pull-up or pull-down for a pin set as input.
type Pull uint8

const (
	Float        Pull = 0 // Let the input float
	Down         Pull = 1 // Apply 50KOhm~60kOhm pull-down
	Up           Pull = 2 // Apply 50kOhm~65kOhm pull-up
	PullNoChange Pull = 3 // Do not change the previous pull resistor setting
)

// Pin is a GPIO number (GPIOnn) on BCM238(5|6|7) or can be aliased by bus
// (P1_xx) or functionality (e.g. SPI0_xx).
//
// As defined at
// https://www.raspberrypi.org/wp-content/uploads/2012/02/BCM2835-ARM-Peripherals.pdf
// page 102 and 103.
type Pin uint8

const (
	INVALID Pin = 255 //
	GROUND  Pin = 254 // Connected to Ground
	V3_3    Pin = 253 // Connected to 3.3v
	V5      Pin = 252 // Connected to 5v
	GPIO0   Pin = 0   // High,  P1_27, I2C_SDA0
	GPIO1   Pin = 1   // High,  P1_28, I2C_SCL0
	GPIO2   Pin = 2   // High,  P1_3,  I2C_SDA1
	GPIO3   Pin = 3   // High,  P1_5,  I2C_SCL1
	GPIO4   Pin = 4   // High,  P1_7,  GPCLK0
	GPIO5   Pin = 5   // High,  P1_29, GPCLK1
	GPIO6   Pin = 6   // High,  P1_31, GPCLK2
	GPIO7   Pin = 7   // High,  P1_26, SPI0_CE1
	GPIO8   Pin = 8   // High,  P1_24, SPI0_CE0
	GPIO9   Pin = 9   // Low,   P1_21, SPI0_MISO
	GPIO10  Pin = 10  // Low,   P1_19, SPI0_MOSI
	GPIO11  Pin = 11  // Low,   P1_23, SPI0_CLK
	GPIO12  Pin = 12  // Low,   P1_32, PWM0
	GPIO13  Pin = 13  // Low,   P1_33, PWM1
	GPIO14  Pin = 14  // Low,   P1_8,  UART_TXD0, UART_TXD1
	GPIO15  Pin = 15  // Low,   P1_10, UART_RXD0, UART_RXD1
	GPIO16  Pin = 16  // Low,   P1_36
	GPIO17  Pin = 17  // Low,   P1_11, SPI1_CE1, IR_IN
	GPIO18  Pin = 18  // Low,   P1_12, IR_OUT
	GPIO19  Pin = 19  // Low,   P1_35
	GPIO20  Pin = 20  // Low,   P1_38
	GPIO21  Pin = 21  // Low,   P1_40
	GPIO22  Pin = 22  // Low,   P1_15
	GPIO23  Pin = 23  // Low,   P1_16
	GPIO24  Pin = 24  // Low,   P1_18
	GPIO25  Pin = 25  // Low,   P1_22
	GPIO26  Pin = 26  // Low,   P1_37
	GPIO27  Pin = 27  // Low,   P1_13
	GPIO28  Pin = 28  // Float, Not connected, SDA0, PCM_CLK
	GPIO29  Pin = 29  // Float, Not connected, SCL0, PCM_FS
	GPIO30  Pin = 30  // Low,   Not connected, PCM_DIN, UART_CTS0, UARTS_CTS1
	GPIO31  Pin = 31  // Low,   Not connected, PCM_DOUT, UART_RTS0, UARTS_RTS1
	GPIO32  Pin = 32  // Low,   Not connected, GPCLK0, UART_TXD0, UARTS_TXD1
	GPIO33  Pin = 33  // Low,   Not connected, UART_RXD0, UARTS_RXD1
	GPIO34  Pin = 34  // High,  Not connected, GPCLK0
	GPIO35  Pin = 35  // High,  Not connected, SPI0_CE1
	GPIO36  Pin = 36  // High,  Not connected, SPI0_CE0, UART_TXD0
	GPIO37  Pin = 37  // Low,   Not connected, SPI0_MISO, UART_RXD0
	GPIO38  Pin = 38  // Low,   Not connected, SPI0_MOSI, UART_RTS0
	GPIO39  Pin = 39  // Low,   Not connected, SPI0_CLK, UART_CTS0
	GPIO40  Pin = 40  // Low,   Audio Right, PWM0_OUT; with low pass filter
	GPIO41  Pin = 41  // Low,   Audio Left, PWM1_OUT; with low pass filter
	GPIO42  Pin = 42  // Low,   Not connected, GPCLK1, SPI2_CLK, UART_RTS1
	GPIO43  Pin = 43  // Low,   Not connected, GPCLK2, SPI2_CE0, UART_CTS1
	GPIO44  Pin = 44  // Float, Not connected, GPCLK1, I2C_SDA0, I2C_SDA1, SPI2_CE1
	GPIO45  Pin = 45  // Float, Audio Left, PWM1_OUT, I2C_SCL0, I2C_SCL1, SPI2_CE2
	GPIO46  Pin = 46  // High,  HDMI hotplug detect
	GPIO47  Pin = 47  // High,  SDCard
	GPIO48  Pin = 48  // High,  SDCard
	GPIO49  Pin = 49  // High,  SDCard
	GPIO50  Pin = 50  // High,  SDCard
	GPIO51  Pin = 51  // High,  SDCard
	GPIO52  Pin = 52  // High,  SDCard
	GPIO53  Pin = 53  // High,  SDCard
)

// Special functions that can be assigned to a GPIO. The values are probed and
// set at runtime. Changing the value of the variables has no effect.
var (
	GPCLK0    Pin = INVALID // GPIO4, GPIO20, GPIO32, GPIO34 (also named GPIO_GCLK)
	GPCLK1    Pin = INVALID // GPIO5, GPIO21
	GPCLK2    Pin = INVALID // GPIO6
	I2C_SCL0  Pin = INVALID // GPIO1, GPIO28, GPIO45
	I2C_SDA0  Pin = INVALID // GPIO0, GPIO29, GPIO44
	I2C_SCL1  Pin = INVALID // GPIO3, GPIO45
	I2C_SDA1  Pin = INVALID // GPIO2, GPIO44
	IR_IN     Pin = INVALID // (any GPIO)
	IR_OUT    Pin = INVALID // (any GPIO)
	PCM_CLK   Pin = INVALID // GPIO18, GPIO28 (I2S)
	PCM_FS    Pin = INVALID // GPIO19, GPIO29
	PCM_DIN   Pin = INVALID // GPIO20, GPIO30
	PCM_DOUT  Pin = INVALID // GPIO21, GPIO31
	PWM0_OUT  Pin = INVALID // GPIO12, GPIO18, GPIO40
	PWM1_OUT  Pin = INVALID // GPIO13, GPIO19, GPIO41, GPIO45
	SPI0_CE0  Pin = INVALID // GPIO8,  GPIO36
	SPI0_CE1  Pin = INVALID // GPIO7,  GPIO35
	SPI0_CLK  Pin = INVALID // GPIO11, GPIO39
	SPI0_MISO Pin = INVALID // GPIO9,  GPIO37
	SPI0_MOSI Pin = INVALID // GPIO10, GPIO38
	SPI1_CE0  Pin = INVALID // GPIO18
	SPI1_CE1  Pin = INVALID // GPIO17
	SPI1_CE2  Pin = INVALID // GPIO16
	SPI1_CLK  Pin = INVALID // GPIO21
	SPI1_MISO Pin = INVALID // GPIO19
	SPI1_MOSI Pin = INVALID // GPIO20
	spi2_miso Pin = INVALID // GPIO40
	spi2_mosi Pin = INVALID // GPIO41
	spi2_clk  Pin = INVALID // GPIO42
	spi2_ce0  Pin = INVALID // GPIO43
	spi2_ce1  Pin = INVALID // GPIO44
	spi2_ce2  Pin = INVALID // GPIO45
	UART_RXD0 Pin = INVALID // GPIO15, GPIO33, GPIO37
	UART_CTS0 Pin = INVALID // GPIO16, GPIO30, GPIO39
	UART_CTS1 Pin = INVALID // GPIO16, GPIO30
	UART_RTS0 Pin = INVALID // GPIO17, GPIO31, GPIO38
	UART_RTS1 Pin = INVALID // GPIO17, GPIO31
	UART_TXD0 Pin = INVALID // GPIO14, GPIO32, GPIO36
	UART_RXD1 Pin = INVALID // GPIO15, GPIO33, GPIO41
	UART_TXD1 Pin = INVALID // GPIO14, GPIO32, GPIO40
)

// Function returns the current GPIO pin function.
func (p Pin) Function() Function {
	// (0x00-0x14)
	return Function((gpioMemory32[p/10] >> ((p % 10) * 3)) & 7)
}

// isClock returns true if the pin is owned an output clock.
func (p Pin) isClock() bool {
	// https://www.raspberrypi.org/wp-content/uploads/2012/02/BCM2835-ARM-Peripherals.pdf
	// Page 105.
	// TODO(maruel): Add code to generate clocks.
	switch p {
	case INVALID:
		return false
	case GPCLK0, GPCLK1, GPCLK2:
		return true
	default:
		return false
	}
}

// isI2C returns true if the pin is owned by the I2C driver.
func (p Pin) isI2C() bool {
	switch p {
	case INVALID:
		return false
	case I2C_SCL0, I2C_SDA0, I2C_SCL1, I2C_SDA1:
		return true
	default:
		return false
	}
}

// isI2S returns true if the pin is owned by the I2S driver.
func (p Pin) isI2S() bool {
	switch p {
	case INVALID:
		return false
	case PCM_CLK, PCM_FS, PCM_DIN, PCM_DOUT:
		return true
	default:
		return false
	}
}

// isIR returns true if the pin is owned by the LIRC driver.
func (p Pin) isIR() bool {
	switch p {
	case INVALID:
		return false
	case IR_IN, IR_OUT:
		return true
	default:
		return false
	}
}

// isPWM returns true if the pin is owned by the PWM driver. By default used
// for audio output.
func (p Pin) isPWM() bool {
	switch p {
	case INVALID:
		return false
	case PWM0_OUT, PWM1_OUT:
		return true
	default:
		return false
	}
}

// isSPI returns true if the pin is owned by the SPI driver.
func (p Pin) isSPI() bool {
	switch p {
	case INVALID:
		return false
	case SPI0_CE0, SPI0_CE1, SPI0_CLK, SPI0_MISO, SPI0_MOSI, SPI1_CE0, SPI1_CE1, SPI1_CE2, SPI1_CLK, SPI1_MISO, SPI1_MOSI:
		return true
	default:
		return false
	}
}

// isUART returns true if the pin is owned by the UART driver.
func (p Pin) isUART() bool {
	switch p {
	case INVALID:
		return false
	case UART_RXD0, UART_TXD0, UART_CTS0, UART_RTS0, UART_RXD1, UART_TXD1, UART_CTS1, UART_RTS1:
		return true
	default:
		return false
	}
}

// TODO(maruel): https://www.kernel.org/doc/Documentation/gpio/sysfs.txt
// Example:
// $ echo "19"  > /sys/class/gpio/export
// $ cat /sys/class/gpio/gpio19/direction
// in
// It's great except that internal pull-up and pull-down are not exposed and
// this is a show stopper.

// In setups a pin as an input.
func (p Pin) In(pull Pull) {
	p.setFunction(In)
	if pull != PullNoChange {
		// https://www.raspberrypi.org/wp-content/uploads/2012/02/BCM2835-ARM-Peripherals.pdf
		// page 101.

		// Set Pull (0x94)
		gpioMemory32[37] = uint32(pull)

		// Datasheet states caller needs to sleep 150 cycles.
		time.Sleep(sleep160cycles)
		offset := 38 + p/32 // (0x98-0x9C)
		gpioMemory32[offset] = 1 << (p % 32)

		time.Sleep(sleep160cycles)
		gpioMemory32[37] = 0
		gpioMemory32[offset] = 0
	}
}

// Read return the current pin level.
func (p Pin) Read() Level {
	// (0x34-0x38)
	return Level((gpioMemory32[13+p/32] & (1 << p)) != 0)
}

// Edge waits until an edge was caught. Doesn't clear previous interruptions.
func (p Pin) Edge(edge Edge) Level {
	// Opportunistically open the gpiofs when needed.
	return Low
}

// ClearEdge clears pending interrupts.
func (p Pin) ClearEdge() {
	/*
		syscall.Ioctl(sysFds[bcmGpioPin], FIONREAD, &count)
		for i := 0; i < count; i++ {
			read(sysFds[bcmGpioPin], &c, 1)
		}
	*/
}

// Out sets a pin as output and sets the initial level.
func (p Pin) Out(level Level) {
	p.setFunction(Out)
	if level == Low {
		p.Low()
	} else {
		p.High()
	}
}

// Low sets a pin already set for output as low. Faster than calling
// Pin.Out(Low).
func (p Pin) Low() {
	// 0x28-0x2C
	gpioMemory32[10+p/32] = 1 << (p & 31)
}

// High sets a pin already set for output as low. Faster than calling
// Pin.Out(High).
func (p Pin) High() {
	// 0x1C-0x20
	gpioMemory32[7+p/32] = 1 << (p & 31)
}

// setFunction changes the GPIO pin function.
//
// TODO(maruel): Refuse to do so when a pin is in altN mode.
//
// TODO(maruel): Update the special function mapping when relevant. Not yet an
// issue since only 'in' and 'out' are used.
func (p Pin) setFunction(f Function) {
	off := p / 10 // (0x00-0x14)
	shift := (p % 10) * 3
	gpioMemory32[off] = (gpioMemory32[off] &^ (7 << shift)) | (uint32(f) << shift)
}

/*
// Close the handle implicitly open by either SetPinPWM or ReleasePinPWM.
//
// It's not required to call it. It doesn't reset the GPIO pin either.
func Close() error {
	err1 := closePiblaster()
	err2 := closeGPIOMem()
	if err1 != nil {
		return err1
	}
	return err2
}
*/

//

// Handle to /dev/gpiomem
var gpioMemory8 []byte

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
var sleep160cycles time.Duration = 220 * time.Nanosecond

func setIfAlt0(p Pin, special *Pin) {
	if p.Function() == Alt0 {
		if *special != INVALID {
			//fmt.Printf("%s and %s have same functionality\n", p, *special)
		}
		*special = p
	}
}

func setIfAlt(p Pin, special0 *Pin, special1 *Pin, special2 *Pin, special3 *Pin, special4 *Pin, special5 *Pin) {
	switch p.Function() {
	case Alt0:
		if special0 != nil {
			if *special0 != INVALID {
				//fmt.Printf("%s and %s have same functionality\n", p, *special0)
			}
			*special0 = p
		}
	case Alt1:
		if special1 != nil {
			if *special1 != INVALID {
				//fmt.Printf("%s and %s have same functionality\n", p, *special1)
			}
			*special1 = p
		}
	case Alt2:
		if special2 != nil {
			if *special2 != INVALID {
				//fmt.Printf("%s and %s have same functionality\n", p, *special2)
			}
			*special2 = p
		}
	case Alt3:
		if special3 != nil {
			if *special3 != INVALID {
				//fmt.Printf("%s and %s have same functionality\n", p, *special3)
			}
			*special3 = p
		}
	case Alt4:
		if special4 != nil {
			if *special4 != INVALID {
				//fmt.Printf("%s and %s have same functionality\n", p, *special4)
			}
			*special4 = p
		}
	case Alt5:
		if special5 != nil {
			if *special5 != INVALID {
				//fmt.Printf("%s and %s have same functionality\n", p, *special5)
			}
			*special5 = p
		}
	}
}

func init() {
	gpioMemory8, globalError = openGPIOMem()
	if globalError != nil {
		return
	}
	gpioMemory32 = unsafeRemap(gpioMemory8)

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
	setIfAlt0(GPIO28, &I2C_SDA0)                                           // Not connected
	setIfAlt0(GPIO29, &I2C_SCL0)                                           // Not connected
	setIfAlt(GPIO30, nil, nil, nil, &UART_CTS0, nil, &UART_CTS1)           // Not connected
	setIfAlt(GPIO31, nil, nil, nil, &UART_RTS0, nil, &UART_RTS1)           // Not connected
	setIfAlt(GPIO32, &GPCLK0, nil, nil, &UART_TXD0, nil, &UART_TXD1)       // Not connected
	setIfAlt(GPIO33, nil, nil, nil, &UART_RXD0, nil, &UART_RXD1)           // Not connected
	setIfAlt0(GPIO34, &GPCLK0)                                             // Not connected
	setIfAlt0(GPIO35, &SPI0_CE1)                                           // Not connected
	setIfAlt(GPIO36, &SPI0_CE0, nil, &UART_TXD0, nil, nil, nil)            // Not connected
	setIfAlt(GPIO37, &SPI0_MISO, nil, &UART_RXD0, nil, nil, nil)           // Not connected
	setIfAlt(GPIO38, &SPI0_MOSI, nil, &UART_RTS0, nil, nil, nil)           // Not connected
	setIfAlt(GPIO39, &SPI0_CLK, nil, &UART_CTS0, nil, nil, nil)            // Not connected
	setIfAlt(GPIO40, &PWM0_OUT, nil, nil, nil, &spi2_miso, &UART_TXD1)     // Connected to audio right
	setIfAlt(GPIO41, &PWM1_OUT, nil, nil, nil, &spi2_mosi, &UART_RXD1)     // Connected to audio left
	setIfAlt(GPIO42, &GPCLK1, nil, nil, nil, &spi2_clk, &UART_RTS1)        // Not connected
	setIfAlt(GPIO43, &GPCLK2, nil, nil, nil, &spi2_ce0, &UART_CTS1)        // Not connected
	setIfAlt(GPIO44, &GPCLK1, &I2C_SDA0, &I2C_SDA1, nil, &spi2_ce1, nil)   // Not connected
	setIfAlt(GPIO45, &PWM1_OUT, &I2C_SCL0, &I2C_SCL1, nil, &spi2_ce2, nil) // Not connected
	// GPIO46-GPIO53 do not have interesting alternate function.

	IR_IN, IR_OUT = getLIRCPins()
}

func openGPIOMem() ([]uint8, error) {
	f, err := os.OpenFile("/dev/gpiomem", os.O_RDWR|os.O_SYNC, 0)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// TODO(maruel): Map PWM, CLK, PADS, TIMER
	return syscall.Mmap(int(f.Fd()), 0, 4096, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
}

func unsafeRemap(i []byte) []uint32 {
	// I/O needs to happen as 32 bits operation, so remap accordingly.
	header := *(*reflect.SliceHeader)(unsafe.Pointer(&i))
	header.Len /= 4
	header.Cap /= 4
	return *(*[]uint32)(unsafe.Pointer(&header))
}

// closeGPIOMem unmaps /dev/gpiomem. Not sure why someone would want to do that.
func closeGPIOMem() error {
	if gpioMemory8 != nil {
		m := gpioMemory8
		gpioMemory8 = nil
		gpioMemory32 = nil
		return syscall.Munmap(m)
	}
	return nil
}
