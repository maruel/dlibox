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

var Pins = [116]Pin{
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
	{number: 10, name: "PC0", defaultPull: host.Float},
	{number: 11, name: "PC1", defaultPull: host.Float},
	{number: 12, name: "PC2", defaultPull: host.Float},
	{number: 13, name: "PC3", defaultPull: host.Up},
	{number: 14, name: "PC4", defaultPull: host.Up},
	{number: 15, name: "PC5", defaultPull: host.Float},
	{number: 16, name: "PC6", defaultPull: host.Up},
	{number: 17, name: "PC7", defaultPull: host.Up},
	{number: 18, name: "PC8", defaultPull: host.Float},
	{number: 19, name: "PC9", defaultPull: host.Float},
	{number: 20, name: "PC10", defaultPull: host.Float},
	{number: 21, name: "PC11", defaultPull: host.Float},
	{number: 22, name: "PC12", defaultPull: host.Float},
	{number: 23, name: "PC13", defaultPull: host.Float},
	{number: 24, name: "PC14", defaultPull: host.Float},
	{number: 25, name: "PC15", defaultPull: host.Float},
	{number: 26, name: "PC16", defaultPull: host.Float},
	{number: 27, name: "PD0", defaultPull: host.Float},
	{number: 28, name: "PD1", defaultPull: host.Float},
	{number: 29, name: "PD2", defaultPull: host.Float},
	{number: 30, name: "PD3", defaultPull: host.Float},
	{number: 31, name: "PD4", defaultPull: host.Float},
	{number: 32, name: "PD5", defaultPull: host.Float},
	{number: 33, name: "PD6", defaultPull: host.Float},
	{number: 34, name: "PD7", defaultPull: host.Float},
	{number: 35, name: "PD8", defaultPull: host.Float},
	{number: 36, name: "PD9", defaultPull: host.Float},
	{number: 37, name: "PD10", defaultPull: host.Float},
	{number: 38, name: "PD11", defaultPull: host.Float},
	{number: 39, name: "PD12", defaultPull: host.Float},
	{number: 40, name: "PD13", defaultPull: host.Float},
	{number: 41, name: "PD14", defaultPull: host.Float},
	{number: 42, name: "PD15", defaultPull: host.Float},
	{number: 43, name: "PD16", defaultPull: host.Float},
	{number: 44, name: "PD17", defaultPull: host.Float},
	{number: 45, name: "PD18", defaultPull: host.Float},
	{number: 46, name: "PD19", defaultPull: host.Float},
	{number: 47, name: "PD20", defaultPull: host.Float},
	{number: 48, name: "PD21", defaultPull: host.Float},
	{number: 49, name: "PD22", defaultPull: host.Float},
	{number: 50, name: "PD23", defaultPull: host.Float},
	{number: 51, name: "PD24", defaultPull: host.Float},
	{number: 52, name: "PE0", defaultPull: host.Float},
	{number: 53, name: "PE1", defaultPull: host.Float},
	{number: 54, name: "PE2", defaultPull: host.Float},
	{number: 55, name: "PE3", defaultPull: host.Float},
	{number: 56, name: "PE4", defaultPull: host.Float},
	{number: 57, name: "PE5", defaultPull: host.Float},
	{number: 58, name: "PE6", defaultPull: host.Float},
	{number: 59, name: "PE7", defaultPull: host.Float},
	{number: 60, name: "PE8", defaultPull: host.Float},
	{number: 61, name: "PE9", defaultPull: host.Float},
	{number: 62, name: "PE10", defaultPull: host.Float},
	{number: 63, name: "PE11", defaultPull: host.Float},
	{number: 64, name: "PE12", defaultPull: host.Float},
	{number: 65, name: "PE13", defaultPull: host.Float},
	{number: 66, name: "PE14", defaultPull: host.Float},
	{number: 67, name: "PE15", defaultPull: host.Float},
	{number: 68, name: "PE16", defaultPull: host.Float},
	{number: 69, name: "PE17", defaultPull: host.Float},
	{number: 70, name: "PF0", defaultPull: host.Float},
	{number: 71, name: "PF1", defaultPull: host.Float},
	{number: 72, name: "PF2", defaultPull: host.Float},
	{number: 73, name: "PF3", defaultPull: host.Float},
	{number: 74, name: "PF4", defaultPull: host.Float},
	{number: 75, name: "PF5", defaultPull: host.Float},
	{number: 76, name: "PF6", defaultPull: host.Float},
	{number: 77, name: "PG0", defaultPull: host.Float},
	{number: 78, name: "PG1", defaultPull: host.Float},
	{number: 79, name: "PG2", defaultPull: host.Float},
	{number: 80, name: "PG3", defaultPull: host.Float},
	{number: 81, name: "PG4", defaultPull: host.Float},
	{number: 82, name: "PG5", defaultPull: host.Float},
	{number: 83, name: "PG6", defaultPull: host.Float},
	{number: 84, name: "PG7", defaultPull: host.Float},
	{number: 85, name: "PG8", defaultPull: host.Float},
	{number: 86, name: "PG9", defaultPull: host.Float},
	{number: 87, name: "PG10", defaultPull: host.Float},
	{number: 88, name: "PG11", defaultPull: host.Float},
	{number: 89, name: "PG12", defaultPull: host.Float},
	{number: 90, name: "PG13", defaultPull: host.Float},
	{number: 91, name: "PH0", defaultPull: host.Float},
	{number: 92, name: "PH1", defaultPull: host.Float},
	{number: 93, name: "PH2", defaultPull: host.Float},
	{number: 94, name: "PH3", defaultPull: host.Float},
	{number: 95, name: "PH4", defaultPull: host.Float},
	{number: 96, name: "PH5", defaultPull: host.Float},
	{number: 97, name: "PH6", defaultPull: host.Float},
	{number: 98, name: "PH7", defaultPull: host.Float},
	{number: 99, name: "PH8", defaultPull: host.Float},
	{number: 100, name: "PH9", defaultPull: host.Float},
	{number: 101, name: "PH10", defaultPull: host.Float},
	{number: 102, name: "PH11", defaultPull: host.Float},
	{number: 103, name: "PL0", defaultPull: host.Float},
	{number: 104, name: "PL1", defaultPull: host.Float},
	{number: 105, name: "PL2", defaultPull: host.Float},
	{number: 106, name: "PL3", defaultPull: host.Float},
	{number: 107, name: "PL4", defaultPull: host.Float},
	{number: 108, name: "PL5", defaultPull: host.Float},
	{number: 109, name: "PL6", defaultPull: host.Float},
	{number: 110, name: "PL7", defaultPull: host.Float},
	{number: 111, name: "PL8", defaultPull: host.Float},
	{number: 112, name: "PL9", defaultPull: host.Float},
	{number: 113, name: "PL10", defaultPull: host.Float},
	{number: 114, name: "PL11", defaultPull: host.Float},
	{number: 115, name: "PL12", defaultPull: host.Float},
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
	PL0  host.PinIO //
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

func init() {
	PB0 = &Pins[0]
	PB1 = &Pins[1]
	PB2 = &Pins[2]
	PB3 = &Pins[3]
	PB4 = &Pins[4]
	PB5 = &Pins[5]
	PB6 = &Pins[6]
	PB7 = &Pins[7]
	PB8 = &Pins[8]
	PB9 = &Pins[9]
	PC0 = &Pins[10]
	PC1 = &Pins[11]
	PC2 = &Pins[12]
	PC3 = &Pins[13]
	PC4 = &Pins[14]
	PC5 = &Pins[15]
	PC6 = &Pins[16]
	PC7 = &Pins[17]
	PC8 = &Pins[18]
	PC9 = &Pins[19]
	PC10 = &Pins[20]
	PC11 = &Pins[21]
	PC12 = &Pins[22]
	PC13 = &Pins[23]
	PC14 = &Pins[24]
	PC15 = &Pins[25]
	PC16 = &Pins[26]
	PD0 = &Pins[27]
	PD1 = &Pins[28]
	PD2 = &Pins[29]
	PD3 = &Pins[30]
	PD4 = &Pins[31]
	PD5 = &Pins[32]
	PD6 = &Pins[33]
	PD7 = &Pins[34]
	PD8 = &Pins[35]
	PD9 = &Pins[36]
	PD10 = &Pins[37]
	PD11 = &Pins[38]
	PD12 = &Pins[39]
	PD13 = &Pins[40]
	PD14 = &Pins[41]
	PD15 = &Pins[42]
	PD16 = &Pins[43]
	PD17 = &Pins[44]
	PD18 = &Pins[45]
	PD19 = &Pins[46]
	PD20 = &Pins[47]
	PD21 = &Pins[48]
	PD22 = &Pins[49]
	PD23 = &Pins[50]
	PD24 = &Pins[51]
	PE0 = &Pins[52]
	PE1 = &Pins[53]
	PE2 = &Pins[54]
	PE3 = &Pins[55]
	PE4 = &Pins[56]
	PE5 = &Pins[57]
	PE6 = &Pins[58]
	PE7 = &Pins[59]
	PE8 = &Pins[60]
	PE9 = &Pins[61]
	PE10 = &Pins[62]
	PE11 = &Pins[63]
	PE12 = &Pins[64]
	PE13 = &Pins[65]
	PE14 = &Pins[66]
	PE15 = &Pins[67]
	PE16 = &Pins[68]
	PE17 = &Pins[69]
	PF0 = &Pins[70]
	PF1 = &Pins[71]
	PF2 = &Pins[72]
	PF3 = &Pins[73]
	PF4 = &Pins[74]
	PF5 = &Pins[75]
	PF6 = &Pins[76]
	PG0 = &Pins[77]
	PG1 = &Pins[78]
	PG2 = &Pins[79]
	PG3 = &Pins[80]
	PG4 = &Pins[81]
	PG5 = &Pins[82]
	PG6 = &Pins[83]
	PG7 = &Pins[84]
	PG8 = &Pins[85]
	PG9 = &Pins[86]
	PG10 = &Pins[87]
	PG11 = &Pins[88]
	PG12 = &Pins[89]
	PG13 = &Pins[90]
	PH0 = &Pins[91]
	PH1 = &Pins[92]
	PH2 = &Pins[93]
	PH3 = &Pins[94]
	PH4 = &Pins[95]
	PH5 = &Pins[96]
	PH6 = &Pins[97]
	PH7 = &Pins[98]
	PH8 = &Pins[99]
	PH9 = &Pins[100]
	PH10 = &Pins[101]
	PH11 = &Pins[102]
	PL0 = &Pins[103]
	PL1 = &Pins[104]
	PL2 = &Pins[105]
	PL3 = &Pins[106]
	PL4 = &Pins[107]
	PL5 = &Pins[108]
	PL6 = &Pins[109]
	PL7 = &Pins[110]
	PL8 = &Pins[111]
	PL9 = &Pins[112]
	PL10 = &Pins[113]
	PL11 = &Pins[114]
	PL12 = &Pins[115]
}
