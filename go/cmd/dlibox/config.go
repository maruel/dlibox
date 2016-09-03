// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/maruel/dlibox/go/anim1d"
	"github.com/pkg/errors"
)

type APA102 struct {
	// Speed of the transfer.
	SPIspeed int64
	// Number of lights controlled by this device. If lower than the actual
	// number of lights, the remaining lights will flash oddly.
	NumberLights   int
	StartupPattern string
}

// Config stores the configuration for this specific host.
type Config struct {
	Alarms   Alarms
	APA102   APA102
	Patterns []string // List of recent patterns. The first is the oldest.
}

const morning = "{\"After\":\"#000000\",\"Before\":{\"After\":\"#ffffff\",\"Before\":{\"After\":\"#ff7f00\",\"Before\":\"#000000\",\"Curve\":\"direct\",\"DurationMS\":6000000,\"OffsetMS\":0,\"_type\":\"Transition\"},\"Curve\":\"direct\",\"DurationMS\":6000000,\"OffsetMS\":6000000,\"_type\":\"Transition\"},\"Curve\":\"direct\",\"DurationMS\":600000,\"OffsetMS\":18000000,\"_type\":\"Transition\"}"

func (c *Config) ResetDefault() {
	*c = Config{
		Alarms: Alarms{
			{
				Enabled: true,
				Hour:    6,
				Minute:  35,
				Days:    Monday | Tuesday | Wednesday | Thursday | Friday,
				Pattern: morning,
			},
			{
				Enabled: true,
				Hour:    6,
				Minute:  55,
				Days:    Saturday | Sunday,
				Pattern: "\"#000000\"",
			},
			{
				Enabled: true,
				Hour:    19,
				Minute:  00,
				Days:    Monday | Tuesday | Wednesday | Thursday | Friday,
				Pattern: "\"#010001\"",
			},
		},
		APA102: APA102{
			SPIspeed:       10000000,
			NumberLights:   150,
			StartupPattern: "\"#000001\"",
		},
		Patterns: []string{
			"{\"_type\":\"Aurore\"}",
			"{\"Child\":{\"Frame\":\"Lff0000ff0000ff0000ff0000ff0000ffffffffffffffffffffffffffffff\",\"_type\":\"Repeated\"},\"MovesPerSec\":6,\"_type\":\"Rotate\"}",
			"{\"Patterns\":[{\"_type\":\"Aurore\"},{\"Seed\":0,\"Stars\":null,\"_type\":\"NightStars\"},{\"AverageDelay\":0,\"Duration\":0,\"_type\":\"WishingStar\"}],\"Weights\":[1,1,1],\"_type\":\"Mixer\"}",
			"{\"Curve\":\"easeinout\",\"DurationShowMS\":1000000,\"DurationTransitionMS\":1000000,\"Patterns\":[\"#ff0000\",\"#00ff00\",\"#0000ff\"],\"_type\":\"Loop\"}",
			"{\"Curve\":\"direct\",\"Left\":\"#000000\",\"Right\":\"#0000ff\",\"_type\":\"Gradient\"}",
			"{\"Curve\":\"direct\",\"Left\":\"#000000\",\"Right\":\"#ff0000\",\"_type\":\"Gradient\"}",
			"{\"Curve\":\"direct\",\"Left\":\"#000000\",\"Right\":\"#00ff00\",\"_type\":\"Gradient\"}",
			"{\"Curve\":\"direct\",\"Left\":\"#000000\",\"Right\":\"#ffffff\",\"_type\":\"Gradient\"}",
			"{\"Child\":\"Lff0000ff0000ee0000dd0000cc0000bb0000aa0000990000880000770000660000550000440000330000220000110000\",\"MovesPerSec\":30,\"_type\":\"PingPong\"}",
			"{\"After\":\"#000000\",\"Before\":{\"After\":\"#ffffff\",\"Before\":{\"After\":\"#ff7f00\",\"Before\":\"#000000\",\"Curve\":\"direct\",\"DurationMS\":0,\"OffsetMS\":0,\"_type\":\"Transition\"},\"Curve\":\"direct\",\"DurationMS\":0,\"OffsetMS\":0,\"_type\":\"Transition\"},\"Curve\":\"direct\",\"DurationMS\":0,\"OffsetMS\":0,\"_type\":\"Transition\"}",
			"\"#000000\"",
			"{\"Child\":\"Lffffff\",\"MovesPerSec\":30,\"_type\":\"PingPong\"}",
			"{\"Curve\":\"easeinout\",\"DurationShowMS\":1000000,\"DurationTransitionMS\":10000000,\"Patterns\":[\"#ff0000\",\"#ff7f00\",\"#ffff00\",\"#00ff00\",\"#0000ff\",\"#4b0082\",\"#8b00ff\"],\"_type\":\"Loop\"}",
			"\"Rainbow\"",
			"{\"Seed\":0,\"Stars\":null,\"_type\":\"NightStars\"}",
			"{\"Child\":\"L010001ff000000ff000000ff\",\"_type\":\"Chronometer\"}",
			morning,
		},
	}
}

