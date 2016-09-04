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
	"github.com/maruel/dlibox/go/pio/buses"
	"github.com/pkg/errors"
)

// Pattern is a JSON encoded pattern.
type Pattern string

// validatePattern verifies that the pattern can be decoded, reencoded and that
// the format is the canonical one.
func (p Pattern) Validate() error {
	var obj anim1d.SPattern
	if err := json.Unmarshal([]byte(p), &obj); err != nil {
		return err
	}
	b, err := obj.MarshalJSON()
	if err == nil && Pattern(b) != p {
		err = fmt.Errorf("pattern not in canonical format: expected %v; got %v", string(b), p)
	}
	return err
}

const morning Pattern = "{\"After\":\"#000000\",\"Before\":{\"After\":\"#ffffff\",\"Before\":{\"After\":\"#ff7f00\",\"Before\":\"#000000\",\"Curve\":\"direct\",\"DurationMS\":6000000,\"OffsetMS\":0,\"_type\":\"Transition\"},\"Curve\":\"direct\",\"DurationMS\":6000000,\"OffsetMS\":6000000,\"_type\":\"Transition\"},\"Curve\":\"direct\",\"DurationMS\":600000,\"OffsetMS\":18000000,\"_type\":\"Transition\"}"

// Patterns is the list of recent patterns. The first is the oldest.
type Patterns []Pattern

func (p *Patterns) ResetDefault() {
	*p = Patterns{
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
	}
}

func (p Patterns) Validate() error {
	for i, s := range p {
		if err := s.Validate(); err != nil {
			return errors.Wrap(err, fmt.Sprintf("can't load recent pattern %d", i))
		}
	}
	return nil
}

// Inject moves the pattern at the top of LRU cache.
func (p *Patterns) Inject(pattern string) {
	for i, old := range *p {
		if old == Pattern(pattern) {
			copy((*p)[i:], (*p)[i+1:])
			*p = (*p)[:len(*p)-1]
			break
		}
	}
	if len(*p) < 25 {
		*p = append(*p, "")
	}
	copy((*p)[1:], *p)
	(*p)[0] = Pattern(pattern)
}

//

// APA102 contains light specific settings.
type APA102 struct {
	// Speed of the transfer.
	SPIspeed int64
	// Number of lights controlled by this device. If lower than the actual
	// number of lights, the remaining lights will flash oddly.
	NumberLights int
	// Pattern when the host starts up.
	StartupPattern Pattern
}

// IR contains InfraRed remote information.
type IR struct {
	Mapping map[buses.Key]Pattern
}

// PIR contains a motion detection behavior.
type PIR struct {
	Pin     string
	Pattern Pattern
}

// Settings is all the host settings.
type Settings struct {
	Alarms Alarms
	APA102 APA102
	IR     IR
	PIR    PIR
}

func (s *Settings) ResetDefault() {
	s.Alarms.ResetDefault()
	s.APA102 = APA102{
		SPIspeed:       10000000,
		NumberLights:   150,
		StartupPattern: "\"#000001\"",
	}
	s.PIR = PIR{
		Pin:     "GPIO4",
		Pattern: "\"#ffffff\"",
	}
	s.IR = IR{
		Mapping: map[buses.Key]Pattern{
			buses.KeyNumeric0: "\"#000000\"",
			buses.Key100Plus:  "\"#ffffff\"",
		},
	}
}

func (s *Settings) Validate() error {
	if err := s.Alarms.Validate(); err != nil {
		return err
	}
	if err := s.APA102.StartupPattern.Validate(); err != nil {
		return errors.Wrap(err, "can't load startup pattern")
	}
	for k, v := range s.IR.Mapping {
		if err := v.Validate(); err != nil {
			return errors.Wrap(err, fmt.Sprintf("can't load pattern for key %s", k))
		}
	}
	if err := s.PIR.Pattern.Validate(); err != nil {
		return errors.Wrap(err, "can't load pattern for PIR")
	}
	return nil
}

// Config contains all the configuration for this specific host.
type Config struct {
	Settings Settings
	Patterns Patterns
}

func (c *Config) ResetDefault() {
	c.Settings.ResetDefault()
	c.Patterns.ResetDefault()
}

func (c *Config) Validate() error {
	if err := c.Settings.Validate(); err != nil {
		return err
	}
	return c.Patterns.Validate()
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
	return c.Validate()
}

func (c *Config) Save(n string) error {
	if err := c.Validate(); err != nil {
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
	if err := c.Settings.Alarms.Reset(p); err != nil {
		return nil
	}
	return p.SetPattern(string(c.Settings.APA102.StartupPattern))
}

func (c *ConfigMgr) Close() error {
	if len(c.path) != 0 {
		return c.Config.Save(c.path)
	}
	return nil
}

//

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
