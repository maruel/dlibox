// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// +build !linux

package gpiomem

import "errors"

// Open is not implemented
func Open() (*Mem, error) {
	return nil, errors.New("not implemented")
}

func (m *Mem) Close() error {
	return nil
}
