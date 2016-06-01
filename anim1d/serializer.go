// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package anim1d

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"image/color"
	"image/png"
	"reflect"
	"strconv"
	"time"
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

// UnmarshalJSON decodes the string "#RRGGBBAA" to the color.
func (c *Color) UnmarshalJSON(d []byte) error {
	var s string
	if err := json.Unmarshal(d, &s); err != nil {
		return err
	}
	c2, err := StringToColor(s)
	if err == nil {
		*c = c2
	}
	return err
}

// MarshalJSON encodes the color as a string "#RRGGBBAA".
func (c *Color) MarshalJSON() ([]byte, error) {
	return json.Marshal(fmt.Sprintf("#%02x%02x%02x%02x", c.R, c.G, c.B, c.A))
}

func (p *SPattern) UnmarshalJSON(b []byte) error {
	if len(b) != 0 && b[0] == '"' {
		// Special case check for Color which is encoded as "#RRGGBBAA" instead of
		// a json dict.
		c := &Color{}
		if err := json.Unmarshal(b, c); err == nil {
			p.Pattern = c
			return err
		}
	}
	tmp, err := jsonUnmarshal(b)
	if err != nil {
		return err
	}
	if len(tmp) == 0 {
		// No error but nothing was present.
		p.Pattern = nil
		return err
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
	p.Pattern, err = Parse(name, b)
	return err
}

func (p *SPattern) MarshalJSON() ([]byte, error) {
	if p.Pattern == nil {
		return []byte("{}"), nil
	}
	b, err := json.Marshal(p.Pattern.(interface{}))
	if err != nil || (len(b) != 0 && b[0] == '"') {
		// Special case check for Color which is encoded as "#RRGGBBAA" instead of
		// a json dict.
		// Also error path.
		return b, err
	}
	tmp, err := jsonUnmarshal(b)
	if err != nil {
		return nil, err
	}
	tmp["_type"] = reflect.TypeOf(p.Pattern).Elem().Name()
	return json.Marshal(tmp)
}

// Parse returns a Pattern object out of the serialized format.
func Parse(name string, data []byte) (Pattern, error) {
	t, ok := serializerLookup[name]
	if !ok {
		return nil, errors.New("pattern not found")
	}

	v := reflect.New(t).Interface()
	if err := json.Unmarshal(data, v); err != nil {
		return nil, err
	}
	return v.(Pattern), nil
}

// Marshal is a shorthand to JSON encode a pattern.
func Marshal(p Pattern) []byte {
	b, err := json.Marshal(&SPattern{p})
	if err != nil {
		panic(err)
	}
	return b
}

// StringToColor converts a #RRGGBBAA encoded string to a Color.
func StringToColor(s string) (Color, error) {
	// Do the parsing manually instead of using a regexp so the code is more
	// portable to C on an ESP8266.
	var c Color
	if len(s) != 9 {
		return c, errors.New("invalid color string")
	}
	if s[0] != '#' {
		return c, errors.New("invalid color string")
	}
	r, err := strconv.ParseUint(s[1:3], 16, 8)
	if err != nil {
		return c, err
	}
	g, err := strconv.ParseUint(s[3:5], 16, 8)
	if err != nil {
		return c, err
	}
	b, err := strconv.ParseUint(s[5:7], 16, 8)
	if err != nil {
		return c, err
	}
	a, err := strconv.ParseUint(s[7:9], 16, 8)
	if err != nil {
		return c, err
	}
	c.R = uint8(r)
	c.G = uint8(g)
	c.B = uint8(b)
	c.A = uint8(a)
	return c, nil
}

// LoadAnimate loads an Animation from a PNG file.
//
// Returns nil if the file can't be found. If vertical is true, rotate the
// image by 90°.
func LoadAnimate(content []byte, frameDuration time.Duration, vertical bool) *Animation {
	img, err := png.Decode(bytes.NewReader(content))
	if err != nil {
		return nil
	}
	bounds := img.Bounds()
	maxY := bounds.Max.Y
	maxX := bounds.Max.X
	if vertical {
		// Invert axes.
		maxY, maxX = maxX, maxY
	}
	buf := make([]Frame, maxY)
	for y := 0; y < maxY; y++ {
		buf[y] = make(Frame, maxX)
		for x := 0; x < maxX; x++ {
			if vertical {
				buf[y][x] = Color(color.NRGBAModel.Convert(img.At(y, x)).(color.NRGBA))
			} else {
				buf[y][x] = Color(color.NRGBAModel.Convert(img.At(x, y)).(color.NRGBA))
			}
		}
	}
	return &Animation{buf, frameDuration}
}
