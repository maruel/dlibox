// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package rejected

import (
	"bytes"
	"fmt"
	"io"

	"github.com/maruel/fastbezier/internal"
)

// Evaluator exposes a function in the uint16 domain.
type Evaluator interface {
	// Eval evaluates a function taking a uint16 as input and returning a
	// uint16.
	//
	// There is no provision for overshooting.
	//
	// Eval should be idempotent.
	Eval(x uint16) uint16
}

// Precise is the precise evaluation of a cubic bezier curve in the uint16
// domain.
type Precise struct {
	x0, y0, x1, y1 float32
}

// MakePrecise returns a precise evaluator.
func MakePrecise(x0, y0, x1, y1 float32) *Precise {
	return &Precise{x0, y0, x1, y1}
}

func (p *Precise) String() string {
	return fmt.Sprintf("Precise{%g, %g, %g, %g}", p.x0, p.y0, p.x1, p.y1)
}

func (p *Precise) Eval(x uint16) uint16 {
	return internal.CubicBezier16(p.x0, p.y0, p.x1, p.y1, x)
}

type point struct {
	x, y uint16
}

// PointsTrimmed is a fast cubic bezier curve evaluator over uint16 that uses a
// lookup table that is quick to generate, at the cost of slightly slower
// evaluation and twice the table size.
//
// It trades off one time initialization and lower CPU utilization over memory
// usage and precision. The precalculated points are distributed along the line
// itself, not on the X axis.
//
// It is constrained over the range [0, 65535] for both x and y. It forces
// points (0, 0) and (65535, 65535).
type PointsTrimmed []point

// MakePointsTrimmed returns a PointsTrimmed object that is a lookup table for
// fast cubic bezier curve evaluation for a bezier curve
// [(0, 0), (x0, y0), (x1, y1), (65535, 65535)].
//
// Memory allocation is 4*(steps-2) bytes.
func MakePointsTrimmed(x0, y0, x1, y1 float32, steps uint16) PointsTrimmed {
	if steps < 3 {
		// Make invalid `steps` value silently work instead of crashing or inducing
		// unnecessary error handling.
		steps = 32
	}

	// Use a fast version that outputs (x, y) values incrementally.
	// https://www.niksula.hut.fi/~hkankaan/Homepages/bezierfast.html
	// Constants
	t := 1. / float32(steps-1)
	t2 := t * t
	const p0X = float32(0)
	const p0Y = float32(0)
	p1X := x0
	p1Y := y0
	p2X := x1
	p2Y := y1
	const p3X = float32(1)
	const p3Y = float32(1)

	// Starting vector.
	fX := p0X
	fY := p0Y

	// First degree derivate: f '(x)*t
	fdX := 3. * (p1X - p0X) * t
	fdY := 3. * (p1Y - p0Y) * t

	// Second degree derivate: f ''(x)*t^2 / 2
	fdd_per_2X := 3. * (p0X - 2.*p1X + p2X) * t2
	fdd_per_2Y := 3. * (p0Y - 2.*p1Y + p2Y) * t2
	fddX := fdd_per_2X + fdd_per_2X
	fddY := fdd_per_2Y + fdd_per_2Y

	fddd_per_2X := 3. * (3.*(p1X-p2X) + p3X - p0X) * t2 * t
	fddd_per_2Y := 3. * (3.*(p1Y-p2Y) + p3Y - p0Y) * t2 * t
	// Third degree derivate: f '''(x)*t^3 / 6
	fddd_per_6X := fddd_per_2X * (1. / 3.)
	fddd_per_6Y := fddd_per_2Y * (1. / 3.)
	fdddX := fddd_per_2X + fddd_per_2X
	fdddY := fddd_per_2Y + fddd_per_2Y

	// Skip first and last points; (0, 0) and (65535, 65535).
	p := make(PointsTrimmed, steps-2)
	for i := range p {
		fX += fdX + fdd_per_2X + fddd_per_6X
		fY += fdY + fdd_per_2Y + fddd_per_6Y
		fdX += fddX + fddd_per_2X
		fdY += fddY + fddd_per_2Y
		fddX += fdddX
		fddY += fdddY
		fdd_per_2X += fddd_per_2X
		fdd_per_2Y += fddd_per_2Y
		p[i] = point{internal.FloatToUint16(fX * 65535.), internal.FloatToUint16(fY * 65535.)}
	}
	return p
}

