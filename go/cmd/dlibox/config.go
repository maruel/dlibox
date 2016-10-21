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

//

// Settings is all the host settings.
type Settings struct {
	Alarms  Alarms
	APA102  APA102
	Button  Button
	Display Display
	IR      IR
	MQTT    MQTT
	Painter Painter
	PIR     PIR
}

func (s *Settings) Lock() {
	s.Alarms.Lock()
	s.APA102.Lock()
	s.Button.Lock()
	s.Display.Lock()
	s.IR.Lock()
	s.MQTT.Lock()
	s.Painter.Lock()
	s.PIR.Lock()
}

func (s *Settings) Unlock() {
	s.PIR.Unlock()
	s.Painter.Unlock()
	s.MQTT.Unlock()
	s.IR.Unlock()
	s.Display.Unlock()
	s.Button.Unlock()
	s.APA102.Unlock()
	s.Alarms.Unlock()
}

func (s *Settings) ResetDefault() {
	s.Alarms.ResetDefault()
	s.APA102.ResetDefault()
	s.Button.ResetDefault()
	s.Display.ResetDefault()
	s.IR.ResetDefault()
	s.MQTT.ResetDefault()
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
	if err := s.MQTT.Validate(); err != nil {
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
	Settings Settings
	// LRU is saved aside Settings because these are not meant to be "updated" by
	// the user, they are a side-effect.
	LRU LRU
}

func (c *Config) Lock() {
	c.Settings.Lock()
	c.LRU.Lock()
}

func (c *Config) Unlock() {
	c.LRU.Unlock()
	c.Settings.Unlock()
}

func (c *Config) ResetDefault() {
	c.Settings.ResetDefault()
	c.LRU.ResetDefault()
}

func (c *Config) Validate() error {
	if err := c.Settings.Validate(); err != nil {
		return err
	}
	return c.LRU.Validate()
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
	c.Lock()
	err = d.Decode(c)
	c.Unlock()
	if err != nil {
		return err
	}
	return c.Validate()
}

func (c *Config) Save(n string) error {
	// There's a window between validating and marshalling where the lock is
	// temporarilly released.
	if err := c.Validate(); err != nil {
		return err
	}
	c.Lock()
	b, err := json.MarshalIndent(c, "", "  ")
	c.Unlock()
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
