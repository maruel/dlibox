// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// InfraRed receiver support through native linux lirc.

package rpi

import (
	"bufio"
	"io/ioutil"
	"net"
	"strconv"
	"strings"
)

// IR is an open port to lirc.
type IR struct {
	c net.Conn
	r *bufio.Reader
}

// MakeIR returns a IR receiver. Returns nil on failure.
func MakeIR() *IR {
	c, err := net.Dial("unix", "/var/run/lirc/lircd")
	if err != nil {
		return nil
	}
	return &IR{c, bufio.NewReader(c)}
}

// Close closes the socket to lirc. It is not a requirement to close before
// process termination.
func (r *IR) Close() error {
	err := r.c.Close()
	r.c = nil
	r.r = nil
	return err
}

// Next blocks until the next IR keypress is received.
//
// Returns the key and the repeat count.
func (r *IR) Next() (string, int, error) {
	// TODO(maruel): handle when isPrefix is set.
	line, _, err := r.r.ReadLine()
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

//

// getLIRCPins queries the kernel module to determine which GPIO pins are taken
// by the driver.
func getLIRCPins() (Pin, Pin) {
	// This is configured in /boot/config.txt as:
	// dtoverlay=gpio_in_pin=23,gpio_out_pin=22
	bytes, err := ioutil.ReadFile("/sys/module/lirc_rpi/parameters/gpio_in_pin")
	if err != nil || len(bytes) == 0 {
		return INVALID, INVALID
	}
	in, err := strconv.Atoi(strings.TrimRight(string(bytes), "\n"))
	if err != nil {
		return INVALID, INVALID
	}
	bytes, err = ioutil.ReadFile("/sys/module/lirc_rpi/parameters/gpio_out_pin")
	if err != nil || len(bytes) == 0 {
		return INVALID, INVALID
	}
	out, err := strconv.Atoi(strings.TrimRight(string(bytes), "\n"))
	if err != nil {
		return INVALID, INVALID
	}
	return Pin(in), Pin(out)
}
