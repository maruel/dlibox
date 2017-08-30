// Copyright 2017 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package msgbus

import (
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestNew_Publish_ephemeral_sync(t *testing.T) {
	b := New()

	if err := b.Publish(Message{Topic: "foo", Payload: []byte("a")}, ExactlyOnce); err != nil {
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
		if err := b.Publish(Message{Topic: "foo", Payload: []byte("b")}, ExactlyOnce); err != nil {
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
	b := New()
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

func TestNew_Subscribe_Close(t *testing.T) {
	b := New()

	if _, err := b.Subscribe("foo", ExactlyOnce); err != nil {
		t.Fatal(err)
	}
	if err := b.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestNew_Publish_Retained(t *testing.T) {
	// Without subscription.
	b := New()
	if err := b.Publish(Message{Topic: "foo", Payload: []byte("yo"), Retained: true}, ExactlyOnce); err != nil {
		t.Fatal(err)
	}
	expected := map[string][]byte{"foo": []byte("yo")}
	if l, err := Retained(b, time.Second, "foo"); err != nil || !reflect.DeepEqual(l, expected) {
		t.Fatal(l, err)
	}
	// Deleted retained message.
	if err := b.Publish(Message{Topic: "foo", Retained: true}, ExactlyOnce); err != nil {
		t.Fatal(err)
	}
	// Just enough time to "wait" but not enough to make this test too slow.
	if l, err := Retained(b, 10*time.Millisecond, "foo"); err != nil || len(l) != 0 {
		t.Fatal(l, err)
	}

	if err := b.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestNew_UnSubscribe_Cap(t *testing.T) {
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
	if !testing.Verbose() {
		log.SetOutput(ioutil.Discard)
		defer log.SetOutput(os.Stderr)
	}
	b := New()

	if b.Publish(Message{Payload: []byte("yo")}, ExactlyOnce) == nil {
		t.Fatal("bad topic")
	}
	if b.Publish(Message{Topic: "#", Payload: make([]byte, 1)}, ExactlyOnce) == nil {
		t.Fatal("topic is query")
	}

	if _, err := b.Subscribe("", ExactlyOnce); err == nil {
		t.Fatal("bad topic")
	}

	b.Unsubscribe("")

	if l, err := Retained(b, time.Second, ""); err == nil || len(l) != 0 {
		t.Fatal("bad topic")
	}

	if err := b.Close(); err != nil {
		t.Fatal(err)
	}
}

//

func TestSubscription(t *testing.T) {
	s := subscription{}
	// Check for s.channel == nil
	s.publish(Message{})

	// First check for s.channel == nil
	s.closeSub()

	// Second check for s.channel == nil is in local_closesub_test.go.
}

func TestNew_dump(t *testing.T) {
	b := New().(*local)
	if err := b.Publish(Message{Topic: "foo", Payload: []byte("yo"), Retained: true}, ExactlyOnce); err != nil {
		t.Fatal(err)
	}
	if _, err := b.Subscribe("bar", ExactlyOnce); err != nil {
		t.Fatal(err)
	}
	if s := b.dump(); s != "Persistent topics:\n- foo: yo\nSubscriptions:\n- bar\n" {
		t.Fatalf("%q", s)
	}
	if err := b.Close(); err != nil {
		t.Fatal(err)
	}
}
