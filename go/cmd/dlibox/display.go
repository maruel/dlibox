// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"image"
	"log"
	"sync"

	"github.com/maruel/dlibox/go/donotuse/conn/i2c"
	"github.com/maruel/dlibox/go/donotuse/devices/ssd1306"
	"github.com/maruel/dlibox/go/donotuse/devices/ssd1306/image1bit"
	"github.com/maruel/dlibox/go/modules"
	"github.com/maruel/dlibox/go/psf"
)

// Display contains small embedded display settings.
type Display struct {
	sync.Mutex
	I2CBus int
}

func (d *Display) ResetDefault() {
	d.Lock()
	d.Unlock()
	d.I2CBus = -1
}

func (d *Display) Validate() error {
	d.Lock()
	d.Unlock()
	return nil
}

func initDisplay(b modules.Bus, config *Display) (*display, error) {
	i2cBus, err := i2c.New(config.I2CBus)
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
	img, err := image1bit.New(image.Rect(0, 0, d.W, d.H))
	if err != nil {
		return nil, err
	}
	f20.Draw(img, 0, 0, image1bit.On, nil, "dlibox!")
	f12.Draw(img, 0, d.H-f12.H-1, image1bit.On, nil, "is awesome")
	if _, err = d.Write(img.Buf); err != nil {
		return nil, err
	}
	c, err := b.Subscribe("display/#", modules.BestEffort)
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
	b   modules.Bus
	img *image1bit.Image
	f12 *psf.Font
	f20 *psf.Font
}

func (d *display) Close() error {
	err := d.b.Unsubscribe("display/#")
	if err != nil {
		log.Printf("failed to unsubscribe: display/#: %v", err)
	}
	return err
}

func (d *display) onMsg(msg modules.Message) {
	switch msg.Topic {
	case "display/settext":
		d.f20.Draw(d.img, 0, 0, image1bit.On, nil, string(msg.Payload))
		if _, err := d.d.Write(d.img.Buf); err != nil {
			log.Printf("display write failure: %# v", msg)
		}
	default:
		log.Printf("display unknown msg: %# v", msg)
	}
}
