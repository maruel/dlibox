// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// pins is a small app to read the function of each pin.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/maruel/dlibox/go/buses/rpi"
)

func makeMapping() ([]string, int) {
	m := make([]string, 256)
	doFunctionalPins(func(name string, p rpi.Pin) {
		m[p] = name
	})
	m[rpi.INVALID] = ""
	max := 0
	for p := rpi.GPIO0; p <= rpi.GPIO53; p++ {
		if len(m[p]) == 0 {
			m[p] = fmt.Sprintf("%s/%s", p.Function().String(), p.ReadInstant())
		}
		if len(m[p]) > max {
			max = len(m[p])
		}
	}
	return m, max
}

func doFunctionalPins(pin func(name string, value rpi.Pin)) {
	pin("GPCLK0", rpi.GPCLK0)
	pin("GPCLK1", rpi.GPCLK1)
	pin("GPCLK2", rpi.GPCLK2)
	pin("I2C_SCL0", rpi.I2C_SCL0)
	pin("I2C_SDA0", rpi.I2C_SDA0)
	pin("I2C_SCL1", rpi.I2C_SCL1)
	pin("I2C_SDA1", rpi.I2C_SDA1)
	pin("IR_IN", rpi.IR_IN)
	pin("IR_OUT", rpi.IR_OUT)
	pin("PCM_CLK", rpi.PCM_CLK)
	pin("PCM_FS", rpi.PCM_FS)
	pin("PCM_DIN", rpi.PCM_DIN)
	pin("PCM_DOUT", rpi.PCM_DOUT)
	pin("PWM0_OUT", rpi.PWM0_OUT)
	pin("PWM1_OUT", rpi.PWM1_OUT)
	pin("SPI0_CE0", rpi.SPI0_CE0)
	pin("SPI0_CE1", rpi.SPI0_CE1)
	pin("SPI0_CLK", rpi.SPI0_CLK)
	pin("SPI0_MISO", rpi.SPI0_MISO)
	pin("SPI0_MOSI", rpi.SPI0_MOSI)
	pin("SPI1_CE0", rpi.SPI1_CE0)
	pin("SPI1_CE1", rpi.SPI1_CE1)
	pin("SPI1_CE2", rpi.SPI1_CE2)
	pin("SPI1_CLK", rpi.SPI1_CLK)
	pin("SPI1_MISO", rpi.SPI1_MISO)
	pin("SPI1_MOSI", rpi.SPI1_MOSI)
	pin("UART_RXD0", rpi.UART_RXD0)
	pin("UART_CTS0", rpi.UART_CTS0)
	pin("UART_CTS1", rpi.UART_CTS1)
	pin("UART_RTS0", rpi.UART_RTS0)
	pin("UART_RTS1", rpi.UART_RTS1)
	pin("UART_TXD0", rpi.UART_TXD0)
	pin("UART_RXD1", rpi.UART_RXD1)
	pin("UART_TXD1", rpi.UART_TXD1)
}

func printFunc(invalid bool) {
	doFunctionalPins(func(name string, value rpi.Pin) {
		if invalid || (value != rpi.INVALID && value.IsConnected()) {
			fmt.Printf("%-9s: %s\n", name, value)
		}
	})
}

func printGPIO(invalid bool, m []string, max int) {
	for p := rpi.GPIO0; p <= rpi.GPIO53; p++ {
		if p.IsConnected() {
			fmt.Printf("%-6s: %s\n", p, m[p])
		} else if invalid {
			fmt.Printf("%-6s: %-*s (not connected)\n", p, max, m[p])
		}
	}
}

