// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Packages the static files in a .go file.
//go:generate go run package.go -out static_files_gen.go ../../../web

// dlibox drives the dlibox LED strip on a Raspberry Pi. It runs a web server
// for remote control.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	_ "net/http/pprof"
	"net/url"
	"os"
	"os/signal"
	"runtime/pprof"
	"strings"
	"syscall"

	"github.com/maruel/interrupt"
)

func mainImpl() error {
	interrupt.HandleCtrlC()
	defer interrupt.Set()
	chanSignal := make(chan os.Signal)
	go func() {
		<-chanSignal
		interrupt.Set()
	}()
	signal.Notify(chanSignal, syscall.SIGTERM)
	log.SetFlags(0)

	cpuprofile := flag.String("cpuprofile", "", "dump CPU profile in file")
	port := flag.Int("port", 80, "HTTP port to listen on")
	mqttHost := flag.String("host", "tcp://dlibox:1833", "MQTT host")
	mqttUser := flag.String("user", "dlibox", "MQTT username")
	mqttPass := flag.String("pass", "dlibox", "MQTT password")
	verbose := flag.Bool("verbose", false, "enable log output")
	flag.Parse()
	if flag.NArg() != 0 {
		return fmt.Errorf("unexpected argument: %s", flag.Args())
	}

	if !*verbose {
		log.SetOutput(ioutil.Discard)
	}

	if *cpuprofile != "" {
		// Run with cpuprofile, then use 'go tool pprof' to analyze it. See
		// http://blog.golang.org/profiling-go-programs for more details.
		f, err := os.Create(*cpuprofile)
		if err != nil {
			return err
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	server := ""
	if len(*mqttHost) != 0 {
		u, err := url.ParseRequestURI(*mqttHost)
		if err != nil {
			return err
		}
		server = strings.SplitN(u.Host, ":", 2)[0]
	}

	// TODO(maruel): Standard way to figure out it's the same host?
	if server == hostname || server == "localhost" || server == "127.0.0.1" {
		return mainController(hostname, *mqttHost, *mqttUser, *mqttPass, *port)
	}
	return mainDevice(hostname, server, *mqttHost, *mqttUser, *mqttPass, *port)
}

func main() {
	if err := mainImpl(); err != nil {
		fmt.Fprintf(os.Stderr, "\ndlibox: %s.\n", err)
		os.Exit(1)
	}
}
