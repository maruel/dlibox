// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Package internal implements common functionality to auto-detect the host.
package internal

import (
	"io/ioutil"
	"strconv"
	"strings"
	"sync"
	"unicode"
)

//

var (
	lock      sync.Mutex
	cpuInfo   map[string]string
	osRelease map[string]string
)

func splitSemiColon(content string) map[string]string {
	// Strictly speaking this format isn't ok, there can be multiple group.
	out := map[string]string{}
	for _, line := range strings.Split(content, "\n") {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		// This format may have space around the ':'.
		key := strings.TrimRightFunc(parts[0], unicode.IsSpace)
		if len(key) == 0 || key[0] == '#' {
			continue
		}
		// Ignore duplicate keys.
		// TODO(maruel): Keep them all.
		if _, ok := out[key]; !ok {
			// Trim on both side, trailing space was observed on "Features" value.
			out[key] = strings.TrimFunc(parts[1], unicode.IsSpace)
		}
	}
	return out
}

func splitStrict(content string) map[string]string {
	out := map[string]string{}
	for _, line := range strings.Split(content, "\n") {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := parts[0]
		if len(key) == 0 || key[0] == '#' {
			continue
		}
		// Overwrite previous key.
		value := parts[1]
		if len(value) > 2 && value[0] == '"' && value[len(value)-1] == '"' {
			// Not exactly 100% right but #closeenough. See for more details
			// https://www.freedesktop.org/software/systemd/man/os-release.html
			var err error
			value, err = strconv.Unquote(value)
			if err != nil {
				continue
			}
		}
		out[key] = value
	}
	return out
}

func makeCPUInfo() map[string]string {
	lock.Lock()
	defer lock.Unlock()
	if cpuInfo == nil {
		if bytes, err := ioutil.ReadFile("/proc/cpuinfo"); err == nil {
			cpuInfo = splitSemiColon(string(bytes))
		} else {
			cpuInfo = map[string]string{}
		}
	}
	return cpuInfo
}

func makeOSRelease() map[string]string {
	lock.Lock()
	defer lock.Unlock()
	if osRelease == nil {
		if bytes, err := ioutil.ReadFile("/etc/os-release"); err == nil {
			osRelease = splitStrict(string(bytes))
		} else {
			osRelease = map[string]string{}
		}
	}
	return osRelease
}
