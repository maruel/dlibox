// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// pins is a small app to read the function of each pin.
package main

import (
	"fmt"
	"os"

	"github.com/maruel/dlibox/go/rpi"
)

func mainImpl() error {
	fmt.Printf("GPCLK0   : %s\n", rpi.GPCLK0)
	fmt.Printf("GPCLK1   : %s\n", rpi.GPCLK1)
	fmt.Printf("GPCLK2   : %s\n", rpi.GPCLK2)
	fmt.Printf("I2C_SCL0 : %s\n", rpi.I2C_SCL0)
	fmt.Printf("I2C_SDA0 : %s\n", rpi.I2C_SDA0)
	fmt.Printf("I2C_SCL1 : %s\n", rpi.I2C_SCL1)
	fmt.Printf("I2C_SDA1 : %s\n", rpi.I2C_SDA1)
	fmt.Printf("IR_IN    : %s\n", rpi.IR_IN)
	fmt.Printf("IR_OUT   : %s\n", rpi.IR_OUT)
	fmt.Printf("PCM_CLK  : %s\n", rpi.PCM_CLK)
	fmt.Printf("PCM_FS   : %s\n", rpi.PCM_FS)
	fmt.Printf("PCM_DIN  : %s\n", rpi.PCM_DIN)
	fmt.Printf("PCM_DOUT : %s\n", rpi.PCM_DOUT)
	fmt.Printf("PWM0_OUT : %s\n", rpi.PWM0_OUT)
	fmt.Printf("PWM1_OUT : %s\n", rpi.PWM1_OUT)
	fmt.Printf("SPI0_CE0 : %s\n", rpi.SPI0_CE0)
	fmt.Printf("SPI0_CE1 : %s\n", rpi.SPI0_CE1)
	fmt.Printf("SPI0_CLK : %s\n", rpi.SPI0_CLK)
	fmt.Printf("SPI0_MISO: %s\n", rpi.SPI0_MISO)
	fmt.Printf("SPI0_MOSI: %s\n", rpi.SPI0_MOSI)
	fmt.Printf("SPI1_CE0 : %s\n", rpi.SPI1_CE0)
	fmt.Printf("SPI1_CE1 : %s\n", rpi.SPI1_CE1)
	fmt.Printf("SPI1_CE2 : %s\n", rpi.SPI1_CE2)
	fmt.Printf("SPI1_CLK : %s\n", rpi.SPI1_CLK)
	fmt.Printf("SPI1_MISO: %s\n", rpi.SPI1_MISO)
	fmt.Printf("SPI1_MOSI: %s\n", rpi.SPI1_MOSI)
	fmt.Printf("UART_RXD0: %s\n", rpi.UART_RXD0)
	fmt.Printf("UART_CTS0: %s\n", rpi.UART_CTS0)
	fmt.Printf("UART_CTS1: %s\n", rpi.UART_CTS1)
	fmt.Printf("UART_RTS0: %s\n", rpi.UART_RTS0)
	fmt.Printf("UART_RTS1: %s\n", rpi.UART_RTS1)
	fmt.Printf("UART_TXD0: %s\n", rpi.UART_TXD0)
	fmt.Printf("UART_RXD1: %s\n", rpi.UART_RXD1)
	fmt.Printf("UART_TXD1: %s\n", rpi.UART_TXD1)
	for p := rpi.GPIO0; p <= rpi.GPIO53; p++ {
		f := p.Function()
		if p.IsConnected() {
			if f == rpi.In || f == rpi.Out {
				fmt.Printf("%-6s: %-4s %s\n", p, f, p.Read())
			} else {
				fmt.Printf("%-6s: %s\n", p, f)
			}
		} else {
			fmt.Printf("%-6s: %-4s (not connected)\n", p, f)
		}
	}
	return nil
}

func main() {
	if err := mainImpl(); err != nil {
		fmt.Fprintf(os.Stderr, "pins: %s.\n", err)
		os.Exit(1)
	}
}
