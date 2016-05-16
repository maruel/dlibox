// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package anim1d

import (
	"encoding/json"
	"errors"
	"reflect"
)

var serializerPatternsLookup map[string]reflect.Type
var serializerMixersLookup map[string]reflect.Type

func init() {
	serializerPatternsLookup = map[string]reflect.Type{
		"StaticColor": reflect.TypeOf(StaticColor{}),
		"PingPong":    reflect.TypeOf(PingPong{}),
		"Animation":   reflect.TypeOf(Animation{}),
		"Rainbow":     reflect.TypeOf(Rainbow{}),
		"Repeated":    reflect.TypeOf(Repeated{}),
		"NightSky":    reflect.TypeOf(NightSky{}),
		"Aurore":      reflect.TypeOf(Aurore{}),
		"NightStars":  reflect.TypeOf(NightStars{}),
		"WishingStar": reflect.TypeOf(WishingStar{}),
		"Gradient":    reflect.TypeOf(Gradient{}),
	}

	// It is a separate because it'll be harder to implement in C++.
	serializerMixersLookup = map[string]reflect.Type{
		"Transition": reflect.TypeOf(Transition{}),
		"Loop":       reflect.TypeOf(Loop{}),
		"Crop":       reflect.TypeOf(Crop{}),
		"Mixer":      reflect.TypeOf(Mixer{}),
		"Scale":      reflect.TypeOf(Scale{}),
	}
}

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

// ParsePattern returns a Pattern object out of the serialized format.
func ParsePattern(name, data string) (Pattern, error) {
	t, ok := serializerPatternsLookup[name]
	if !ok {
		return nil, errors.New("pattern not found")
	}

	v := reflect.New(t).Interface()
	if err := json.Unmarshal([]byte(data), v); err != nil {
		return nil, err
	}
	return v.(Pattern), nil
}

// ParseMixer returns a Pattern object out of the serialized format. It accepts
// mixers as input.
func ParseMixer(name, data string) (Pattern, error) {
	t, ok := serializerMixersLookup[name]
	if !ok {
		return nil, errors.New("mixer not found")
	}

	// TODO(maruel): Implement json.Unmarshaler?
	v := reflect.New(t).Interface()
	if err := json.Unmarshal([]byte(data), v); err != nil {
		return nil, err
	}
	return v.(Pattern), nil
}
