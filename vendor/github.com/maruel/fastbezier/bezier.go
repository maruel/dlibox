// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package fastbezier

import (
	"bytes"
	"fmt"
	"io"

	"github.com/maruel/fastbezier/internal"
)

// LUT is a fast cubic bezier curve evaluator over uint16 that uses a lookup
// table.
//
// Values are constrained in the range [0, 65535] for both x and y. It forces
// points (0, 0) and (65535, 65535).
type LUT []uint16

// Make returns a LUT object.
//
// Memory allocation is 2*(steps+1) bytes.
//
// It is only useful when the table is going to be stored as a precalculated
// table. Otherwise it is preferable to use `MakeFast`.
func Make(x0, y0, x1, y1 float32, steps uint16) LUT {
	if steps < 3 {
		// Make invalid `steps` value silently work instead of crashing or inducing
		// unnecessary error handling.
		steps = 32
	}

	// TODO(maruel): Overshoot the curve inversion points to reduce the worst
	// case error.

	stepsm1 := 1. / float32(steps-1)
	l := make(LUT, steps, steps+1)
	for i := range l {
		l[i] = internal.FloatToUint16(internal.CubicBezier(x0, y0, x1, y1, float32(i)*stepsm1) * 65535.)
	}
	// Adds a second 65535 to speed up Eval(); otherwise x==65535 has to be
	// special cased which slows it down.
	l = append(l, 65535)
	return l
}

// MakeFast returns a LUT object that is slightly less precise but takes half of
// the time to generate than `Make`.
//
// With default steps, max error is 109 instead of 104, generally in the range
// of a delta 10 higher than with `Make`.
func MakeFast(x0, y0, x1, y1 float32, steps uint16) LUT {
	if steps < 3 {
		// Make invalid `steps` value silently work instead of crashing or inducing
		// unnecessary error handling.
		steps = 32
	}
	stepsm1 := 1. / float32(steps-1)
	l := make(LUT, steps, steps+1)

	// TODO(maruel): Overshoot the curve inversion points to reduce the worst
	// case error.

	// Use a fast version that outputs (x, y) values incrementally. Use a 2x
	// resolution to get a good enough precision, especially a curvature
	// inversion points.
	// https://www.niksula.hut.fi/~hkankaan/Homepages/bezierfast.html
	// Constants
	stepsInc := 2 * steps
	t := 1. / float32(stepsInc-1)
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

	j := uint16(1)
	fJ := stepsm1
	var fX1, fY1 float32
	for i := uint16(0); j < steps; i++ {
		fX += fdX + fdd_per_2X + fddd_per_6X
		fY += fdY + fdd_per_2Y + fddd_per_6Y
		fdX += fddX + fddd_per_2X
		fdY += fddY + fddd_per_2Y
		fddX += fdddX
		fddY += fdddY
		fdd_per_2X += fddd_per_2X
		fdd_per_2Y += fddd_per_2Y
		if fX > fJ {
			a := fY1 * (fX - fJ)
			b := fY * (fJ - fX1)
			y := (a + b) / (fX - fX1)
			l[j] = internal.FloatToUint16(y * 65535.)
			j++
			fJ += stepsm1
		}
		fX1 = fX
		fY1 = fY
	}
	// Adds a second 65535 to speed up Eval(); otherwise x==65535 has to be
	// special cased which slows it down.
	l = append(l, 65535)
	return l
}

func (l LUT) String() string {
	b := bytes.NewBufferString("LUT{")
	steps := len(l) - 2
	for i, y := range l {
		x := i * 65535 / steps
		fmt.Fprintf(b, "(%d, %d)", x, y)
		if i == steps {
			break
		}
		io.WriteString(b, ", ")
	}
	io.WriteString(b, "}")
	return b.String()
}

func (l LUT) Eval(x uint16) uint16 {
	steps := uint32(len(l) - 2)
	x32 := uint32(x)
	index := x32 * steps / 65535
	nextX := (index + 1) * 65535 / steps
	baseX := index * 65535 / steps
	a := uint32(l[index]) * (nextX - x32)
	b := uint32(l[index+1]) * (x32 - baseX)
	return uint16((a + b) / (nextX - baseX))
}
