// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package allwinner

import (
	"errors"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/maruel/dlibox/go/pio/host"
	"github.com/maruel/dlibox/go/pio/host/internal/gpiomem"
	"github.com/maruel/dlibox/go/pio/host/internal/sysfs"
)

var Pins = [10]Pin{
	{number: 0, name: "PB0", defaultPull: host.Float},
	{number: 1, name: "PB1", defaultPull: host.Float},
	{number: 2, name: "PB2", defaultPull: host.Float},
	{number: 3, name: "PB3", defaultPull: host.Float},
	{number: 4, name: "PB4", defaultPull: host.Float},
	{number: 5, name: "PB5", defaultPull: host.Float},
	{number: 6, name: "PB6", defaultPull: host.Float},
	{number: 7, name: "PB7", defaultPull: host.Float},
	{number: 8, name: "PB8", defaultPull: host.Float},
	{number: 9, name: "PB9", defaultPull: host.Float},
	// TODO(maruel): add rest, only a few PC pins have a default pull up.
}

// Functional is pins.Functional on this CPU.
var Functional = map[string]host.Pin{}

// Page 23~24
// Each pin supports 6 functions.

type Pin struct {
	number      int
	name        string
	defaultPull host.Pull
	edge        *sysfs.Pin // Mutable, set once, then never set back to nil
}

