// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package internal

import (
	"io/ioutil"
	"log"
	"strings"
)

func init() {
	CPUInfo = map[string]string{}
	bytes, err := ioutil.ReadFile("/proc/cpuinfo")
	if err != nil {
		log.Printf("Failed to read /proc/cpuinfo: %s", err)
	}
	for _, line := range strings.Split(string(bytes), "\n") {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		// Ignore information for other processors than the #0.
		if len(CPUInfo[key]) == 0 {
			CPUInfo[key] = strings.TrimSpace(parts[1])
		}
	}
}
