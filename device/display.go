// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package device

import (
	"log"

	"github.com/maruel/dlibox/nodes"
	"github.com/maruel/msgbus"
	"github.com/maruel/psf"
	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/devices/ssd1306"
	"periph.io/x/periph/devices/ssd1306/image1bit"
)

type displayDev struct {
	NodeBase
	Cfg *nodes.Display

	d   *ssd1306.Dev
	img *image1bit.VerticalLSB
	f12 *psf.Font
	f20 *psf.Font
}

func (d *displayDev) init(b msgbus.Bus) error {
	i2cBus, err := i2creg.Open(d.Cfg.I2C.ID)
	if err != nil {
		return err
	}
	opts := ssd1306.DefaultOpts
	opts.W = d.Cfg.W
	opts.H = d.Cfg.H
	d.d, err = ssd1306.NewI2C(i2cBus, &opts)
	if err != nil {
		return err
	}
	d.f12, err = psf.Load("Terminus12x6")
	if err != nil {
		return err
	}
	d.f20, err = psf.Load("Terminus20x10")
	if err != nil {
		return err
	}
	bounds := d.d.Bounds()
	d.img = image1bit.NewVerticalLSB(bounds)
	d.f20.Draw(d.img, 0, 0, image1bit.On, nil, "dlibox!")
	d.f12.Draw(d.img, 0, bounds.Dy()-d.f12.H-1, image1bit.On, nil, "is awesome")
	if _, err = d.d.Write(d.img.Pix); err != nil {
		return err
	}
	c, err := b.Subscribe("#", msgbus.ExactlyOnce)
	if err != nil {
		return err
	}
	go func() {
		for msg := range c {
			d.onMsg(msg)
		}
	}()
	return nil
}

/*
func (d *displayDev) Close() error {
	d.b.Unsubscribe("display/#")
	return nil
}
*/

func (d *displayDev) onMsg(msg msgbus.Message) {
	switch msg.Topic {
	case "settext":
		d.f20.Draw(d.img, 0, 0, image1bit.On, nil, string(msg.Payload))
		if _, err := d.d.Write(d.img.Pix); err != nil {
			log.Printf("display write failure: %# v", msg)
		}
	default:
		log.Printf("display unknown msg: %# v", msg)
	}
}
