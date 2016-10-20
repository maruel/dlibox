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
	"sync"

	"github.com/maruel/dlibox/go/anim1d"
	"github.com/maruel/dlibox/go/donotuse/conn/ir"
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

// LRU is the list of recent patterns. The first is the oldest.
type LRU struct {
	Max      int
	Patterns []Pattern
}

func (l *LRU) ResetDefault() {
	*l = LRU{
		Max: 25,
		Patterns: []Pattern{
			"{\"_type\":\"Aurore\"}",
			"{\"Child\":{\"Frame\":\"Lff0000ff0000ff0000ff0000ff0000ffffffffffffffffffffffffffffff\",\"_type\":\"Repeated\"},\"MovePerHour\":21600,\"_type\":\"Rotate\"}",
			"{\"Patterns\":[{\"_type\":\"Aurore\"},{\"Seed\":0,\"Stars\":null,\"_type\":\"NightStars\"},{\"AverageDelay\":0,\"Duration\":0,\"_type\":\"WishingStar\"}],\"_type\":\"Add\"}",
			"{\"Curve\":\"easeinout\",\"DurationShowMS\":1000000,\"DurationTransitionMS\":1000000,\"Patterns\":[\"#ff0000\",\"#00ff00\",\"#0000ff\"],\"_type\":\"Loop\"}",
			"{\"Curve\":\"direct\",\"Left\":\"#000000\",\"Right\":\"#0000ff\",\"_type\":\"Gradient\"}",
			"{\"Curve\":\"direct\",\"Left\":\"#000000\",\"Right\":\"#ff0000\",\"_type\":\"Gradient\"}",
			"{\"Curve\":\"direct\",\"Left\":\"#000000\",\"Right\":\"#00ff00\",\"_type\":\"Gradient\"}",
			"{\"Curve\":\"direct\",\"Left\":\"#000000\",\"Right\":\"#ffffff\",\"_type\":\"Gradient\"}",
			"{\"Child\":\"Lff0000ff0000ee0000dd0000cc0000bb0000aa0000990000880000770000660000550000440000330000220000110000\",\"MovePerHour\":108000,\"_type\":\"PingPong\"}",
			"{\"After\":\"#000000\",\"Before\":{\"After\":\"#ffffff\",\"Before\":{\"After\":\"#ff7f00\",\"Before\":\"#000000\",\"Curve\":\"direct\",\"DurationMS\":0,\"OffsetMS\":0,\"_type\":\"Transition\"},\"Curve\":\"direct\",\"DurationMS\":0,\"OffsetMS\":0,\"_type\":\"Transition\"},\"Curve\":\"direct\",\"DurationMS\":0,\"OffsetMS\":0,\"_type\":\"Transition\"}",
			"\"#000000\"",
			"{\"Child\":\"Lffffff\",\"MovePerHour\":108000,\"_type\":\"PingPong\"}",
			"{\"Curve\":\"easeinout\",\"DurationShowMS\":1000000,\"DurationTransitionMS\":10000000,\"Patterns\":[\"#ff0000\",\"#ff7f00\",\"#ffff00\",\"#00ff00\",\"#0000ff\",\"#4b0082\",\"#8b00ff\"],\"_type\":\"Loop\"}",
			"\"Rainbow\"",
			"{\"Seed\":0,\"Stars\":null,\"_type\":\"NightStars\"}",
			"{\"Child\":\"L010001ff000000ff000000ff\",\"_type\":\"Chronometer\"}",
			morning,
		},
	}
}

func (l *LRU) Validate() error {
	for i, s := range l.Patterns {
		if err := s.Validate(); err != nil {
			return errors.Wrap(err, fmt.Sprintf("can't load recent pattern %d", i))
		}
	}
	return nil
}

// Inject moves the pattern at the top of LRU cache.
func (l *LRU) Inject(pattern string) {
	for i, old := range l.Patterns {
		if old == Pattern(pattern) {
			copy(l.Patterns[i:], l.Patterns[i+1:])
			l.Patterns = l.Patterns[:len(l.Patterns)-1]
			break
		}
	}
	if len(l.Patterns) < l.Max {
		l.Patterns = append(l.Patterns, "")
	}
	copy(l.Patterns[1:], l.Patterns)
	l.Patterns[0] = Pattern(pattern)
}

//

// APA102 contains light specific settings.
type APA102 struct {
	// Speed of the transfer.
	SPIspeed int64
	// Number of lights controlled by this device. If lower than the actual
	// number of lights, the remaining lights will flash oddly.
	NumberLights int
}

