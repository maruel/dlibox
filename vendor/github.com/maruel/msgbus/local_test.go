// Copyright 2017 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package msgbus

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestNew_Publish_ephemeral_sync(t *testing.T) {
	t.Parallel()
	b := New()

	if err := b.Publish(Message{Topic: "foo", Payload: []byte("a")}, ExactlyOnce, false); err != nil {
		t.Fatal(err)
	}

	c, err := b.Subscribe("foo", ExactlyOnce)
	if err != nil {
		t.Fatal(err)
	}

	select {
	case v := <-c:
		t.Fatal(v)
	default:
	}

	go func() {
		if err := b.Publish(Message{Topic: "foo", Payload: []byte("b")}, ExactlyOnce, false); err != nil {
			t.Fatal(err)
		}
	}()
	if v := <-c; v.Topic != "foo" || string(v.Payload) != "b" {
		t.Fatalf("%s != foo; %q != b", v.Topic, string(v.Payload))
	}
	b.Unsubscribe("foo")
	b.Unsubscribe("foo")

	if err := b.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestNew_Publish_ephemeral_async(t *testing.T) {
	t.Parallel()
	b := New()
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

func TestNew_Subscribe_Close(t *testing.T) {
	t.Parallel()
	b := New()

	if _, err := b.Subscribe("foo", ExactlyOnce); err != nil {
		t.Fatal(err)
	}
	if err := b.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestNew_Publish_Retained(t *testing.T) {
	t.Parallel()
	// Without subscription.
	b := New()
	if err := b.Publish(Message{Topic: "foo", Payload: make([]byte, 1)}, ExactlyOnce, true); err != nil {
		t.Fatal(err)
	}
	if l, err := b.Retained("foo"); err != nil || len(l) != 1 {
		t.Fatal(l, err)
	}
	// Deleted retained message.
	if err := b.Publish(Message{Topic: "foo"}, ExactlyOnce, true); err != nil {
		t.Fatal(err)
	}
	if l, err := b.Retained("foo"); err != nil || len(l) != 0 {
		t.Fatal(l, err)
	}

	if err := b.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestNew_UnSubscribe_Cap(t *testing.T) {
	t.Parallel()
	b := New()

	for i := 0; i < 17; i++ {
		if _, err := b.Subscribe("foo", ExactlyOnce); err != nil {
			t.Fatal(err)
		}
	}
	b.Unsubscribe("foo")

	if _, err := b.Subscribe("bar", ExactlyOnce); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 16; i++ {
		if _, err := b.Subscribe("foo", ExactlyOnce); err != nil {
			t.Fatal(err)
		}
	}
	b.Unsubscribe("foo")

	if err := b.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestNew_Err(t *testing.T) {
	// Can't be t.Parallel() due to log.SetOutput().
	if !testing.Verbose() {
		log.SetOutput(ioutil.Discard)
		defer log.SetOutput(os.Stderr)
	}
	b := New()

	if err := b.Publish(Message{Payload: make([]byte, 1)}, ExactlyOnce, false); err == nil {
		t.Fatal("bad topic")
	}
	if err := b.Publish(Message{Topic: "#", Payload: make([]byte, 1)}, ExactlyOnce, false); err == nil {
		t.Fatal("topic is query")
	}

	if _, err := b.Subscribe("", ExactlyOnce); err == nil {
		t.Fatal("bad topic")
	}

	b.Unsubscribe("")

	if l, err := b.Retained(""); err == nil || len(l) != 0 {
		t.Fatal("bad topic")
	}

	if err := b.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestNew_subscription(t *testing.T) {
	t.Parallel()
	s := subscription{}
	// Check for s.channel == nil
	s.publish(Message{})

	// First check for s.channel == nil
	s.closeSub()

	// Second check for s.channel == nil is in local_closesub_test.go.
}
