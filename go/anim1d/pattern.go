// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package anim1d

// Pattern is a interface to draw an animated line.
type Pattern interface {
	// Render fills the buffer with the image at this time frame.
	//
	// The image should be derived from timeMS, which is the time since this
	// pattern was started.
	//
	// Calling Render() with a nil pattern is valid. Patterns should be callable
	// without crashing with an object initialized with default values.
	//
	// timeMS will cycle after 49.7 days. The reason it's not using time.Duration
	// is that int64 calculation on ARM is very slow and abysmal on xtensa, which
	// this code is transpiled to.
	Render(pixels Frame, timeMS uint32)
}
