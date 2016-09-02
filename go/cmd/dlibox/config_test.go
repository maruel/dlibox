// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"testing"

	"github.com/maruel/ut"
)

func TestInject(t *testing.T) {
	var config Config
	config.Patterns = []string{"first", "second", "third"}
	prev := make([]string, len(config.Patterns))
	copy(prev, config.Patterns)
	config.Inject("new")
	ut.AssertEqual(t, []string{"new", "first", "second", "third"}, config.Patterns)

	config.Inject("second")
	ut.AssertEqual(t, []string{"second", "new", "first", "third"}, config.Patterns)
}