func (p PointsTrimmed) String() string {
	b := bytes.NewBufferString("PointsTrimmed{(0, 0), ")
	for _, point := range p {
		fmt.Fprintf(b, "(%d, %d)", point.x, point.y)
		io.WriteString(b, ", ")
	}
	io.WriteString(b, "(65535, 65535)}")
	return b.String()
}

func (p PointsTrimmed) Eval(x uint16) uint16 {
	// For very short table, it's faster to do a linear search. The trade off is
	// CPU architecture dependent.
	for i, p1 := range p {
		if p1.x > x {
			if i == 0 {
				return uint16(uint32(p1.y) * uint32(x) / uint32(p1.x))
			}
			p0 := p[i-1]
			a := uint32(p0.y) * uint32(p1.x-x)
			b := uint32(p1.y) * uint32(x-p0.x)
			return uint16((a + b) / uint32(p1.x-p0.x))
		}
	}
	p0 := p[len(p)-1]
	a := uint32(p0.y) * uint32(65535-x)
	b := uint32(65535) * uint32(x-p0.x)
	return uint16((a + b) / uint32(65535-p0.x))
}

// PointsFull is the version of Points with (0, 0) and (65535, 65535) in the
// table for faster evaluation.
type PointsFull []point

func MakePointsFull(x0, y0, x1, y1 float32, steps uint16) PointsFull {
	if steps < 3 {
		// Make invalid `steps` value silently work instead of crashing or inducing
		// unnecessary error handling.
		steps = 32
	}

	// Use a fast version that outputs (x, y) values incrementally.
	// https://www.niksula.hut.fi/~hkankaan/Homepages/bezierfast.html
	// Constants
	t := 1. / float32(steps-1)
	t2 := t * t
	const p0X = float32(0)
	const p0Y = float32(0)
	p1X := x0
	p1Y := y0
	p2X := x1
	p2Y := y1
	const p3X = float32(1)
	const p3Y = float32(1)

	// Starting vector.
	fX := p0X
	fY := p0Y

	// First degree derivate: f '(x)*t
	fdX := 3. * (p1X - p0X) * t
	fdY := 3. * (p1Y - p0Y) * t

	// Second degree derivate: f ''(x)*t^2 / 2
	fdd_per_2X := 3. * (p0X - 2.*p1X + p2X) * t2
	fdd_per_2Y := 3. * (p0Y - 2.*p1Y + p2Y) * t2
	fddX := fdd_per_2X + fdd_per_2X
	fddY := fdd_per_2Y + fdd_per_2Y

	fddd_per_2X := 3. * (3.*(p1X-p2X) + p3X - p0X) * t2 * t
	fddd_per_2Y := 3. * (3.*(p1Y-p2Y) + p3Y - p0Y) * t2 * t
	// Third degree derivate: f '''(x)*t^3 / 6
	fddd_per_6X := fddd_per_2X * (1. / 3.)
	fddd_per_6Y := fddd_per_2Y * (1. / 3.)
	fdddX := fddd_per_2X + fddd_per_2X
	fdddY := fddd_per_2Y + fddd_per_2Y

	p := make(PointsFull, steps)
	for i := range p {
		p[i] = point{internal.FloatToUint16(fX * 65535.), internal.FloatToUint16(fY * 65535.)}
		if i == len(p)-1 {
			break
		}
		fX += fdX + fdd_per_2X + fddd_per_6X
		fY += fdY + fdd_per_2Y + fddd_per_6Y
		fdX += fddX + fddd_per_2X
		fdY += fddY + fddd_per_2Y
		fddX += fdddX
		fddY += fdddY
		fdd_per_2X += fddd_per_2X
		fdd_per_2Y += fddd_per_2Y
	}
	return p
}

func (p PointsFull) String() string {
	b := bytes.NewBufferString("PointsFull{")
	for i, point := range p {
		fmt.Fprintf(b, "(%d, %d)", point.x, point.y)
		if i != len(p)-1 {
			io.WriteString(b, ", ")
		}
	}
	io.WriteString(b, "}")
	return b.String()
}

func (p PointsFull) Eval(x uint16) uint16 {
	// For very short table, it's faster to do a linear search. The trade off is
	// CPU architecture dependent.
	// TODO(maruel): Implement binary search and benchmark.
	for i, p1 := range p {
		if p1.x > x {
			p0 := p[i-1]
			a := uint32(p0.y) * uint32(p1.x-x)
			b := uint32(p1.y) * uint32(x-p0.x)
			return uint16((a + b) / uint32(p1.x-p0.x))
		}
	}
	return 65535
}

