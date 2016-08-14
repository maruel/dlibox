// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// pins is a small app to read the function of each pin.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/maruel/dlibox/go/rpi"
)

func printPin(name string, value rpi.Pin, showInvalid bool) {
	if showInvalid || value != rpi.INVALID {
		fmt.Printf("%-9s: %s\n", name, value)
	}
}

func mainImpl() error {
	all := flag.Bool("a", false, "print everything")
	fun := flag.Bool("f", false, "print functional pins (e.g. I2C_SCL1)")
	gpio := flag.Bool("g", false, "print GPIO pins (e.g. GPIO1) (default)")
	hardware := flag.Bool("h", false, "print hardware pins (e.g. P1_1)")
	invalid := flag.Bool("n", false, "show not connected/INVALID pins")
	flag.Parse()
	if *all {
		*fun = true
		*gpio = true
		*hardware = true
		*invalid = true
	} else if !*fun && !*gpio && !*hardware {
		*gpio = true
	}

	if *fun {
		printPin("GPCLK0", rpi.GPCLK0, *invalid)
		printPin("GPCLK1", rpi.GPCLK1, *invalid)
		printPin("GPCLK2", rpi.GPCLK2, *invalid)
		printPin("I2C_SCL0", rpi.I2C_SCL0, *invalid)
		printPin("I2C_SDA0", rpi.I2C_SDA0, *invalid)
		printPin("I2C_SCL1", rpi.I2C_SCL1, *invalid)
		printPin("I2C_SDA1", rpi.I2C_SDA1, *invalid)
		printPin("IR_IN", rpi.IR_IN, *invalid)
		printPin("IR_OUT", rpi.IR_OUT, *invalid)
		printPin("PCM_CLK", rpi.PCM_CLK, *invalid)
		printPin("PCM_FS", rpi.PCM_FS, *invalid)
		printPin("PCM_DIN", rpi.PCM_DIN, *invalid)
		printPin("PCM_DOUT", rpi.PCM_DOUT, *invalid)
		printPin("PWM0_OUT", rpi.PWM0_OUT, *invalid)
		printPin("PWM1_OUT", rpi.PWM1_OUT, *invalid)
		printPin("SPI0_CE0", rpi.SPI0_CE0, *invalid)
		printPin("SPI0_CE1", rpi.SPI0_CE1, *invalid)
		printPin("SPI0_CLK", rpi.SPI0_CLK, *invalid)
		printPin("SPI0_MISO", rpi.SPI0_MISO, *invalid)
		printPin("SPI0_MOSI", rpi.SPI0_MOSI, *invalid)
		printPin("SPI1_CE0", rpi.SPI1_CE0, *invalid)
		printPin("SPI1_CE1", rpi.SPI1_CE1, *invalid)
		printPin("SPI1_CE2", rpi.SPI1_CE2, *invalid)
		printPin("SPI1_CLK", rpi.SPI1_CLK, *invalid)
		printPin("SPI1_MISO", rpi.SPI1_MISO, *invalid)
		printPin("SPI1_MOSI", rpi.SPI1_MOSI, *invalid)
		printPin("UART_RXD0", rpi.UART_RXD0, *invalid)
		printPin("UART_CTS0", rpi.UART_CTS0, *invalid)
		printPin("UART_CTS1", rpi.UART_CTS1, *invalid)
		printPin("UART_RTS0", rpi.UART_RTS0, *invalid)
		printPin("UART_RTS1", rpi.UART_RTS1, *invalid)
		printPin("UART_TXD0", rpi.UART_TXD0, *invalid)
		printPin("UART_RXD1", rpi.UART_RXD1, *invalid)
		printPin("UART_TXD1", rpi.UART_TXD1, *invalid)
	}
	if *gpio {
		for p := rpi.GPIO0; p <= rpi.GPIO53; p++ {
			f := p.Function()
			if p.IsConnected() {
				if f == rpi.In || f == rpi.Out {
					fmt.Printf("%-6s: %-4s %s\n", p, f, p.ReadInstant())
				} else {
					fmt.Printf("%-6s: %s\n", p, f)
				}
			} else if *invalid {
				fmt.Printf("%-6s: %-4s (not connected)\n", p, f)
			}
		}
	}
	if *hardware {
		fmt.Printf("P1_1        : %s\n", rpi.P1_1)
		fmt.Printf("P1_2        : %s\n", rpi.P1_2)
		fmt.Printf("P1_3        : %s\n", rpi.P1_3)
		fmt.Printf("P1_4        : %s\n", rpi.P1_4)
		fmt.Printf("P1_5        : %s\n", rpi.P1_5)
		fmt.Printf("P1_6        : %s\n", rpi.P1_6)
		fmt.Printf("P1_7        : %s\n", rpi.P1_7)
		fmt.Printf("P1_8        : %s\n", rpi.P1_8)
		fmt.Printf("P1_9        : %s\n", rpi.P1_9)
		fmt.Printf("P1_10       : %s\n", rpi.P1_10)
		fmt.Printf("P1_11       : %s\n", rpi.P1_11)
		fmt.Printf("P1_12       : %s\n", rpi.P1_12)
		fmt.Printf("P1_13       : %s\n", rpi.P1_13)
		fmt.Printf("P1_14       : %s\n", rpi.P1_14)
		fmt.Printf("P1_15       : %s\n", rpi.P1_15)
		fmt.Printf("P1_16       : %s\n", rpi.P1_16)
		fmt.Printf("P1_17       : %s\n", rpi.P1_17)
		fmt.Printf("P1_18       : %s\n", rpi.P1_18)
		fmt.Printf("P1_19       : %s\n", rpi.P1_19)
		fmt.Printf("P1_20       : %s\n", rpi.P1_20)
		fmt.Printf("P1_21       : %s\n", rpi.P1_21)
		fmt.Printf("P1_22       : %s\n", rpi.P1_22)
		fmt.Printf("P1_23       : %s\n", rpi.P1_23)
		fmt.Printf("P1_24       : %s\n", rpi.P1_24)
		fmt.Printf("P1_25       : %s\n", rpi.P1_25)
		fmt.Printf("P1_26       : %s\n", rpi.P1_26)
		if rpi.P1_27 != rpi.INVALID || *invalid {
			fmt.Printf("P1_27       : %s\n", rpi.P1_27)
			fmt.Printf("P1_28       : %s\n", rpi.P1_28)
			fmt.Printf("P1_29       : %s\n", rpi.P1_29)
			fmt.Printf("P1_30       : %s\n", rpi.P1_30)
			fmt.Printf("P1_31       : %s\n", rpi.P1_31)
			fmt.Printf("P1_32       : %s\n", rpi.P1_32)
			fmt.Printf("P1_33       : %s\n", rpi.P1_33)
			fmt.Printf("P1_34       : %s\n", rpi.P1_34)
			fmt.Printf("P1_35       : %s\n", rpi.P1_35)
			fmt.Printf("P1_36       : %s\n", rpi.P1_36)
			fmt.Printf("P1_37       : %s\n", rpi.P1_37)
			fmt.Printf("P1_38       : %s\n", rpi.P1_38)
			fmt.Printf("P1_39       : %s\n", rpi.P1_39)
			fmt.Printf("P1_40       : %s\n", rpi.P1_40)
		}
		if rpi.P5_1 != rpi.INVALID || *invalid {
			fmt.Printf("P5_1        : %s\n", rpi.P5_1)
			fmt.Printf("P5_2        : %s\n", rpi.P5_2)
			fmt.Printf("P5_3        : %s\n", rpi.P5_3)
			fmt.Printf("P5_4        : %s\n", rpi.P5_4)
			fmt.Printf("P5_5        : %s\n", rpi.P5_5)
			fmt.Printf("P5_6        : %s\n", rpi.P5_6)
			fmt.Printf("P5_7        : %s\n", rpi.P5_7)
			fmt.Printf("P5_8        : %s\n", rpi.P5_8)
		}
		fmt.Printf("AUDIO_LEFT  : %s\n", rpi.AUDIO_LEFT)
		fmt.Printf("AUDIO_RIGHT : %s\n", rpi.AUDIO_RIGHT)
		fmt.Printf("HDMI_HOTPLUG: %s\n", rpi.HDMI_HOTPLUG_DETECT)
	}
	return nil
}

func main() {
	if err := mainImpl(); err != nil {
		fmt.Fprintf(os.Stderr, "pins: %s.\n", err)
		os.Exit(1)
	}
}
