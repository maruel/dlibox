// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// pioselftest verifies that basic GPIO pin functionality work.
package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/maruel/dlibox/go/pio/host"
	"github.com/maruel/dlibox/go/pio/host/sysfs"
	"github.com/maruel/dlibox/go/pio/protocols/gpio"
)

const (
	// Sleep for a short delay as there can be some capacitance on the line,
	// requiring a few CPU cycles before the input stabilizes to the new value.
	shortDelay = time.Nanosecond

	// Purely to help diagnose issues.
	longDelay = 2 * time.Second
)

func getPin(s string, useSysfs bool) (gpio.PinIO, error) {
	number, err := strconv.Atoi(s)
	if err != nil {
		return nil, err
	}
	var p gpio.PinIO
	if useSysfs {
		p, err = sysfs.PinByNumber(number)
		if err != nil {
			return nil, err
		}
	} else {
		p = gpio.ByNumber(number)
	}
	if p == nil {
		return nil, errors.New("invalid pin number")
	}
	return p, nil
}

func slowSleep(do bool) {
	if do {
		fmt.Printf("  Sleep(%s)\n", longDelay)
		time.Sleep(longDelay)
	}
}

func waitChan(c <-chan gpio.Level) (gpio.Level, bool) {
	select {
	case i := <-c:
		return i, true
	case <-time.After(time.Second):
		return gpio.Low, false
	}
}

func doEdgesInner(p1, p2 gpio.PinIO, slow bool, c <-chan gpio.Level) error {
	time.Sleep(shortDelay)

	fmt.Printf("    %s.Out(Low)\n", p2)
	if err := p2.Out(gpio.Low); err != nil {
		return err
	}
	if l, ok := waitChan(c); !ok {
		return errors.New("edge didn't trigger")
	} else {
		fmt.Printf("    %s <- %s\n", l, p1)
		if l != gpio.Low {
			return fmt.Errorf("%s: expected Low, got %s", p1, l)
		}
	}

	slowSleep(slow)
	fmt.Printf("    %s.Out(High)\n", p2)
	if err := p2.Out(gpio.High); err != nil {
		return err
	}
	if l, ok := waitChan(c); !ok {
		return errors.New("edge didn't trigger")
	} else {
		fmt.Printf("    %s <- %s\n", l, p1)
		if l != gpio.High {
			return fmt.Errorf("%s: expected High, got %s", p1, l)
		}
	}

	slowSleep(slow)
	fmt.Printf("    %s.Out(Low)\n", p2)
	if err := p2.Out(gpio.Low); err != nil {
		return err
	}
	if l, ok := waitChan(c); !ok {
		return errors.New("edge didn't trigger")
	} else {
		fmt.Printf("    %s <- %s\n", l, p1)
		if l != gpio.Low {
			return fmt.Errorf("%s: expected Low, got %s", p1, l)
		}
	}
	return nil
}

func doEdges(p1, p2 gpio.PinIO, slow bool) error {
	fmt.Printf("  Testing edges\n")
	slowSleep(slow)
	fmt.Printf("    %s.Edges()\n", p1)
	c, err := p1.Edges()
	if err != nil {
		return err
	}
	// Create an inner function instead of using defer to simplify debugging.
	err = doEdgesInner(p1, p2, slow, c)
	fmt.Printf("    %s.DisableEdges()\n", p1)
	p1.DisableEdges()
	return err
}

func doPull(p1, p2 gpio.PinIO, slow bool) error {
	fmt.Printf("  Testing pull resistor\n")
	// p1 is float.
	slowSleep(slow)
	fmt.Printf("    %s.In(Down)\n", p2)
	if err := p2.In(gpio.Down); err != nil {
		return err
	}
	time.Sleep(shortDelay)
	fmt.Printf("    -> %s: %s\n    -> %s: %s\n", p1, p1.Function(), p2, p2.Function())
	if p1.Read() != gpio.Low {
		return errors.New("read pull down failure")
	}

	slowSleep(slow)
	fmt.Printf("    %s.In(Up)\n", p2)
	if err := p2.In(gpio.Up); err != nil {
		return err
	}
	time.Sleep(shortDelay)
	fmt.Printf("    -> %s: %s\n    -> %s: %s\n", p1, p1.Function(), p2, p2.Function())
	if p1.Read() != gpio.High {
		return errors.New("read pull up failure")
	}
	return nil
}

