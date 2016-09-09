// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// pins is a small app to read the function of each pin.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/maruel/dlibox/go/pio/host"
	"github.com/maruel/dlibox/go/pio/host/bcm283x"
	"github.com/maruel/dlibox/go/pio/host/cpu"
	"github.com/maruel/dlibox/go/pio/host/rpi"
)

// makeMapping returns a map between the pin name with its functionality and
// the pin, by number.
func makeMapping() ([]string, int) {
	m := make([]string, len(host.AllPins))
	max := 0
	for i, p := range host.AllPins {
		if b, ok := p.(*bcm283x.Pin); ok {
			// TODO(maruel): When function is Alt, should put the actual function.
			m[i] = fmt.Sprintf("%s/%s", b.Function(), b.Read())
		} else {
			m[i] = p.String()
		}
		if len(m[i]) > max {
			max = len(m[i])
		}
	}
	return m, max
}

func doFunctionalPins(pin func(name string, value host.Pin)) {
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
}

func printFunc(invalid bool) {
	doFunctionalPins(func(name string, value host.Pin) {
		p, _ := value.(*bcm283x.Pin)
		if invalid || (p != nil && rpi.IsConnected(p)) {
			fmt.Printf("%-9s: %s\n", name, value)
		}
	})
}

func printGPIO(invalid bool, m []string, max int) {
	for i, p := range host.AllPins {
		if rpi.IsConnected(p) {
			fmt.Printf("%-6s: %s\n", p, m[i])
		} else if invalid {
			fmt.Printf("%-6s: %-*s (not connected)\n", p, max, m[i])
		}
	}
}

func printPin(invalid bool, m []string, max int, hdr string, pos1 int, pin1, pin2 host.Pin) {
	name1 := ""
	if n := pin1.Number(); n >= 0 {
		name1 = m[n]
	}
	name2 := ""
	if n := pin2.Number(); n >= 0 {
		name2 = m[n]
	}
	fmt.Printf("%3s %*s %6s %2d x x %2d  %-6s %s\n", hdr, max, name1, pin1, pos1, pos1+1, pin2, name2)
}

func printHardware(invalid bool, m []string, max int) {
	fmt.Print("Header    Func  Name  Pos Pos  Name   Func\n")
	printPin(invalid, m, max, "P1:", 1, rpi.P1_1, rpi.P1_2)
	printPin(invalid, m, max, "", 3, rpi.P1_3, rpi.P1_4)
	printPin(invalid, m, max, "", 5, rpi.P1_5, rpi.P1_6)
	printPin(invalid, m, max, "", 7, rpi.P1_7, rpi.P1_8)
	printPin(invalid, m, max, "", 9, rpi.P1_9, rpi.P1_10)
	printPin(invalid, m, max, "", 11, rpi.P1_11, rpi.P1_12)
	printPin(invalid, m, max, "", 13, rpi.P1_13, rpi.P1_14)
	printPin(invalid, m, max, "", 15, rpi.P1_15, rpi.P1_16)
	printPin(invalid, m, max, "", 17, rpi.P1_17, rpi.P1_18)
	printPin(invalid, m, max, "", 19, rpi.P1_19, rpi.P1_20)
	printPin(invalid, m, max, "", 21, rpi.P1_21, rpi.P1_22)
	printPin(invalid, m, max, "", 23, rpi.P1_23, rpi.P1_24)
	printPin(invalid, m, max, "", 25, rpi.P1_25, rpi.P1_26)
	if rpi.IsConnected(rpi.P1_27) || invalid {
		printPin(invalid, m, max, "", 27, rpi.P1_27, rpi.P1_28)
		printPin(invalid, m, max, "", 29, rpi.P1_29, rpi.P1_30)
		printPin(invalid, m, max, "", 31, rpi.P1_31, rpi.P1_32)
		printPin(invalid, m, max, "", 33, rpi.P1_33, rpi.P1_34)
		printPin(invalid, m, max, "", 35, rpi.P1_35, rpi.P1_36)
		printPin(invalid, m, max, "", 37, rpi.P1_37, rpi.P1_38)
		printPin(invalid, m, max, "", 39, rpi.P1_39, rpi.P1_40)
	}
	if rpi.IsConnected(rpi.P5_1) || invalid {
		fmt.Print("\n")
		printPin(invalid, m, max, "P5:", 1, rpi.P5_1, rpi.P5_2)
		printPin(invalid, m, max, "", 3, rpi.P5_3, rpi.P5_4)
		printPin(invalid, m, max, "", 5, rpi.P5_5, rpi.P5_6)
		printPin(invalid, m, max, "", 7, rpi.P5_7, rpi.P5_8)
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
		fmt.Printf("Version: %d  MaxSpeed: %dMhz\n", rpi.Version, cpu.MaxSpeed/1000000)
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
