// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package sysfs

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/maruel/dlibox/go/pio"
	"github.com/maruel/dlibox/go/pio/protocols/gpio"
)

// Pins is all the pins exported by GPIO sysfs.
var Pins map[int]*Pin

// PinByNumber returns a *Pin for the pin number, if any.
func PinByNumber(i int) (*Pin, error) {
	p, ok := Pins[i]
	if !ok {
		return nil, errors.New("invalid pin number")
	}
	if err := p.open(); err != nil {
		return nil, err
	}
	return p, nil
}

// Pin represents one GPIO pin as found by sysfs.
type Pin struct {
	number int
	name   string
	root   string // Something like /sys/class/gpio/gpio%d/

	lock       sync.Mutex
	direction  int             // Cache of the last known direction
	fDirection *os.File        // handle to /sys/class/gpio/gpio*/direction; never closed
	fEdge      *os.File        // handle to /sys/class/gpio/gpio*/edge; never closed
	fValue     *os.File        // handle to /sys/class/gpio/gpio*/value; never closed
	epollFd    int             // Never closed
	event      event           // Initialized once
	edges      chan gpio.Level // Closed when edges are terminated
	wg         sync.WaitGroup  // Set when Edges() is running
}

func (p *Pin) String() string {
	return p.name
}

func (p *Pin) Number() int {
	return p.number
}

func (p *Pin) Function() string {
	p.lock.Lock()
	defer p.lock.Unlock()
	// TODO(maruel): There's an internal bug which causes p.direction to be invalid (!?)
	// Need to figure it out ASAP.
	if err := p.open(); err != nil {
		return err.Error()
	}
	var buf [4]byte
	if _, err := p.fDirection.Seek(0, 0); err != nil {
		return err.Error()
	}
	if _, err := p.fDirection.Read(buf[:]); err != nil {
		return err.Error()
	}
	if buf[0] == 'i' && buf[1] == 'n' {
		p.direction = dIn
	} else if buf[0] == 'o' && buf[1] == 'u' && buf[2] == 't' {
		p.direction = dOut
	}
	if p.direction == dIn {
		return "In/" + p.Read().String()
	} else if p.direction == dOut {
		return "Out/" + p.Read().String()
	}
	return "N/A"
}

// In setups a pin as an input.
func (p *Pin) In(pull gpio.Pull) error {
	if pull != gpio.PullNoChange && pull != gpio.Float {
		return errors.New("not implemented")
	}
	p.lock.Lock()
	defer p.lock.Unlock()
	if p.direction == dIn {
		return nil
	}
	if err := p.open(); err != nil {
		return err
	}
	if _, err := p.fDirection.Seek(0, 0); err != nil {
		return err
	}
	if _, err := p.fDirection.Write([]byte("in")); err != nil {
		return err
	}
	p.direction = dIn
	return nil
}

func (p *Pin) Read() gpio.Level {
	var buf [2]byte
	if _, err := p.fValue.Seek(0, 0); err != nil {
		// Error.
		//fmt.Printf("%s: %v", p, err)
		return gpio.Low
	}
	if _, err := p.fValue.Read(buf[:]); err != nil {
		// Error.
		//fmt.Printf("%s: %v", p, err)
		return gpio.Low
	}
	if buf[0] == '0' {
		return gpio.Low
	}
	if buf[0] == '1' {
		return gpio.High
	}
	// Error.
	return gpio.Low
}

// Edges creates a edge detection loop and implements gpio.PinIn.
//
// It is the function that opens the gpio sysfs file handle for /edge.
func (p *Pin) Edges() (<-chan gpio.Level, error) {
	p.lock.Lock()
	defer p.lock.Unlock()
	if p.edges != nil {
		return nil, errors.New("must call DisableEdges() between Edges()")
	}
	var err error
	if p.fEdge == nil {
		p.fEdge, err = os.OpenFile(p.root+"edge", os.O_RDWR|os.O_APPEND, 0600)
		if err != nil {
			return nil, err
		}
	}
	if p.epollFd == 0 {
		if p.epollFd, err = p.event.makeEvent(p.fValue); err != nil {
			return nil, err
		}
	}
	if _, err := p.fEdge.Seek(0, 0); err != nil {
		return nil, err
	}
	if _, err = p.fEdge.Write([]byte("both")); err != nil {
		return nil, err
	}
	p.edges = make(chan gpio.Level)
	p.wg.Add(1)
	var started sync.WaitGroup
	started.Add(1)
	go p.edgeLoop(&started)
	started.Wait()
	return p.edges, nil
}

// DisableEdges stops a previous Edges() call.
func (p *Pin) DisableEdges() {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.disableEdge()
}

func (p *Pin) Pull() gpio.Pull {
	return gpio.PullNoChange
}

// Out sets a pin as output.
func (p Pin) Out(l gpio.Level) error {
	p.lock.Lock()
	defer p.lock.Unlock()
	if p.direction != dOut {
		if err := p.open(); err != nil {
			return err
		}
		// Cancel any outstanding edge detection.
		p.disableEdge()

		// "To ensure glitch free operation, values "low" and "high" may be written
		// to configure the GPIO as an output with that initial value."
		if _, err := p.fDirection.Seek(0, 0); err != nil {
			return err
		}
		var d []byte
		if l == gpio.Low {
			d = []byte("low")
		} else {
			d = []byte("high")
		}
		if _, err := p.fDirection.Write(d); err != nil {
			return err
		}
		p.direction = dOut
		return nil
	}
	if _, err := p.fValue.Seek(0, 0); err != nil {
		return nil
	}
	var d [2]byte
	if l == gpio.Low {
		d[0] = '0'
	} else {
		d[0] = '1'
	}
	_, err := p.fValue.Write(d[:])
	return err
}

