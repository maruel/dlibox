// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package msgbus

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestQOS_String(t *testing.T) {
	data := []struct {
		v        QOS
		expected string
	}{
		{BestEffort, "BestEffort"},
		{QOS(-1), "QOS(-1)"},
	}
	for _, line := range data {
		if actual := line.v.String(); actual != line.expected {
			t.Fatalf("%q != %q", actual, line.expected)
		}
	}
}

func TestLog(t *testing.T) {
	if !testing.Verbose() {
		log.SetOutput(ioutil.Discard)
		defer log.SetOutput(os.Stderr)
	}
	b := Log(New())
	c, err := b.Subscribe("foo", BestEffort)
	if err != nil {
		t.Fatal(err)
	}
	if err := b.Publish(Message{Topic: "foo", Payload: make([]byte, 1)}, BestEffort); err != nil {
		t.Fatal(err)
	}
	if i := <-c; i.Topic != "foo" {
		t.Fatalf("%s != foo", i.Topic)
	}
	b.Unsubscribe("foo")
	if err := b.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestLog_Subscribe(t *testing.T) {
	if !testing.Verbose() {
		log.SetOutput(ioutil.Discard)
		defer log.SetOutput(os.Stderr)
	}
	b := Log(New())
	if _, err := b.Subscribe("", BestEffort); err == nil {
		t.Fatal("bad topic")
	}
	if err := b.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestRebasePub(t *testing.T) {
	b := RebasePub(New(), "foo")
	c, err := b.Subscribe("foo/bar", BestEffort)
	if err != nil {
		t.Fatal(err)
	}
	if err := b.Publish(Message{Topic: "bar", Payload: []byte("yo"), Retained: true}, BestEffort); err != nil {
		t.Fatal(err)
	}
	if i := <-c; i.Topic != "foo/bar" {
		t.Fatalf("%s != foo/bar", i.Topic)
	}
	b.Unsubscribe("foo")
	expected := map[string][]byte{"foo/bar": []byte("yo")}
	if l, err := Retained(b, time.Second, "foo/bar"); err != nil || !reflect.DeepEqual(l, expected) {
		t.Fatal(l, err)
	}
	if s := b.(fmt.Stringer).String(); s != "LocalBus/foo/" {
		t.Fatal(s)
	}
	if err := b.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestRebasePub_err(t *testing.T) {
	if !testing.Verbose() {
		log.SetOutput(ioutil.Discard)
		defer log.SetOutput(os.Stderr)
	}
	if RebasePub(New(), "a\000") != nil {
		t.Fatal("bad topic")
	}
	if RebasePub(New(), "#") != nil {
		t.Fatal("can't use a query")
	}
}

func TestRebaseSub(t *testing.T) {
	b := RebaseSub(New(), "foo")
	c, err := b.Subscribe("bar", BestEffort)
	if err != nil {
		t.Fatal(err)
	}
	if err := b.Publish(Message{Topic: "foo/bar", Payload: []byte("yo"), Retained: true}, BestEffort); err != nil {
		t.Fatal(err)
	}
	if i := <-c; i.Topic != "bar" {
		t.Fatalf("%s != bar", i.Topic)
	}
	b.Unsubscribe("bar")
	expected := map[string][]byte{"bar": []byte("yo")}
	if l, err := Retained(b, time.Second, "bar"); err != nil || !reflect.DeepEqual(l, expected) {
		t.Fatal(l, err)
	}
	if s := b.(fmt.Stringer).String(); s != "LocalBus/foo/" {
		t.Fatal(s)
	}
	if err := b.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestRebaseSub_err(t *testing.T) {
	if !testing.Verbose() {
		log.SetOutput(ioutil.Discard)
		defer log.SetOutput(os.Stderr)
	}
	if RebaseSub(New(), "a\000") != nil {
		t.Fatal("bad topic")
	}
	if RebaseSub(New(), "#") != nil {
		t.Fatal("can't use a query")
	}
}

func TestRebaseSub_root(t *testing.T) {
	b := RebaseSub(New(), "foo")
	c, err := b.Subscribe("//bar", BestEffort)
	if err != nil {
		t.Fatal(err)
	}
	if err := b.Publish(Message{Topic: "bar", Payload: []byte("yo"), Retained: true}, BestEffort); err != nil {
		t.Fatal(err)
	}
	if i := <-c; i.Topic != "bar" {
		t.Fatalf("%s != bar", i.Topic)
	}
	b.Unsubscribe("//bar")
	expected := map[string][]byte{"bar": []byte("yo")}
	if l, err := Retained(b, time.Second, "//bar"); err != nil || !reflect.DeepEqual(l, expected) {
		t.Fatal(l, err)
	}
	if err := b.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestRebaseSub_Err(t *testing.T) {
	b := RebaseSub(New(), "foo")
	if _, err := b.Subscribe("#/a", BestEffort); err == nil {
		t.Fatal("bad topic")
	}
	if _, err := Retained(b, time.Second, "#/a"); err == nil {
		t.Fatal("bad topic")
	}
}

func TestRetained(t *testing.T) {
	if _, err := Retained(nil, 0, "#"); err == nil {
		t.Fatal("can't use query")
	}
	if _, err := Retained(nil, 0, "a", "a"); err == nil {
		t.Fatal("can't use same topic twice")
	}
}

//

func TestParseTopicGood(t *testing.T) {
	data := []string{
		"$",
		"/",
		"//",
		"+/+/+",
		"+/+/#",
		"sport/tennis/+",
		"sport/tennis/#",
		"sport/tennis/player1",
		"sport/tennis/player1/ranking",
		"sport/tennis/+/score/wimbledon",
		strings.Repeat("a", 65535),
	}
	for i, line := range data {
		if _, err := parseTopic(line); err != nil {
			t.Fatalf("%d: parseTopic(%#v) returned %v", i, line, err)
		}
	}
}

func TestParseTopicBad(t *testing.T) {
	data := []string{
		"",
		"sport/tennis#",
		"sport/tennis/#/ranking",
		"sport/tennis+",
		"sport/#/tennis",
		strings.Repeat("a", 65536),
	}
	for i, line := range data {
		if _, err := parseTopic(line); err == nil {
			t.Fatalf("%d: parseTopic(%#v) returned non nil", i, line)
		}
	}
}

func TestMatchSuccess(t *testing.T) {
	data := [][2]string{
		{"sport/tennis/#", "sport/tennis"},
		{"sport/tennis/#", "sport/tennis/player1"},
		{"sport/tennis/#", "sport/tennis/player1/ranking"},
		{"sport/tennis/+", "sport/tennis"},
		{"sport/tennis/+", "sport/tennis/player1"},
		{"sport/+/player1", "sport/tennis/player1"},
	}
	for i, line := range data {
		q, err := parseTopic(line[0])
		if err != nil {
			t.Fatalf("%d: %#v.match(%#v): %v", i, line[0], line[1], err)
		}
		if !q.match(line[1]) {
			t.Fatalf("%d: %#v.match(%#v) returned false", i, line[0], line[1])
		}
	}
}

func TestMatchFail(t *testing.T) {
	data := [][2]string{
		{"sport/tennis/#", ""},
		{"sport/tennis/#", "sport"},
		{"sport/tennis/#", "sport/badminton"},
		{"sport/tennis/+", "sport"},
		{"sport/tennis/+", "sport/badminton"},
		{"sport/+/player1", "sport/tennis/player2"},
		{"sport/tennis", "sport/tennis/ball"},
		{"sport/tennis/ball", "sport/tennis"},
		{"sport/tennis/#", "sport/tennis/$player1"}, // 4.7.2
		{"+/tennis/", "sport/tennis/$player1"},      // 4.7.2
	}
	for i, line := range data {
		q, err := parseTopic(line[0])
		if err != nil {
			t.Fatalf("%d: %#v.match(%#v): %v", i, line[0], line[1], err)
		}
		if q.match(line[1]) {
			t.Fatalf("%d: %#v.match(%#v) returned true", i, line[0], line[1])
		}
	}
}
