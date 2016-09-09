// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package bcm283x

import (
	"errors"
	"fmt"
	"time"

	"github.com/maruel/dlibox/go/pio/host"
	"github.com/maruel/dlibox/go/pio/host/cpu"
	"github.com/maruel/dlibox/go/pio/host/internal/gpiomem"
	"github.com/maruel/dlibox/go/pio/host/ir"
)

// Function specifies the active functionality of a pin. The alternative
// function is GPIO pin dependent.
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

const functionName = "InOutAlt5Alt4Alt0Alt1Alt2Alt3"

var functionIndex = [...]uint8{0, 2, 5, 9, 13, 17, 21, 25, 29}

func (i Function) String() string {
	if i >= Function(len(functionIndex)-1) {
		return fmt.Sprintf("Function(%d)", i)
	}
	return functionName[functionIndex[i]:functionIndex[i+1]]
}

// Pin is a GPIO number (GPIOnn) on BCM238(5|6|7). If you search for pin per
// their position on the P1 header, look at ../rpi package.
//
// Pin implements host.Pin.
type Pin uint8

const (
	INVALID Pin = 255 //
	GROUND  Pin = 254 // Connected to Ground
	V3_3    Pin = 253 // Connected to 3.3v
	V5      Pin = 252 // Connected to 5v
	GPIO0   Pin = 0   // High,  I2C_SDA0
	GPIO1   Pin = 1   // High,  I2C_SCL0
	GPIO2   Pin = 2   // High,  I2C_SDA1
	GPIO3   Pin = 3   // High,  I2C_SCL1
	GPIO4   Pin = 4   // High,  GPCLK0
	GPIO5   Pin = 5   // High,  GPCLK1
	GPIO6   Pin = 6   // High,  GPCLK2
	GPIO7   Pin = 7   // High,  SPI0_CE1
	GPIO8   Pin = 8   // High,  SPI0_CE0
	GPIO9   Pin = 9   // Low,   SPI0_MISO
	GPIO10  Pin = 10  // Low,   SPI0_MOSI
	GPIO11  Pin = 11  // Low,   SPI0_CLK
	GPIO12  Pin = 12  // Low,   PWM0_OUT
	GPIO13  Pin = 13  // Low,   PWM1_OUT
	GPIO14  Pin = 14  // Low,   UART_TXD0, UART_TXD1
	GPIO15  Pin = 15  // Low,   UART_RXD0, UART_RXD1
	GPIO16  Pin = 16  // Low,   UART_CTS0, SPI1_CE2, UART_CTS1
	GPIO17  Pin = 17  // Low,   UART_RTS0, SPI1_CE1, UART_RTS1
	GPIO18  Pin = 18  // Low,   PCM_CLK, SPI1_CE0, PWM0_OUT
	GPIO19  Pin = 19  // Low,   PCM_FS, SPI1_MISO, PWM1_OUT
	GPIO20  Pin = 20  // Low,   PCM_DIN, SPI1_MOSI, GPCLK0
	GPIO21  Pin = 21  // Low,   PCM_DOUT, SPI1_CLK, GPCLK1
	GPIO22  Pin = 22  // Low,
	GPIO23  Pin = 23  // Low,
	GPIO24  Pin = 24  // Low,
	GPIO25  Pin = 25  // Low,
	GPIO26  Pin = 26  // Low,
	GPIO27  Pin = 27  // Low,
	GPIO28  Pin = 28  // Float, I2C_SDA0, PCM_CLK
	GPIO29  Pin = 29  // Float, I2C_SCL0, PCM_FS
	GPIO30  Pin = 30  // Low,   PCM_DIN, UART_CTS0, UARTS_CTS1
	GPIO31  Pin = 31  // Low,   PCM_DOUT, UART_RTS0, UARTS_RTS1
	GPIO32  Pin = 32  // Low,   GPCLK0, UART_TXD0, UARTS_TXD1
	GPIO33  Pin = 33  // Low,   UART_RXD0, UARTS_RXD1
	GPIO34  Pin = 34  // High,  GPCLK0
	GPIO35  Pin = 35  // High,  SPI0_CE1
	GPIO36  Pin = 36  // High,  SPI0_CE0, UART_TXD0
	GPIO37  Pin = 37  // Low,   SPI0_MISO, UART_RXD0
	GPIO38  Pin = 38  // Low,   SPI0_MOSI, UART_RTS0
	GPIO39  Pin = 39  // Low,   SPI0_CLK, UART_CTS0
	GPIO40  Pin = 40  // Low,   PWM0_OUT, SPI2_MISO, UART_TXD1
	GPIO41  Pin = 41  // Low,   PWM1_OUT, SPI2_MOSI, UART_RXD1
	GPIO42  Pin = 42  // Low,   GPCLK1, SPI2_CLK, UART_RTS1
	GPIO43  Pin = 43  // Low,   GPCLK2, SPI2_CE0, UART_CTS1
	GPIO44  Pin = 44  // Float, GPCLK1, I2C_SDA0, I2C_SDA1, SPI2_CE1
	GPIO45  Pin = 45  // Float, PWM1_OUT, I2C_SCL0, I2C_SCL1, SPI2_CE2
	GPIO46  Pin = 46  // High,
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
	GPCLK0    host.Pin = INVALID // GPIO4, GPIO20, GPIO32, GPIO34 (also named GPIO_GCLK)
	GPCLK1    host.Pin = INVALID // GPIO5, GPIO21, GPIO42, GPIO44
	GPCLK2    host.Pin = INVALID // GPIO6, GPIO43
	I2C_SCL0  host.Pin = INVALID // GPIO1, GPIO29, GPIO45
	I2C_SDA0  host.Pin = INVALID // GPIO0, GPIO28, GPIO44
	I2C_SCL1  host.Pin = INVALID // GPIO3, GPIO45
	I2C_SDA1  host.Pin = INVALID // GPIO2, GPIO44
	IR_IN     host.Pin = INVALID // (any GPIO)
	IR_OUT    host.Pin = INVALID // (any GPIO)
	PCM_CLK   host.Pin = INVALID // GPIO18, GPIO28 (I2S)
	PCM_FS    host.Pin = INVALID // GPIO19, GPIO29
	PCM_DIN   host.Pin = INVALID // GPIO20, GPIO30
	PCM_DOUT  host.Pin = INVALID // GPIO21, GPIO31
	PWM0_OUT  host.Pin = INVALID // GPIO12, GPIO18, GPIO40
	PWM1_OUT  host.Pin = INVALID // GPIO13, GPIO19, GPIO41, GPIO45
	SPI0_CE0  host.Pin = INVALID // GPIO8,  GPIO36
	SPI0_CE1  host.Pin = INVALID // GPIO7,  GPIO35
	SPI0_CLK  host.Pin = INVALID // GPIO11, GPIO39
	SPI0_MISO host.Pin = INVALID // GPIO9,  GPIO37
	SPI0_MOSI host.Pin = INVALID // GPIO10, GPIO38
	SPI1_CE0  host.Pin = INVALID // GPIO18
	SPI1_CE1  host.Pin = INVALID // GPIO17
	SPI1_CE2  host.Pin = INVALID // GPIO16
	SPI1_CLK  host.Pin = INVALID // GPIO21
	SPI1_MISO host.Pin = INVALID // GPIO19
	SPI1_MOSI host.Pin = INVALID // GPIO20
	SPI2_MISO host.Pin = INVALID // GPIO40
	SPI2_MOSI host.Pin = INVALID // GPIO41
	SPI2_CLK  host.Pin = INVALID // GPIO42
	SPI2_CE0  host.Pin = INVALID // GPIO43
	SPI2_CE1  host.Pin = INVALID // GPIO44
	SPI2_CE2  host.Pin = INVALID // GPIO45
	UART_RXD0 host.Pin = INVALID // GPIO15, GPIO33, GPIO37
	UART_CTS0 host.Pin = INVALID // GPIO16, GPIO30, GPIO39
	UART_CTS1 host.Pin = INVALID // GPIO16, GPIO30
	UART_RTS0 host.Pin = INVALID // GPIO17, GPIO31, GPIO38
	UART_RTS1 host.Pin = INVALID // GPIO17, GPIO31
	UART_TXD0 host.Pin = INVALID // GPIO14, GPIO32, GPIO36
	UART_RXD1 host.Pin = INVALID // GPIO15, GPIO33, GPIO41
	UART_TXD1 host.Pin = INVALID // GPIO14, GPIO32, GPIO40
)

