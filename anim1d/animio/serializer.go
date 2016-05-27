// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package animio

import (
	"encoding/json"
	"errors"
	"reflect"

	"github.com/maruel/dotstar/anim1d"
)

// List all known patterns and mixers that can be instantiated.
var serializerLookup map[string]reflect.Type
var knownPatterns = []anim1d.Pattern{
	// Patterns
	&anim1d.StaticColor{},
	&anim1d.PingPong{},
	&anim1d.Animation{},
	&anim1d.Rainbow{},
	&anim1d.Repeated{},
	&anim1d.NightSky{},
	&anim1d.Aurore{},
	&anim1d.NightStars{},
	&anim1d.WishingStar{},
	&anim1d.Gradient{},
	// Mixers
	&anim1d.Transition{},
	&anim1d.Loop{},
	&anim1d.Crop{},
	&anim1d.Mixer{},
	&anim1d.Scale{},
}

func init() {
	serializerLookup = make(map[string]reflect.Type, len(knownPatterns))
	for _, i := range knownPatterns {
		r := reflect.TypeOf(i).Elem()
		serializerLookup[r.Name()] = r
	}
}

// TODO(maruel): Use the lookup when unserializing.

/*
func (p *Pattern) UnmarshalJSON(b []byte) error {
	tmp := map[string]string{}
	json.Unmarshal(b, tmp)
	name, ok := tmp["type"]
	if !ok {
		return errors.New("bad json data")
	}
	var err error
	*p, err = ParsePattern(name, string(b))
	return err
}
*/

// Parse returns a Pattern object out of the serialized format.
func Parse(name, data string) (anim1d.Pattern, error) {
	t, ok := serializerLookup[name]
	if !ok {
		return nil, errors.New("pattern not found")
	}

	v := reflect.New(t).Interface()
	if err := json.Unmarshal([]byte(data), v); err != nil {
		return nil, err
	}
	return v.(anim1d.Pattern), nil
}
