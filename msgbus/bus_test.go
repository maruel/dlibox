// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package msgbus

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestQOS_String(t *testing.T) {
	t.Parallel()
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
	// Can't be t.Parallel() due to log.SetOutput().
	if !testing.Verbose() {
		log.SetOutput(ioutil.Discard)
		defer log.SetOutput(os.Stderr)
	}
	b := Log(New())
	c, err := b.Subscribe("foo", BestEffort)
	if err != nil {
		t.Fatal(err)
	}
	if err := b.Publish(Message{Topic: "foo", Payload: make([]byte, 1)}, BestEffort, false); err != nil {
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
	// Can't be t.Parallel() due to log.SetOutput().
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
	if err := b.Publish(Message{Topic: "bar", Payload: make([]byte, 1)}, BestEffort, true); err != nil {
		t.Fatal(err)
	}
	if i := <-c; i.Topic != "foo/bar" {
		t.Fatalf("%s != foo/bar", i.Topic)
	}
	b.Unsubscribe("foo")
	if l, err := b.Retained("foo/bar"); err != nil || len(l) != 1 || l[0].Topic != "foo/bar" {
		t.Fatal(l, err)
	}
	if err := b.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestRebaseSub(t *testing.T) {
	b := RebaseSub(New(), "foo")
	c, err := b.Subscribe("bar", BestEffort)
	if err != nil {
		t.Fatal(err)
	}
	if err := b.Publish(Message{Topic: "foo/bar", Payload: make([]byte, 1)}, BestEffort, true); err != nil {
		t.Fatal(err)
	}
	if i := <-c; i.Topic != "bar" {
		t.Fatalf("%s != bar", i.Topic)
	}
	b.Unsubscribe("bar")
	if l, err := b.Retained("bar"); err != nil || len(l) != 1 || l[0].Topic != "bar" {
		t.Fatal(l, err)
	}
	if err := b.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestRebaseSub_root(t *testing.T) {
	b := RebaseSub(New(), "foo")
	c, err := b.Subscribe("//bar", BestEffort)
	if err != nil {
		t.Fatal(err)
	}
	if err := b.Publish(Message{Topic: "bar", Payload: make([]byte, 1)}, BestEffort, true); err != nil {
		t.Fatal(err)
	}
	if i := <-c; i.Topic != "bar" {
		t.Fatalf("%s != bar", i.Topic)
	}
	b.Unsubscribe("//bar")
	/*
		if l, err := b.Retained("//bar"); err != nil || len(l) != 1 || l[0].Topic != "bar" {
			t.Fatal(l, err)
		}
	*/
	if err := b.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestRebaseSub_Err(t *testing.T) {
	b := RebaseSub(New(), "foo")
	if _, err := b.Subscribe("#/a", BestEffort); err == nil {
		t.Fatal("bad topic")
	}
	if _, err := b.Retained("#/a"); err == nil {
		t.Fatal("bad topic")
	}
}

//

func TestParseTopicGood(t *testing.T) {
	t.Parallel()
	data := []string{
		"//",
		"+/+/+",
		"+/+/#",
		"sport/tennis/+",
		"sport/tennis/#",
		"sport/tennis/player1",
		"sport/tennis/player1/ranking",
		"sport/tennis/+/score/wimbledon",
	}
	for i, line := range data {
		if parseTopic(line) == nil {
			t.Fatalf("%d: parseTopic(%#v) returned nil", i, line)
		}
	}
}

func TestParseTopicBad(t *testing.T) {
	t.Parallel()
	data := []string{
		"",
		"$",
		"sport/tennis#",
		"sport/tennis/#/ranking",
		"sport/tennis+",
		"sport/#/tennis",
	}
	for i, line := range data {
		if parseTopic(line) != nil {
			t.Fatalf("%d: parseTopic(%#v) returned non nil", i, line)
		}
	}
}

func TestMatchSuccess(t *testing.T) {
	t.Parallel()
	data := [][2]string{
		{"sport/tennis/#", "sport/tennis"},
		{"sport/tennis/#", "sport/tennis/player1"},
		{"sport/tennis/#", "sport/tennis/player1/ranking"},
		{"sport/tennis/+", "sport/tennis"},
		{"sport/tennis/+", "sport/tennis/player1"},
		{"sport/+/player1", "sport/tennis/player1"},
	}
	for i, line := range data {
		q := parseTopic(line[0])
		if !q.match(line[1]) {
			t.Fatalf("%d: %#v.match(%#v) returned false", i, line[0], line[1])
		}
	}
}

func TestMatchFail(t *testing.T) {
	t.Parallel()
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
		q := parseTopic(line[0])
		if q.match(line[1]) {
			t.Fatalf("%d: %#v.match(%#v) returned true", i, line[0], line[1])
		}
	}
}
