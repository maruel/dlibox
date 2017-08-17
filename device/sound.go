// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package device

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/maruel/dlibox/nodes/sound"
	"github.com/maruel/msgbus"
)

type soundDev struct {
	cfg  *sound.Dev
	b    msgbus.Bus
	root string
}

func initSound(b msgbus.Bus, cfg *sound.Dev) (*soundDev, error) {
	s := &soundDev{b: b, cfg: cfg}
	c, err := b.Subscribe("sound", msgbus.BestEffort)
	if err != nil {
		return nil, err
	}
	go func() {
		for {
			for msg := range c {
				s.onMsg(msg)
			}
		}
	}()
	return s, nil
}

func (s *soundDev) Close() error {
	s.b.Unsubscribe("sound")
	return nil
}

func (s *soundDev) onMsg(m msgbus.Message) {
	p := filepath.Join(s.root, filepath.Base(string(m.Payload))+".wav")
	if _, err := os.Stat(p); err != nil {
		log.Printf("sound: file not present: %s", p)
		return
	}
	play(s.cfg.DeviceID, p)
}

// play plays a wav.
//
// If a sound was already playing, the command will fail, ignore.
func play(card, path string) {
	c := exec.Command("aplay", "-D", card, path)
	go c.Run()
}
