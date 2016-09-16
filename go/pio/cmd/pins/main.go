// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// pins is a small app to read the function of each pin.
package main

import (
	"flag"
	"fmt"
	"os"
	"sort"

	"github.com/maruel/dlibox/go/pio/host"
	"github.com/maruel/dlibox/go/pio/host/cpu"
	"github.com/maruel/dlibox/go/pio/host/pins"
	"github.com/maruel/dlibox/go/pio/host/rpi"
)

func getMaxName() int {
	max := 0
	for _, p := range pins.All {
		if l := len(p.String()); l > max {
			max = l
		}
	}
	return max
}

func getMaxFn() int {
	max := 0
	for _, p := range pins.All {
		if l := len(p.Function()); l > max {
			max = l
		}
	}
	return max
}

func doFunctionalPins(pin func(name string, value host.Pin)) {
	/*
		// TODO(maruel): Migrate this to host too.
		pin("GPCLK0", bcm283x.GPCLK0)
		pin("GPCLK1", bcm283x.GPCLK1)
		pin("GPCLK2", bcm283x.GPCLK2)
		pin("I2C_SCL0", bcm283x.I2C_SCL0)
		pin("I2C_SDA0", bcm283x.I2C_SDA0)
		pin("I2C_SCL1", bcm283x.I2C_SCL1)
		pin("I2C_SDA1", bcm283x.I2C_SDA1)
		pin("IR_IN", bcm283x.IR_IN)
		pin("IR_OUT", bcm283x.IR_OUT)
		pin("PCM_CLK", bcm283x.PCM_CLK)
		pin("PCM_FS", bcm283x.PCM_FS)
		pin("PCM_DIN", bcm283x.PCM_DIN)
		pin("PCM_DOUT", bcm283x.PCM_DOUT)
		pin("PWM0_OUT", bcm283x.PWM0_OUT)
		pin("PWM1_OUT", bcm283x.PWM1_OUT)
		pin("SPI0_CE0", bcm283x.SPI0_CE0)
		pin("SPI0_CE1", bcm283x.SPI0_CE1)
		pin("SPI0_CLK", bcm283x.SPI0_CLK)
		pin("SPI0_MISO", bcm283x.SPI0_MISO)
		pin("SPI0_MOSI", bcm283x.SPI0_MOSI)
		pin("SPI1_CE0", bcm283x.SPI1_CE0)
		pin("SPI1_CE1", bcm283x.SPI1_CE1)
		pin("SPI1_CE2", bcm283x.SPI1_CE2)
		pin("SPI1_CLK", bcm283x.SPI1_CLK)
		pin("SPI1_MISO", bcm283x.SPI1_MISO)
		pin("SPI1_MOSI", bcm283x.SPI1_MOSI)
		pin("UART_RXD0", bcm283x.UART_RXD0)
		pin("UART_CTS0", bcm283x.UART_CTS0)
		pin("UART_CTS1", bcm283x.UART_CTS1)
		pin("UART_RTS0", bcm283x.UART_RTS0)
		pin("UART_RTS1", bcm283x.UART_RTS1)
		pin("UART_TXD0", bcm283x.UART_TXD0)
		pin("UART_RXD1", bcm283x.UART_RXD1)
		pin("UART_TXD1", bcm283x.UART_TXD1)
	*/
}

func printFunc(invalid bool) {
	doFunctionalPins(func(name string, value host.Pin) {
		/*
			p, _ := value.(*bcm283x.Pin)
			if invalid || (p != nil && rpi.IsConnected(p)) {
				fmt.Printf("%-9s: %s\n", name, value)
			}
		*/
	})
}

func printGPIO(invalid bool, maxName, maxFn int) {
	ids := make([]int, 0, len(pins.All))
	for i := range pins.All {
		ids = append(ids, i)
	}
	sort.Ints(ids)
	for _, id := range ids {
		p := pins.All[id]
		if rpi.IsConnected(p) {
			fmt.Printf("%-*s: %s\n", maxName, p, p.Function())
		} else if invalid {
			fmt.Printf("%-*s: %-*s (not connected)\n", maxName, p, maxFn, p.Function())
		}
	}
}

