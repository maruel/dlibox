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
// The connection timeouts are fine tuned for a LAN. It will likely fail on a
// slower connection or when used over the internet.
//
// will is the message to send if the connection is not closed correctly; when
// Close() is not called.
//
// order determines is messages are processed in order or not. Out of order
// processing means that a subscription will not be blocked by another one that
// fails to process its queue in time.
//
// This main purpose of this library is to create a layer that is simpler, more
// usable and more Go-idiomatic than paho.mqtt.golang.
func NewMQTT(server, clientID, user, password string, will Message, order bool) (Bus, error) {
	opts := mqtt.NewClientOptions().AddBroker(server)
	opts.ClientID = clientID
	// Default 10min is too slow.
	opts.MaxReconnectInterval = 30 * time.Second
	opts.Order = order
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

// mqttBus main purpose is to hide the complex thing that paho.mqtt.golang is.
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

// Close gracefully closes the connection to the server.
//
// Waits 1s for the connection to terminate correctly. If this function is not
// called, the will message in NewMQTT() will be activated.
func (m *mqttBus) Close() error {
	m.client.Disconnect(1000)
	m.client = nil
	return nil
}

func (m *mqttBus) Publish(msg Message, qos QOS) error {
	// Quick local check.
	p, err := parseTopic(msg.Topic)
	if err != nil {
		return err
	}
	if p.isQuery() {
		return errors.New("cannot publish to a topic query")
	}
	token := m.client.Publish(msg.Topic, byte(qos), msg.Retained, msg.Payload)
	if qos > BestEffort {
		token.Wait()
	}
	return token.Error()
}

func (m *mqttBus) Subscribe(topicQuery string, qos QOS) (<-chan Message, error) {
	// Quick local check.
	if _, err := parseTopic(topicQuery); err != nil {
		return nil, err
	}

	c := make(chan Message)
	token := m.client.Subscribe(topicQuery, byte(qos), func(client mqtt.Client, msg mqtt.Message) {
		c <- Message{Topic: msg.Topic(), Payload: msg.Payload(), Retained: msg.Retained()}
	})
	token.Wait()
	return c, token.Error()
}

func (m *mqttBus) Unsubscribe(topicQuery string) {
	// Quick local check.
	if _, err := parseTopic(topicQuery); err != nil {
		log.Printf("%s.Unsubscribe(%s): %v", m, topicQuery, err)
		return
	}

	token := m.client.Unsubscribe(topicQuery)
	token.Wait()
	if err := token.Error(); err != nil {
		log.Printf("%s.Unsubscribe(%s): %v", m, topicQuery, err)
	}
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
