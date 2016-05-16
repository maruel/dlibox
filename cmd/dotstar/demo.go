// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

/*
	demoMode := flag.Bool("demo", false, "enable cycling through a few animations as a demo")

	if *demoMode {
		go func() {
			patterns := []struct {
				d int
				p anim1d.Pattern
			}{
				{3, registry.Patterns["Rainbow static"]},
				{10, registry.Patterns["Glow rainbow"]},
				{10, registry.Patterns["Étoile floue"]},
				{7, registry.Patterns["Canne de Noël"]},
				{7, registry.Patterns["K2000"]},
				{5, registry.Patterns["Ping pong"]},
				{5, registry.Patterns["Glow"]},
				{5, registry.Patterns["Glow gris"]},
			}
			i := 0
			p.SetPattern(patterns[i].p)
			delay := time.Duration(patterns[i].d) * time.Second
			for {
				select {
				case <-time.After(delay):
					i = (i + 1) % len(patterns)
					p.SetPattern(patterns[i].p)
					delay = time.Duration(patterns[i].d) * time.Second
				case <-interrupt.Channel:
					return
				}
			}
		}()
		properties = append(properties, "demo")
	}
*/
