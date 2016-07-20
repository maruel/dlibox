// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"testing"

	"github.com/maruel/ut"
)

func TestConfig(t *testing.T) {
	c := Config{}
	c.ResetDefault()
	ut.AssertEqual(t, nil, c.verify())
}