func printHardware(invalid bool, m []string, max int) {
	fmt.Print("Header    Func  Name  Pos Pos  Name   Func\n")
	fmt.Printf("P1: %*s %6s  1 x x 2  %-6s %s\n", max, m[rpi.P1_1], rpi.P1_1, rpi.P1_2, m[rpi.P1_2])
	fmt.Printf("    %*s %6s  3 x x 4  %-6s %s\n", max, m[rpi.P1_3], rpi.P1_3, rpi.P1_4, m[rpi.P1_4])
	fmt.Printf("    %*s %6s  5 x x 6  %-6s %s\n", max, m[rpi.P1_5], rpi.P1_5, rpi.P1_6, m[rpi.P1_6])
	fmt.Printf("    %*s %6s  7 x x 8  %-6s %s\n", max, m[rpi.P1_7], rpi.P1_7, rpi.P1_8, m[rpi.P1_8])
	fmt.Printf("    %*s %6s  9 x x 10 %-6s %s\n", max, m[rpi.P1_9], rpi.P1_9, rpi.P1_10, m[rpi.P1_10])
	fmt.Printf("    %*s %6s 11 x x 12 %-6s %s\n", max, m[rpi.P1_11], rpi.P1_11, rpi.P1_12, m[rpi.P1_12])
	fmt.Printf("    %*s %6s 13 x x 14 %-6s %s\n", max, m[rpi.P1_13], rpi.P1_13, rpi.P1_14, m[rpi.P1_14])
	fmt.Printf("    %*s %6s 15 x x 16 %-6s %s\n", max, m[rpi.P1_15], rpi.P1_15, rpi.P1_16, m[rpi.P1_16])
	fmt.Printf("    %*s %6s 17 x x 18 %-6s %s\n", max, m[rpi.P1_17], rpi.P1_17, rpi.P1_18, m[rpi.P1_18])
	fmt.Printf("    %*s %6s 19 x x 20 %-6s %s\n", max, m[rpi.P1_19], rpi.P1_19, rpi.P1_20, m[rpi.P1_20])
	fmt.Printf("    %*s %6s 21 x x 22 %-6s %s\n", max, m[rpi.P1_21], rpi.P1_21, rpi.P1_22, m[rpi.P1_22])
	fmt.Printf("    %*s %6s 23 x x 24 %-6s %s\n", max, m[rpi.P1_23], rpi.P1_23, rpi.P1_24, m[rpi.P1_24])
	fmt.Printf("    %*s %6s 25 x x 26 %-6s %s\n", max, m[rpi.P1_25], rpi.P1_25, rpi.P1_26, m[rpi.P1_26])
	if rpi.P1_27 != rpi.INVALID || invalid {
		fmt.Printf("    %*s %6s 27 x x 28 %-6s %s\n", max, m[rpi.P1_27], rpi.P1_27, rpi.P1_28, m[rpi.P1_28])
		fmt.Printf("    %*s %6s 29 x x 30 %-6s %s\n", max, m[rpi.P1_29], rpi.P1_29, rpi.P1_30, m[rpi.P1_30])
		fmt.Printf("    %*s %6s 31 x x 32 %-6s %s\n", max, m[rpi.P1_31], rpi.P1_31, rpi.P1_32, m[rpi.P1_32])
		fmt.Printf("    %*s %6s 33 x x 34 %-6s %s\n", max, m[rpi.P1_33], rpi.P1_33, rpi.P1_34, m[rpi.P1_34])
		fmt.Printf("    %*s %6s 35 x x 36 %-6s %s\n", max, m[rpi.P1_35], rpi.P1_35, rpi.P1_36, m[rpi.P1_36])
		fmt.Printf("    %*s %6s 37 x x 38 %-6s %s\n", max, m[rpi.P1_37], rpi.P1_37, rpi.P1_38, m[rpi.P1_38])
		fmt.Printf("    %*s %6s 39 x x 40 %-6s %s\n", max, m[rpi.P1_39], rpi.P1_39, rpi.P1_40, m[rpi.P1_40])
	}
	if rpi.P5_1 != rpi.INVALID || invalid {
		fmt.Print("\n")
		fmt.Printf("P5: %*s %6s 1 x x 2 %-6s %s\n", max, m[rpi.P5_1], rpi.P5_2, rpi.P5_1, m[rpi.P5_2])
		fmt.Printf("    %*s %6s 3 x x 4 %-6s %s\n", max, m[rpi.P5_3], rpi.P5_4, rpi.P5_3, m[rpi.P5_4])
		fmt.Printf("    %*s %6s 5 x x 6 %-6s %s\n", max, m[rpi.P5_5], rpi.P5_6, rpi.P5_5, m[rpi.P5_6])
		fmt.Printf("    %*s %6s 7 x x 8 %-6s %s\n", max, m[rpi.P5_7], rpi.P5_8, rpi.P5_7, m[rpi.P5_8])
	}
	fmt.Print("\n")
	fmt.Printf("AUDIO_LEFT  : %s\n", rpi.AUDIO_LEFT)
	fmt.Printf("AUDIO_RIGHT : %s\n", rpi.AUDIO_RIGHT)
	fmt.Printf("HDMI_HOTPLUG: %s\n", rpi.HDMI_HOTPLUG_DETECT)
}

func mainImpl() error {
	all := flag.Bool("a", false, "print everything")
	fun := flag.Bool("f", false, "print functional pins (e.g. I2C_SCL1)")
	gpio := flag.Bool("g", false, "print GPIO pins (e.g. GPIO1) (default)")
	hardware := flag.Bool("h", false, "print hardware pins (e.g. P1_1)")
	info := flag.Bool("i", false, "show general information")
	invalid := flag.Bool("n", false, "show not connected/INVALID pins")
	flag.Parse()
	if *all {
		*fun = true
		*gpio = true
		*hardware = true
		*info = true
		*invalid = true
	} else if !*fun && !*gpio && !*hardware && !*info {
		*gpio = true
	}

	if *info {
		fmt.Printf("Version: %d  MaxSpeed: %dMhz\n", rpi.Version, rpi.MaxSpeed/1000000)
	}
	m, max := makeMapping()
	if *fun {
		printFunc(*invalid)
	}
	if *gpio {
		printGPIO(*invalid, m, max)
	}
	if *hardware {
		printHardware(*invalid, m, max)
	}
	return nil
}

func main() {
	if err := mainImpl(); err != nil {
		fmt.Fprintf(os.Stderr, "pins: %s.\n", err)
		os.Exit(1)
	}
}
