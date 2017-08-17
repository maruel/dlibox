// Copyright 2017 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Package shared contains shared code between the controller and the device(s).
package shared

import (
	"errors"
	"fmt"
	"net"
	"os"
	"os/user"

	"github.com/maruel/dlibox/msgbus"
	"periph.io/x/periph"
)

// InitState initializes the MQTT node state.
func InitState(bus msgbus.Bus, state *periph.State) {
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
}

var (
	hostname string
	home     string
)

// Home returns the home directory even when cross compiled and panics on
// failure.
//
// When cross compiling, user.Current() fails.
func Home() string {
	if home != "" {
		return home
	}
	u, err := user.Current()
	if err == nil {
		home = u.HomeDir
		return home
	}
	if home = os.Getenv("HOME"); len(home) != 0 {
		return home
	}
	panic(errors.New("can't find HOME"))
}

// Hostname is like os.Hostname() but caches the value and panics on failure.
func Hostname() string {
	if hostname != "" {
		return hostname
	}
	var err error
	if hostname, err = os.Hostname(); err != nil {
		panic(err)
	}
	return hostname
}
