// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package controller

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"github.com/maruel/dlibox/controller/alarm"
	"github.com/maruel/dlibox/controller/rules"
	"github.com/maruel/dlibox/nodes"
	"github.com/maruel/dlibox/shared"
)

// config contains all the configuration that the user can specify.
type config struct {
	// Not stored in MQTT
	Alarms alarm.Config
	Rules  rules.Rules

	// Stored in MQTT as nodes.Nodes
	Devices map[nodes.ID]*nodes.Dev
}

// db is all the settings and values that are persisted on disk.
type db struct {
	mu     sync.Mutex
	Config config
	// AnimLRU is saved outside of Config because these are not meant to be
	// "updated" by the user, they are a side-effect.
	AnimLRU animLRU
}

func (d *db) load(n string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	f, err := os.Open(n)
	if err != nil {
		if os.IsNotExist(err) {
			// Ignore if the file is not present.
			return nil
		}
		return err
	}
	defer f.Close()
	j := json.NewDecoder(f)
	j.UseNumber()
	return j.Decode(d)
}

func (d *db) save(n string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	b, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		return err
	}
	f, err := os.Create(n)
	if err != nil {
		return err
	}
	if _, err = f.Write(append(b, '\n')); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}

type dbMgr struct {
	db
	path string
}

func (d *dbMgr) Load() error {
	d.path = filepath.Join(shared.Home(), "dlibox.json")
	return d.db.load(d.path)
}

func (d *dbMgr) Close() error {
	if len(d.path) != 0 {
		return d.db.save(d.path)
	}
	return nil
}
