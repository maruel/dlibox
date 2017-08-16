// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package device

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/maruel/dlibox/modules/shared"
	"github.com/maruel/dlibox/msgbus"
)

type Sound struct {
	sync.Mutex
	Card string
	Root string
}

func (s *Sound) ResetDefault() {
	s.Lock()
	defer s.Unlock()
	s.Card = "hw:1,0"
	s.Root = shared.Home()
}

func (s *Sound) Validate() error {
	s.Lock()
	defer s.Unlock()
	return nil
}

func initSound(b msgbus.Bus, config *Sound) (*sound, error) {
	s := &sound{b: b, config: config}
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

// play plays a wav.
//
// If a sound was already playing, the command will fail, ignore.
func play(card, path string) {
	c := exec.Command("aplay", "-D", card, path)
	go c.Run()
}

type sound struct {
	b      msgbus.Bus
	config *Sound
}

func (s *sound) Close() error {
	s.b.Unsubscribe("sound")
	return nil
}

func (s *sound) onMsg(m msgbus.Message) {
	s.config.Lock()
	defer s.config.Unlock()
	p := filepath.Join(s.config.Root, filepath.Base(string(m.Payload))+".wav")
	if _, err := os.Stat(p); err != nil {
		log.Printf("sound: file not present: %s", p)
		return
	}
	play(s.config.Card, p)
}
