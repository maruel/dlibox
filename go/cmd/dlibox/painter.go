// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"github.com/google/pio/devices"
	"github.com/maruel/dlibox/go/anim1d"
)

func initPainter(leds devices.Display, fps int, config *Painter) (*anim1d.Painter, error) {
	p := anim1d.NewPainter(leds, fps)
	if len(config.Last) != 0 {
		if err := p.SetPattern(string(config.Last)); err != nil {
			return nil, err
		}
	}
	if len(config.Startup) != 0 {
		if err := p.SetPattern(string(config.Startup)); err != nil {
			return nil, err
		}
	}
	return p, nil
}
