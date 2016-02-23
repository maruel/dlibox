// Copyright 2015 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

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

func read(name string) []byte {
	return []byte(staticFiles[name])
}

var staticFiles = map[string]string{
{{range $key, $value := .}}	{{$key}}: {{$value}},
{{end}}}
`))

func walk(path string, info os.FileInfo, err error) error {
	if info.Name()[0] == '.' {
		return nil
	}
	if info.IsDir() {
		return nil
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	contents[strconv.Quote(info.Name())] = strconv.Quote(string(data))
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
		if err := filepath.Walk(inputDir, walk); err != nil {
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
