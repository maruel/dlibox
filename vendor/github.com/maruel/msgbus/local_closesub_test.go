// Copyright 2017 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// +build !race

package msgbus

import (
	"sync"
	"testing"
)

func TestSubscription_closeSub(t *testing.T) {
	s := subscription{}
	// Second check for s.channel == nil.
	// This one is not deterministic but as long as we can cover this often
	// enough, it's fine. It is abuse a race condition so this test cannot be run
	// under the race detector.
	s.channel = make(chan Message)
	s.mu.RLock()

	ready := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(1)
	defer wg.Wait()
	go func() {
		defer wg.Done()
		ready <- struct{}{}
		s.closeSub()
	}()

	<-ready
	// Cheat. This triggers the race detector (for obvious reason) but permits to
	// cover a very hard to reach line.
	s.channel = nil
	s.mu.RUnlock()
}
