// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Package ir implements InfraRed receiver support through native linux lirc.
//
// See http://www.lirc.org/ for details about daemon configuration.
package ir

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"

	"github.com/maruel/dlibox/go/pio/buses"
)

// Pins queries the kernel module to determine which GPIO pins are taken by
// the driver.
//
// The return values can be converted to bcm238x.Pin. Return (-1, -1) on
// failure.
func Pins() (int, int) {
	// This is configured in /boot/config.txt as:
	// dtoverlay=gpio_in_pin=23,gpio_out_pin=22
	bytes, err := ioutil.ReadFile("/sys/module/lirc_rpi/parameters/gpio_in_pin")
	if err != nil || len(bytes) == 0 {
		return -1, -1
	}
	in, err := strconv.Atoi(strings.TrimRight(string(bytes), "\n"))
	if err != nil {
		return -1, -1
	}
	bytes, err = ioutil.ReadFile("/sys/module/lirc_rpi/parameters/gpio_out_pin")
	if err != nil || len(bytes) == 0 {
		return -1, -1
	}
	out, err := strconv.Atoi(strings.TrimRight(string(bytes), "\n"))
	if err != nil {
		return -1, -1
	}
	return in, out
}

// Bus is an open port to lirc.
type Bus struct {
	w           net.Conn
	c           chan buses.Message
	lock        sync.Mutex
	list        map[string][]string // list of remotes and associated keys
	pendingList map[string][]string // list of remotes and associated keys being created.
}

// Make returns a IR receiver / emitter handle.
func Make() (*Bus, error) {
	w, err := net.Dial("unix", "/var/run/lirc/lircd")
	if err != nil {
		return nil, err
	}
	b := &Bus{w: w, c: make(chan buses.Message), list: map[string][]string{}}
	// Inconditionally retrieve the list of all known buttons at start.
	if _, err := w.Write([]byte("LIST\n")); err != nil {
		w.Close()
		return nil, err
	}
	go b.loop(bufio.NewReader(w))
	return b, nil
}

// Close closes the socket to lirc. It is not a requirement to close before
// process termination.
func (b *Bus) Close() error {
	return b.w.Close()
}

// Emit implements buses.IR.
func (b *Bus) Emit(remote, button string) error {
	// http://www.lirc.org/html/lircd.html#lbAH
	_, err := fmt.Fprintf(b.w, "SEND_ONCE %s %s", remote, button)
	return err
}

// Channel implements buses.IR.
func (b *Bus) Channel() <-chan buses.Message {
	return b.c
}

// Codes returns all the known codes.
//
// Empty if the list was not retrieved yet.
func (b *Bus) Codes() map[string][]string {
	b.lock.Lock()
	defer b.lock.Unlock()
	return b.list
}

func (b *Bus) loop(r *bufio.Reader) {
	defer func() {
		close(b.c)
		b.c = nil
	}()
	for {
		line, err := read(r)
		if line == "BEGIN" {
			err = b.readData(r)
		} else if len(line) != 0 {
			// Format is: <code> <repeat count> <button name> <remote control name>
			// http://www.lirc.org/html/lircd.html#lbAG
			if parts := strings.SplitN(line, " ", 5); len(parts) != 4 {
				log.Printf("ir: corrupted line: #v", line)
			} else {
				if i, err2 := strconv.Atoi(parts[1]); err2 != nil {
					log.Printf("ir: corrupted line: #v", line)
				} else {
					b.c <- buses.Message{parts[2], i + 1, parts[3]}
				}
			}
		}
		if err != nil {
			break
		}
	}
}

func (b *Bus) readData(r *bufio.Reader) error {
	// Format is:
	// BEGIN
	// <original command>
	// SUCCESS
	// DATA
	// <number of entries 1 based>
	// <entries>
	// ...
	// END
	cmd, err := read(r)
	if err != nil {
		return err
	}
	switch cmd {
	case "SIGHUP":
		_, err = b.w.Write([]byte("LIST\n"))
	default:
		// In case of any error, ignore the rest.
		line, err := read(r)
		if err != nil {
			return err
		}
		if line != "SUCCESS" {
			log.Printf("ir: unexpected line: %v, expected SUCCESS", line)
			return nil
		}
		if line, err = read(r); err != nil {
			return err
		}
		if line != "DATA" {
			log.Printf("ir: unexpected line: %v, expected DATA", line)
			return nil
		}
		if line, err = read(r); err != nil {
			return err
		}
		nbLines, err := strconv.Atoi(line)
		if err != nil {
			return err
		}
		list := make([]string, nbLines)
		for i := 0; i < nbLines; i++ {
			if list[i], err = read(r); err != nil {
				return err
			}
		}
		switch {
		case cmd == "LIST":
			// Request the codes for each remote.
			b.pendingList = map[string][]string{}
			for _, l := range list {
				if _, ok := b.pendingList[l]; ok {
					log.Printf("ir: unexpected %s", cmd)
				} else {
					b.pendingList[l] = []string{}
					if _, err = fmt.Fprintf(b.w, "LIST %s\n", l); err != nil {
						return err
					}
				}
			}
		case strings.HasPrefix(line, "LIST "):
			if b.pendingList == nil {
				log.Printf("ir: unexpected %s", cmd)
			} else {
				remote := cmd[5:]
				b.pendingList[remote] = list
				all := true
				for _, v := range b.pendingList {
					if len(v) == 0 {
						all = false
						break
					}
				}
				if all {
					b.lock.Lock()
					b.list = b.pendingList
					b.pendingList = nil
					b.lock.Unlock()
				}
			}
		default:
		}
	}
	line, err := read(r)
	if err != nil {
		return err
	}
	if line != "END" {
		log.Printf("ir: unexpected line: %v, expected END", line)
	}
	return nil
}

func read(r *bufio.Reader) (string, error) {
	raw, err := r.ReadBytes('\n')
	if err != nil {
		return "", err
	}
	if len(raw) != 0 {
		raw = raw[:len(raw)-1]
	}
	return string(raw), nil
}

var _ buses.IR = &Bus{}
