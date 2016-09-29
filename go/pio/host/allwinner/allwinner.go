// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package allwinner

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/maruel/dlibox/go/pio"
	"github.com/maruel/dlibox/go/pio/host/internal"
	"github.com/maruel/dlibox/go/pio/host/internal/gpiomem"
	"github.com/maruel/dlibox/go/pio/host/sysfs"
	"github.com/maruel/dlibox/go/pio/protocols/analog"
	"github.com/maruel/dlibox/go/pio/protocols/gpio"
	"github.com/maruel/dlibox/go/pio/protocols/pins"
)

var (
	X32KFOUT gpio.PinIO   = &pins.BasicPin{Name: "X32KFOUT"} // Clock output of 32Khz crystal
	KEY_ADC  analog.PinIO = &pins.BasicPin{Name: "KEY_ADC"}  // 6 bits resolution ADC for key application; can work up to 250Hz conversion rate; reference voltage is 2.0V
	EAROUTP  analog.PinIO = &pins.BasicPin{Name: "EAROUTP"}  // Earpiece amplifier negative differential output
	EAROUTN  analog.PinIO = &pins.BasicPin{Name: "EAROUTN"}  // Earpiece amplifier positive differential output
)

// 0x24/4 = 9

var Pins = []Pin{
	{index: 0, group: 9 * 1, offset: 0, name: "PB0", defaultPull: gpio.Float},
	{index: 1, group: 9 * 1, offset: 1, name: "PB1", defaultPull: gpio.Float},
	{index: 2, group: 9 * 1, offset: 2, name: "PB2", defaultPull: gpio.Float},
	{index: 3, group: 9 * 1, offset: 3, name: "PB3", defaultPull: gpio.Float},
	{index: 4, group: 9 * 1, offset: 4, name: "PB4", defaultPull: gpio.Float},
	{index: 5, group: 9 * 1, offset: 5, name: "PB5", defaultPull: gpio.Float},
	{index: 6, group: 9 * 1, offset: 6, name: "PB6", defaultPull: gpio.Float},
	{index: 7, group: 9 * 1, offset: 7, name: "PB7", defaultPull: gpio.Float},
	{index: 8, group: 9 * 1, offset: 8, name: "PB8", defaultPull: gpio.Float},
	{index: 9, group: 9 * 1, offset: 9, name: "PB9", defaultPull: gpio.Float},
	{index: 10, group: 9 * 2, offset: 0, name: "PC0", defaultPull: gpio.Float},
	{index: 11, group: 9 * 2, offset: 1, name: "PC1", defaultPull: gpio.Float},
	{index: 12, group: 9 * 2, offset: 2, name: "PC2", defaultPull: gpio.Float},
	{index: 13, group: 9 * 2, offset: 3, name: "PC3", defaultPull: gpio.Up},
	{index: 14, group: 9 * 2, offset: 4, name: "PC4", defaultPull: gpio.Up},
	{index: 15, group: 9 * 2, offset: 5, name: "PC5", defaultPull: gpio.Float},
	{index: 16, group: 9 * 2, offset: 6, name: "PC6", defaultPull: gpio.Up},
	{index: 17, group: 9 * 2, offset: 7, name: "PC7", defaultPull: gpio.Up},
	{index: 18, group: 9 * 2, offset: 8, name: "PC8", defaultPull: gpio.Float},
	{index: 19, group: 9 * 2, offset: 9, name: "PC9", defaultPull: gpio.Float},
	{index: 20, group: 9 * 2, offset: 10, name: "PC10", defaultPull: gpio.Float},
	{index: 21, group: 9 * 2, offset: 11, name: "PC11", defaultPull: gpio.Float},
	{index: 22, group: 9 * 2, offset: 12, name: "PC12", defaultPull: gpio.Float},
	{index: 23, group: 9 * 2, offset: 13, name: "PC13", defaultPull: gpio.Float},
	{index: 24, group: 9 * 2, offset: 14, name: "PC14", defaultPull: gpio.Float},
	{index: 25, group: 9 * 2, offset: 15, name: "PC15", defaultPull: gpio.Float},
	{index: 26, group: 9 * 2, offset: 16, name: "PC16", defaultPull: gpio.Float},
	{index: 27, group: 9 * 3, offset: 0, name: "PD0", defaultPull: gpio.Float},
	{index: 28, group: 9 * 3, offset: 1, name: "PD1", defaultPull: gpio.Float},
	{index: 29, group: 9 * 3, offset: 2, name: "PD2", defaultPull: gpio.Float},
	{index: 30, group: 9 * 3, offset: 3, name: "PD3", defaultPull: gpio.Float},
	{index: 31, group: 9 * 3, offset: 4, name: "PD4", defaultPull: gpio.Float},
	{index: 32, group: 9 * 3, offset: 5, name: "PD5", defaultPull: gpio.Float},
	{index: 33, group: 9 * 3, offset: 6, name: "PD6", defaultPull: gpio.Float},
	{index: 34, group: 9 * 3, offset: 7, name: "PD7", defaultPull: gpio.Float},
	{index: 35, group: 9 * 3, offset: 8, name: "PD8", defaultPull: gpio.Float},
	{index: 36, group: 9 * 3, offset: 9, name: "PD9", defaultPull: gpio.Float},
	{index: 37, group: 9 * 3, offset: 10, name: "PD10", defaultPull: gpio.Float},
	{index: 38, group: 9 * 3, offset: 11, name: "PD11", defaultPull: gpio.Float},
	{index: 39, group: 9 * 3, offset: 12, name: "PD12", defaultPull: gpio.Float},
	{index: 40, group: 9 * 3, offset: 13, name: "PD13", defaultPull: gpio.Float},
	{index: 41, group: 9 * 3, offset: 14, name: "PD14", defaultPull: gpio.Float},
	{index: 42, group: 9 * 3, offset: 15, name: "PD15", defaultPull: gpio.Float},
	{index: 43, group: 9 * 3, offset: 16, name: "PD16", defaultPull: gpio.Float},
	{index: 44, group: 9 * 3, offset: 17, name: "PD17", defaultPull: gpio.Float},
	{index: 45, group: 9 * 3, offset: 18, name: "PD18", defaultPull: gpio.Float},
	{index: 46, group: 9 * 3, offset: 19, name: "PD19", defaultPull: gpio.Float},
	{index: 47, group: 9 * 3, offset: 20, name: "PD20", defaultPull: gpio.Float},
	{index: 48, group: 9 * 3, offset: 21, name: "PD21", defaultPull: gpio.Float},
	{index: 49, group: 9 * 3, offset: 22, name: "PD22", defaultPull: gpio.Float},
	{index: 50, group: 9 * 3, offset: 23, name: "PD23", defaultPull: gpio.Float},
	{index: 51, group: 9 * 3, offset: 24, name: "PD24", defaultPull: gpio.Float},
	{index: 52, group: 9 * 4, offset: 0, name: "PE0", defaultPull: gpio.Float},
	{index: 53, group: 9 * 4, offset: 1, name: "PE1", defaultPull: gpio.Float},
	{index: 54, group: 9 * 4, offset: 2, name: "PE2", defaultPull: gpio.Float},
	{index: 55, group: 9 * 4, offset: 3, name: "PE3", defaultPull: gpio.Float},
	{index: 56, group: 9 * 4, offset: 4, name: "PE4", defaultPull: gpio.Float},
	{index: 57, group: 9 * 4, offset: 5, name: "PE5", defaultPull: gpio.Float},
	{index: 58, group: 9 * 4, offset: 6, name: "PE6", defaultPull: gpio.Float},
	{index: 59, group: 9 * 4, offset: 7, name: "PE7", defaultPull: gpio.Float},
	{index: 60, group: 9 * 4, offset: 8, name: "PE8", defaultPull: gpio.Float},
	{index: 61, group: 9 * 4, offset: 9, name: "PE9", defaultPull: gpio.Float},
	{index: 62, group: 9 * 4, offset: 10, name: "PE10", defaultPull: gpio.Float},
	{index: 63, group: 9 * 4, offset: 11, name: "PE11", defaultPull: gpio.Float},
	{index: 64, group: 9 * 4, offset: 12, name: "PE12", defaultPull: gpio.Float},
	{index: 65, group: 9 * 4, offset: 13, name: "PE13", defaultPull: gpio.Float},
	{index: 66, group: 9 * 4, offset: 14, name: "PE14", defaultPull: gpio.Float},
	{index: 67, group: 9 * 4, offset: 15, name: "PE15", defaultPull: gpio.Float},
	{index: 68, group: 9 * 4, offset: 16, name: "PE16", defaultPull: gpio.Float},
	{index: 69, group: 9 * 4, offset: 17, name: "PE17", defaultPull: gpio.Float},
	{index: 70, group: 9 * 5, offset: 0, name: "PF0", defaultPull: gpio.Float},
	{index: 71, group: 9 * 5, offset: 1, name: "PF1", defaultPull: gpio.Float},
	{index: 72, group: 9 * 5, offset: 2, name: "PF2", defaultPull: gpio.Float},
	{index: 73, group: 9 * 5, offset: 3, name: "PF3", defaultPull: gpio.Float},
	{index: 74, group: 9 * 5, offset: 4, name: "PF4", defaultPull: gpio.Float},
	{index: 75, group: 9 * 5, offset: 5, name: "PF5", defaultPull: gpio.Float},
	{index: 76, group: 9 * 5, offset: 6, name: "PF6", defaultPull: gpio.Float},
	{index: 77, group: 9 * 6, offset: 0, name: "PG0", defaultPull: gpio.Float},
	{index: 78, group: 9 * 6, offset: 1, name: "PG1", defaultPull: gpio.Float},
	{index: 79, group: 9 * 6, offset: 2, name: "PG2", defaultPull: gpio.Float},
	{index: 80, group: 9 * 6, offset: 3, name: "PG3", defaultPull: gpio.Float},
	{index: 81, group: 9 * 6, offset: 4, name: "PG4", defaultPull: gpio.Float},
	{index: 82, group: 9 * 6, offset: 5, name: "PG5", defaultPull: gpio.Float},
	{index: 83, group: 9 * 6, offset: 6, name: "PG6", defaultPull: gpio.Float},
	{index: 84, group: 9 * 6, offset: 7, name: "PG7", defaultPull: gpio.Float},
	{index: 85, group: 9 * 6, offset: 8, name: "PG8", defaultPull: gpio.Float},
	{index: 86, group: 9 * 6, offset: 9, name: "PG9", defaultPull: gpio.Float},
	{index: 87, group: 9 * 6, offset: 10, name: "PG10", defaultPull: gpio.Float},
	{index: 88, group: 9 * 6, offset: 11, name: "PG11", defaultPull: gpio.Float},
	{index: 89, group: 9 * 6, offset: 12, name: "PG12", defaultPull: gpio.Float},
	{index: 90, group: 9 * 6, offset: 13, name: "PG13", defaultPull: gpio.Float},
	{index: 91, group: 9 * 7, offset: 0, name: "PH0", defaultPull: gpio.Float},
	{index: 92, group: 9 * 7, offset: 1, name: "PH1", defaultPull: gpio.Float},
	{index: 93, group: 9 * 7, offset: 2, name: "PH2", defaultPull: gpio.Float},
	{index: 94, group: 9 * 7, offset: 3, name: "PH3", defaultPull: gpio.Float},
	{index: 95, group: 9 * 7, offset: 4, name: "PH4", defaultPull: gpio.Float},
	{index: 96, group: 9 * 7, offset: 5, name: "PH5", defaultPull: gpio.Float},
	{index: 97, group: 9 * 7, offset: 6, name: "PH6", defaultPull: gpio.Float},
	{index: 98, group: 9 * 7, offset: 7, name: "PH7", defaultPull: gpio.Float},
	{index: 99, group: 9 * 7, offset: 8, name: "PH8", defaultPull: gpio.Float},
	{index: 100, group: 9 * 7, offset: 9, name: "PH9", defaultPull: gpio.Float},
	{index: 101, group: 9 * 7, offset: 10, name: "PH10", defaultPull: gpio.Float},
	{index: 102, group: 9 * 7, offset: 11, name: "PH11", defaultPull: gpio.Float},
}

