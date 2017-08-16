// Copyright 2017 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/maruel/dlibox/go/msgbus"
	"github.com/maruel/serve-dir/loghttp"
	"periph.io/x/periph"
	"periph.io/x/periph/host"
)

// mainDevice is the main function when running as a device (a node).
func mainDevice(hostname, server, mqttHost, mqttUser, mqttPass string, port int) error {
	// Initialize periph.
	state, err := host.Init()
	if err != nil {
		return err
	}

	/*
		_, err = initDisplay(bus, &config.Settings.Display)
		if err != nil {
			// Non-fatal.
			log.Printf("Display not connected: %v", err)
		}

		leds, err := initLEDs(bus, *fake, &config.Settings.APA102)
		if err != nil {
			// Non-fatal.
			log.Printf("LEDs: %v", err)
		} else if leds != nil {
			defer leds.Close()
			p, err := initPainter(bus, leds, leds.fps, &config.Settings.Painter, &config.LRU)
			if err != nil {
				return err
			}
			defer p.Close()

		}
		h, err := initHalloween(bus, &config.Settings.Halloween)
		if err != nil {
			// Non-fatal.
			log.Printf("Halloween: %v", err)
		} else if h != nil {
			defer h.Close()
		}

		if err = initButton(bus, &config.Settings.Button); err != nil {
			// Non-fatal.
			log.Printf("Button not connected: %v", err)
		}

		if err = initIR(bus, &config.Settings.IR); err != nil {
			// Non-fatal.
			log.Printf("IR not connected: %v", err)
		}

		if err = initPIR(bus, &config.Settings.PIR); err != nil {
			// Non-fatal.
			log.Printf("PIR not connected: %v", err)
		}

		if err = alarm.Init(bus, &config.Settings.Alarms); err != nil {
			return err
		}

		s, err := initSound(bus, &config.Settings.Sound)
		if err != nil {
			// Non-fatal.
			log.Printf("Sound failed: %v", err)
		} else if s != nil {
			defer s.Close()
		}
	*/

	if port != 0 {
		ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			return err
		}
		if server != "" {
			http.DefaultServeMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				http.Redirect(w, r, "http://"+server, 302)
			})
		}
		server := http.Server{
			Addr:           ln.Addr().String(),
			Handler:        &loghttp.Handler{Handler: http.DefaultServeMux},
			ReadTimeout:    60 * time.Second,
			WriteTimeout:   60 * time.Second,
			MaxHeaderBytes: 1 << 16,
		}

		go server.Serve(ln)
		log.Printf("Visit: http://%s:%d/debug/pprof for debugging", hostname, port)
	}

	root := "dlibox/" + hostname
	bus := msgbus.New()
	if mqttHost != "" {
		server, err := msgbus.NewMQTT(mqttHost, hostname, mqttUser, mqttPass, msgbus.Message{root + "/$online", []byte("false")})
		if err != nil {
			// TODO(maruel): Have it continuously try to automatically reconnect.
			log.Printf("Failed to connect to server: %v", err)
		} else {
			bus = server
		}
	}

	// TODO(maruel): Uses modified Homie convention. The main modification is
	// that it is the controller that decides which nodes to expose ($nodes), not
	// the device itself. This means no need for configuration on the device
	// itself, assuming all devices run all the same code.
	// https://github.com/marvinroger/homie#device-attributes
	bus = msgbus.RebaseSub(msgbus.RebasePub(bus, root), root)
	initState(bus, state)
	return watchFile()
}

func initState(bus msgbus.Bus, state *periph.State) {
	// TODO(maruel): We need to get the MAC and IP of the network that is UP. In
	// the case where multiple networks are up, too bad.
	// TODO(maruel): Use priority eth > wlan > blan
	var ip net.IP
	var mask []byte
	ifaces, _ := net.Interfaces()
	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
				mask = v.Mask
			case *net.IPAddr:
				//ip = v.IP
			}
		}
	}

	bus.Publish(msgbus.Message{"$localip", []byte(ip.String())}, msgbus.MinOnce, true)
	bus.Publish(msgbus.Message{"$mac", mask}, msgbus.MinOnce, true)
	bus.Publish(msgbus.Message{"$implementation", []byte("dlibox")}, msgbus.MinOnce, true)
	if state != nil {
		bus.Publish(msgbus.Message{"$implementation/periph/state", []byte(fmt.Sprintf("%v", state))}, msgbus.MinOnce, true)
	}
	bus.Publish(msgbus.Message{"$online", []byte("true")}, msgbus.MinOnce, true)
}
