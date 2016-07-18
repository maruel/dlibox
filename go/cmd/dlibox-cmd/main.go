// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// dlibox-cmd is meant to run on a host to query via mDNS and MQTT the current
// dlibox instances.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/hashicorp/mdns"
)

func getInterfaces() ([]net.Interface, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	var out []net.Interface
	for _, i := range interfaces {
		if i.Flags&net.FlagUp == 0 || i.Flags&net.FlagLoopback != 0 || i.Flags&net.FlagMulticast == 0 {
			continue
		}
		out = append(out, i)
	}
	return out, err
}

// query prints out all the mDNS hosts advertized on the local network.
func query() error {
	var wgI sync.WaitGroup
	var wgE sync.WaitGroup
	const timeout = 5 * time.Second
	ifs, err := getInterfaces()
	if err != nil {
		return nil
	}

	entries := make(chan *mdns.ServiceEntry)
	wgE.Add(1)
	go func() {
		defer wgE.Done()
		for e := range entries {
			fmt.Printf("%-40s  %-15s  %s  %-5d\n", e.Name, e.AddrV4, e.AddrV6, e.Port)
			for _, i := range e.InfoFields {
				fmt.Printf("  - %s\n", i)
			}
		}
	}()

	errs := make([]error, len(ifs))
	for i := range ifs {
		wgI.Add(1)
		go func(i int) {
			defer wgI.Done()
			params := mdns.QueryParam{
				Service:   "_dlibox._tcp",
				Domain:    "local",
				Timeout:   timeout,
				Interface: &ifs[i],
				Entries:   entries,
			}
			errs[i] = mdns.Query(&params)
		}(i)
	}
	wgI.Wait()
	close(entries)
	wgE.Wait()
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}

func mainImpl() error {
	verbose := flag.Bool("verbose", false, "enable log output")
	flag.Parse()
	if flag.NArg() != 0 {
		return fmt.Errorf("unexpected argument: %s", flag.Args())
	}

	if !*verbose {
		log.SetOutput(ioutil.Discard)
	}

	return query()
}

func main() {
	if err := mainImpl(); err != nil {
		fmt.Fprintf(os.Stderr, "\ndlibox-cmd: %s.\n", err)
		os.Exit(1)
	}
}