var functional = map[string]gpio.PinIO{
	/*
		"AIF2_BCLK":    pins.INVALID,
		"AIF2_DIN":     pins.INVALID,
		"AIF2_DOUT":    pins.INVALID,
		"AIF2_SYNC":    pins.INVALID,
		"AIF3_BCLK":    pins.INVALID,
		"AIF3_DIN":     pins.INVALID,
		"AIF3_DOUT":    pins.INVALID,
		"AIF3_SYNC":    pins.INVALID,
		"CCIR_CLK":     pins.INVALID,
		"CCIR_D0":      pins.INVALID,
		"CCIR_D1":      pins.INVALID,
		"CCIR_D2":      pins.INVALID,
		"CCIR_D3":      pins.INVALID,
		"CCIR_D4":      pins.INVALID,
		"CCIR_D5":      pins.INVALID,
		"CCIR_D6":      pins.INVALID,
		"CCIR_D7":      pins.INVALID,
		"CCIR_DE":      pins.INVALID,
		"CCIR_HSYNC":   pins.INVALID,
		"CCIR_VSYNC":   pins.INVALID,
		"CSI_D0":       pins.INVALID,
		"CSI_D1":       pins.INVALID,
		"CSI_D2":       pins.INVALID,
		"CSI_D3":       pins.INVALID,
		"CSI_D4":       pins.INVALID,
		"CSI_D5":       pins.INVALID,
		"CSI_D6":       pins.INVALID,
		"CSI_D7":       pins.INVALID,
		"CSI_HSYNC":    pins.INVALID,
		"CSI_MCLK":     pins.INVALID,
		"CSI_PCLK":     pins.INVALID,
		"CSI_SCL":      pins.INVALID,
		"CSI_SDA":      pins.INVALID,
		"CSI_VSYNC":    pins.INVALID,
	*/
	"I2C0_SCL":  pins.INVALID,
	"I2C0_SDA":  pins.INVALID,
	"I2C1_SCL":  pins.INVALID,
	"I2C1_SDA":  pins.INVALID,
	"I2C2_SCL":  pins.INVALID,
	"I2C2_SDA":  pins.INVALID,
	"I2S0_MCLK": pins.INVALID,
	/*
		"JTAG_CK0":     pins.INVALID,
		"JTAG_CK1":     pins.INVALID,
		"JTAG_DI0":     pins.INVALID,
		"JTAG_DI1":     pins.INVALID,
		"JTAG_DO0":     pins.INVALID,
		"JTAG_DO1":     pins.INVALID,
		"JTAG_MS0":     pins.INVALID,
		"JTAG_MS1":     pins.INVALID,
		"LCD_CLK":      pins.INVALID,
		"LCD_D10":      pins.INVALID,
		"LCD_D11":      pins.INVALID,
		"LCD_D12":      pins.INVALID,
		"LCD_D13":      pins.INVALID,
		"LCD_D14":      pins.INVALID,
		"LCD_D15":      pins.INVALID,
		"LCD_D18":      pins.INVALID,
		"LCD_D19":      pins.INVALID,
		"LCD_D2":       pins.INVALID,
		"LCD_D20":      pins.INVALID,
		"LCD_D21":      pins.INVALID,
		"LCD_D22":      pins.INVALID,
		"LCD_D23":      pins.INVALID,
		"LCD_D3":       pins.INVALID,
		"LCD_D4":       pins.INVALID,
		"LCD_D5":       pins.INVALID,
		"LCD_D6":       pins.INVALID,
		"LCD_D7":       pins.INVALID,
		"LCD_DE":       pins.INVALID,
		"LCD_HSYNC":    pins.INVALID,
		"LCD_VSYNC":    pins.INVALID,
		"LVDS_VN0":     pins.INVALID,
		"LVDS_VN1":     pins.INVALID,
		"LVDS_VN2":     pins.INVALID,
		"LVDS_VN3":     pins.INVALID,
		"LVDS_VNC":     pins.INVALID,
		"LVDS_VP0":     pins.INVALID,
		"LVDS_VP1":     pins.INVALID,
		"LVDS_VP2":     pins.INVALID,
		"LVDS_VP3":     pins.INVALID,
		"LVDS_VPC":     pins.INVALID,
		"MDC":          pins.INVALID,
		"MDIO":         pins.INVALID,
		"MIC_CLK":      pins.INVALID,
		"MIC_DATA":     pins.INVALID,
		"NAND_ALE":     pins.INVALID,
		"NAND_CE0":     pins.INVALID,
		"NAND_CE1":     pins.INVALID,
		"NAND_CLE":     pins.INVALID,
		"NAND_DQ0":     pins.INVALID,
		"NAND_DQ1":     pins.INVALID,
		"NAND_DQ2":     pins.INVALID,
		"NAND_DQ3":     pins.INVALID,
		"NAND_DQ4":     pins.INVALID,
		"NAND_DQ5":     pins.INVALID,
		"NAND_DQ6":     pins.INVALID,
		"NAND_DQ7":     pins.INVALID,
		"NAND_DQS":     pins.INVALID,
		"NAND_RB0":     pins.INVALID,
		"NAND_RB1":     pins.INVALID,
		"NAND_RE":      pins.INVALID,
		"NAND_WE":      pins.INVALID,
		"OWA_OUT":      pins.INVALID,
	*/
	"PCM0_BCLK":    pins.INVALID,
	"PCM0_DIN":     pins.INVALID,
	"PCM0_DOUT":    pins.INVALID,
	"PCM0_SYNC":    pins.INVALID,
	"PCM1_BCLK":    pins.INVALID,
	"PCM1_DIN":     pins.INVALID,
	"PCM1_DOUT":    pins.INVALID,
	"PCM1_SYNC":    pins.INVALID,
	"PLL_LOCK_DBG": pins.INVALID,
	"PWM0":         pins.INVALID,
	/*
		"RGMII_CLKI":   pins.INVALID,
		"RGMII_RXCK":   pins.INVALID,
		"RGMII_RXCT":   pins.INVALID,
		"RGMII_RXD0":   pins.INVALID,
		"RGMII_RXD1":   pins.INVALID,
		"RGMII_RXD2":   pins.INVALID,
		"RGMII_RXD3":   pins.INVALID,
		"RGMII_RXER":   pins.INVALID,
		"RGMII_TXCK":   pins.INVALID,
		"RGMII_TXCT":   pins.INVALID,
		"RGMII_TXD0":   pins.INVALID,
		"RGMII_TXD1":   pins.INVALID,
		"RGMII_TXD2":   pins.INVALID,
		"RGMII_TXD3":   pins.INVALID,
		"SDC0_CLK":     pins.INVALID,
		"SDC0_CMD":     pins.INVALID,
		"SDC0_D0":      pins.INVALID,
		"SDC0_D1":      pins.INVALID,
		"SDC0_D2":      pins.INVALID,
		"SDC0_D3":      pins.INVALID,
		"SDC1_CLK":     pins.INVALID,
		"SDC1_CMD":     pins.INVALID,
		"SDC1_D0":      pins.INVALID,
		"SDC1_D1":      pins.INVALID,
		"SDC1_D2":      pins.INVALID,
		"SDC1_D3":      pins.INVALID,
		"SDC2_CLK":     pins.INVALID,
		"SDC2_CMD":     pins.INVALID,
		"SDC2_D0":      pins.INVALID,
		"SDC2_D1":      pins.INVALID,
		"SDC2_D2":      pins.INVALID,
		"SDC2_D3":      pins.INVALID,
		"SDC2_D4":      pins.INVALID,
		"SDC2_D5":      pins.INVALID,
		"SDC2_D6":      pins.INVALID,
		"SDC2_D7":      pins.INVALID,
		"SDC2_DS":      pins.INVALID,
		"SDC2_RST":     pins.INVALID,
		"SIM_CLK":      pins.INVALID,
		"SIM_DATA":     pins.INVALID,
		"SIM_DET":      pins.INVALID,
		"SIM_PWREN":    pins.INVALID,
		"SIM_RST":      pins.INVALID,
		"SIM_VPPEN":    pins.INVALID,
		"SIM_VPPPP":    pins.INVALID,
	*/
	"SPI0_CLK":  pins.INVALID,
	"SPI0_CS0":  pins.INVALID,
	"SPI0_MISO": pins.INVALID,
	"SPI0_MOSI": pins.INVALID,
	"SPI1_CLK":  pins.INVALID,
	"SPI1_CS0":  pins.INVALID,
	"SPI1_MISO": pins.INVALID,
	"SPI1_MOSI": pins.INVALID,
	/*
		"S_CIR_RX":  pins.INVALID,
		"S_I2C_CSK": pins.INVALID,
		"S_I2C_SCL": pins.INVALID,
		"S_I2C_SDA": pins.INVALID,
		"S_I2C_SDA": pins.INVALID,
		"S_JTAG_CK": pins.INVALID,
		"S_JTAG_DI": pins.INVALID,
		"S_JTAG_DO": pins.INVALID,
		"S_JTAG_MS": pins.INVALID,
		"S_PWM":     pins.INVALID,
		"S_RSB_SCK": pins.INVALID,
		"S_RSB_SDA": pins.INVALID,
		"S_UART_RX": pins.INVALID,
		"S_UART_TX": pins.INVALID,
		"TS_CLK":    pins.INVALID,
		"TS_D0":     pins.INVALID,
		"TS_D1":     pins.INVALID,
		"TS_D2":     pins.INVALID,
		"TS_D3":     pins.INVALID,
		"TS_D4":     pins.INVALID,
		"TS_D5":     pins.INVALID,
		"TS_D6":     pins.INVALID,
		"TS_D7":     pins.INVALID,
		"TS_DVLD":   pins.INVALID,
		"TS_ERR":    pins.INVALID,
		"TS_SYNC":   pins.INVALID,
	*/
	"UART0_RX":  pins.INVALID,
	"UART0_TX":  pins.INVALID,
	"UART1_CTS": pins.INVALID,
	"UART1_RTS": pins.INVALID,
	"UART1_RX":  pins.INVALID,
	"UART1_TX":  pins.INVALID,
	"UART2_CTS": pins.INVALID,
	"UART2_RTS": pins.INVALID,
	"UART2_RX":  pins.INVALID,
	"UART2_TX":  pins.INVALID,
	"UART3_CTS": pins.INVALID,
	"UART3_RTS": pins.INVALID,
	"UART3_RX":  pins.INVALID,
	"UART3_TX":  pins.INVALID,
	"UART4_CTS": pins.INVALID,
	"UART4_RTS": pins.INVALID,
	"UART4_RX":  pins.INVALID,
	"UART4_TX":  pins.INVALID,
}

