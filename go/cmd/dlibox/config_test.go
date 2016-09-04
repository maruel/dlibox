// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"testing"

	"github.com/maruel/ut"
)

func Test(t *testing.T) {
	var config Config
	config.ResetDefault()
	ut.AssertEqual(t, nil, config.Validate())
}

func TestInject(t *testing.T) {
	t.Parallel()
	var config Config
	config.Patterns = Patterns{"first", "second", "third"}
	prev := make(Patterns, len(config.Patterns))
	copy(prev, config.Patterns)
	config.Patterns.Inject("new")
	ut.AssertEqual(t, Patterns{"new", "first", "second", "third"}, config.Patterns)

	config.Patterns.Inject("second")
	ut.AssertEqual(t, Patterns{"second", "new", "first", "third"}, config.Patterns)
}
