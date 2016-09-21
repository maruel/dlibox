// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// pioselftest is a small app to verify that basic GPIO pin functionality work.
package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/maruel/dlibox/go/pio/host"
	"github.com/maruel/dlibox/go/pio/host/hal/pins"
)

func getPin(s string) (host.PinIO, error) {
	number, err := strconv.Atoi(s)
	if err != nil {
		return nil, err
	}
	p := pins.ByNumber(number)
	if p == nil {
		return nil, errors.New("invalid pin number")
	}
	return p, nil
}

const shortDelay = time.Microsecond
const longDelay = 2 * time.Second

func slowSleep(do bool) {
	if do {
		fmt.Printf("  Sleep(%s)\n", longDelay)
		time.Sleep(longDelay)
	}
}

func waitChan(c <-chan host.Level) (host.Level, bool) {
	select {
	case i := <-c:
		return i, true
	case <-time.After(time.Second):
		return host.Low, false
	}
}

func doEdges(p1, p2 host.PinIO, slow bool) error {
	slowSleep(slow)
	fmt.Printf("  %s.Edges()\n", p1)
	c, err := p1.Edges()
	if err != nil {
		return err
	}
	defer p1.DisableEdges()
	time.Sleep(shortDelay)

	fmt.Printf("  %s.Set(Low)\n", p2)
	p2.Set(host.Low)
	if l, ok := waitChan(c); !ok {
		return errors.New("edge didn't trigger")
	} else if l != host.Low {
		return fmt.Errorf("expected Low, got %s", l)
	}

	slowSleep(slow)
	fmt.Printf("  %s.Set(High)\n", p2)
	p2.Set(host.High)
	if l, ok := waitChan(c); !ok {
		return errors.New("edge didn't trigger")
	} else if l != host.High {
		return fmt.Errorf("expected High, got %s", l)
	}

	slowSleep(slow)
	fmt.Printf("  %s.Set(Low)\n", p2)
	p2.Set(host.Low)
	if l, ok := waitChan(c); !ok {
		return errors.New("edge didn't trigger")
	} else if l != host.Low {
		return fmt.Errorf("expected Low, got %s", l)
	}
	return nil
}

func doCycle(p1, p2 host.PinIO, noEdge, noPull, slow bool) error {
	// Do a 'shortDelay' sleep between writting and reading because there can be
	// propagation delay in the wire.
	//
	// Random observation, needs to be confirmed:
	// On A64, on some pin the pull resistor is low and can give a 3.3v/2 output
	// when crossing an output at high.
	fmt.Printf("%s -> %s\n", p1, p2)
	pull := host.Float
	if noPull {
		pull = host.PullNoChange
	}
	fmt.Printf("  %s.In(%s)\n", p1, pull)
	if err := p1.In(pull); err != nil {
		return err
	}
	fmt.Printf("  %s.Out()\n", p2)
	if err := p2.Out(); err != nil {
		return err
	}
	fmt.Printf("  %s.Set(Low)\n", p2)
	p2.Set(host.Low)
	time.Sleep(shortDelay)
	fmt.Printf("  -> %s: %s\n- %s: %s\n", p1, p1.Function(), p2, p2.Function())
	if p1.Read() != host.Low {
		return errors.New("read low failure")
	}

	slowSleep(slow)
	fmt.Printf("  %s.Set(High)\n", p2)
	p2.Set(host.High)
	time.Sleep(shortDelay)
	fmt.Printf("  -> %s: %s\n- %s: %s\n", p1, p1.Function(), p2, p2.Function())
	if p1.Read() != host.High {
		return errors.New("read high failure")
	}

	if !noEdge {
		if err := doEdges(p1, p2, slow); err != nil {
			return err
		}
	}

	if !noPull {
		// p1 is float.
		slowSleep(slow)
		fmt.Printf("  %s.In(Down)\n", p2)
		if err := p2.In(host.Down); err != nil {
			return err
		}
		time.Sleep(shortDelay)
		fmt.Printf("  -> %s: %s\n- %s: %s\n", p1, p1.Function(), p2, p2.Function())
		if p1.Read() != host.Low {
			return errors.New("read pull down failure")
		}

		slowSleep(slow)
		fmt.Printf("  %s.In(Up)\n", p2)
		if err := p2.In(host.Up); err != nil {
			return err
		}
		time.Sleep(shortDelay)
		fmt.Printf("  -> %s: %s\n- %s: %s\n", p1, p1.Function(), p2, p2.Function())
		if p1.Read() != host.High {
			return errors.New("read pull up failure")
		}
	}
	return nil
}

func mainImpl() error {
	noEdge := flag.Bool("e", false, "no edge test, necessary when testing without sysfs")
	// This flag should be set automatically when sysfs gpio is detected.
	noPull := flag.Bool("p", false, "no pull test, necessary when testing sysfs gpio")
	slow := flag.Bool("s", false, "slow; insert a second between each step")
	flag.Parse()

	if flag.NArg() != 2 {
		return errors.New("specify the two pins to use; they must be connected together")
	}
	// On Allwinner CPUs, it's a good idea to test specifically the PLx pins,
	// since they use a different register memory block than groups PB to PH.
	p1, err := getPin(flag.Args()[0])
	if err != nil {
		return err
	}
	p2, err := getPin(flag.Args()[1])
	if err != nil {
		return err
	}
	fmt.Printf("Using pins and their current state:\n- %s: %s\n- %s: %s\n\n", p1, p1.Function(), p2, p2.Function())
	if err := doCycle(p1, p2, *noEdge, *noPull, *slow); err != nil {
		return err
	}
	fmt.Print("\n")
	return doCycle(p2, p1, *noEdge, *noPull, *slow)
}

func main() {
	if err := mainImpl(); err != nil {
		fmt.Fprintf(os.Stderr, "pioselfcheck: %s.\n", err)
		os.Exit(1)
	}
}