type Pin struct {
	index       uint8      // only used to lookup in mapping
	group       uint8      // as per register offset calculation
	offset      uint8      // as per register offset calculation
	name        string     // name as per datasheet
	defaultPull gpio.Pull  // default pull at startup
	edge        *sysfs.Pin // Mutable, set once, then never set back to nil
}

// http://forum.pine64.org/showthread.php?tid=474
// about number calculation.
var (
	PB0  gpio.PinIO = &Pins[0]   // 32
	PB1  gpio.PinIO = &Pins[1]   // 33
	PB2  gpio.PinIO = &Pins[2]   // 34
	PB3  gpio.PinIO = &Pins[3]   // 35
	PB4  gpio.PinIO = &Pins[4]   // 36
	PB5  gpio.PinIO = &Pins[5]   // 37
	PB6  gpio.PinIO = &Pins[6]   // 38
	PB7  gpio.PinIO = &Pins[7]   // 39
	PB8  gpio.PinIO = &Pins[8]   // 40
	PB9  gpio.PinIO = &Pins[9]   // 41
	PC0  gpio.PinIO = &Pins[10]  //
	PC1  gpio.PinIO = &Pins[11]  //
	PC2  gpio.PinIO = &Pins[12]  //
	PC3  gpio.PinIO = &Pins[13]  //
	PC4  gpio.PinIO = &Pins[14]  //
	PC5  gpio.PinIO = &Pins[15]  //
	PC6  gpio.PinIO = &Pins[16]  //
	PC7  gpio.PinIO = &Pins[17]  //
	PC8  gpio.PinIO = &Pins[18]  //
	PC9  gpio.PinIO = &Pins[19]  //
	PC10 gpio.PinIO = &Pins[20]  //
	PC11 gpio.PinIO = &Pins[21]  //
	PC12 gpio.PinIO = &Pins[22]  //
	PC13 gpio.PinIO = &Pins[23]  //
	PC14 gpio.PinIO = &Pins[24]  //
	PC15 gpio.PinIO = &Pins[25]  //
	PC16 gpio.PinIO = &Pins[26]  //
	PD0  gpio.PinIO = &Pins[27]  //
	PD1  gpio.PinIO = &Pins[28]  //
	PD2  gpio.PinIO = &Pins[29]  //
	PD3  gpio.PinIO = &Pins[30]  //
	PD4  gpio.PinIO = &Pins[31]  //
	PD5  gpio.PinIO = &Pins[32]  //
	PD6  gpio.PinIO = &Pins[33]  //
	PD7  gpio.PinIO = &Pins[34]  //
	PD8  gpio.PinIO = &Pins[35]  //
	PD9  gpio.PinIO = &Pins[36]  //
	PD10 gpio.PinIO = &Pins[37]  //
	PD11 gpio.PinIO = &Pins[38]  //
	PD12 gpio.PinIO = &Pins[39]  //
	PD13 gpio.PinIO = &Pins[40]  //
	PD14 gpio.PinIO = &Pins[41]  //
	PD15 gpio.PinIO = &Pins[42]  //
	PD16 gpio.PinIO = &Pins[43]  //
	PD17 gpio.PinIO = &Pins[44]  //
	PD18 gpio.PinIO = &Pins[45]  //
	PD19 gpio.PinIO = &Pins[46]  //
	PD20 gpio.PinIO = &Pins[47]  //
	PD21 gpio.PinIO = &Pins[48]  //
	PD22 gpio.PinIO = &Pins[49]  //
	PD23 gpio.PinIO = &Pins[50]  //
	PD24 gpio.PinIO = &Pins[51]  //
	PE0  gpio.PinIO = &Pins[52]  //
	PE1  gpio.PinIO = &Pins[53]  //
	PE2  gpio.PinIO = &Pins[54]  //
	PE3  gpio.PinIO = &Pins[55]  //
	PE4  gpio.PinIO = &Pins[56]  //
	PE5  gpio.PinIO = &Pins[57]  //
	PE6  gpio.PinIO = &Pins[58]  //
	PE7  gpio.PinIO = &Pins[59]  //
	PE8  gpio.PinIO = &Pins[60]  //
	PE9  gpio.PinIO = &Pins[61]  //
	PE10 gpio.PinIO = &Pins[62]  //
	PE11 gpio.PinIO = &Pins[63]  //
	PE12 gpio.PinIO = &Pins[64]  //
	PE13 gpio.PinIO = &Pins[65]  //
	PE14 gpio.PinIO = &Pins[66]  //
	PE15 gpio.PinIO = &Pins[67]  //
	PE16 gpio.PinIO = &Pins[68]  //
	PE17 gpio.PinIO = &Pins[69]  //
	PF0  gpio.PinIO = &Pins[70]  //
	PF1  gpio.PinIO = &Pins[71]  //
	PF2  gpio.PinIO = &Pins[72]  //
	PF3  gpio.PinIO = &Pins[73]  //
	PF4  gpio.PinIO = &Pins[74]  //
	PF5  gpio.PinIO = &Pins[75]  //
	PF6  gpio.PinIO = &Pins[76]  //
	PG0  gpio.PinIO = &Pins[77]  // 192
	PG1  gpio.PinIO = &Pins[78]  // 193
	PG2  gpio.PinIO = &Pins[79]  // 194
	PG3  gpio.PinIO = &Pins[80]  // 195
	PG4  gpio.PinIO = &Pins[81]  // 196
	PG5  gpio.PinIO = &Pins[82]  // 197
	PG6  gpio.PinIO = &Pins[83]  // 198
	PG7  gpio.PinIO = &Pins[84]  // 199
	PG8  gpio.PinIO = &Pins[85]  // 200
	PG9  gpio.PinIO = &Pins[86]  // 201
	PG10 gpio.PinIO = &Pins[87]  // 202
	PG11 gpio.PinIO = &Pins[88]  // 203
	PG12 gpio.PinIO = &Pins[89]  // 204
	PG13 gpio.PinIO = &Pins[90]  // 205
	PH0  gpio.PinIO = &Pins[91]  // 224
	PH1  gpio.PinIO = &Pins[92]  // 225
	PH2  gpio.PinIO = &Pins[93]  // 226
	PH3  gpio.PinIO = &Pins[94]  // 227
	PH4  gpio.PinIO = &Pins[95]  // 228
	PH5  gpio.PinIO = &Pins[96]  // 229
	PH6  gpio.PinIO = &Pins[97]  // 230
	PH7  gpio.PinIO = &Pins[98]  // 232
	PH8  gpio.PinIO = &Pins[99]  // 233
	PH9  gpio.PinIO = &Pins[100] // 234
	PH10 gpio.PinIO = &Pins[101] // 235
	PH11 gpio.PinIO = &Pins[102] //
)

