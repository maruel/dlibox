// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package msgbus

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	mqtt "github.com/maruel/paho.mqtt.golang"
)

// NewMQTT returns an initialized active MQTT connection.
//
// The connection timeouts are fine tuned for a LAN.
//
// This main purpose of this library is to hide the horror that
// paho.mqtt.golang is.
func NewMQTT(server, clientID, user, password string, will Message) (Bus, error) {
	opts := mqtt.NewClientOptions().AddBroker(server)
	opts.ClientID = clientID
	// Default 10min is too slow.
	opts.MaxReconnectInterval = 30 * time.Second
	// Global ordering flag.
	// opts.Order = false
	opts.Username = user
	opts.Password = password
	if len(will.Topic) != 0 {
		opts.SetBinaryWill(will.Topic, will.Payload, byte(ExactlyOnce), true)
	}
	m := &mqttBus{server: server}
	opts.OnConnect = m.onConnect
	opts.OnConnectionLost = m.onConnectionLost
	opts.DefaultPublishHandler = m.unexpectedMessage
	m.client = mqtt.NewClient(opts)
	token := m.client.Connect()
	token.Wait()
	if err := token.Error(); err != nil {
		return nil, err
	}
	return m, nil
}

//

// mqttBus main purpose is to hide the horror that paho.mqtt.golang is.
//
// This Bus is thread safe.
type mqttBus struct {
	client mqtt.Client
	server string

	mu               sync.Mutex
	disconnectedOnce bool
}

func (m *mqttBus) String() string {
	return fmt.Sprintf("MQTT{%s}", m.server)
}

func (m *mqttBus) Close() error {
	m.client.Disconnect(500)
	m.client = nil
	return nil
}

func (m *mqttBus) Publish(msg Message, qos QOS, retained bool) error {
	// Quick local check.
	p := parseTopic(msg.Topic)
	if p == nil || p.isQuery() {
		return errors.New("invalid topic")
	}
	token := m.client.Publish(msg.Topic, byte(qos), retained, msg.Payload)
	if qos > BestEffort {
		token.Wait()
	}
	return token.Error()
}

func (m *mqttBus) Subscribe(topicQuery string, qos QOS) (<-chan Message, error) {
	// Quick local check.
	p := parseTopic(topicQuery)
	if p == nil {
		return nil, errors.New("invalid topic")
	}

	c := make(chan Message)
	token := m.client.Subscribe(topicQuery, byte(qos), func(client mqtt.Client, msg mqtt.Message) {
		c <- Message{msg.Topic(), msg.Payload()}
	})
	token.Wait()
	return c, token.Error()
}

func (m *mqttBus) Unsubscribe(topicQuery string) {
	// Quick local check.
	p := parseTopic(topicQuery)
	if p == nil {
		log.Printf("%s.Unsubscribe(%s): invalid topic", m, topicQuery)
		return
	}

	token := m.client.Unsubscribe(topicQuery)
	token.Wait()
	if err := token.Error(); err != nil {
		log.Printf("%s.Unsubscribe(%s): %v", m, topicQuery, err)
	}
}

func (m *mqttBus) Retained(topicQuery string) ([]Message, error) {
	// Quick local check.
	p := parseTopic(topicQuery)
	if p == nil {
		return nil, errors.New("invalid topic")
	}

	// Do a quick Subscribe(), retrieve all retained messages until one with
	// !msg.Retained() or a timeout then Unsubscribe.
	c := make(chan Message)
	token := m.client.Subscribe(topicQuery, byte(ExactlyOnce), func(client mqtt.Client, msg mqtt.Message) {
		// TODO(maruel): This assumes that retained messages are sent first by the
		// broker. This is likely not true.
		if msg.Retained() {
			c <- Message{msg.Topic(), msg.Payload()}
		}
	})
	if err := token.Error(); err != nil {
		// TODO(maruel): This will leak the channel.
		return nil, err
	}
	var out []Message
	// TODO(maruel): This is crappy.
	for loop := true; loop; {
		after := time.After(1 * time.Second)
		select {
		case i := <-c:
			out = append(out, i)
		case <-after:
			loop = false
		}
	}
	m.Unsubscribe(topicQuery)
	return out, nil
}

func (m *mqttBus) unexpectedMessage(c mqtt.Client, msg mqtt.Message) {
	log.Printf("%s: Unexpected message %s", m, msg.Topic())
}

func (m *mqttBus) onConnect(c mqtt.Client) {
	m.mu.Lock()
	d := m.disconnectedOnce
	m.mu.Unlock()
	if d {
		log.Printf("%s: connected", m)
	}
}

func (m *mqttBus) onConnectionLost(c mqtt.Client, err error) {
	log.Printf("%s: connection lost: %v", m, err)
	m.mu.Lock()
	m.disconnectedOnce = true
	m.mu.Unlock()
}

var _ Bus = &mqttBus{}
