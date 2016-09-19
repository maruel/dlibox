// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package sysfs

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/maruel/dlibox/go/pio/host"
)

// GetPin returns a *Pin for the pin number, if any.
func GetPin(i int) (*Pin, error) {
	if err := Init(); err != nil {
		return nil, err
	}
	lock.Lock()
	defer lock.Unlock()
	p, ok := Pins[i]
	if !ok {
		return nil, errors.New("invalid pin number")
	}
	if err := p.openValue(); err != nil {
		return nil, err
	}
	return p, nil
}

type Pin struct {
	number int
	name   string
	root   string // Something like /sys/class/gpio/gpio%d/

	lock       sync.Mutex
	fDirection *os.File
	fEdge      *os.File
	fValue     *os.File // handle to /sys/class/gpio/gpio*/value.
	epollFd    int
	event      event
}

func (p *Pin) String() string {
	return p.name
}

func (p *Pin) Number() int {
	return p.number
}

func (p *Pin) Function() string {
	if err := p.openValue(); err != nil {
		return err.Error()
	}
	if err := p.openDirection(); err != nil {
		return err.Error()
	}
	var buf [4]byte
	p.lock.Lock()
	defer p.lock.Unlock()
	if _, err := p.fDirection.Seek(0, 0); err != nil {
		return err.Error()
	}
	if _, err := p.fDirection.Read(buf[:]); err != nil {
		return err.Error()
	}
	if buf[0] == 'i' && buf[1] == 'n' {
		return "In/" + p.Read().String()
	} else if buf[0] == 'o' && buf[1] == 'u' && buf[2] == 't' {
		return "Out/" + p.Read().String()
	}
	return "N/A"
}

// In setups a pin as an input.
func (p *Pin) In(pull host.Pull) error {
	if err := p.setIn(true); err != nil {
		return err
	}
	if pull != host.PullNoChange && pull != host.Float {
		return errors.New("not implemented")
	}
	return nil
}

func (p *Pin) Read() host.Level {
	var buf [2]byte
	if _, err := p.fValue.Seek(0, 0); err != nil {
		// Error.
		fmt.Printf("%s: %v", p, err)
		return host.Low
	}
	if _, err := p.fValue.Read(buf[:]); err != nil {
		// Error.
		fmt.Printf("%s: %v", p, err)
		return host.Low
	}
	if buf[0] == '0' {
		return host.Low
	}
	if buf[0] == '1' {
		return host.High
	}
	// Error.
	return host.Low
}

// Edges creates a edge detection loop and implements host.PinIn.
func (p *Pin) Edges() (chan host.Level, error) {
	last := p.Read()
	if err := p.setEdge(true); err != nil {
		return nil, err
	}
	c := make(chan host.Level)
	go func() {
		defer close(c)
		var b [1]byte
		for {
			if _, err := p.fValue.Seek(0, 0); err != nil {
				return
			}
			for {
				p.lock.Lock()
				ep := p.epollFd
				v := p.fValue
				p.lock.Unlock()
				if ep == 0 {
					return
				}
				if nr, err := p.event.wait(ep); err != nil {
					return
				} else if nr < 1 {
					continue
				}
				if _, err := v.Seek(0, 0); err != nil {
					return
				}
				if _, err := v.Read(b[:]); err != nil {
					return
				}
				break
			}
			// Make sure to ignore spurious wake up.
			if b[0] == '1' {
				if last != host.High {
					c <- host.High
					last = host.High
				}
			} else {
				if last != host.Low {
					c <- host.Low
					last = host.Low
				}
			}
		}
	}()
	return c, nil
}

// DisableEdge stops edges.
func (p *Pin) DisableEdge() {
	p.setEdge(false)
}

// Out sets a pin as output
func (p Pin) Out() error {
	return p.setIn(false)
}

func (p Pin) Set(l host.Level) {
	p.lock.Lock()
	defer p.lock.Unlock()
	var d [2]byte
	if l == host.Low {
		d[0] = '0'
	} else {
		d[0] = '1'
	}
	if _, err := p.fValue.Seek(0, 0); err != nil {
		// Error.
	}
	if _, err := p.fValue.Write(d[:]); err != nil {
		// Error.
	}
}

//