// PinIO implementation.

// Number implements gpio.PinIO
//
// It returns the GPIO pin number as represented by gpio sysfs.
func (p *Pin) Number() int {
	return int(p.group/9)*32 + int(p.offset)
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
		if s := mapping[p.index][0]; len(s) != 0 {
			return s
		}
		return "<Alt1>"
	case alt2:
		if s := mapping[p.index][1]; len(s) != 0 {
			return s
		}
		return "<Alt2>"
	case alt3:
		if s := mapping[p.index][2]; len(s) != 0 {
			return s
		}
		return "<Alt3>"
	case alt4:
		if s := mapping[p.index][3]; len(s) != 0 {
			return s
		}
		return "<Alt4>"
	case alt5:
		if s := mapping[p.index][4]; len(s) != 0 {
			return s
		}
		return "<Alt5>"
	case disabled:
		return "<Disabled>"
	default:
		return "<Internal error>"
	}
}

func (p *Pin) In(pull gpio.Pull) error {
	if gpioMemory == nil {
		return errors.New("subsystem not initialized")
	}
	if !p.setFunction(in) {
		return fmt.Errorf("failed to set pin %s as input", p.name)
	}
	if pull == gpio.PullNoChange {
		return nil
	}
	off := p.group + 7 + p.offset/16
	shift := 2 * (p.offset % 16)
	// Do it in a way that is concurrent safe.
	// Pn_PULL  n*0x24+0x1C Port n Pull Register (n from 1(B) to 7(H))
	gpioMemory[off] &^= 3 << shift
	switch pull {
	case gpio.Down:
		gpioMemory[off] = 2 << shift
	case gpio.Up:
		gpioMemory[off] = 1 << shift
	default:
	}
	return nil
}

