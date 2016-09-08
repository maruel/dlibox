// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package host

import "fmt"

// Level is the level of the pin: Low or High.
type Level bool

const (
	// Low represents 0v.
	Low Level = false
	// High represents Vin, generally 3.3v or 5v.
	High Level = true
)

func (l Level) String() string {
	if l == Low {
		return "Low"
	}
	return "High"
}

// Edge specifies the processor to generate an interrupt based on edge
// detection.
type Edge uint8

const (
	EdgeNone Edge = Edge(0)
	Rising   Edge = Edge(1)
	Falling  Edge = Edge(2)
	EdgeBoth Edge = Edge(3)
)

const edgeName = "NoneRisingFallingBoth"

var edgeIndex = [...]uint8{0, 4, 10, 17, 21}

func (i Edge) String() string {
	if i >= Edge(len(edgeIndex)-1) {
		return fmt.Sprintf("Edge(%d)", i)
	}
	return edgeName[edgeIndex[i]:edgeIndex[i+1]]
}

// Pull specifies the internal pull-up or pull-down for a pin set as input.
//
// The pull resistor stays set even after the processor shuts down. It is not
// possible to 'read back' what value was specified for each pin.
type Pull uint8

const (
	Float        Pull = 0 // Let the input float
	Down         Pull = 1 // Apply pull-down; for a bcm283x, the resistor is 50KOhm~60kOhm
	Up           Pull = 2 // Apply pull-up; for a bcm283x, the resistor is 50kOhm~65kOhm
	PullNoChange Pull = 3 // Do not change the previous pull resistor setting
)

const pullName = "FloatDownUpPullNoChange"

var pullIndex = [...]uint8{0, 5, 9, 11, 23}

func (i Pull) String() string {
	if i >= Pull(len(pullIndex)-1) {
		return fmt.Sprintf("Pull(%d)", i)
	}
	return pullName[pullIndex[i]:pullIndex[i+1]]
}

// Pin is a generic GPIO pin. It supports both input and output.
type Pin interface {
	PinIn
	PinOut
}

// PinIn is an input GPIO pin.
type PinIn interface {
	// In setups a pin as an input.
	In(pull Pull, edge Edge) error
	// ReadInstant return the current pin level.
	//
	// Behavior is undefined if pin is set as Output.
	ReadInstant() Level
	// ReadEdge waits until a edge detection occured and returns the pin level
	// read.
	//
	// Behavior is undefined if pin is set as Output or as EdgeNone.
	ReadEdge() Level
}

// PinOut is an output GPIO pin.
type PinOut interface {
	// Out sets a pin as output. The caller should immediately call SetLow() or
	// SetHigh() afterward.
	Out() error
	// SetLow sets a pin already set for output as Low.
	//
	// Behavior is undefined if Out() wasn't used before.
	SetLow()
	// SetHigh sets a pin already set for output as High.
	//
	// Behavior is undefined if Out() wasn't used before.
	SetHigh()
}