// openValue opens the gpio sysfs handle to /value.
func (p *Pin) openValue() error {
	p.lock.Lock()
	defer p.lock.Unlock()
	var err error
	if p.fValue == nil {
		// Ignore the failure unless this is a permission failure. Exporting a pin
		// that is already exported causes a write failure.
		_, err = exportHandle.Write([]byte(fmt.Sprintf("%d\n", p.number)))
		if os.IsPermission(err) {
			return err
		}
		// There's a race condition where the file may be created but udev is still
		// running the script to make it readable to the current user. It's simpler
		// to just loop a little as if /export is accessible, it doesn't make sense
		// that gpioN/value doesn't become accessible eventually.
		timeout := 5 * time.Second
		for start := time.Now(); time.Since(start) < timeout; {
			p.fValue, err = os.OpenFile(p.root+"value", os.O_RDWR, 0600)
			// The virtual file creation is synchronous when writing to /export for
			// udev rule execution is asynchronous.
			if err == nil || !os.IsPermission(err) {
				break
			}
		}
	}
	return err
}

// openDirection opens the gpio sysfs handle to /direction.
//
// Assumes openValue() succeeded before.
func (p *Pin) openDirection() error {
	p.lock.Lock()
	defer p.lock.Unlock()
	if p.fDirection == nil {
		var err error
		if p.fDirection, err = os.OpenFile(p.root+"direction", os.O_RDWR|os.O_APPEND, 0600); err != nil {
			return err
		}
	}
	return nil
}

func (p *Pin) setIn(asIn bool) error {
	if err := p.openValue(); err != nil {
		return err
	}
	if err := p.openDirection(); err != nil {
		return err
	}
	var b [3]byte
	if asIn {
		b[0] = 'i'
		b[1] = 'n'
	} else {
		b[0] = 'o'
		b[1] = 'u'
		b[2] = 't'
	}
	p.lock.Lock()
	defer p.lock.Unlock()
	if _, err := p.fDirection.Seek(0, 0); err != nil {
		return err
	}
	if _, err := p.fDirection.Write(b[:]); err != nil {
		return err
	}
	return nil
}

// setEdge changes the edge detection setting for the pin.
//
// It is the function that opens the gpio sysfs file handle.
func (p *Pin) setEdge(enable bool) error {
	p.lock.Lock()
	defer p.lock.Unlock()
	if !enable {
		if p.fEdge != nil {
			if _, err := p.fEdge.Seek(0, 0); err != nil {
				return err
			}
			_, err := p.fEdge.Write([]byte("none"))
			return err
		}
		return nil
	}
	var err error
	if p.fEdge == nil {
		p.fEdge, err = os.OpenFile(p.root+"edge", os.O_RDWR|os.O_APPEND, 0600)
		if err != nil {
			return err
		}
	}
	if p.epollFd == 0 {
		if p.epollFd, err = p.event.makeEvent(p.fValue); err != nil {
			return err
		}
	}
	if _, err := p.fEdge.Seek(0, 0); err != nil {
		return err
	}
	_, err = p.fEdge.Write([]byte("both"))
	return err
}

//

var (
	lock         sync.Mutex
	exportHandle io.WriteCloser // handle to /sys/class/gpio/export
	Pins         map[int]*Pin   // Pins is all the pins exported by GPIO sysfs.
)

func initLinux() error {
	lock.Lock()
	defer lock.Unlock()
	if exportHandle != nil {
		return nil
	}
	f, err := os.OpenFile("/sys/class/gpio/export", os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	items, err := ioutil.ReadDir("/sys/class/gpio/")
	if err != nil {
		f.Close()
		return err
	}
	// There is host that use non-continuous pin numbering.
	Pins = map[int]*Pin{}
	for _, item := range items {
		name := item.Name()
		if !strings.HasPrefix(name, "gpiochip") {
			continue
		}
		if err := exportGPIOChip("/sys/class/gpio/" + name + "/"); err != nil {
			f.Close()
			return err
		}
	}
	exportHandle = f
	return nil
}

func readInt(path string) (int, error) {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return 0, err
	}
	if len(raw) == 0 || raw[len(raw)-1] != '\n' {
		return 0, errors.New("invalid value")
	}
	return strconv.Atoi(string(raw[:len(raw)-1]))
}

func exportGPIOChip(path string) error {
	base, err := readInt(path + "base")
	if err != nil {
		return err
	}
	number, err := readInt(path + "ngpio")
	if err != nil {
		return err
	}
	// TODO(maruel): The chip driver may lie and lists GPIO pins that cannot be
	// exported. The only way to know about it is to export it before opening.
	for i := base; i < base+number; i++ {
		Pins[i] = &Pin{number: i, name: fmt.Sprintf("GPIO%d", i), root: fmt.Sprintf("/sys/class/gpio/gpio%d/", i)}
	}
	return nil
}