func (p *Pin) Read() gpio.Level {
	// Pn_DAT  n*0x24+0x10  Port n Data Register (n from 1(B) to 7(H))
	return gpio.Level(gpioMemory[p.group+4]&(1<<p.offset) != 0)
}

// Edges creates a edge detection loop and implements gpio.PinIn.
//
// This requires opening a gpio sysfs file handle. The pin will be exported at
// /sys/class/gpio/gpio*/. Note that the pin will not be unexported at
// shutdown.
//
// Not all pins support edge detection Allwinner processors!
func (p *Pin) Edges() (<-chan gpio.Level, error) {
	switch p.group {
	case 1 * 9, 6 * 9, 7 * 9:
		// TODO(maruel): Some pins do not support Alt5 in these groups.
	default:
		return nil, errors.New("only groups PB, PG, PH (and PL if available) support edge based triggering")
	}
	// This is a race condition but this is fine; at worst PinByNumber() is called
	// twice but it is guaranteed to return the same value. p.edge is never set
	// to nil.
	if p.edge == nil {
		var err error
		if p.edge, err = sysfs.PinByNumber(p.Number()); err != nil {
			return nil, err
		}
	}
	if err := p.edge.In(gpio.PullNoChange); err != nil {
		return nil, err
	}
	return p.edge.Edges()
}