func (c *Config) Load(n string) error {
	f, err := os.Open(n)
	if err != nil {
		if os.IsNotExist(err) {
			// Ignore if the file is not present.
			return nil
		}
		return err
	}
	defer f.Close()
	d := json.NewDecoder(f)
	d.UseNumber()
	if err := d.Decode(c); err != nil {
		return err
	}
	return c.verify()
}

func (c *Config) Save(n string) error {
	if err := c.verify(); err != nil {
		return err
	}
	b, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	f, err := os.Create(n)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(append(b, '\n'))
	return err
}

// Inject moves the pattern at the top of LRU cache.
func (c *Config) Inject(pattern string) {
	for i, old := range c.Patterns {
		if old == pattern {
			copy(c.Patterns[i:], c.Patterns[i+1:])
			c.Patterns = c.Patterns[:len(c.Patterns)-1]
			break
		}
	}
	if len(c.Patterns) < 25 {
		c.Patterns = append(c.Patterns, "")
	}
	copy(c.Patterns[1:], c.Patterns)
	c.Patterns[0] = pattern
}

// check verifies that the pattern can be decoded, reencoded and that the
// format is the canonical one.
func check(s string) error {
	var p anim1d.SPattern
	if err := json.Unmarshal([]byte(s), &p); err != nil {
		return err
	}
	b, err := p.MarshalJSON()
	if err != nil {
		return err
	}
	if string(b) != s {
		return fmt.Errorf("pattern not in canonical format: expected %v; got %v", string(b), s)
	}
	return nil
}

func (c *Config) verify() error {
	if err := check(c.APA102.StartupPattern); err != nil {
		return errors.Wrap(err, "can't load startup pattern")
	}
	for _, a := range c.Alarms {
		if err := check(a.Pattern); err != nil {
			return errors.Wrap(err, fmt.Sprintf("can't load pattern for alarm %s", a))
		}
	}
	for i, s := range c.Patterns {
		if err := check(s); err != nil {
			return errors.Wrap(err, fmt.Sprintf("can't load recent pattern %d", i))
		}
	}
	return nil
}

type ConfigMgr struct {
	Config
	path string
}

func (c *ConfigMgr) Load() error {
	home, err := getHome()
	if err != nil {
		return err
	}
	c.path = filepath.Join(home, "dlibox.json")
	return c.Config.Load(c.path)
}

func (c *ConfigMgr) Init(p *anim1d.Painter) error {
	if err := c.Alarms.Reset(p); err != nil {
		return nil
	}
	return p.SetPattern(c.APA102.StartupPattern)
}

func (c *ConfigMgr) Close() error {
	if len(c.path) != 0 {
		return c.Save(c.path)
	}
	return nil
}

// getHome returns the home directory even when cross compiled.
//
// When cross compiling, user.Current() fails.
func getHome() (string, error) {
	if u, err := user.Current(); err == nil {
		return u.HomeDir, nil
	}
	if s := os.Getenv("HOME"); len(s) != 0 {
		return s, nil
	}
	return "", errors.New("can't find HOME")
}
