// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Package ir implements InfraRed receiver support through native linux lirc.
//
// See http://www.lirc.org/ for details about daemon configuration.
package ir

import (
	"bufio"
	"io/ioutil"
	"net"
	"strconv"
	"strings"
)

// Bus is an open port to lirc.
type Bus struct {
	c net.Conn
	r *bufio.Reader
}

// Make returns a IR receiver / emitter handle. Returns nil on failure.
func Make() *Bus {
	c, err := net.Dial("unix", "/var/run/lirc/lircd")
	if err != nil {
		return nil
	}
	return &Bus{c, bufio.NewReader(c)}
}

// Close closes the socket to lirc. It is not a requirement to close before
// process termination.
func (b *Bus) Close() error {
	err := b.c.Close()
	b.c = nil
	b.r = nil
	return err
}

// Next blocks until the next IR keypress is received.
//
// Returns the key and the repeat count.
func (b *Bus) Next() (string, int, error) {
	// TODO(maruel): handle when isPrefix is set.
	line, _, err := b.r.ReadLine()
	if err != nil {
		return "", 0, err
	}
	// Format is: <code> <repeat count> <button name> <remote control name>
	// http://www.lirc.org/html/lircd.html#lbAG
	parts := strings.SplitN(string(line), " ", 5)
	// Ignore corrupted output and notifications like "BEGIN", "SIGHUP" and "END".
	if len(parts) != 5 {
		return "", 0, nil
	}
	i, err := strconv.Atoi(parts[0])
	if err != nil {
		return "", 0, nil
	}
	return parts[3], i, nil
}

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
