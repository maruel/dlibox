// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Package anim1d draws 1D (line) animations that are stateless.
//
// That is, they can be played forward, backward, on multiple devices
// synchronized with only clock being synchronized.
//
// It contains all the building blocks to create animations. All the animations
// are designed to be stateless and serializable so multiple devices can
// seamless synchronize.
package anim1d