func (a *APA102) ResetDefault() {
	*a = APA102{
		SPIspeed:     10000000,
		NumberLights: 150,
	}
}

func (a *APA102) Validate() error {
	return nil
}

// Button contains settings for controlling the lights through a button.
type Button struct {
	PinNumber int
}

func (b *Button) ResetDefault() {
	*b = Button{}
}

func (b *Button) Validate() error {
	return nil
}

// Display contains small embedded display settings.
type Display struct {
}

func (d *Display) ResetDefault() {
	*d = Display{}
}

func (d *Display) Validate() error {
	return nil
}

// IR contains InfraRed remote information.
type IR struct {
	Mapping map[ir.Key]Pattern // TODO(maruel): We may actually do something more complex than just set a pattern.
}

func (i *IR) ResetDefault() {
	*i = IR{
		Mapping: map[ir.Key]Pattern{
			ir.KEY_NUMERIC_0: "\"#000000\"",
			ir.KEY_100PLUS:   "\"#ffffff\"",
		},
	}
}

func (i *IR) Validate() error {
	for k, v := range i.Mapping {
		if err := v.Validate(); err != nil {
			return errors.Wrap(err, fmt.Sprintf("can't load pattern for key %s", k))
		}
	}
	return nil
}

// Painter contains settings about patterns.
type Painter struct {
	Named   map[string]Pattern // Patterns that are 'named'.
	Startup Pattern            // Startup pattern to use. If not set, use Last.
	Last    Pattern            // Last pattern used.
}

func (p *Painter) ResetDefault() {
	*p = Painter{
		Named:   map[string]Pattern{},
		Startup: "\"#010001\"",
		Last:    "\"#010001\"",
	}
}

func (p *Painter) Validate() error {
	for k, v := range p.Named {
		if err := v.Validate(); err != nil {
			return errors.Wrap(err, fmt.Sprintf("can't load pattern %s", k))
		}
	}
	if len(p.Startup) != 0 {
		if err := p.Startup.Validate(); err != nil {
			return errors.Wrap(err, "can't load pattern for Last")
		}
	}
	if len(p.Last) != 0 {
		return p.Last.Validate()
	}
	return nil
}

// PIR contains a motion detection behavior.
type PIR struct {
	Pin     int
	Pattern Pattern // TODO(maruel): We may actually do something more complex than just set a pattern.
}

func (p *PIR) ResetDefault() {
	*p = PIR{
		Pin:     -1,
		Pattern: "\"#ffffff\"",
	}
}

func (p *PIR) Validate() error {
	if err := p.Pattern.Validate(); err != nil {
		return errors.Wrap(err, "can't load pattern for PIR")
	}
	return nil
}

// Settings is all the host settings.
type Settings struct {
	Alarms  Alarms
	APA102  APA102
	Button  Button
	Display Display
	IR      IR
	Painter Painter
	PIR     PIR
}

func (s *Settings) ResetDefault() {
	s.Alarms.ResetDefault()
	s.APA102.ResetDefault()
	s.Button.ResetDefault()
	s.Display.ResetDefault()
	s.IR.ResetDefault()
	s.Painter.ResetDefault()
	s.PIR.ResetDefault()
}

func (s *Settings) Validate() error {
	if err := s.Alarms.Validate(); err != nil {
		return err
	}
	if err := s.APA102.Validate(); err != nil {
		return err
	}
	if err := s.Button.Validate(); err != nil {
		return err
	}
	if err := s.Display.Validate(); err != nil {
		return err
	}
	if err := s.IR.Validate(); err != nil {
		return err
	}
	if err := s.Painter.Validate(); err != nil {
		return err
	}
	if err := s.PIR.Validate(); err != nil {
		return err
	}
	return nil
}

// Config contains all the configuration for this specific host.
type Config struct {
	mu       sync.Mutex
	Settings Settings
	// LRU is saved aside Settings because these are not meant to be "updated" by
	// the user, they are a side-effect.
	LRU LRU
}

func (c *Config) ResetDefault() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Settings.ResetDefault()
	c.LRU.ResetDefault()
}

func (c *Config) Validate() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.validate()
}

func (c *Config) validate() error {
	if err := c.Settings.Validate(); err != nil {
		return err
	}
	return c.LRU.Validate()
}

func (c *Config) Load(n string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
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
	return c.validate()
}

func (c *Config) Save(n string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if err := c.validate(); err != nil {
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

func (c *ConfigMgr) Init(p *painter) error {
	// TODO(maruel): Use module.Bus instead.
	if err := c.Settings.Alarms.Reset(p); err != nil {
		return nil
	}
	return nil
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
