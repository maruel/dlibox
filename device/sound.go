// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package device

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/maruel/dlibox/nodes"
	"github.com/maruel/msgbus"
)

type soundDev struct {
	NodeBase
	Cfg  *nodes.Sound
	root string
}

func (s *soundDev) init(b msgbus.Bus) error {
	c, err := b.Subscribe("sound", msgbus.ExactlyOnce)
	if err != nil {
		return err
	}
	go func() {
		for {
			for msg := range c {
				s.onMsg(msg)
			}
		}
	}()
	return nil
}

func (s *soundDev) onMsg(m msgbus.Message) {
	p := filepath.Join(s.root, filepath.Base(string(m.Payload))+".wav")
	if _, err := os.Stat(p); err != nil {
		log.Printf("%s: file not present: %s", s, p)
		return
	}
	play(s.Cfg.DeviceID, p)
}

// play plays a wav.
//
// If a sound was already playing, the command will fail, ignore.
func play(card, path string) {
	c := exec.Command("aplay", "-D", card, path)
	go c.Run()
}
