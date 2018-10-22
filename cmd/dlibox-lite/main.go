// Copyright 2018 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// dlibox-lite is an home automation system to be used with home assistant and
// similar MQTT based systems.
package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"runtime/pprof"
	"strconv"
	"strings"
	"syscall"

	"github.com/maruel/dlibox/device"
	"github.com/maruel/dlibox/shared"
	"github.com/maruel/interrupt"
	"github.com/maruel/msgbus"
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
	config := flag.String("config", "", "device configuration file")
	mqttHost := flag.String("mqtt", "tcp://dlibox:1883", "MQTT host in the form tcp://user:pass@host:port")
	flag.Parse()
	if flag.NArg() != 0 {
		return fmt.Errorf("unexpected argument: %s", flag.Args())
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
	serverID := ""
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
	serverID = parts[0]
	if i, err := strconv.Atoi(parts[1]); i == 0 || err != nil {
		return errors.New("mqtt port is required")
	}
	if len(u.Path) != 0 {
		return errors.New("mqtt must not include a path")
	}
	if len(u.RawQuery) != 0 {
		return errors.New("mqtt must not not include a query argument")
	}

	root := "dlibox/" + shared.Hostname()
	will := msgbus.Message{Topic: root + "/$online", Payload: []byte("false"), Retained: true}
	usr := ""
	pwd := ""
	if u.User != nil {
		usr = u.User.Username()
		pwd, _ = u.User.Password()
	}
	server := u.Scheme + "://" + u.Host
	clientID := shared.Hostname()

	log.Printf("MQTT:%s ClientID:%s User:%s Pass:%s", server, clientID, usr, pwd)
	mqttServer, err := msgbus.NewMQTT(server, clientID, usr, pwd, will, true)
	if err != nil {
		// TODO(maruel): Have it continuously try to automatically reconnect.
		log.Printf("Failed to connect to server: %v", err)
	} else {
		bus = mqttServer
	}
	bus = msgbus.Log(bus)
	return device.Main(serverID, bus, *port)
}

func main() {
	if err := mainImpl(); err != nil {
		fmt.Fprintf(os.Stderr, "\ndlibox-lite: %s.\n", err)
		os.Exit(1)
	}
}
