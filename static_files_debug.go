// Copyright 2015 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// +build debug

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var paths []string

func init() {
	gopath := strings.Split(os.Getenv("GOPATH"), string(os.PathListSeparator))[0]
	paths = []string{
		filepath.Join(gopath, "src", "github.com", "maruel", "dotstar", "web", "static"),
		filepath.Join(gopath, "src", "github.com", "maruel", "dotstar", "images"),
	}
}

func read(name string) []byte {
	for _, p := range paths {
		if content, err := ioutil.ReadFile(filepath.Join(p, name)); err == nil {
			return content
		}
	}
	panic(fmt.Errorf("failed to find %s", name))
}
