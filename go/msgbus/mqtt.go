// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package msgbus

import (
	"errors"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// NewMQTT returns an initialized active MQTT connection.
//
// The connection timeouts are fine tuned for a LAN.
//
// This main purpose of this library is to hide the horror that
// paho.mqtt.golang is.
func NewMQTT(server, clientID, user, password string) (Bus, error) {
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
	return &mqttBus{client: client}, nil
	/*
		// Subscribe to all messages and filter locally. This causes a huge amount
		// of unnecessary traffic since it effectively acts as a local broker.
		m := &mqttBus{client: client}
		token := m.client.Subscribe("#", byte(BestEffort), func(client mqtt.Client, msgQ mqtt.Message) {
			msg := Message{msgQ.Topic(), msgQ.Payload()}
			m.mu.Lock()
			defer m.mu.Unlock()
			for i := range m.subscribers {
				if m.subscribers[i].topic_query.match(msg.Topic) {
					m.subscribers[i].publish(msg)
				}
			}
		})
		token.Wait()
		if err := token.Error(); err != nil {
			return nil, err
		}
		return m, nil
	*/
}

//

// mqttBus main purpose is to hide the horror that paho.mqtt.golang is.
//
// This Bus is thread safe.
type mqttBus struct {
	client mqtt.Client

	// For local brokerage:
	//mu          sync.Mutex
	//subscribers []*subscription
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

func (m *mqttBus) Subscribe(topic_query string, qos QOS) (<-chan Message, error) {
	// Quick local check.
	p := parseTopic(topic_query)
	if p == nil {
		return nil, errors.New("invalid topic")
	}

	c := make(chan Message)
	token := m.client.Subscribe(topic_query, byte(qos), func(client mqtt.Client, msg mqtt.Message) {
		c <- Message{msg.Topic(), msg.Payload()}
	})
	token.Wait()
	return c, token.Error()
	/*
		c := make(chan Message)
		m.mu.Lock()
		defer m.mu.Unlock()
		m.subscribers = append(m.subscribers, &subscription{topic_query: p, channel: c})
		return c, nil
	*/
}

func (m *mqttBus) Unsubscribe(topic_query string) {
	// Quick local check.
	p := parseTopic(topic_query)
	if p == nil {
		// Invalid topic.
		return
	}

	token := m.client.Unsubscribe(topic_query)
	token.Wait()
	// token.Error() is lost.
	/*
		m.mu.Lock()
		defer m.mu.Unlock()
		for i := range m.subscribers {
			if m.subscribers[i].topic_query.isEqual(p) {
				m.subscribers[i].closeSub()
				copy(m.subscribers[i:], m.subscribers[i+1:])
				m.subscribers = m.subscribers[:len(m.subscribers)-1]
				return
			}
		}
	*/
}

func (m *mqttBus) Retained(topic_query string) ([]Message, error) {
	// Quick local check.
	p := parseTopic(topic_query)
	if p == nil {
		return nil, errors.New("invalid topic")
	}

	// TODO(maruel): It looks it needs to do a quick Subscribe + poll every
	// messages until one with !msg.Retained() or a timeout then Unsubscribe.
	return nil, errors.New("implement me")
}

var _ Bus = &mqttBus{}
