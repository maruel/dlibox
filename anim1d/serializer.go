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
	&Frame{},
	&Rainbow{},
	&Repeated{},
	&NightSky{},
	&Aurore{},
	&NightStars{},
	&WishingStar{},
	&Gradient{},
	// Mixers
	&Transition{},
	&Cycle{},
	&Loop{},
	&Rotate{},
	&PingPong{},
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

// UnmarshalJSON decodes the string "#RRGGBB" to the color.
//
// If unmarshalling fails, 'c' is not touched.
func (c *Color) UnmarshalJSON(d []byte) error {
	var s string
	if err := json.Unmarshal(d, &s); err != nil {
		return err
	}
	if len(s) == 0 || s[0] != '#' {
		return errors.New("invalid color string")
	}
	c2, err := stringToColor(s[1:])
	if err == nil {
		*c = c2
	}
	return err
}

// MarshalJSON encodes the color as a string "#RRGGBB".
func (c *Color) MarshalJSON() ([]byte, error) {
	return json.Marshal(fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B))
}

// UnmarshalJSON decodes the string "LRRGGBB..." to the colors.
//
// If unmarshalling fails, 'f' is not touched.
func (f *Frame) UnmarshalJSON(d []byte) error {
	var s string
	if err := json.Unmarshal(d, &s); err != nil {
		return err
	}
	if len(s) == 0 || (len(s)-1)%6 != 0 || s[0] != 'L' {
		return errors.New("invalid frame string")
	}
	l := (len(s) - 1) / 6
	f2 := make(Frame, l)
	for i := 0; i < l; i++ {
		var err error
		if f2[i], err = stringToColor(s[1+i*6 : 1+(i+1)*6]); err != nil {
			return err
		}
	}
	*f = f2
	return nil
}

// MarshalJSON encodes the frame as a string "LRRGGBB...".
func (f Frame) MarshalJSON() ([]byte, error) {
	out := bytes.Buffer{}
	out.Grow(1 + 6*len(f))
	out.WriteByte('L')
	for _, c := range f {
		fmt.Fprintf(&out, "%02x%02x%02x", c.R, c.G, c.B)
	}
	return json.Marshal(out.String())
}

// UnmarshalJSON decodes a Pattern.
//
// It knows how to decode Color, Frame or other arbitrary Pattern.
//
// If unmarshalling fails, 'f' is not touched.
func (p *SPattern) UnmarshalJSON(b []byte) error {
	if len(b) > 2 && b[0] == '"' {
		// Special case check for Color which is encoded as "#RRGGBB" and Frame
		// which is encoded as "|LRRGGBB..." instead of a json dict.
		if b[1] == '#' {
			c := &Color{}
			if err := json.Unmarshal(b, c); err == nil {
				p.Pattern = c
				return err
			}
		} else if b[1] == 'L' {
			var f Frame
			if err := json.Unmarshal(b, &f); err == nil {
				p.Pattern = f
				return err
			}
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
	p2, err := Parse(name, b)
	if err == nil {
		p.Pattern = p2
	}
	return err
}

func (p *SPattern) MarshalJSON() ([]byte, error) {
	if p.Pattern == nil {
		return []byte("{}"), nil
	}
	b, err := json.Marshal(p.Pattern)
	if err != nil || (len(b) != 0 && b[0] == '"') {
		// Special case check for Color which is encoded as "#RRGGBB" instead of a
		// json dict.
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

// stringToColor converts a "RRGGBB" encoded string to a Color.
func stringToColor(s string) (Color, error) {
	// Do the parsing manually instead of using a regexp so the code is more
	// portable to C on an ESP8266.
	var c Color
	if len(s) != 6 {
		return c, errors.New("invalid color string")
	}
	r, err := strconv.ParseUint(s[0:2], 16, 8)
	if err != nil {
		return c, err
	}
	g, err := strconv.ParseUint(s[2:4], 16, 8)
	if err != nil {
		return c, err
	}
	b, err := strconv.ParseUint(s[4:6], 16, 8)
	if err != nil {
		return c, err
	}
	c.R = uint8(r)
	c.G = uint8(g)
	c.B = uint8(b)
	return c, nil
}

// LoadPNG loads a PNG file and creates a Cycle out of the lines.
//
// If vertical is true, rotate the image by 90°.
func LoadPNG(content []byte, frameDuration time.Duration, vertical bool) *Cycle {
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
	}
	for y := 0; y < maxY; y++ {
		for x := 0; x < maxX; x++ {
			c1 := color.NRGBAModel.Convert(img.At(x, y)).(color.NRGBA)
			c := Color{c1.R, c1.G, c1.B}
			if vertical {
				buf[x][y] = c
			} else {
				buf[y][x] = c
			}
		}
	}
	children := make([]SPattern, maxY)
	for i, p := range buf {
		children[i].Pattern = p
	}
	return &Cycle{children, frameDuration}
}
