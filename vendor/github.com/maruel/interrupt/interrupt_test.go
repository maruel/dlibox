// Copyright 2014 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package interrupt

import (
	"testing"
)

func TestSet(t *testing.T) {
	if IsSet() {
		t.Fatal("IsSet() should be false")
	}
	select {
	case <-Channel:
		t.Fatal("Channel should not trigger")
	default:
	}

	HandleCtrlC()

	if IsSet() {
		t.Fatal("IsSet() should be false")
	}
	select {
	case <-Channel:
		t.Fatal("Channel should not trigger")
	default:
	}

	Set()

	for i := 0; i < 2; i++ {
		if !IsSet() {
			t.Fatal("IsSet() should be true")
		}
		x, ok := <-Channel
		if !x {
			t.Fatal("Channel should send true")
		}
		if !ok {
			t.Fatal("Channel should be open")
		}
	}
}
