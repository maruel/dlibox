// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Package drivers is essentially a registry of drivers.
//
// Every device driver should register itself in their package init() function
// by calling drivers.Register().
//
// The user call drivers.Init() on startup to initialize all the registered
// drivers in the correct order all at once.
package drivers
