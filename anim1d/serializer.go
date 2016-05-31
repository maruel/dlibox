// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package anim1d

import (
	"bytes"
	"encoding/json"
	"errors"
	"reflect"
)

// List all known patterns and mixers that can be instantiated.
var serializerLookup map[string]reflect.Type

var knownPatterns = []Pattern{
	// Patterns
	&Color{},
	&PingPong{},
	&Animation{},
	&Rainbow{},
	&Repeated{},
	&NightSky{},
	&Aurore{},
	&NightStars{},
	&WishingStar{},
	&Gradient{},
	// Mixers
	&Transition{},
	&Loop{},
	&Crop{},
	&Mixer{},
	&Scale{},
}

func init() {
	serializerLookup = make(map[string]reflect.Type, len(knownPatterns))
	for _, i := range knownPatterns {
		r := reflect.TypeOf(i).Elem()
		serializerLookup[r.Name()] = r
	}
}

// SPattern is a Pattern that can be serialized.
//
// It is only meant to be used in mixers.
type SPattern struct {
	Pattern
}

// jsonUnmarshal unmarshals data into a map of interface{} without mangling
// int64.
func jsonUnmarshal(b []byte) (map[string]interface{}, error) {
	var tmp map[string]interface{}
	d := json.NewDecoder(bytes.NewReader(b))
	d.UseNumber()
	err := d.Decode(&tmp)
	return tmp, err
}

func (p *SPattern) UnmarshalJSON(b []byte) error {
	tmp, err := jsonUnmarshal(b)
	if err != nil {
		return nil
	}
	if len(tmp) == 0 {
		p.Pattern = nil
		return nil
	}
	n, ok := tmp["_type"]
	if !ok {
		return errors.New("bad json data")
	}
	name, ok := n.(string)
	if !ok {
		return errors.New("bad json data")
	}
	// _type will be ignored.
	p.Pattern, err = Parse(name, string(b))
	return err
}

func (p *SPattern) MarshalJSON() ([]byte, error) {
	if p.Pattern == nil {
		return []byte("{}"), nil
	}
	b, err := json.Marshal(p.Pattern.(interface{}))
	if err != nil {
		return nil, err
	}
	tmp, err := jsonUnmarshal(b)
	if err != nil {
		return nil, err
	}
	tmp["_type"] = reflect.TypeOf(p.Pattern).Elem().Name()
	return json.Marshal(tmp)
}

// Parse returns a Pattern object out of the serialized format.
func Parse(name, data string) (Pattern, error) {
	t, ok := serializerLookup[name]
	if !ok {
		return nil, errors.New("pattern not found")
	}

	v := reflect.New(t).Interface()
	if err := json.Unmarshal([]byte(data), v); err != nil {
		return nil, err
	}
	return v.(Pattern), nil
}
