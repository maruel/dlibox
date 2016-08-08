// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// InfraRed receiver support through native linux lirc.

package rpi

import (
	"bufio"
	"net"
	"strconv"
	"strings"
)

// IR is an open port to lirc.
type IR struct {
	c net.Conn
	r *bufio.Reader
}

// MakeIR returns a IR receiver.
//
// Returns nil on failure.
//
// You can try to 'irrecord -a -d /var/run/lirc/lircd ~/lircd.conf' and identify
// a few key, then 'grep -R '<hex value>' /usr/share/lirc/remotes/' on one of
// the key you found. Then copy the file found or your generated config file to
// /etc/lirc/lircd.conf.
//
// To debug, use either 'irw' or 'nc -U /var/run/lirc/lircd' and use your IR
// remote to test if lirc is correctly configured.
//
// http://alexba.in/blog/2013/01/06/setting-up-lirc-on-the-raspberrypi/ is a
// good starter guide.
func MakeIR() *IR {
	c, err := net.Dial("unix", "/var/run/lirc/lircd")
	if err != nil {
		return nil
	}
	return &IR{c, bufio.NewReader(c)}
}

func (r *IR) Close() error {
	err := r.c.Close()
	r.c = nil
	r.r = nil
	return err
}

// Next blocks until the next IR keypress is received.
//
// Returns the key and the repeat count.
//
// Use one of the following to list all the valid key names:
//
//     irrecord -l
//     grep -hoER '(BTN|KEY)_\w+' /usr/share/lirc/remotes | sort | uniq | less
//
// http://www.lirc.org/api-docs/html/input__map_8inc_source.html
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
