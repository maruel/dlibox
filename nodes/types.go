// Copyright 2017 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package nodes

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"periph.io/x/periph/conn/physic"
)

// NodeCfg is the descriptor of a node.
type NodeCfg interface {
	Validator
	toProperties() map[ID]Property
}

type I2CRef struct {
	ID string
}

type SPIRef struct {
	ID string
	Hz physic.Frequency
}

// Anim1D is an APA102 LED strip.
type Anim1D struct {
	APA102 bool
	I2C    I2CRef
	SPI    SPIRef
	// NumberLights is the number of lights controlled by this device. If lower
	// than the actual number of lights, the remaining lights will flash oddly.
	NumberLights int
	FPS          int
}

// Validate implements Validator.
func (a *Anim1D) Validate() error {
	if !a.APA102 {
		return errors.New("anim1d: only APA102 is supported")
	}
	if len(a.I2C.ID) != 0 {
		if len(a.SPI.ID) != 0 || a.SPI.Hz != 0 {
			return errors.New("anim1d: can't use both I2C and SPI")
		}
	} else {
		if len(a.SPI.ID) == 0 {
			return errors.New("anim1d: SPI.ID is required")
		}
		if a.SPI.Hz < physic.KiloHertz {
			return errors.New("anim1d: SPI.Hz is required")
		}
	}
	if a.NumberLights <= 0 || a.NumberLights > 1000000 {
		return errors.New("anim1d: NumberLights is required")
	}
	if a.FPS <= 0 || a.FPS > 240 {
		return errors.New("anim1d: FPS is required")
	}
	return nil
}

func (a *Anim1D) toProperties() map[ID]Property {
	return map[ID]Property{
		"anim1d": {
			DataType: "string",
			Settable: true,
		},
	}
}

// Button represents a physical GPIO input pin of type Button.
type Button struct {
	Pin string
}

// Validate implements Validator.
func (b *Button) Validate() error {
	if len(b.Pin) == 0 {
		return errors.New("button: Pin is required")
	}
	return nil
}

func (b *Button) toProperties() map[ID]Property {
	return map[ID]Property{
		"button": {DataType: "boolean"},
	}
}

// Display is an ssd1306 display.
//
// TODO(maruel): make it more generic so other kind of display are supported.
type Display struct {
	SSD1306 bool
	I2C     struct {
		ID string
	}
	W, H int
}

// Validate implements Validator.
func (d *Display) Validate() error {
	if !d.SSD1306 {
		return errors.New("display: only SSD1306 is supported")
	}
	if len(d.I2C.ID) == 0 {
		return errors.New("display: I2C.ID is required")
	}
	if d.W == 0 {
		return errors.New("display: W is required")
	}
	if d.H == 0 {
		return errors.New("display: H is required")
	}
	return nil
}

func (d *Display) toProperties() map[ID]Property {
	return map[ID]Property{
		"markee": {
			DataType: "string",
			Settable: true,
		},
		"content": {
			DataType: "string",
			Settable: true,
		},
	}
}

// IR is an InfraRed Remote receiver.
//
// In practice, only lirc is supported.
type IR struct {
}

// Validate implements Validator.
func (i *IR) Validate() error {
	return nil
}

func (i *IR) toProperties() map[ID]Property {
	return map[ID]Property{
		"ir": {DataType: "string"},
	}
}

// PIR represents a GPIO physical pin that is connected to a motion detector.
type PIR struct {
	Pin string
}

// Validate implements Validator.
func (p *PIR) Validate() error {
	if len(p.Pin) == 0 {
		return errors.New("pir: Pin is required")
	}
	return nil
}

func (p *PIR) toProperties() map[ID]Property {
	return map[ID]Property{
		"pir": {DataType: "boolean"},
	}
}

// Sound is a sound output device.
type Sound struct {
	DeviceID string // Empty to use the default sound card.
}

// Validate implements Validator.
func (s *Sound) Validate() error {
	return nil
}

func (s *Sound) toProperties() map[ID]Property {
	return map[ID]Property{
		// An URL.
		"speakers": {
			DataType: "string",
			Settable: true,
		},
	}
}

//

// Type is all the known node types.
type Type string

// Validate implements Validator.
func (t Type) Validate() error {
	if _, ok := TypesMap[t]; !ok {
		return fmt.Errorf("unknown type %q", t)
	}
	return nil
}

var knownTypes = []Validator{
	&Anim1D{},
	&Button{},
	&Display{},
	&IR{},
	&PIR{},
	&Sound{},
}

// TypesMap is the list of known types and their associated name.
var TypesMap map[Type]reflect.Type

func init() {
	TypesMap = map[Type]reflect.Type{}
	for _, t := range knownTypes {
		r := reflect.TypeOf(t).Elem()
		TypesMap[Type(strings.ToLower(r.Name()))] = r
	}
}