//

// open opens the gpio sysfs handle to /value and direction.
//
// lock must be held.
func (p *Pin) open() error {
	if exportHandle == nil {
		return errors.New("sysfs gpio is not initialized")
	}
	if p.fDirection != nil {
		return nil
	}
	// Ignore the failure unless this is a permission failure. Exporting a pin
	// that is already exported causes a write failure.
	_, err := exportHandle.Write([]byte(fmt.Sprintf("%d\n", p.number)))
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
	p.fDirection, err = os.OpenFile(p.root+"direction", os.O_RDWR|os.O_APPEND, 0600)
	if err != nil {
		p.fValue.Close()
		p.fValue = nil
	}
	return err
}

// disableEdge disable the edge detection setting for the pin, if any.
func (p *Pin) disableEdge() error {
	if p.direction != dIn {
		return errors.New("pin wasn't set as input first")
	}
	if p.edges != nil {
		// Drain it if needed. This works because p.edges is not buffered.
		select {
		case <-p.edges:
		default:
		}
		// Only after it is safe to close.
		close(p.edges)
		p.edges = nil
		if _, err := p.fEdge.Seek(0, 0); err != nil {
			return err
		}
		if _, err := p.fEdge.Write([]byte("none")); err != nil {
			return err
		}
	}
	return nil
}

func (p *Pin) edgeLoop(started *sync.WaitGroup) {
	defer p.wg.Done()
	last := p.Read()
	started.Done()
	var b [1]byte
	for {
		for {
			p.lock.Lock()
			c := p.edges
			p.lock.Unlock()
			if c == nil {
				return
			}
			if _, err := p.fValue.Seek(0, 0); err != nil {
				log.Printf("edgeLoop() ended: %v\n", err)
				return
			}
			if nr, err := p.event.wait(p.epollFd); err != nil {
				log.Printf("edgeLoop() ended: %v\n", err)
				return
			} else if nr < 1 {
				continue
			}
			if _, err := p.fValue.Read(b[:]); err != nil {
				log.Printf("edgeLoop() ended: %v\n", err)
				return
			}
			break
		}
		p.lock.Lock()
		c := p.edges
		p.lock.Unlock()
		if c == nil {
			return
		}
		// Make sure to ignore spurious wake up.
		if b[0] == '1' {
			if last != gpio.High {
				c <- gpio.High
				last = gpio.High
			}
		} else {
			if last != gpio.Low {
				c <- gpio.Low
				last = gpio.Low
			}
		}
	}
}

//

var exportHandle io.Writer // handle to /sys/class/gpio/export

const (
	dUnknown = 0
	dIn      = 1
	dOut     = 2
)

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

// driverGPIO implements pio.Driver.
type driverGPIO struct {
}

func (d *driverGPIO) String() string {
	return "sysfs-gpio"
}

func (d *driverGPIO) Type() pio.Type {
	// It intentionally load later than processors.
	return pio.Pins
}

func (d *driverGPIO) Prerequisites() []string {
	return nil
}

// Init initializes GPIO sysfs handling code.
//
// Uses gpio sysfs as described at
// https://www.kernel.org/doc/Documentation/gpio/sysfs.txt
//
// GPIO sysfs is often the only way to do edge triggered interrupts. Doing this
// requires cooperation from a driver in the kernel.
//
// The main drawback of GPIO sysfs is that it doesn't expose internal pull
// resistor and it is much slower than using memory mapped hardware registers.
func (d *driverGPIO) Init() (bool, error) {
	f, err := os.OpenFile("/sys/class/gpio/export", os.O_WRONLY, 0600)
	if err != nil {
		return true, err
	}
	items, err := ioutil.ReadDir("/sys/class/gpio/")
	if err != nil {
		f.Close()
		return true, err
	}
	// There are hosts that use non-continuous pin numbering.
	Pins = map[int]*Pin{}
	for _, item := range items {
		name := item.Name()
		if !strings.HasPrefix(name, "gpiochip") {
			continue
		}
		if err := d.exportGPIOChip("/sys/class/gpio/" + name + "/"); err != nil {
			f.Close()
			return true, err
		}
	}
	exportHandle = f
	return true, nil
}

func (d *driverGPIO) exportGPIOChip(path string) error {
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
		if _, ok := Pins[i]; ok {
			return fmt.Errorf("found two pins with number %d", i)
		}
		p := &Pin{
			number: i,
			name:   fmt.Sprintf("GPIO%d", i),
			root:   fmt.Sprintf("/sys/class/gpio/gpio%d/", i),
		}
		Pins[i] = p
		// Intentionally ignore any error. In that case, a Processor driver already
		// registered the pin.
		_ = gpio.Register(p)
		// We cannot use gpio.MapFunction() since there is no API to determine this.
	}
	return nil
}
