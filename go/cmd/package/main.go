// Copyright 2015 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// package is an internal tool to build the statically embedded files in
// dlibox.
package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"text/template"
)

var tmpl = template.Must(template.New("tmpl").Parse(`// Automatically generated file. Do not edit!
// Generated with "go run package/main.go"

// +build !debug

package main

const cacheControl30d = "Cache-Control:public, max-age=259200" // 30d
const cacheControl5m = "Cache-Control:public, max-age=300" // 5m

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
{{range .}}	{{.Name}}: {{.Content}},
{{end}}}
`))

type file struct {
	Name    string
	Content string
}

type files []file

func (f files) Len() int           { return len(f) }
func (f files) Less(i, j int) bool { return f[i].Name < f[j].Name }
func (f files) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }

type context struct {
	basePath string
	files    files
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
	c.files = append(c.files, file{strconv.Quote(p), strconv.Quote(string(data))})
	return nil
}

func mainImpl() error {
	outputFile := flag.String("out", "", "outputfile")
	flag.Parse()
	if flag.NArg() == 0 {
		return errors.New("Usage: package -out [output file] [input dir ...]")
	}
	c := &context{}
	for _, inputDir := range flag.Args() {
		inputDir, err := filepath.Abs(inputDir)
		if err != nil {
			return err
		}
		c.basePath = inputDir
		if err := filepath.Walk(inputDir, c.walk); err != nil {
			return err
		}
	}
	f, err := os.Create(*outputFile)
	if err != nil {
		return err
	}
	defer f.Close()
	sort.Sort(c.files)
	if err := tmpl.Execute(f, c.files); err != nil {
		return err
	}
	return exec.Command("gofmt", "-w", "-s", *outputFile).Run()
}

func main() {
	if err := mainImpl(); err != nil {
		fmt.Fprintf(os.Stderr, "\npackage: %s.\n", err)
		os.Exit(1)
	}
}
