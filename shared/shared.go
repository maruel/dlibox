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
	"strings"

	"github.com/maruel/dlibox/msgbus"
	"periph.io/x/periph"
)

// InitState initializes the MQTT node state.
func InitState(bus msgbus.Bus, state *periph.State) {
	ip := ""
	mac := ""
	ifaces, _ := net.Interfaces()
	for _, i := range ifaces {
		if i.Flags&net.FlagUp == 0 || i.Flags&net.FlagLoopback != 0 {
			continue
		}
		if strings.HasPrefix(i.Name, "virbr") || strings.HasPrefix(i.Name, "docker") {
			continue
		}
		addrs, _ := i.Addrs()
		mac = i.HardwareAddr.String()
		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				if v.IP.IsLoopback() || v.IP.IsUnspecified() {
					continue
				}
				ip = v.IP.String()
				goto done
			}
		}
	}
done:

	bus.Publish(msgbus.Message{"$localip", []byte(ip)}, msgbus.MinOnce, true)
	bus.Publish(msgbus.Message{"$mac", []byte(mac)}, msgbus.MinOnce, true)
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
