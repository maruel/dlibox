// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// dlibox is an home automation system.
package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	_ "net/http/pprof"
	"net/url"
	"os"
	"os/signal"
	"runtime/pprof"
	"strconv"
	"strings"
	"syscall"

	"github.com/maruel/dlibox/go/modules/controller"
	"github.com/maruel/dlibox/go/modules/device"
	"github.com/maruel/dlibox/go/modules/shared"
	"github.com/maruel/dlibox/go/msgbus"
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
	mqttHost := flag.String("mqtt", "tcp://dlibox:1833", "MQTT host in the form tcp://user:pass@host:port")
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

	bus := msgbus.New()
	server := ""
	isController := false
	if len(*mqttHost) != 0 {
		u, err := url.ParseRequestURI(*mqttHost)
		if err != nil {
			return err
		}
		parts := strings.SplitN(u.Host, ":", 2)
		if len(parts) != 2 {
			return errors.New("mqtt port is required")
		}
		if len(parts[0]) == 0 {
			return errors.New("mqtt host is required")
		}
		server = parts[0]
		if i, err := strconv.Atoi(parts[1]); i == 0 || err != nil {
			return errors.New("mqtt port is required")
		}
		if len(u.Path) != 0 {
			return errors.New("mqtt must not include a path")
		}
		if len(u.RawQuery) != 0 {
			return errors.New("mqtt must not not include a query argument")
		}

		// TODO(maruel): Standard way to figure out it's the same host? Likely by
		// resolving the IP.
		isController := server == shared.Hostname() || server == "localhost" || server == "127.0.0.1"
		root := "dlibox"
		if !isController {
			root += shared.Hostname()
		}
		will := msgbus.Message{root + "/$online", []byte("false")}
		usr := ""
		pwd := ""
		if u.User != nil {
			usr = u.User.Username()
			pwd, _ = u.User.Password()
		}
		mqttServer, err := msgbus.NewMQTT(server, shared.Hostname(), usr, pwd, will)
		if err != nil {
			// TODO(maruel): Have it continuously try to automatically reconnect.
			log.Printf("Failed to connect to server: %v", err)
		} else {
			bus = mqttServer
		}
		// Everything is under the namespace "dlibox/"
		bus = msgbus.RebasePub(msgbus.RebaseSub(bus, "dlibox"), "dlibox")
	}

	if isController {
		return controller.Main(bus, *port)
	}
	return device.Main(server, bus, *port)
}

func main() {
	if err := mainImpl(); err != nil {
		fmt.Fprintf(os.Stderr, "\ndlibox: %s.\n", err)
		os.Exit(1)
	}
}