// TableTrimmed is a fast cubic bezier curve evaluator over uint16 that uses a
// lookup table that is slower to generate compared to Points, with the benefit
// of faster evaluation and half the table size.
//
// It trades off one time initialization and CPU utilization over memory usage
// and precision.
//
// Each point is spaced uniformly across the X axis.
//
// It is constrained over the range [0, 65535] for both x and y. It forces
// points (0, 0) and (65535, 65535).
type TableTrimmed []uint16

// MakeTableTrimmed returns a TableTrimmed object that contains a lookup table
// for fast cubic bezier curve evaluation for the curve
// [(0, 0), (x0, y0), (x1, y1), (65535, 65535)].
//
// This function is slower than MakePoints but Eval() calls are faster.
// TableTrimmed uses half of the memory of Points. It is meant to be used when
// a significant number of calls to Eval() will be done.
//
// Memory allocation is 2*(steps-2) bytes.
func MakeTableTrimmed(x0, y0, x1, y1 float32, steps uint16) TableTrimmed {
	if steps < 3 {
		// Make invalid `steps` value silently work instead of crashing or inducing
		// unnecessary error handling.
		steps = 32
	}
	stepsm1 := 1. / float32(steps-1)
	// 0 and 65535 are omitted.
	t := make(TableTrimmed, steps-2)
	for i := range t {
		t[i] = internal.FloatToUint16(internal.CubicBezier(x0, y0, x1, y1, float32(i+1)*stepsm1) * 65535.)
	}
	return t
}

func (t TableTrimmed) String() string {
	b := bytes.NewBufferString("TableTrimmed{(0, 0), ")
	steps := len(t) + 1
	for i, y := range t {
		x := (i + 1) * 65535 / steps
		fmt.Fprintf(b, "(%d, %d)", x, y)
		io.WriteString(b, ", ")
	}
	io.WriteString(b, "(65535, 65535)}")
	return b.String()
}

func (t TableTrimmed) Eval(x uint16) uint16 {
	// Points 0 and 65535 are omitted from the table.
	steps := uint32(len(t) + 1)
	x32 := uint32(x)
	switch index := x32 * steps / 65535; index {
	case 0:
		// The first point (0, 0) is not stored.
		return uint16(uint32(t[0]) * x32 / (65535 / steps))
	case steps - 1:
		// The last point (65535, 65535) is not stored.
		baseX := index * 65535 / steps
		a := uint32(t[len(t)-1]) * (65535 - x32)
		b := uint32(65535) * (x32 - baseX)
		return uint16((a + b) / (65535 - baseX))
	case steps:
		// For x==65535.
		return 65535
	default:
		nextX := (index + 1) * 65535 / steps
		baseX := index * 65535 / steps
		a := uint32(t[index-1]) * (nextX - x32)
		b := uint32(t[index]) * (x32 - baseX)
		return uint16((a + b) / (nextX - baseX))
	}
}

// TableFull is the version of Table with (0, 0) and (65535, 65535) in the
// table for faster evaluation.
type TableFull []uint16

func MakeTableFull(x0, y0, x1, y1 float32, steps uint16) TableFull {
	if steps < 3 {
		// Make invalid `steps` value silently work instead of crashing or inducing
		// unnecessary error handling.
		steps = 32
	}
	// Adds a second 65535 to speed up Eval(); otherwise x==65535 has to be
	// special cased which slows Eval() down.
	stepsm1 := 1. / float32(steps-1)
	t := make(TableFull, steps, steps+1)
	for i := range t {
		t[i] = internal.FloatToUint16(internal.CubicBezier(x0, y0, x1, y1, float32(i)*stepsm1) * 65535.)
	}
	t = append(t, 65535)
	return t
}

func (t TableFull) String() string {
	b := bytes.NewBufferString("TableFull{")
	steps := len(t) - 2
	for i, y := range t {
		if i == len(t)-1 {
			break
		}
		x := i * 65535 / steps
		fmt.Fprintf(b, "(%d, %d)", x, y)
		if i != steps {
			io.WriteString(b, ", ")
		}
	}
	io.WriteString(b, "}")
	return b.String()
}

func (t TableFull) Eval(x uint16) uint16 {
	steps := uint32(len(t) - 2)
	x32 := uint32(x)
	index := x32 * steps / 65535
	nextX := (index + 1) * 65535 / steps
	baseX := index * 65535 / steps
	a := uint32(t[index]) * (nextX - x32)
	b := uint32(t[index+1]) * (x32 - baseX)
	return uint16((a + b) / (nextX - baseX))
}