// http://forum.pine64.org/showthread.php?tid=474
// about number calculation.
var (
	PB0  host.PinIO // Z UART2_TX, -, JTAG_MS0, -, PB_EINT0
	PB1  host.PinIO // Z, UART2_RX, -, JTAG_CK0, SIM_PWREN, PB_EINT1
	PB2  host.PinIO //
	PB3  host.PinIO //
	PB4  host.PinIO //
	PB5  host.PinIO //
	PB6  host.PinIO // AIF2_DOUT, PCM0_DOUT, -, SIM_RST
	PB7  host.PinIO // AIF2_DIN, PCM0_DIN, -, SIM_DET
	PB8  host.PinIO //
	PB9  host.PinIO //
	PC0  host.PinIO //
	PC1  host.PinIO //
	PC2  host.PinIO //
	PC3  host.PinIO //
	PC4  host.PinIO //
	PC5  host.PinIO //
	PC6  host.PinIO //
	PC7  host.PinIO //
	PC8  host.PinIO //
	PC9  host.PinIO //
	PC10 host.PinIO //
	PC11 host.PinIO //
	PC12 host.PinIO //
	PC13 host.PinIO //
	PC14 host.PinIO //
	PC15 host.PinIO //
	PC16 host.PinIO //
	PD0  host.PinIO //
	PD1  host.PinIO //
	PD2  host.PinIO //
	PD3  host.PinIO //
	PD4  host.PinIO //
	PD5  host.PinIO //
	PD6  host.PinIO //
	PD7  host.PinIO //
	PD8  host.PinIO //
	PD9  host.PinIO //
	PD10 host.PinIO //
	PD11 host.PinIO //
	PD12 host.PinIO //
	PD13 host.PinIO //
	PD14 host.PinIO //
	PD15 host.PinIO //
	PD16 host.PinIO //
	PD17 host.PinIO //
	PD18 host.PinIO //
	PD19 host.PinIO //
	PD20 host.PinIO //
	PD21 host.PinIO //
	PD22 host.PinIO //
	PD23 host.PinIO //
	PD24 host.PinIO //
	PE0  host.PinIO //
	PE1  host.PinIO //
	PE2  host.PinIO //
	PE3  host.PinIO //
	PE4  host.PinIO //
	PE5  host.PinIO //
	PE6  host.PinIO //
	PE7  host.PinIO //
	PE8  host.PinIO //
	PE9  host.PinIO //
	PE10 host.PinIO //
	PE11 host.PinIO //
	PE12 host.PinIO //
	PE13 host.PinIO //
	PE14 host.PinIO //
	PE15 host.PinIO //
	PE16 host.PinIO //
	PE17 host.PinIO //
	PF0  host.PinIO //
	PF1  host.PinIO //
	PF2  host.PinIO //
	PF3  host.PinIO //
	PF4  host.PinIO //
	PF5  host.PinIO //
	PF6  host.PinIO //
	PG0  host.PinIO //
	PG1  host.PinIO //
	PG2  host.PinIO //
	PG3  host.PinIO //
	PG4  host.PinIO //
	PG5  host.PinIO //
	PG6  host.PinIO //
	PG7  host.PinIO //
	PG8  host.PinIO //
	PG9  host.PinIO //
	PG10 host.PinIO //
	PG11 host.PinIO //
	PG12 host.PinIO //
	PG13 host.PinIO //
	PH0  host.PinIO //
	PH1  host.PinIO //
	PH2  host.PinIO //
	PH3  host.PinIO //
	PH4  host.PinIO //
	PH5  host.PinIO //
	PH6  host.PinIO //
	PH7  host.PinIO //
	PH8  host.PinIO //
	PH9  host.PinIO //
	PH10 host.PinIO //
	PH11 host.PinIO //
	PL1  host.PinIO //
	PL2  host.PinIO //
	PL3  host.PinIO //
	PL4  host.PinIO //
	PL5  host.PinIO //
	PL6  host.PinIO //
	PL7  host.PinIO //
	PL8  host.PinIO //
	PL9  host.PinIO //
	PL10 host.PinIO //
	PL11 host.PinIO //
	PL12 host.PinIO //
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

func (p *Pin) Function() string {
	// TODO(maruel): Add; can be disabled, unlike broadcom.
	return ""
}

func (p *Pin) In(pull host.Pull) error {
	return errors.New("implement me")
}

func (p *Pin) Read() host.Level {
	return host.Low
}

func (p *Pin) Edges() (chan host.Level, error) {
	return nil, errors.New("implement me")
}

func (p *Pin) Out() error {
	return errors.New("implement me")
}

func (p *Pin) Set(l host.Level) {
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
	return nil
}

//

// mapping excludes functions in and out.
// Datasheet, page 23.
// http://files.pine64.org/doc/datasheet/pine64/A64_Datasheet_V1.1.pdf
var mapping = [][5]string{
	{"UART2_TX", "", "JTAG_MS0", "", "PB_EINT0"},                    // PB0
	{"UART2_RX", "", "JTAG_CK0", "SIM_PWREN", "PB_EINT1"},           // PB1
	{"UART2_RTS", "", "JTAG_DO0", "SIM_VPPEN", "PB_EINT2"},          // PB2
	{"UART2_CTS", "I2S0_MCLK", "JTAG_DI0", "SIM_VPPPP", "PB_EINT3"}, // PB3
	{"AIF2_SYNC", "PCM0_SYNC", "", "SIM_CLK", "PB_EINT4"},           // PB4
	{"AIF2_BCLK", "PCM0_BCLK", "", "SIM_DATA", "PB_EINT5"},          // PB5
	{"AIF2_DOUT", "PCM0_DOUT", "", "SIM_RST", "PB_EINT6"},           // PB6
	{"AIF2_DIN", "PCM0_DIN", "", "SIM_DET", "PB_EINT7"},             // PB7
	{"", "", "", "UART0_TX", "PB_EINT8"},                            // PB8
	{"", "", "", "UART0_RX", "PB_EINT9"},                            // PB9
	// TODO(maruel): Add the rest.
}

// getBaseAddress queries the virtual file system to retrieve the base address
// of the GPIO registers.
//
// Defaults to 0x01C20800 as per datasheet if could query the file system.
func getBaseAddress() uint64 {
	base := uint64(0x01C20800)
	link, err := os.Readlink("/sys/bus/platform/drivers/sun50i-pinctrl/driver")
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

// http://files.pine64.org/doc/datasheet/pine64/A64_Datasheet_V1.1.pdf
// Pulls

// http://files.pine64.org/doc/datasheet/pine64/Allwinner_A64_User_Manual_V1.0.pdf
// Page 376 GPIO PB to PH:
// 1 to 7 means PB to PH.
// Pn_CFG0 n*0x24+0x00  Port n Configure Register 0 (n from 1 to 7)
// Pn_CFG1 n*0x24+0x04  Port n Configure Register 1 (n from 1 to 7)
// Pn_CFG2 n*0x24+0x08  Port n Configure Register 2 (n from 1 to 7)
// Pn_CFG3 n*0x24+0x0C  Port n Configure Register 3 (n from 1 to 7)
// Pn_DAT  n*0x24+0x10  Port n Data Register (n from 1 to 7)
// Pn_DRV0 n*0x24+0x14  Port n Multi-Driving Register 0 (n from 1 to 7)
// Pn_DRV1 n*0x24+0x18  Port n Multi-Driving Register 1 (n from 1 to 7)
// Pn_PUL0 n*0x24+0x1C  Port n Pull  Register 0 (n from 1 to 7)
// Pn_PUL1 n*0x24+0x20  Port n Pull Register 1 (n from 1 to 7)
var gpioMemory32 []uint32

// Page 73 for memory mapping overview.
// Page 194 for PWM.
// Page 230 for crypto engine.
// Page 278 audio including ADC.
// Page 376 GPIO PB to PH
// Page 410 GPIO PL
// Page 536 I²C (TWI)
// Page 545 SPI
// Page 560 UART
// Page 621 I2S/PCM

var _ host.PinIO = &Pin{}
