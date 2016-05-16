// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

// github.com/stianeikeland/go-rpio

/*
func listenToPin(pinNumber int, p *anim1d.Painter, r *anim1d.PatternRegistry) {
	pin := rpio.Pin(pinNumber)
	pin.Input()
	pin.PullUp()
	last := rpio.High
	names := make([]string, 0, len(r.Patterns))
	for n := range r.Patterns {
		names = append(names, n)
	}
	sort.Strings(names)
	index := 0
	for {
		// Types of press:
		// - Short press (<2s)
		// - 2s press
		// - 4s press
		// - double-click (incompatible with repeated short press)
		//
		// Functions:
		// - Bonne nuit
		// - Next / Prev
		// - Éteindre (longer press après bonne nuit?)
		if state := pin.Read(); state != last {
			last = state
			if state == rpio.Low {
				index = (index + 1) % len(names)
				// TODO(maruel): Data race.
				p.SetPattern(r.Patterns[names[index]])
			}
		}
		select {
		case <-interrupt.Channel:
			return
		case <-time.After(time.Millisecond):
		}
	}
}

	pinNumber := flag.Int("pin", 0, "pin to listen to")

	if *pinNumber != 0 {
		// Open and map memory to access gpio, check for errors
		if err := rpio.Open(); err != nil {
			return err
		}
		defer rpio.Close()
		go listenToPin(*pinNumber, p, registry)
	}
*/
