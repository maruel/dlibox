// Copyright 2017 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package msgbus

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"

	mqtt "github.com/maruel/paho.mqtt.golang"
)

func TestNewMQTT_fail(t *testing.T) {
	_, err := NewMQTT("", "client", "user", "pass", Message{Topic: "status", Payload: []byte("dead")})
	if err == nil {
		t.Fatal("invalid host")
	}
}

func TestMQTT(t *testing.T) {
	// Can't be t.Parallel() due to log.SetOutput().
	if !testing.Verbose() {
		log.SetOutput(ioutil.Discard)
		defer log.SetOutput(os.Stderr)
	}
	c := clientFake{}
	m := mqttBus{client: &c}

	if err := m.Publish(Message{Topic: "a/#/b"}, BestEffort, false); err == nil {
		t.Fatal("bad topic")
	}
	if _, err := m.Subscribe("a/#/b", BestEffort); err == nil {
		t.Fatal("bad topic")
	}
	m.Unsubscribe("a/#/b")
	if l, err := m.Retained("a/#/b"); err == nil || len(l) != 0 {
		t.Fatal("bad topic")
	}
	if err := m.Close(); err != nil {
		t.Fatal(err)
	}
}

//

type clientFake struct {
	err error
}

func (c *clientFake) IsConnected() bool {
	return true
}

func (c *clientFake) Connect() mqtt.Token {
	return &tokenFake{err: c.err}
}

func (c *clientFake) Disconnect(quiesce uint) {
}

func (c *clientFake) Publish(topic string, qos byte, retained bool, payload interface{}) mqtt.Token {
	return &tokenFake{err: c.err}
}

func (c *clientFake) Subscribe(topic string, qos byte, callback mqtt.MessageHandler) mqtt.Token {
	return &tokenFake{err: c.err}
}

func (c *clientFake) SubscribeMultiple(filters map[string]byte, callback mqtt.MessageHandler) mqtt.Token {
	return &tokenFake{err: c.err}
}

func (c *clientFake) Unsubscribe(topics ...string) mqtt.Token {
	return &tokenFake{err: c.err}
}

func (c *clientFake) AddRoute(topic string, callback mqtt.MessageHandler) {
}

type tokenFake struct {
	mqtt.UnsubscribeToken // to get flowComplete()
	err                   error
}

func (t *tokenFake) Wait() bool {
	return true
}
func (t *tokenFake) WaitTimeout(time.Duration) bool {
	return true
}

func (t *tokenFake) Error() error {
	return t.err
}