func (p *Pin) DisableEdges() {
	if p.edge != nil {
		p.edge.DisableEdges()
	}
}

func (p *Pin) Pull() gpio.Pull {
	off := p.group + 7 + p.offset/16
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
	// Pn_DAT  n*0x24+0x10  Port n Data Register (n from 1(B) to 7(H))
	if l {
		gpioMemory[p.group+4] |= bit
	} else {
		gpioMemory[p.group+4] &^= bit
	}
	return nil
}

//

// function returns the current GPIO pin function.
func (p *Pin) function() function {
	if gpioMemory == nil {
		return disabled
	}
	off := p.group + p.offset/8
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
	} else if actual == in {
		p.DisableEdges()
	}
	off := p.group + p.offset/8
	shift := 4 * (p.offset % 8)
	mask := uint32(disabled) << shift
	v := (uint32(f) << shift) ^ mask
	// First disable, then setup. This is concurrent safe.
	// Pn_CFGx n*0x24+0x0x  Port n Configure Register x (n from 1(B) to 7(H))
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
	alt5     function = 6 // often interrupt based edge detection as input
	disabled function = 7
)

// http://files.pine64.org/doc/datasheet/pine64/Allwinner_A64_User_Manual_V1.0.pdf
// Page 376 GPIO PB to PH.
var gpioMemory []uint32

