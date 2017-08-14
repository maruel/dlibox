// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"log"
	"sync"

	"github.com/maruel/dlibox/go/msgbus"
	"github.com/maruel/psf"
	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/devices/ssd1306"
	"periph.io/x/periph/devices/ssd1306/image1bit"
)

// Display contains small embedded display settings.
type Display struct {
	sync.Mutex
	I2CBus string
}

func (d *Display) ResetDefault() {
	d.Lock()
	d.Unlock()
	d.I2CBus = ""
}

func (d *Display) Validate() error {
	d.Lock()
	d.Unlock()
	return nil
}

func initDisplay(b msgbus.Bus, config *Display) (*display, error) {
	i2cBus, err := i2creg.Open(config.I2CBus)
	if err != nil {
		return nil, err
	}
	d, err := ssd1306.NewI2C(i2cBus, 128, 64, false)
	if err != nil {
		return nil, err
	}
	f12, err := psf.Load("Terminus12x6")
	if err != nil {
		return nil, err
	}
	f20, err := psf.Load("Terminus20x10")
	if err != nil {
		return nil, err
	}
	bounds := d.Bounds()
	img := image1bit.NewVerticalLSB(bounds)
	f20.Draw(img, 0, 0, image1bit.On, nil, "dlibox!")
	f12.Draw(img, 0, bounds.Dy()-f12.H-1, image1bit.On, nil, "is awesome")
	if _, err = d.Write(img.Pix); err != nil {
		return nil, err
	}
	c, err := b.Subscribe("display/#", msgbus.BestEffort)
	if err != nil {
		return nil, err
	}
	disp := &display{d, b, img, f12, f20}
	go func() {
		for msg := range c {
			disp.onMsg(msg)
		}
	}()
	return disp, nil
}

type display struct {
	d   *ssd1306.Dev
	b   msgbus.Bus
	img *image1bit.VerticalLSB
	f12 *psf.Font
	f20 *psf.Font
}

func (d *display) Close() error {
	d.b.Unsubscribe("display/#")
	return nil
}

func (d *display) onMsg(msg msgbus.Message) {
	switch msg.Topic {
	case "display/settext":
		d.f20.Draw(d.img, 0, 0, image1bit.On, nil, string(msg.Payload))
		if _, err := d.d.Write(d.img.Pix); err != nil {
			log.Printf("display write failure: %# v", msg)
		}
	default:
		log.Printf("display unknown msg: %# v", msg)
	}
}
