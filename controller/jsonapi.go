// Copyright 2018 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/maruel/anim1d"
	"github.com/maruel/msgbus"
)

// jsonAPI contains the global state/caches for the JSON API.
type jsonAPI struct {
	hostname string
	b        msgbus.Bus
	l        io.WriterTo
	db       *db
}

func (j *jsonAPI) init(hostname string, b msgbus.Bus, d *db, l io.WriterTo) {
	j.hostname = hostname
	j.b = b
	j.l = l
	j.db = d
}

// getAPIs returns the JSON API handlers.
func (j *jsonAPI) getAPIs() []apiHandler {
	return []apiHandler{
		{"/api/dlibox/v1/pattern/list", j.apiPatternList},
		{"/api/dlibox/v1/pattern/get", j.apiPatternGet},
		{"/api/dlibox/v1/pattern/set", j.apiPatternSet},
		{"/api/dlibox/v1/publish", j.apiPublish},
		{"/api/dlibox/v1/server/state", j.apiServerState},
		{"/api/dlibox/v1/settings/get", j.apiSettingGet},
		{"/api/dlibox/v1/settings/set", j.apiSettingSet},
	}
}

// /api/dlibox/v1/pattern/list

func (j *jsonAPI) apiPatternList() ([]pattern, int) {
	j.db.AnimLRU.Lock()
	defer j.db.AnimLRU.Unlock()
	// TODO(maruel): We must make a copy or figure out a way for the lock to live
	// through the JSON encoding.
	return j.db.AnimLRU.Patterns, 200
}

// /api/dlibox/v1/pattern/get

func (j *jsonAPI) apiPatternGet(name string) (pattern, int) {
	//j.db.Config.Painter.Lock()
	//defer j.db.Config.Painter.Unlock()
	/*
		l := j.db.Config.Painter.Last
		if l == "" {
			l = j.db.Config.Painter.Startup
		}
		return l
	*/
	return pattern(""), 200
}

// /api/dlibox/v1/pattern/set

// TODO(maruel): Accept interface{}
func (j *jsonAPI) apiPatternSet(raw []byte) (interface{}, int) {
	var obj anim1d.SPattern
	if err := json.Unmarshal(raw, &obj); err != nil {
		log.Printf("web: invalid JSON pattern: %v", err)
		return map[string]string{"error": err.Error()}, 400
	}
	// Reencode in canonical format to send it back to the user.
	raw, err := obj.MarshalJSON()
	if err != nil {
		log.Printf("web: internal error: %v", err)
		return map[string]string{"error": err.Error()}, 500
	}
	if err := j.b.Publish(msgbus.Message{Topic: "painter/setuser", Payload: raw}, msgbus.ExactlyOnce); err != nil {
		log.Printf("web: failed to publish: %v", err)
		return map[string]string{"error": fmt.Sprintf("failed to publish: %v", err)}, 500
	}
	return raw, 200
}

// /api/dlibox/v1/publish

func (j *jsonAPI) apiPublish(state string) (map[string]string, int) {
	/*
		if !State(state).Valid() {
			return map[string]string{"error": "state is invalid"}, 400
		}
	*/
	if err := j.b.Publish(msgbus.Message{Topic: "//dlibox/halloween/state", Payload: []byte(state)}, msgbus.ExactlyOnce); err != nil {
		log.Printf("web: failed to publish: %v", err)
		return map[string]string{"error": fmt.Sprintf("failed to publish: %v", err)}, 500
	}
	return map[string]string{"ok": "1"}, 200
}

// /api/dlibox/v1/settings/get

func (j *jsonAPI) apiSettingGet() (interface{}, int) {
	// TODO(maruel): Lock.
	//j.db.Config.Lock()
	//defer j.db.Config.Unlock()
	return j.db.Config, 200
}

// /api/dlibox/v1/settings/set

func (j *jsonAPI) apiSettingSet(settings config) (config, int) {
	// TODO(maruel): Lock.
	j.db.Config = settings
	// Serialize it again to return the canonical form.
	return settings, 200
}

// /api/periph/v1/server/state

type serverStateOut struct {
	Hostname string
}

func (j *jsonAPI) apiServerState() (*serverStateOut, int) {
	out := &serverStateOut{
		Hostname: j.hostname,
	}
	return out, 200
}
