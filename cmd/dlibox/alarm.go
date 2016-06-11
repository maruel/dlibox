// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"log"
	"time"

	"github.com/maruel/dlibox-go/anim1d"
)

type WeekdayBit int

const (
	Sunday WeekdayBit = 1 << iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

func (w WeekdayBit) IsEnabledFor(d time.Weekday) bool {
	return (w & WeekdayBit(1<<uint(d))) != 0
}

type Alarm struct {
	Enabled bool
	Hour    int
	Minute  int
	Days    WeekdayBit
	Pattern string // JSON serialized pattern.
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

func (a *Alarm) Reset(p *anim1d.Painter) error {
	if a.timer != nil {
		a.timer.Stop()
		a.timer = nil
	}
	// TODO(maruel): Make sure all alarms are valid.
	now := time.Now()
	if next := a.Next(now); !next.IsZero() {
		a.timer = time.AfterFunc(next.Sub(now), func() {
			if err := p.SetPattern(a.Pattern); err != nil {
				log.Printf("failed to unmarshal pattern %q", a.Pattern)
			}
			a.Reset(p)
		})
	}
	return nil
}

type Alarms []Alarm

func (a Alarms) Reset(p *anim1d.Painter) error {
	var err error
	for _, a := range a {
		if err1 := a.Reset(p); err1 != nil {
			err = err1
		}
	}
	return err
}
