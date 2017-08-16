// Copyright 2015 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// To use this file without having to constantly run go generate, use:
//    go install -tags debug

// +build debug

package controller

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const cacheControl30d = "Cache-Control:no-cache, no-store"
const cacheControl5m = "Cache-Control:no-cache, no-store"

func mustRead(name string) []byte {
	if content := read(name); content != nil {
		return content
	}
	panic(fmt.Errorf("failed to find %s", name))
}

func read(name string) []byte {
	gopath := strings.Split(os.Getenv("GOPATH"), string(os.PathListSeparator))[0]
	staticPaths := []string{
		filepath.Join(gopath, "src", "github.com", "maruel", "dlibox", "web"),
		filepath.Join(gopath, "src", "github.com", "maruel", "dlibox", "cmd", "dlibox", "images"),
	}
	for _, p := range staticPaths {
		if content, err := ioutil.ReadFile(filepath.Join(p, name)); err == nil {
			return content
		}
	}
	return nil
}
