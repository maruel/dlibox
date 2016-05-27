// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package animio

import (
	"image/color"
	"testing"

	"github.com/maruel/ut"
)

func TestParse(t *testing.T) {
	// TODO(maruel): Test all known patterns and mixers.
	p, err := Parse("StaticColor", "{\"c\":{\"r\":1, \"g\":2, \"b\":3, \"a\":4}}")
	ut.AssertEqual(t, nil, err)
	b := [1]color.NRGBA{}
	p.NextFrame(b[:], 0)
	ut.AssertEqual(t, color.NRGBA{1, 2, 3, 4}, b[0])
}
