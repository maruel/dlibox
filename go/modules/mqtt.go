// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package modules

import (
	"errors"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// MQTT main purpose is to hide the horror that paho.mqtt.golang is.
//
// It intentionally has a much simpler surface.
type MQTT struct {
	client mqtt.Client
}

func Make(server, clientID, user, password string) (*MQTT, error) {
	opts := mqtt.NewClientOptions().AddBroker(server).SetClientID(clientID)
	// Use lower timeouts than the defaults since they are high and the current
	// assumption is local network.
	opts.SetConnectTimeout(5 * time.Second)
	opts.SetKeepAlive(4 * time.Second)
	opts.SetPingTimeout(2 * time.Second)
	if len(user) != 0 {
		opts.SetUsername(user)
	}
	if len(password) != 0 {
		opts.SetPassword(password)
	}
	// TODO(maruel): opts.SetTLSConfig()
	// https://github.com/eclipse/paho.mqtt.golang/blob/master/samples/ssl.go
	// TODO(maruel): opts.SetBinaryWill()
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}
	return &MQTT{client}, nil
}

func (m *MQTT) Close() error {
	m.client.Disconnect(500)
	m.client = nil
	return nil
}

func (m *MQTT) Publish(msg *Message, qos QOS, retained bool) error {
	// Make it back synchronous.
	token := m.client.Publish(msg.Topic, byte(qos), retained, msg.Payload)
	token.Wait()
	return token.Error()
}

func (m *MQTT) Subscribe(topic string, qos QOS) (<-chan *Message, error) {
	c := make(chan *Message)
	token := m.client.Subscribe(topic, byte(qos), func(client mqtt.Client, msg mqtt.Message) {
		c <- &Message{msg.Topic(), msg.Payload()}
	})
	return c, token.Error()
}

func (m *MQTT) Unsubscribe(topic string) error {
	token := m.client.Unsubscribe(topic)
	token.Wait()
	return token.Error()
}

func (m *MQTT) Get(topic string, qos QOS) ([]*Message, error) {
	// TODO(maruel): It looks it needs to do a quick Subscribe + poll every
	// messages until one with !msg.Retained() or a timeout then Unsubscribe.
	return nil, errors.New("implement me")
}

var _ Bus = &MQTT{}