// Number implements host.Pin
func (p Pin) Number() int {
	return int(p)
}

// Function returns the current GPIO pin function.
func (p Pin) Function() Function {
	if gpioMemory32 == nil {
		return Alt5
	}
	// 0x00    RW   GPIO Function Select 0 (GPIO0-9)
	// 0x04    RW   GPIO Function Select 1 (GPIO10-19)
	// 0x08    RW   GPIO Function Select 2 (GPIO20-29)
	// 0x0C    RW   GPIO Function Select 3 (GPIO30-39)
	// 0x10    RW   GPIO Function Select 4 (GPIO40-49)
	// 0x14    RW   GPIO Function Select 5 (GPIO50-53)
	return Function((gpioMemory32[p/10] >> ((p % 10) * 3)) & 7)
}

// In setups a pin as an input and implements host.PinIn.
//
// Specifying a value for pull other than host.PullNoChange causes this
// function to be slightly slower (about 1ms).
//
// Will fail if requesting to change a pin that is set to special functionality.
func (p Pin) In(pull host.Pull) error {
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
		offset := 38 + p/32
		gpioMemory32[offset] = 1 << (p % 32)

		time.Sleep(sleep160cycles)
		gpioMemory32[37] = 0
		gpioMemory32[offset] = 0
	}
	return nil
}

// Read return the current pin level and implements host.PinIn.
//
// This function is very fast. It works even if the pin is set as output.
func (p Pin) Read() host.Level {
	if gpioMemory32 == nil {
		return host.Low
	}
	// 0x34    R    GPIO Pin Level 0 (GPIO0-31)
	// 0x38    R    GPIO Pin Level 1 (GPIO32-53)
	return host.Level((gpioMemory32[13+p/32] & (1 << (p & 31))) != 0)
}

