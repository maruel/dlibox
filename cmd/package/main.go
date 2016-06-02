// Copyright 2015 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// package is an internal tool to build the statically embedded files in
// dotstar.
package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"text/template"
)

var contents map[string]string

var tmpl = template.Must(template.New("tmpl").Parse(`// Automatically generated file. Do not edit!
// Generated with "go run package/main.go"

// +build !debug

package main

func mustRead(name string) []byte {
	return []byte(staticFiles[name])
}

func read(name string) []byte {
	if content, ok := staticFiles[name]; ok {
		return []byte(content)
	}
	return nil
}

var staticFiles = map[string]string{
{{range $key, $value := .}}	{{$key}}: {{$value}},
{{end}}}
`))

type context struct {
	basePath string
}

func (c *context) walk(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	name := info.Name()
	if name[0] == '.' {
		return nil
	}
	if name[0] >= 'A' && name[0] <= 'Z' {
		// Ignore uppercase filename.
		return nil
	}
	if info.IsDir() {
		return nil
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	p := path[len(c.basePath)+1:]
	contents[strconv.Quote(p)] = strconv.Quote(string(data))
	return nil
}

func mainImpl() error {
	outputFile := flag.String("out", "", "outputfile")
	flag.Parse()
	if flag.NArg() == 0 {
		return errors.New("Usage: package -out [output file] [input dir ...]")
	}
	contents = map[string]string{}
	for _, inputDir := range flag.Args() {
		inputDir, err := filepath.Abs(inputDir)
		if err != nil {
			return err
		}
		c := &context{inputDir}
		if err := filepath.Walk(inputDir, c.walk); err != nil {
			return err
		}
	}
	f, err := os.Create(*outputFile)
	if err != nil {
		return err
	}
	defer f.Close()
	return tmpl.Execute(f, contents)
}

func main() {
	if err := mainImpl(); err != nil {
		fmt.Fprintf(os.Stderr, "\npackage: %s.\n", err)
		os.Exit(1)
	}
}
