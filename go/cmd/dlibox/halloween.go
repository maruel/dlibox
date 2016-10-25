// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"errors"
	"log"
	"reflect"
	"sync"
	"time"

	"github.com/maruel/dlibox/go/modules"
)

type Halloween struct {
	sync.Mutex
	Enabled bool
	Modes   map[string]State
	Cmds    map[State][]Command
}

func (h *Halloween) ResetDefault() {
	h.Lock()
	defer h.Unlock()
	h.Enabled = false
}

func (h *Halloween) Validate() error {
	h.Lock()
	defer h.Unlock()
	return nil
}

func merge(chans ...<-chan modules.Message) <-chan modules.Message {
	out := make(chan modules.Message)
	c := make([]reflect.SelectCase, len(chans))
	for i := range chans {
		c[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(chans[i])}
	}
	go func() {
		defer close(out)
		for {
			i, msg, ok := reflect.Select(c)
			if !ok {
				if len(c) == 1 {
					break
				}
				if i != len(c)-1 {
					copy(c[i:], c[i+1:])
				}
				c = c[:len(c)-1]
				continue
			}
			out <- msg.Interface().(modules.Message)
		}
	}()
	return out
}

func initHalloween(b modules.Bus, config *Halloween) (*halloween, error) {
	if !config.Enabled {
		return nil, errors.New("not the controller")
	}
	if config.Modes == nil {
		return nil, errors.New("halloween Modes is missing")
	}
	if config.Cmds == nil {
		return nil, errors.New("halloween Cmds is missing")
	}

	h := &halloween{b: b, config: config, timer: time.NewTimer(0)}
	c1, err := b.Subscribe("//dlibox/+/pir", modules.ExactlyOnce)
	if err != nil {
		return nil, err
	}
	c2, err := b.Subscribe("//dlibox/halloween/#", modules.ExactlyOnce)
	if err != nil {
		return nil, err
	}
	c := merge(c1, c2)
	go func() {
		for {
			for msg := range c {
				h.onMsg(msg)
			}
		}
	}()
	return h, nil
}

// State is the state machine for the incoming children.
type State string

const (
	Idle     State = "idle"
	Incoming State = "incoming"
	Balcon   State = "balcon"
	Back     State = "back"
)

type halloween struct {
	b         modules.Bus
	config    *Halloween
	state     State
	lastLaugh time.Time
	timer     *time.Timer
}

func (h *halloween) Close() error {
	var err error
	if err1 := h.b.Unsubscribe("//dlibox/+/pir"); err1 != nil {
		log.Printf("failed to unsubscribe: //dlibox/+/pir: %v", err1)
		err = err1
	}
	if err1 := h.b.Unsubscribe("//dlibox/halloween/#"); err1 != nil {
		log.Printf("failed to unsubscribe: //dlibox/halloween/#: %v", err1)
		err = err1
	}
	return err
}

func (h *halloween) onMsg(m modules.Message) {
	h.config.Lock()
	defer h.config.Unlock()
	if s, ok := h.config.Modes[m.Topic]; ok {
		if h.state == s {
			// Reset the timer.
			//h.timer.Reset(d)
			return
		}
		h.state = s
		for _, cmd := range h.config.Cmds[s] {
			if err := h.b.Publish(cmd.ToMsg(), modules.ExactlyOnce, false); err != nil {
				log.Printf("halloween: %s: %v", s, cmd)
			}
		}
		return
	}
	switch {
	default:
	}
}