func printPin(invalid bool, maxName, maxFn int, hdr string, pos1 int, pin1, pin2 host.Pin) {
	fmt.Printf("%3s %*s %*s %2d x x %2d  %-*s %s\n", hdr, maxFn, pin1.Function(), maxName, pin1, pos1, pos1+1, maxName, pin2, pin2.Function())
}

func printHardware(invalid bool, maxName, maxFn int) {
	// TODO(maruel): Remove the raspbianism from here.
	fmt.Print("Header    Func  Name  Pos Pos  Name   Func\n")
	printPin(invalid, maxName, maxFn, "P1:", 1, rpi.P1_1, rpi.P1_2)
	printPin(invalid, maxName, maxFn, "", 3, rpi.P1_3, rpi.P1_4)
	printPin(invalid, maxName, maxFn, "", 5, rpi.P1_5, rpi.P1_6)
	printPin(invalid, maxName, maxFn, "", 7, rpi.P1_7, rpi.P1_8)
	printPin(invalid, maxName, maxFn, "", 9, rpi.P1_9, rpi.P1_10)
	printPin(invalid, maxName, maxFn, "", 11, rpi.P1_11, rpi.P1_12)
	printPin(invalid, maxName, maxFn, "", 13, rpi.P1_13, rpi.P1_14)
	printPin(invalid, maxName, maxFn, "", 15, rpi.P1_15, rpi.P1_16)
	printPin(invalid, maxName, maxFn, "", 17, rpi.P1_17, rpi.P1_18)
	printPin(invalid, maxName, maxFn, "", 19, rpi.P1_19, rpi.P1_20)
	printPin(invalid, maxName, maxFn, "", 21, rpi.P1_21, rpi.P1_22)
	printPin(invalid, maxName, maxFn, "", 23, rpi.P1_23, rpi.P1_24)
	printPin(invalid, maxName, maxFn, "", 25, rpi.P1_25, rpi.P1_26)
	if rpi.IsConnected(rpi.P1_27) || invalid {
		printPin(invalid, maxName, maxFn, "", 27, rpi.P1_27, rpi.P1_28)
		printPin(invalid, maxName, maxFn, "", 29, rpi.P1_29, rpi.P1_30)
		printPin(invalid, maxName, maxFn, "", 31, rpi.P1_31, rpi.P1_32)
		printPin(invalid, maxName, maxFn, "", 33, rpi.P1_33, rpi.P1_34)
		printPin(invalid, maxName, maxFn, "", 35, rpi.P1_35, rpi.P1_36)
		printPin(invalid, maxName, maxFn, "", 37, rpi.P1_37, rpi.P1_38)
		printPin(invalid, maxName, maxFn, "", 39, rpi.P1_39, rpi.P1_40)
	}
	if rpi.IsConnected(rpi.P5_1) || invalid {
		fmt.Print("\n")
		printPin(invalid, maxName, maxFn, "P5:", 1, rpi.P5_1, rpi.P5_2)
		printPin(invalid, maxName, maxFn, "", 3, rpi.P5_3, rpi.P5_4)
		printPin(invalid, maxName, maxFn, "", 5, rpi.P5_5, rpi.P5_6)
		printPin(invalid, maxName, maxFn, "", 7, rpi.P5_7, rpi.P5_8)
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

	// Explicitly initialize to catch any error.
	if err := pins.Init(); err != nil {
		return err
	}
	if *info {
		fmt.Printf("Version: %d  MaxSpeed: %dMhz\n", rpi.Version, cpu.MaxSpeed/1000000)
	}
	maxName := getMaxName()
	maxFn := getMaxFn()
	if *fun {
		printFunc(*invalid)
	}
	if *gpio {
		printGPIO(*invalid, maxName, maxFn)
	}
	if *hardware {
		printHardware(*invalid, maxName, maxFn)
	}
	return nil
}

func main() {
	if err := mainImpl(); err != nil {
		fmt.Fprintf(os.Stderr, "pins: %s.\n", err)
		os.Exit(1)
	}
}