// Out sets a pin as output and implements host.PinOut. The caller should
// immediately call SetLow() or SetHigh() afterward.
//
// Will fail if requesting to change a pin that is set to special functionality.
func (p Pin) Out() error {
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
func (p Pin) Set(l host.Level) {
	if gpioMemory32 != nil {
		// 0x1C    W    GPIO Pin Output Set 0 (GPIO0-31)
		// 0x20    W    GPIO Pin Output Set 1 (GPIO32-53)
		base := 7 + p/32
		if l == host.Low {
			// 0x28    W    GPIO Pin Output Clear 0 (GPIO0-31)
			// 0x2C    W    GPIO Pin Output Clear 1 (GPIO32-53)
			base += 3
		}
		gpioMemory32[base] = 1 << (p & 31)
	}
}

// setFunction changes the GPIO pin function.
//
// Returns false if the pin was in AltN. Only accepts In and Out
func (p Pin) setFunction(f Function) bool {
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
	off := p / 10
	shift := (p % 10) * 3
	gpioMemory32[off] = (gpioMemory32[off] &^ (7 << shift)) | (uint32(f) << shift)
	return true
}

// GetPin returns a pin from its name.
func GetPin(name string) Pin {
	switch name {
	case "GPIO0":
		return GPIO0
	case "GPIO1":
		return GPIO1
	case "GPIO2":
		return GPIO2
	case "GPIO3":
		return GPIO3
	case "GPIO4":
		return GPIO4
	case "GPIO5":
		return GPIO5
	case "GPIO6":
		return GPIO6
	case "GPIO7":
		return GPIO7
	case "GPIO8":
		return GPIO8
	case "GPIO9":
		return GPIO9
	case "GPIO10":
		return GPIO10
	case "GPIO11":
		return GPIO11
	case "GPIO12":
		return GPIO12
	case "GPIO13":
		return GPIO13
	case "GPIO14":
		return GPIO14
	case "GPIO15":
		return GPIO15
	case "GPIO16":
		return GPIO16
	case "GPIO17":
		return GPIO17
	case "GPIO18":
		return GPIO18
	case "GPIO19":
		return GPIO19
	case "GPIO20":
		return GPIO20
	case "GPIO21":
		return GPIO21
	case "GPIO22":
		return GPIO22
	case "GPIO23":
		return GPIO23
	case "GPIO24":
		return GPIO24
	case "GPIO25":
		return GPIO25
	case "GPIO26":
		return GPIO26
	case "GPIO27":
		return GPIO27
	case "GPIO28":
		return GPIO28
	case "GPIO29":
		return GPIO29
	case "GPIO30":
		return GPIO30
	case "GPIO31":
		return GPIO31
	case "GPIO32":
		return GPIO32
	case "GPIO33":
		return GPIO33
	case "GPIO34":
		return GPIO34
	case "GPIO35":
		return GPIO35
	case "GPIO36":
		return GPIO36
	case "GPIO37":
		return GPIO37
	case "GPIO38":
		return GPIO38
	case "GPIO39":
		return GPIO39
	case "GPIO40":
		return GPIO40
	case "GPIO41":
		return GPIO41
	case "GPIO42":
		return GPIO42
	case "GPIO43":
		return GPIO43
	case "GPIO44":
		return GPIO44
	case "GPIO45":
		return GPIO45
	case "GPIO46":
		return GPIO46
	case "GPIO47":
		return GPIO47
	case "GPIO48":
		return GPIO48
	case "GPIO49":
		return GPIO49
	case "GPIO50":
		return GPIO50
	case "GPIO51":
		return GPIO51
	case "GPIO52":
		return GPIO52
	case "GPIO53":
		return GPIO53
	case "GROUND":
		return GROUND
	case "V3_3":
		return V3_3
	case "V5":
		return V5
	default:
		return INVALID
	}
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

func setIfAlt0(p Pin, special *host.Pin) {
	if p.Function() == Alt0 {
		if *special != INVALID {
			//fmt.Printf("%s and %s have same functionality\n", p, *special)
		}
		*special = p
	}
}

func setIfAlt(p Pin, special0 *host.Pin, special1 *host.Pin, special2 *host.Pin, special3 *host.Pin, special4 *host.Pin, special5 *host.Pin) {
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
				//log.Printf("%s and %s have same functionality\n", p, *special2)
			}
			*special2 = p
		}
	case Alt3:
		if special3 != nil {
			if *special3 != INVALID {
				//log.Printf("%s and %s have same functionality\n", p, *special3)
			}
			*special3 = p
		}
	case Alt4:
		if special4 != nil {
			if *special4 != INVALID {
				//log.Printf("%s and %s have same functionality\n", p, *special4)
			}
			*special4 = p
		}
	case Alt5:
		if special5 != nil {
			if *special5 != INVALID {
				//log.Printf("%s and %s have same functionality\n", p, *special5)
			}
			*special5 = p
		}
	}
}

func init() {
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
		IR_IN = Pin(in)
	}
	if out != -1 {
		IR_OUT = Pin(out)
	}
}
