// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"time"

	"github.com/maruel/dotstar/anim1d"
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
	Pattern anim1d.Pattern
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

func (a *Alarm) SetupTimer(now time.Time, p *anim1d.Painter) *time.Timer {
	if next := a.Next(now); next.IsZero() {
		var t *time.Timer
		t = time.AfterFunc(next.Sub(now), func() {
			p.SetPattern(a.Pattern)
			// Rearm on the next event.
			now = time.Now()
			if next = a.Next(now); next.IsZero() {
				t.Reset(next.Sub(now))
			}
		})
		return t
	}
	return nil
}

type Alarms []Alarm

func (a Alarms) AlarmLoop(p *anim1d.Painter) {
	now := time.Now()
	var timers []*time.Timer
	for _, a := range a {
		if t := a.SetupTimer(now, p); t != nil {
			timers = append(timers, t)
		}
	}
	for _, t := range timers {
		t.Stop()
	}
}