// Page 73 for memory mapping overview.
// Page 194 for PWM.
// Page 230 for crypto engine.
// Page 278 audio including ADC.
// Page 376 GPIO PB to PH
// Page 410 GPIO PL
// Page 536 I²C (I2C)
// Page 545 SPI
// Page 560 UART
// Page 621 I2S/PCM

var _ gpio.PinIO = &Pin{}

// mapping excludes functions in and out.
// Datasheet, page 23.
// http://files.pine64.org/doc/datasheet/pine64/A64_Datasheet_V1.1.pdf
//
// - The datasheet uses TWI instead of I2C but I renamed for consistency.
// - AIF is audio interface, i.e. to connect to S/PDIF
// - RGMII means Reduced gigabit media-independent interface
// - SDC means SDCard?
// - NAND is for NAND flash controller
// - CSI and CCI are for video capture
var mapping = [103][5]string{
	{"UART2_TX", "", "JTAG_MS0", "", "PB_EINT0"},                    // PB0
	{"UART2_RX", "", "JTAG_CK0", "SIM_PWREN", "PB_EINT1"},           // PB1
	{"UART2_RTS", "", "JTAG_DO0", "SIM_VPPEN", "PB_EINT2"},          // PB2
	{"UART2_CTS", "I2S0_MCLK", "JTAG_DI0", "SIM_VPPPP", "PB_EINT3"}, // PB3
	{"AIF2_SYNC", "PCM0_SYNC", "", "SIM_CLK", "PB_EINT4"},           // PB4
	{"AIF2_BCLK", "PCM0_BCLK", "", "SIM_DATA", "PB_EINT5"},          // PB5
	{"AIF2_DOUT", "PCM0_DOUT", "", "SIM_RST", "PB_EINT6"},           // PB6
	{"AIF2_DIN", "PCM0_DIN", "", "SIM_DET", "PB_EINT7"},             // PB7
	{"", "", "UART0_TX", "", "PB_EINT8"},                            // PB8
	{"", "", "UART0_RX", "", "PB_EINT9"},                            // PB9
	{"NAND_WE", "", "SPI0_MOSI"},                                    // PC0
	{"NAND_ALE", "SDC2_DS", "SPI0_MISO"},                            // PC1
	{"NAND_CLE", "", "SPI0_CLK"},                                    // PC2
	{"NAND_CE1", "", "SPI0_CS0"},                                    // PC3
	{"NAND_CE0"},                                                    // PC4
	{"NAND_RE", "SDC2_CLK"},                                         // PC5
	{"NAND_RB0", "SDC2_CMD"},                                        // PC6
	{"NAND_RB1"},                                                    // PC7
	{"NAND_DQ0", "SDC2_D0"},                                         // PC8
	{"NAND_DQ1", "SDC2_D1"},                                         // PC9
	{"NAND_DQ2", "SDC2_D2"},                                         // PC10
	{"NAND_DQ3", "SDC2_D3"},                                         // PC11
	{"NAND_DQ4", "SDC2_D4"},                                         // PC12
	{"NAND_DQ5", "SDC2_D5"},                                         // PC13
	{"NAND_DQ6", "SDC2_D6"},                                         // PC14
	{"NAND_DQ7", "SDC2_D7"},                                         // PC15
	{"NAND_DQS", "SDC2_RST"},                                        // PC16
	{"LCD_D2", "UART3_TX", "SPI1_CS0", "CCIR_CLK"},                  // PD0
	{"LCD_D3", "UART3_RX", "SPI1_CLK", "CCIR_DE"},                   // PD1
	{"LCD_D4", "UART4_TX", "SPI1_MOSI", "CCIR_HSYNC"},               // PD2
	{"LCD_D5", "UART4_RX", "SPI1_MISO", "CCIR_VSYNC"},               // PD3
	{"LCD_D6", "UART4_RTS", "", "CCIR_D0"},                          // PD4
	{"LCD_D7", "UART4_CTS", "", "CCIR_D1"},                          // PD5
	{"LCD_D10", "", "", "CCIR_D2"},                                  // PD6
	{"LCD_D11", "", "", "CCIR_D3"},                                  // PD7
	{"LCD_D12", "", "RGMII_RXD3", "CCIR_D4"},                        // PD8
	{"LCD_D13", "", "RGMII_RXD2", "CCIR_D5"},                        // PD9
	{"LCD_D14", "", "RGMII_RXD1"},                                   // PD10
	{"LCD_D15", "", "RGMII_RXD0"},                                   // PD11
	{"LCD_D18", "LVDS_VP0", "RGMII_RXCK"},                           // PD12
	{"LCD_D19", "LVDS_VN0", "RGMII_RXCT"},                           // PD13
	{"LCD_D20", "LVDS_VP1", "RGMII_RXER"},                           // PD14
	{"LCD_D21", "LVDS_VN1", "RGMII_TXD3", "CCIR_D6"},                // PD15
	{"LCD_D22", "LVDS_VP2", "RGMII_TXD2", "CCIR_D7"},                // PD16
	{"LCD_D23", "LVDS_VN2", "RGMII_TXD1"},                           // PD17
	{"LCD_CLK", "LVDS_VPC", "RGMII_TXD0"},                           // PD18
	{"LCD_DE", "LVDS_VNC", "RGMII_TXCK"},                            // PD19
	{"LCD_HSYNC", "LVDS_VP3", "RGMII_TXCT"},                         // PD20
	{"LCD_VSYNC", "LVDS_VN3", "RGMII_CLKI"},                         // PD21
	{"PWM0", "", "MDC"},                                             // PD22
	{"", "", "MDIO"},                                                // PD23
	{"", ""},                                                        // PD24
	{"CSI_PCLK", "", "TS_CLK"},                                      // PE0
	{"CSI_MCLK", "", "TS_ERR"},                                      // PE1
	{"CSI_HSYNC", "", "TS_SYNC"},                                    // PE2
	{"CSI_VSYNC", "", "TS_DVLD"},                                    // PE3
	{"CSI_D0", "", "TS_D0"},                                         // PE4
	{"CSI_D1", "", "TS_D1"},                                         // PE5
	{"CSI_D2", "", "TS_D2"},                                         // PE6
	{"CSI_D3", "", "TS_D3"},                                         // PE7
	{"CSI_D4", "", "TS_D4"},                                         // PE8
	{"CSI_D5", "", "TS_D5"},                                         // PE9
	{"CSI_D6", "", "TS_D6"},                                         // PE10
	{"CSI_D7", "", "TS_D7"},                                         // PE11
	{"CSI_SCK"},                                                     // PE12
	{"CSI_SDA"},                                                     // PE13
	{"PLL_LOCK_DBG", "I2C2_SCL"},                                    // PE14
	{"", "I2C2_SDA"},                                                // PE15
	{},                                                              // PE16
	{"", ""},                                                        // PE17
	{"SDC0_D1", "JTAG_MS1"},                                         // PF0
	{"SDC0_D0", "JTAG_DI1"},                                         // PF1
	{"SDC0_CLK", "UART0_TX"},                                        // PF2
	{"SDC0_CMD", "JTAG_DO1"},                                        // PF3
	{"SDC0_D3", "UART0_RX"},                                         // PF4
	{"SDC0_D2", "JTAG_CK1"},                                         // PF5
	{"", "", ""},                                                    // PF6
	{"SDC1_CLK", "", "", "", "PG_EINT0"},                            // PG0
	{"SDC1_CMD", "", "", "", "PG_EINT1"},                            // PG1
	{"SDC1_D0", "", "", "", "PG_EINT2"},                             // PG2
	{"SDC1_D1", "", "", "", "PG_EINT3"},                             // PG3
	{"SDC1_D2", "", "", "", "PG_EINT4"},                             // PG4
	{"SDC1_D3", "", "", "", "PG_EINT5"},                             // PG5
	{"UART1_TX", "", "", "", "PG_EINT6"},                            // PG6
	{"UART1_RX", "", "", "", "PG_EINT7"},                            // PG7
	{"UART1_RTS", "", "", "", "PG_EINT8"},                           // PG8
	{"UART1_CTS", "", "", "", "PG_EINT9"},                           // PG9
	{"AIF3_SYNC", "PCM1_SYNC", "", "", "PG_EINT10"},                 // PG10
	{"AIF3_BCLK", "PCM1_BCLK", "", "", "PG_EINT11"},                 // PG11
	{"AIF3_DOUT", "PCM1_DOUT", "", "", "PG_EINT12"},                 // PG12
	{"AIF3_DIN", "PCM1_DIN", "", "", "PG_EINT13"},                   // PG13
	{"I2C0_SCL", "", "", "", "PH_EINT0"},                            // PH0
	{"I2C0_SDA", "", "", "", "PH_EINT1"},                            // PH1
	{"I2C1_SCL", "", "", "", "PH_EINT2"},                            // PH2
	{"I2C1_SDA", "", "", "", "PH_EINT3"},                            // PH3
	{"UART3_TX", "", "", "", "PH_EINT4"},                            // PH4
	{"UART3_RX", "", "", "", "PH_EINT5"},                            // PH5
	{"UART3_RTS", "", "", "", "PH_EINT6"},                           // PH6
	{"UART3_CTS", "", "", "", "PH_EINT7"},                           // PH7
	{"OWA_OUT", "", "", "", "PH_EINT8"},                             // PH8
	{"", "", "", "", "PH_EINT9"},                                    // PH9
	{"MIC_CLK", "", "", "", "PH_EINT10"},                            // PH10
	{"MIC_DATA", "", "", "", "PH_EINT11"},                           // PH11
}

// getBaseAddress queries the virtual file system to retrieve the base address
// of the GPIO registers for GPIO pins in groups PB to PH.
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

// driver implements pio.Driver.
type driver struct {
}

func (d *driver) String() string {
	return "allwinner"
}

func (d *driver) Type() pio.Type {
	return pio.Processor
}

func (d *driver) Prerequisites() []string {
	return nil
}

func (d *driver) Init() (bool, error) {
	if !internal.IsAllwinner() {
		return false, nil
	}
	mem, err := gpiomem.OpenMem(getBaseAddress())
	if err != nil {
		return true, err
	}
	gpioMemory = mem.Uint32
	for i := range Pins {
		if err := gpio.Register(&Pins[i]); err != nil {
			return true, err
		}
		if f := Pins[i].Function(); f[:2] != "In" && f[:3] != "Out" {
			functional[f] = &Pins[i]
		}
	}
	for f, p := range functional {
		gpio.MapFunction(f, p)
	}
	return true, nil
}

var _ pio.Driver = &driver{}
