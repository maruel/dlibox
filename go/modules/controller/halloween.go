// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package controller

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/maruel/dlibox/go/modules"
	"github.com/maruel/dlibox/go/msgbus"
)

type Halloween struct {
	sync.Mutex
	Enabled   bool
	Modes     map[string]State
	Cmds      map[State][]modules.Command
	IdleAfter int // seconds
}

func (h *Halloween) ResetDefault() {
	h.Lock()
	defer h.Unlock()
	h.Enabled = false
	h.IdleAfter = 15
}

func (h *Halloween) Validate() error {
	h.Lock()
	defer h.Unlock()
	for k, v := range h.Modes {
		if !v.Valid() {
			return fmt.Errorf("halloween: Modes[%q] has invalid state %q", k, v)
		}
	}
	for k := range h.Cmds {
		if !k.Valid() {
			return fmt.Errorf("halloween: Cmds[%q] is an invalid state", k)
		}
	}
	return nil
}

func initHalloween(b msgbus.Bus, config *Halloween) (*halloween, error) {
	if !config.Enabled {
		return nil, errors.New("not the controller")
	}
	if config.Modes == nil {
		return nil, errors.New("halloween Modes is missing")
	}
	if config.Cmds == nil {
		return nil, errors.New("halloween Cmds is missing")
	}

	h := &halloween{
		b:      b,
		config: config,
		state:  Idle,
	}
	// Listen to all messages, since we don't know the one that could be keys in
	// the config. Technically we know but it's easier to just get them all.
	// Revisit this decision if it becomes a problem.
	c, err := b.Subscribe("//#", msgbus.BestEffort)
	if err != nil {
		return nil, err
	}
	go func() {
		for {
			for msg := range c {
				h.onMsg(msg)
			}
		}
	}()
	// Trigger Idle on startup.
	h.publishState(h.state)
	return h, nil
}

// State is the state machine for the incoming children.
type State string

const (
	// Idle is the animation while nothing in happening.
	Idle State = "idle"
	// Incoming is when little monsters (children) are walking in front of the
	// house.
	Incoming State = "incoming"
	// Porch is when the children are in front of the door.
	Porch State = "porch"
)

func (s State) Valid() bool {
	switch s {
	case Idle, Incoming, Porch:
		return true
	}
	return false
}

type halloween struct {
	b         msgbus.Bus
	config    *Halloween
	state     State
	timerIdle *time.Timer
}

func (h *halloween) Close() error {
	h.config.Lock()
	defer h.config.Unlock()
	if h.timerIdle != nil {
		h.timerIdle.Stop()
	}
	h.b.Unsubscribe("//#")
	return nil
}

func (h *halloween) onMsg(m msgbus.Message) {
	h.config.Lock()
	defer h.config.Unlock()
	if s, ok := h.config.Modes[m.Topic]; ok {
		if h.state == s {
			// Didn't change state, do not trigger anything.
			return
		}
		if s == Porch && h.state == Incoming {
			// Ignore, we'll wait for going back to idle first.
			return
		}
		if s != Idle {
			// Reset the timer. Note that the timer is only armed when the switch is
			// triggered by Modes. If someone sends a state change manually via
			// mosquitto_pub -t dlibox/halloween/state, the timer will not be armed.
			if h.timerIdle != nil {
				h.timerIdle.Stop()
			}
			if h.config.IdleAfter != 0 {
				h.timerIdle = time.AfterFunc(time.Duration(h.config.IdleAfter)*time.Second, h.setIdle)
			}
		}
		// Broadcast the new state. onMsg() will be called again with this state.
		h.publishState(s)
		return
	}

	if m.Topic == "dlibox/halloween/state" {
		s := State(m.Payload)
		if !s.Valid() {
			log.Printf("halloween: state is invalid: %q", s)
			return
		}
		h.state = s
		for _, cmd := range h.config.Cmds[h.state] {
			// TODO(maruel): Run them in parallel.
			if err := h.b.Publish(cmd.ToMsg(), msgbus.BestEffort, false); err != nil {
				log.Printf("halloween: %s->%v: %v", h.state, cmd, err)
			}
		}
	}
}

func (h *halloween) setIdle() {
	h.config.Lock()
	defer h.config.Unlock()
	if h.state == Idle {
		return
	}
	log.Printf("halloween: going back idle")
	h.publishState(Idle)
}

func (h *halloween) publishState(s State) {
	if err := h.b.Publish(msgbus.Message{"//dlibox/halloween/state", []byte(s)}, msgbus.BestEffort, true); err != nil {
		log.Printf("halloween: failed to publish state: %v", err)
	}
}
