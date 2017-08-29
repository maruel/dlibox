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
	"strings"
	"time"
)

const rainbowKey = "Rainbow"
const randKey = "rand"

// patternsLookup lists all known patterns that can be instantiated.
var patternsLookup map[string]reflect.Type

var knownPatterns = []Pattern{
	// Patterns
	&Color{},
	&Frame{},
	&Rainbow{},
	&Repeated{},
	&Aurore{},
	&NightStars{},
	&Lightning{},
	&WishingStar{},
	// Mixers
	&Gradient{},
	&Split{},
	&Transition{},
	&Loop{},
	&Chronometer{},
	&Rotate{},
	&PingPong{},
	&Crop{},
	&Subset{},
	&Dim{},
	&Add{},
	&Scale{},
}

// valuesLookup lists all the known values that can be instantiated.
var valuesLookup map[string]reflect.Type

var knownValues = []Value{
	new(Const),
	new(Percent),
	&OpAdd{},
	&OpMod{},
	&OpStep{},
	&Rand{},
}

func init() {
	patternsLookup = make(map[string]reflect.Type, len(knownPatterns))
	for _, i := range knownPatterns {
		r := reflect.TypeOf(i).Elem()
		patternsLookup[r.Name()] = r
	}
	valuesLookup = make(map[string]reflect.Type, len(knownValues))
	for _, i := range knownValues {
		r := reflect.TypeOf(i).Elem()
		valuesLookup[r.Name()] = r
	}
}

// SPattern

// SPattern is a Pattern that can be serialized.
//
// It is only meant to be used in mixers.
type SPattern struct {
	Pattern
}

// Render implements Pattern.
func (s *SPattern) Render(pixels Frame, timeMS uint32) {
	if s.Pattern == nil {
		return
	}
	s.Pattern.Render(pixels, timeMS)
}

// UnmarshalJSON decodes a Pattern.
//
// It knows how to decode Color, Frame or other arbitrary Pattern.
//
// If unmarshalling fails, 's' is not touched.
func (s *SPattern) UnmarshalJSON(b []byte) error {
	// Try to decode first as a string, then as a dict. Not super efficient but
	// it works.
	if p2, err := parsePatternString(b); err == nil {
		s.Pattern = p2
		return nil
	}
	o, err := jsonUnmarshalWithType(b, patternsLookup, nil)
	if err != nil {
		return err
	}
	if o == nil {
		s.Pattern = nil
	} else {
		s.Pattern = o.(Pattern)
	}
	return nil
}

// UnmarshalJSON decodes the string "#RRGGBB" to the color.
//
// If unmarshalling fails, 'c' is not touched.
func (c *Color) UnmarshalJSON(b []byte) error {
	s, err := jsonUnmarshalString(b)
	if err != nil {
		return err
	}
	return c.FromString(s)
}

// MarshalJSON encodes the color as a string "#RRGGBB".
func (c *Color) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

// UnmarshalJSON decodes the string "LRRGGBB..." to the colors.
//
// If unmarshalling fails, 'f' is not touched.
func (f *Frame) UnmarshalJSON(b []byte) error {
	s, err := jsonUnmarshalString(b)
	if err != nil {
		return err
	}
	return f.FromString(s)
}

// MarshalJSON encodes the frame as a string "LRRGGBB...".
func (f Frame) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.String())
}

// UnmarshalJSON decodes the string "Rainbow" to the rainbow.
func (r *Rainbow) UnmarshalJSON(b []byte) error {
	s, err := jsonUnmarshalString(b)
	if err != nil {
		return err
	}
	if s != rainbowKey {
		return errors.New("invalid color string")
	}
	return err
}

// MarshalJSON encodes the rainbow as a string "Rainbow".
func (r *Rainbow) MarshalJSON() ([]byte, error) {
	return json.Marshal(rainbowKey)
}

// MarshalJSON includes the additional key "_type" to help with unmarshalling.
func (s *SPattern) MarshalJSON() ([]byte, error) {
	if s.Pattern == nil {
		return []byte("{}"), nil
	}
	return jsonMarshalWithType(s.Pattern)
}

// LoadPNG loads a PNG file and creates a Loop out of the lines.
//
// If vertical is true, rotate the image by 90°.
func LoadPNG(content []byte, frameDuration time.Duration, vertical bool) *Loop {
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
	return &Loop{
		Patterns: children,
		ShowMS:   uint32(frameDuration / time.Millisecond),
	}
}

//

// parsePatternString returns a Pattern object out of the serialized JSON
// string.
func parsePatternString(b []byte) (Pattern, error) {
	s, err := jsonUnmarshalString(b)
	if err != nil {
		return nil, err
	}
	// Could try to do one after the other? It's kind of a hack at the moment.
	if len(s) != 0 {
		switch s[0] {
		case '#':
			// "#RRGGBB"
			c := &Color{}
			err := json.Unmarshal(b, c)
			return c, err
		case 'L':
			// "LRRGGBBRRGGBB..."
			var f Frame
			err := json.Unmarshal(b, &f)
			return f, err
		case rainbowKey[0]:
			// "Rainbow"
			r := &Rainbow{}
			err := json.Unmarshal(b, r)
			return r, err
		}
	}
	return nil, errors.New("unrecognized pattern string")
}

