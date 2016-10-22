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

	h := &halloween{config: config}
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

// state is the state machine for the incoming children.
type state int

const (
	idle state = iota
	incoming
	balcon
	back
)

type halloween struct {
	b      modules.Bus
	config *Halloween
	state  state
	timer  *time.Timer
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
	switch m.Topic {
	default:
	}
}
