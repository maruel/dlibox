// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Package alarm defines the alarms based on the time of the day.
package alarm

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/maruel/dlibox/controller/rules"
	"github.com/maruel/msgbus"
)

// WeekdayBit is a bitmask of each day.
type WeekdayBit int

// Week days.
const (
	Sunday WeekdayBit = 1 << iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
	lastDay
)

const weekdayLetters = "SMTWTFS"

// IsEnabledFor returns true if the bitmask is set for this week day.
func (w WeekdayBit) IsEnabledFor(d time.Weekday) bool {
	return (w & WeekdayBit(1<<uint(d))) != 0
}

func (w WeekdayBit) String() string {
	var out [7]rune
	for i := uint(0); i < 7; i++ {
		if (w & (1 << i)) != 0 {
			out[i] = rune(weekdayLetters[i])
		} else {
			out[i] = '•'
		}
	}
	return string(out[:])
}

// Alarm represents a single alarm.
type Alarm struct {
	Enabled bool
	Hour    int
	Minute  int
	Days    WeekdayBit
	Cmd     rules.Command
	timer   *time.Timer
}

// Next returns when the next trigger should be according to the alarm
// schedule.
//
// Return 0 if not enabled.
func (a *Alarm) Next(now time.Time) time.Time {
	if a.Enabled && a.Days != 0 {
		out := time.Date(now.Year(), now.Month(), now.Day(), a.Hour, a.Minute, 0, 0, now.Location())
		if a.Days.IsEnabledFor(now.Weekday()) {
			if now.Hour() < a.Hour || (now.Hour() == a.Hour && now.Minute() < a.Minute) {
				return out
			}
		}
		for i := 1; i < 8; i++ {
			out = out.Add(24 * time.Hour)
			if a.Days.IsEnabledFor(out.Weekday()) {
				return out
			}
		}
		panic("unexpected code")
	}
	return time.Time{}
}

// Reset reinitializes with a message bus.
func (a *Alarm) Reset(b msgbus.Bus) error {
	if a.timer != nil {
		a.timer.Stop()
		a.timer = nil
	}
	now := time.Now()
	if next := a.Next(now); !next.IsZero() {
		a.timer = time.AfterFunc(next.Sub(now), func() {
			if err := b.Publish(a.Cmd.ToMsg(), msgbus.BestEffort, false); err != nil {
				log.Printf("failed to publish command %v", a.Cmd)
			}
			a.Reset(b)
		})
	}
	return nil
}

// Validate confirms the settings are valid.
func (a *Alarm) Validate() error {
	if a.Days >= lastDay {
		return errors.New("invalid days")
	}
	if a.Hour < 0 || a.Hour >= 24 {
		return errors.New("invalid hour")
	}
	if a.Minute < 0 || a.Minute >= 60 {
		return errors.New("invalid minute")
	}
	return a.Cmd.Validate()
}

func (a *Alarm) String() string {
	out := fmt.Sprintf("%02d:%02d (%s)", a.Hour, a.Minute, a.Days)
	if !a.Enabled {
		out += " (disabled)"
	}
	return out
}

// Config is what should be serialized.
type Config struct {
	Alarms map[string]*Alarm
}

// Init initializes the timers.
func Init(b msgbus.Bus, config *Config) error {
	var err error
	for _, a := range config.Alarms {
		if err1 := a.Reset(b); err1 != nil {
			err = err1
		}
	}
	return err
}

// ResetDefault initializes the default alarms.
func (c *Config) ResetDefault() {
	c.Alarms = map[string]*Alarm{
		"Morning weekdays": {
			Enabled: true,
			Hour:    6,
			Minute:  35,
			Days:    Monday | Tuesday | Wednesday | Thursday | Friday,
			Cmd:     rules.Command{Topic: "painter/setautomated", Payload: "#010203"},
		},
		"Monring weekends": {
			Enabled: true,
			Hour:    6,
			Minute:  55,
			Days:    Saturday | Sunday,
			Cmd:     rules.Command{Topic: "painter/setautomated", Payload: "\"#000000\""},
		},
		"Evening weekdays": {
			Enabled: true,
			Hour:    19,
			Minute:  00,
			Days:    Monday | Tuesday | Wednesday | Thursday | Friday,
			Cmd:     rules.Command{Topic: "painter/setautomated", Payload: "\"#010001\""},
		},
	}
}

// Validate confirms the settings are valid.
func (c *Config) Validate() error {
	for name, a := range c.Alarms {
		if len(name) == 0 {
			return errors.New("alarm without a name")
		}
		if err := a.Validate(); err != nil {
			return fmt.Errorf("can't validate alarm %s: %v", name, err)
		}
	}
	return nil
}