// SValue

// SValue is the serializable version of Value.
type SValue struct {
	Value
}

// Eval implements Value.
func (s *SValue) Eval(timeMS uint32, l int) int32 {
	if s.Value == nil {
		return 0
	}
	return s.Value.Eval(timeMS, l)
}

// UnmarshalJSON decodes a Value.
//
// It knows how to decode Const or other arbitrary Value.
//
// If unmarshalling fails, 'f' is not touched.
func (s *SValue) UnmarshalJSON(b []byte) error {
	// Try to decode first as a int, then as a string, then as a dict. Not super
	// efficient but it works.
	if c, err := jsonUnmarshalInt32(b); err == nil {
		s.Value = Const(c)
		return nil
	}
	if v, err := jsonUnmarshalString(b); err == nil {
		// It could be either a Percent or a Rand.
		if v == randKey {
			s.Value = &Rand{}
			return nil
		}
		if strings.HasPrefix(v, "+") {
			var o OpAdd
			if err := o.UnmarshalJSON(b); err == nil {
				s.Value = &o
			}
			return err
		}
		if strings.HasPrefix(v, "-") {
			var o OpAdd
			if err := o.UnmarshalJSON(b); err == nil {
				o.AddMS = -o.AddMS
				s.Value = &o
			}
			return err
		}
		if strings.HasPrefix(v, "%") {
			var o OpMod
			if err := o.UnmarshalJSON(b); err == nil {
				s.Value = &o
			}
			return err
		}
		if strings.HasSuffix(v, "%") {
			var p Percent
			if err := p.UnmarshalJSON(b); err == nil {
				s.Value = &p
			}
			return err
		}
		return fmt.Errorf("unknown value %q", v)
	}
	o, err := jsonUnmarshalWithType(b, valuesLookup, nil)
	if err != nil {
		return err
	}
	s.Value = o.(Value)
	return nil
}

// MarshalJSON includes the additional key "_type" to help with unmarshalling.
func (s *SValue) MarshalJSON() ([]byte, error) {
	if s.Value == nil {
		// nil value marshals to the constant 0.
		return []byte("0"), nil
	}
	return jsonMarshalWithType(s.Value)
}

// UnmarshalJSON decodes the int to the const.
//
// If unmarshalling fails, 'c' is not touched.
func (c *Const) UnmarshalJSON(b []byte) error {
	i, err := jsonUnmarshalInt32(b)
	if err != nil {
		return err
	}
	*c = Const(i)
	return err
}

// MarshalJSON encodes the const as a int.
func (c *Const) MarshalJSON() ([]byte, error) {
	return json.Marshal(int(*c))
}

// UnmarshalJSON decodes the percent in the form of a string.
//
// If unmarshalling fails, 'p' is not touched.
func (p *Percent) UnmarshalJSON(b []byte) error {
	s, err := jsonUnmarshalString(b)
	if err != nil {
		return err
	}
	if !strings.HasSuffix(s, "%") {
		return errors.New("percent must end with %")
	}
	f, err := strconv.ParseFloat(s[:len(s)-1], 32)
	if err == nil {
		// Convert back to fixed point.
		*p = Percent(int32(f * 655.36))
	}
	return err
}

// MarshalJSON encodes the percent as a string.
func (p *Percent) MarshalJSON() ([]byte, error) {
	return json.Marshal(strconv.FormatFloat(float64(*p)/655.36, 'g', 4, 32) + "%")
}

// UnmarshalJSON decodes the add in the form of a string.
//
// If unmarshalling fails, 'o' is not touched.
func (o *OpAdd) UnmarshalJSON(b []byte) error {
	s, err := jsonUnmarshalString(b)
	if err != nil {
		return err
	}
	i := int64(0)
	if strings.HasPrefix(s, "+") {
		i, err = strconv.ParseInt(s[1:], 10, 32)
	} else if strings.HasPrefix(s, "-") {
		i, err = strconv.ParseInt(s, 10, 32)
	} else {
		return errors.New("add: must start with + or -")
	}
	if err == nil {
		o.AddMS = int32(i)
	}
	if i < 0 {
		return errors.New("add: value must be positive")
	}
	return err
}

// MarshalJSON encodes the add as a string.
func (o *OpAdd) MarshalJSON() ([]byte, error) {
	if o.AddMS >= 0 {
		return json.Marshal("+" + strconv.FormatInt(int64(o.AddMS), 10))
	}
	return json.Marshal(strconv.FormatInt(int64(o.AddMS), 10))
}

