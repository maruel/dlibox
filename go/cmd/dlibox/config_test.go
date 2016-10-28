// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"testing"

	"github.com/maruel/ut"
)

func TestConfig(t *testing.T) {
	var config Config
	config.ResetDefault()
	ut.AssertEqual(t, nil, config.autoFix())
}

func TestInject(t *testing.T) {
	t.Parallel()
	var config Config
	config.LRU = LRU{Max: 3, Patterns: []Pattern{"first", "second", "third"}}
	prev := make([]Pattern, len(config.LRU.Patterns))
	copy(prev, config.LRU.Patterns)
	config.LRU.Inject("new")
	ut.AssertEqual(t, []Pattern{"new", "first", "second"}, config.LRU.Patterns)

	config.LRU.Inject("second")
	ut.AssertEqual(t, []Pattern{"second", "new", "first"}, config.LRU.Patterns)
}
