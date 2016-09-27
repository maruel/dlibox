// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package drivers

import (
	"errors"
	"testing"
)

func TestInitSimple(t *testing.T) {
	initTest([]Driver{
		&driver{
			name:    "CPU",
			t:       Processor,
			prereqs: nil,
			ok:      true,
			err:     nil,
		},
	})
	if len(allDrivers[Processor]) != 1 {
		t.Fatalf("%v", allDrivers)
	}
	if len(byName) != 1 {
		t.Fatalf("%v", byName)
	}
	actual, errs := Init()
	if len(errs) != 0 {
		t.Fatalf("%v", errs)
	}
	if len(actual) != 1 {
		t.Fatalf("%v", actual)
	}
}

func TestInitSkip(t *testing.T) {
	initTest([]Driver{
		&driver{
			name:    "CPU",
			t:       Processor,
			prereqs: nil,
			ok:      false,
			err:     nil,
		},
	})
	actual, errs := Init()
	if len(errs) != 0 {
		t.Fatalf("%v", errs)
	}
	if len(actual) != 0 {
		t.Fatalf("%v", actual)
	}
}

func TestInitErr(t *testing.T) {
	initTest([]Driver{
		&driver{
			name:    "CPU",
			t:       Processor,
			prereqs: nil,
			ok:      true,
			err:     errors.New("oops"),
		},
	})
	actual, errs := Init()
	if len(errs) != 1 {
		t.Fatalf("%v", errs)
	}
	if len(actual) != 0 {
		t.Fatalf("%v", actual)
	}
}

func TestInitBadOrder(t *testing.T) {
	initTest([]Driver{
		&driver{
			name:    "CPU",
			t:       Processor,
			prereqs: []string{"Board"},
			ok:      true,
			err:     nil,
		},
		&driver{
			name:    "Board",
			t:       Pins,
			prereqs: nil,
			ok:      true,
			err:     nil,
		},
	})
	actual, errs := Init()
	if len(errs) != 1 {
		t.Fatalf("%v", errs)
	}
	if len(actual) != 0 {
		t.Fatalf("%v", actual)
	}
}

func TestInitMissing(t *testing.T) {
	initTest([]Driver{
		&driver{
			name:    "CPU",
			t:       Processor,
			prereqs: []string{"Board"},
			ok:      true,
			err:     nil,
		},
	})
	actual, errs := Init()
	if len(errs) != 1 {
		t.Fatalf("%v", errs)
	}
	if len(actual) != 0 {
		t.Fatalf("%v", actual)
	}
}

func TestExplodeStagesSimple(t *testing.T) {
	d := []Driver{
		&driver{
			name:    "CPU",
			t:       Processor,
			prereqs: nil,
			ok:      true,
			err:     nil,
		},
	}
	initTest(d)
	actual, err := explodeStages(d)
	if len(actual) != 1 || len(actual[0]) != 1 {
		t.Fatalf("%v", actual)
	}
	if err != nil {
		t.Fatalf("%v", err)
	}
}

func TestExplodeStages1Dep(t *testing.T) {
	// This explodes the stage into two.
	d := []Driver{
		&driver{
			name:    "CPU-specialized",
			t:       Processor,
			prereqs: []string{"CPU-generic"},
			ok:      true,
			err:     nil,
		},
		&driver{
			name:    "CPU-generic",
			t:       Processor,
			prereqs: nil,
			ok:      true,
			err:     nil,
		},
	}
	initTest(d)
	actual, err := explodeStages(d)
	if len(actual) != 2 || len(actual[0]) != 1 || len(actual[1]) != 1 {
		t.Fatalf("%v", actual)
	}
	if err != nil {
		t.Fatalf("%v", err)
	}
}

func TestExplodeStages3Dep(t *testing.T) {
	// This explodes the stage into 3 due to diamond shaped DAG.
	d := []Driver{
		&driver{
			name:    "base2",
			t:       Processor,
			prereqs: []string{"root"},
			ok:      true,
			err:     nil,
		},
		&driver{
			name:    "base1",
			t:       Processor,
			prereqs: []string{"root"},
			ok:      true,
			err:     nil,
		},
		&driver{
			name:    "root",
			t:       Processor,
			prereqs: nil,
			ok:      true,
			err:     nil,
		},
		&driver{
			name:    "super",
			t:       Processor,
			prereqs: []string{"base1", "base2"},
			ok:      true,
			err:     nil,
		},
	}
	initTest(d)
	actual, err := explodeStages(d)
	if len(actual) != 3 || len(actual[0]) != 1 || len(actual[1]) != 2 || len(actual[2]) != 1 {
		t.Fatalf("%v", actual)
	}
	if err != nil {
		t.Fatalf("%v", err)
	}
}

//

func reset() {
	for i := range allDrivers {
		allDrivers[i] = nil
	}
	byName = map[string]Driver{}
	actualDrivers = nil
}

func initTest(drivers []Driver) {
	reset()
	for _, d := range drivers {
		Register(d)
	}
}

type driver struct {
	name    string
	t       Type
	prereqs []string
	ok      bool
	err     error
}

func (d *driver) String() string {
	return d.name
}

func (d *driver) Type() Type {
	return d.t
}

func (d *driver) Prerequisites() []string {
	return d.prereqs
}

func (d *driver) Init() (bool, error) {
	return d.ok, d.err
}
