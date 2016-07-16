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
	"os"
	"sync"
	"time"

	"github.com/hashicorp/mdns"
)

// query prints out all the mDNS hosts advertized on the local network.
func query() error {
	var wg sync.WaitGroup
	entries := make(chan *mdns.ServiceEntry)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for e := range entries {
			fmt.Printf("%-40s  %-15s  %s  %-5d\n", e.Name, e.AddrV4, e.AddrV6, e.Port)
			for _, i := range e.InfoFields {
				fmt.Printf("  - %s\n", i)
			}
		}
	}()
	params := mdns.QueryParam{
		//Service: "dlibox",
		Domain:  "local",
		Timeout: 10 * time.Second,
		Entries: entries,
	}
	if err := mdns.Query(&params); err != nil {
		return err
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