func doCycle(p1, p2 gpio.PinIO, noEdge, noPull, slow bool) error {
	// Do a 'shortDelay' sleep between writting and reading because there can be
	// propagation delay in the wire.
	//
	// Random observation, needs to be confirmed:
	// On A64, on some pin the pull resistor is low and can give a 3.3v/2 output
	// when crossing an output at high.
	fmt.Printf("Testing %s -> %s\n", p2, p1)
	fmt.Printf("  Testing base functionality\n")
	pull := gpio.Float
	if noPull {
		pull = gpio.PullNoChange
	}
	fmt.Printf("    %s.In(%s)\n", p1, pull)
	if err := p1.In(pull); err != nil {
		return err
	}
	fmt.Printf("    %s.Out(Low)\n", p2)
	if err := p2.Out(gpio.Low); err != nil {
		return err
	}
	time.Sleep(shortDelay)
	fmt.Printf("    -> %s: %s\n    -> %s: %s\n", p1, p1.Function(), p2, p2.Function())
	if l := p1.Read(); l != gpio.Low {
		return fmt.Errorf("%s: expected to read Low but got %s", p1, l)
	}

	slowSleep(slow)
	fmt.Printf("    %s.Out(High)\n", p2)
	if err := p2.Out(gpio.High); err != nil {
		return err
	}
	time.Sleep(shortDelay)
	fmt.Printf("    -> %s: %s\n    -> %s: %s\n", p1, p1.Function(), p2, p2.Function())
	if l := p1.Read(); l != gpio.High {
		return fmt.Errorf("%s: expected to read High but got %s", p1, l)
	}

	if !noEdge {
		if err := doEdges(p1, p2, slow); err != nil {
			return err
		}
	}

	if !noPull {
		if err := doPull(p1, p2, slow); err != nil {
			return err
		}
	}
	return nil
}

func mainImpl() error {
	noEdge := flag.Bool("e", false, "no edge test, necessary when testing without sysfs")
	slow := flag.Bool("s", false, "slow; insert a second between each step")
	useSysfs := flag.Bool("sysfs", false, "force the use of sysfs")
	flag.Parse()

	if flag.NArg() != 2 {
		return errors.New("specify the two pins to use; they must be connected together")
	}

	state, err := host.Init()
	if err != nil {
		return err
	}
	var finalErr error
	if len(state.Failed) != 0 {
		fmt.Printf("Drivers failed to load:\n")
		for _, f := range state.Failed {
			fmt.Printf("  - %s: %v\n", f.D, f.Err)
		}
		finalErr = errors.New("please fix the drivers. Do you need to run as root?")
	}
	fmt.Printf("Using drivers:\n")
	for _, driver := range state.Loaded {
		fmt.Printf("  - %s\n", driver.String())
	}
	fmt.Printf("Drivers skipped:\n")
	for _, driver := range state.Skipped {
		fmt.Printf("  - %s\n", driver.String())
	}

	// On Allwinner CPUs, it's a good idea to test specifically the PLx pins,
	// since they use a different register memory block than groups PB to PH.
	p1, err := getPin(flag.Args()[0], *useSysfs)
	if err != nil {
		return err
	}
	p2, err := getPin(flag.Args()[1], *useSysfs)
	if err != nil {
		return err
	}

	// Disable pull testing when using sysfs.
	_, noPull := p1.(*sysfs.Pin)
	if noPull {
		fmt.Printf("Skipping input pull resistor on sysfs\n")
	}

	fmt.Printf("Using pins and their current state:\n- %s: %s\n- %s: %s\n\n", p1, p1.Function(), p2, p2.Function())
	err = doCycle(p1, p2, *noEdge, noPull, *slow)
	if err == nil {
		err = doCycle(p2, p1, *noEdge, noPull, *slow)
	}
	if err2 := p1.In(gpio.PullNoChange); err2 != nil {
		fmt.Printf("(Exit) Failed to reset %s as input: %s\n", p1, err2)
	}
	if err2 := p2.In(gpio.PullNoChange); err2 != nil {
		fmt.Printf("(Exit) Failed to reset %s as input: %s\n", p1, err2)
	}
	if err == nil {
		return finalErr
	}
	return err
}

func main() {
	if err := mainImpl(); err != nil {
		fmt.Fprintf(os.Stderr, "pioselfcheck: %s.\n", err)
		os.Exit(1)
	}
}
