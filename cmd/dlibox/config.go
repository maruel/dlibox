// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"errors"
	"os"
	"os/user"
	"path/filepath"

	"github.com/maruel/dlibox-go/anim1d"
)

type APA102 struct {
	// Speed of the transfer.
	SPIspeed int64
	// Number of lights controlled by this device. If lower than the actual
	// number of lights, the remaining lights will flash oddly.
	NumberLights int
}

// Config stores the configuration for this specific host.
type Config struct {
	Alarms   Alarms
	APA102   APA102
	Patterns []string // List of recent patterns.
}

func (c *Config) ResetDefault() {
	*c = Config{
		//"{\"Duration\":600000000000,\"After\":\"#00000000\",\"Offset\":1800000000000,\"Before\":{\"Duration\":600000000000,\"After\":\"#ffffffff\",\"Offset\":600000000000,\"Before\":{\"Duration\":600000000000,\"After\":\"#ff7f00ff\",\"Offset\":0,\"Before\":\"#00000000\",\"Transition\":\"linear\",\"_type\":\"Transition\"},\"Transition\":\"linear\",\"_type\":\"Transition\"},\"Transition\":\"linear\",\"_type\":\"Transition\"}",
		Alarms: Alarms{
			{
				Enabled: true,
				Hour:    6,
				Minute:  55,
				Days:    Monday | Tuesday | Wednesday | Thursday | Friday,
				Pattern: "\"#FFFFFF\"",
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
			SPIspeed:     10000000,
			NumberLights: 150,
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
	return d.Decode(c)
}

func (c *Config) Save(n string) error {
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
	return c.Alarms.Reset(p)
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