// UnmarshalJSON decodes the mod in the form of a string.
//
// If unmarshalling fails, 'o' is not touched.
func (o *OpMod) UnmarshalJSON(b []byte) error {
	s, err := jsonUnmarshalString(b)
	if err != nil {
		return err
	}
	if !strings.HasPrefix(s, "%") {
		return errors.New("mod: must start with %")
	}
	i, err := strconv.ParseInt(s[1:], 10, 32)
	if err == nil {
		o.TickMS = int32(i)
	}
	if i < 0 {
		return errors.New("mod: value must be positive")
	}
	return err
}

// MarshalJSON encodes the mod as a string.
func (o *OpMod) MarshalJSON() ([]byte, error) {
	return json.Marshal("%" + strconv.FormatInt(int64(o.TickMS), 10))
}

// UnmarshalJSON decodes the string to the rand.
//
// If unmarshalling fails, 'r' is not touched.
func (r *Rand) UnmarshalJSON(b []byte) error {
	s, err := jsonUnmarshalString(b)
	if err == nil {
		// Shortcut.
		if s != randKey {
			return errors.New("invalid format")
		}
		r.TickMS = 0
		return nil
	}
	// SValue.UnmarshalJSON would handle it but implement it here so calling
	// UnmarshalJSON on a concrete instance still work. The issue is that we do
	// not want to recursively call ourselves so create a temporary type.
	type tmpRand Rand
	var r2 tmpRand
	if err := json.Unmarshal(b, &r2); err != nil {
		return err
	}
	*r = Rand(r2)
	return nil
}

// MarshalJSON encodes the rand as a string.
func (r *Rand) MarshalJSON() ([]byte, error) {
	if r.TickMS == 0 {
		// Shortcut.
		return json.Marshal(randKey)
	}
	type tmpRand Rand
	r2 := tmpRand(*r)
	return jsonMarshalWithTypeName(r2, "Rand")
}

// UnmarshalJSON is because MovePerHour is a superset of SValue.
func (m *MovePerHour) UnmarshalJSON(b []byte) error {
	var s SValue
	if err := s.UnmarshalJSON(b); err != nil {
		return err
	}
	*m = MovePerHour(s)
	return nil
}

// MarshalJSON is because MovePerHour is a superset of SValue.
func (m *MovePerHour) MarshalJSON() ([]byte, error) {
	s := SValue{m.Value}
	return s.MarshalJSON()
}

// General

// jsonUnmarshalDict unmarshals data into a map of interface{} without mangling
// int64.
func jsonUnmarshalDict(b []byte) (map[string]interface{}, error) {
	var tmp map[string]interface{}
	d := json.NewDecoder(bytes.NewReader(b))
	d.UseNumber()
	err := d.Decode(&tmp)
	return tmp, err
}

func jsonUnmarshalInt32(b []byte) (int32, error) {
	var i int32
	err := json.Unmarshal(b, &i)
	return i, err
}

func jsonUnmarshalString(b []byte) (string, error) {
	var s string
	err := json.Unmarshal(b, &s)
	return s, err
}

func jsonMarshalWithType(v interface{}) ([]byte, error) {
	t := reflect.TypeOf(v)
	switch t.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Ptr, reflect.Slice:
		return jsonMarshalWithTypeName(v, t.Elem().Name())
	default:
		return jsonMarshalWithTypeName(v, t.Name())
	}
}

func jsonMarshalWithTypeName(v interface{}, name string) ([]byte, error) {
	b, err := json.Marshal(v)
	if err != nil || (len(b) != 0 && b[0] != '{') {
		// Special case check for custom marshallers that do not encode as a dict.
		return b, err
	}
	// Inject "_type".
	tmp, err := jsonUnmarshalDict(b)
	if err != nil {
		return nil, err
	}
	tmp["_type"] = name
	return json.Marshal(tmp)
}

func jsonUnmarshalWithType(b []byte, lookup map[string]reflect.Type, null interface{}) (interface{}, error) {
	tmp, err := jsonUnmarshalDict(b)
	if err != nil {
		return nil, err
	}
	if len(tmp) == 0 {
		// No error but nothing was present. Treat "{}" as equivalent encoding for
		// null.
		return null, nil
	}
	n, ok := tmp["_type"]
	if !ok {
		return nil, errors.New("missing value type")
	}
	name, ok := n.(string)
	if !ok {
		return nil, errors.New("invalid value type")
	}
	// "_type" will be ignored, no need to reencode the dict to json.
	return parseDictToType(name, b, lookup)
}

// parseDictToType decodes an object out of the serialized JSON dict.
func parseDictToType(name string, b []byte, lookup map[string]reflect.Type) (interface{}, error) {
	t, ok := lookup[name]
	if !ok {
		return nil, fmt.Errorf("type %#v not found", name)
	}
	v := reflect.New(t).Interface()
	if err := json.Unmarshal(b, v); err != nil {
		return nil, err
	}
	return v, nil
}
