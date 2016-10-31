// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package modules

import (
	"io"
	"log"
	"strings"
)

// QOS defines the quality of service to use when publishing and subscribing to
// messages.
type QOS int

const (
	// The broker/client will deliver the message once, with no confirmation.
	BestEffort QOS = 0
	// The broker/client will deliver the message at least once, with
	// confirmation required.
	MinOnce QOS = 1
	// The broker/client will deliver the message exactly once by using a four
	// step handshake.
	ExactlyOnce QOS = 2
)

// Message represents a single message to a single topic.
type Message struct {
	Topic   string
	Payload []byte
}

// Bus represents a minimal publisher-subscriber bus.
//
// The topics are expected to use the MQTT definition. Mosquitto has a nice doc
// about this: https://mosquitto.org/man/mqtt-7.html
//
// There is a pure local-only implementation and another using a real MQTT
// server.
type Bus interface {
	io.Closer

	// Publish publishes a message to a topic.
	//
	// If msg.Payload is empty, the topic is deleted if it was retained.
	Publish(msg Message, qos QOS, retained bool) error

	// Subscribe sends back updates to this topic query.
	//
	// The format must be according to MQTT.
	Subscribe(topic string, qos QOS) (<-chan Message, error)
	// Unsubscribe removes a previous subscription.
	//
	// It is an error to unsubscribe from a non-subscribed topic.
	Unsubscribe(topic string) error
	// Get retrieves matching messages for a retained topic query.
	Get(topic string, qos QOS) ([]Message, error)
}

// Logging logs all messages published on a Bus.
func Logging(b Bus) Bus {
	return &logging{b}
}

// Rebase rebases a Bus for all topics.
//
// User can then use:
//  - "../" to backtrack closer to root
//  - "//" to ignore the root
func Rebase(b Bus, root string) Bus {
	if len(root) != 0 && root[len(root)-1] != '/' {
		root += "/"
	}
	return &rebasePublisher{&rebaseSubscriber{b, root}, root}
}

// RebasePublisher rebases a Bus when publishing messages.
//
// Messages retrieved are unaffected.
func RebasePublisher(b Bus, root string) Bus {
	if len(root) != 0 && root[len(root)-1] != '/' {
		root += "/"
	}
	return &rebasePublisher{b, root}
}

// RebaseSubscriber rebases a Bus when subscribing or getting topics.
//
// Messages published are unaffected.
func RebaseSubscriber(b Bus, root string) Bus {
	if len(root) != 0 && root[len(root)-1] != '/' {
		root += "/"
	}
	return &rebaseSubscriber{b, root}
}

//

type logging struct {
	Bus
}

func (l *logging) Publish(msg Message, qos QOS, retained bool) error {
	log.Printf("Publish({%s, %q}, %s, %t)", msg.Topic, string(msg.Payload), qos, retained)
	return l.Bus.Publish(msg, qos, retained)
}

func (l *logging) Subscribe(topic string, qos QOS) (<-chan Message, error) {
	log.Printf("Subscribe(%s, %s)", topic, qos)
	c, err := l.Bus.Subscribe(topic, qos)
	if err != nil {
		return c, err
	}
	c2 := make(chan Message)
	go func() {
		defer close(c2)
		for msg := range c {
			log.Printf("<- Message{%s, %q}", msg.Topic, string(msg.Payload))
			c2 <- msg
		}
	}()
	return c2, nil
}

type rebasePublisher struct {
	Bus
	root string
}

func (r *rebasePublisher) Publish(msg Message, qos QOS, retained bool) error {
	msg.Topic = mergeTopic(r.root, msg.Topic)
	return r.Bus.Publish(msg, qos, retained)
}

type rebaseSubscriber struct {
	Bus
	root string
}

func (r *rebaseSubscriber) Subscribe(topic string, qos QOS) (<-chan Message, error) {
	if strings.HasPrefix(topic, "//") {
		return r.Bus.Subscribe(topic[2:], qos)
	}
	// TODO(maruel): Support mergeTopic().
	actual := r.root + topic
	c, err := r.Bus.Subscribe(actual, qos)
	p := parseTopic(actual)
	if err != nil {
		return c, err
	}
	c2 := make(chan Message)
	offset := len(r.root)
	go func() {
		defer close(c2)
		// Translate the topics.
		for msg := range c {
			if !p.match(msg.Topic) {
				// TODO(maruel): There's a bug when subscribing over MQTT.
				log.Printf("bus: unexpected topic prefix %q, expected %q", msg.Topic, actual)
			} else {
				c2 <- Message{msg.Topic[offset:], msg.Payload}
			}
		}
	}()
	return c2, nil
}

func (r *rebaseSubscriber) Unsubscribe(topic string) error {
	if strings.HasPrefix(topic, "//") {
		return r.Bus.Unsubscribe(topic[2:])
	}
	// TODO(maruel): Support mergeTopic().
	return r.Bus.Unsubscribe(r.root + topic)
}

func (r *rebaseSubscriber) Get(topic string, qos QOS) ([]Message, error) {
	// TODO(maruel): Support mergeTopic().
	msgs, err := r.Bus.Get(r.root+topic, qos)
	if err != nil {
		return msgs, err
	}
	offset := len(r.root)
	for i := range msgs {
		msgs[i].Topic = msgs[i].Topic[offset:]
	}
	return msgs, err
}

func mergeTopic(root, topic string) string {
	if strings.HasPrefix(topic, "//") {
		return topic[2:]
	}
	for strings.HasPrefix(topic, "../") {
		if len(root) == 0 {
			panic(topic)
		}
		i := strings.LastIndexByte(root, '/')
		root = root[:i]
		topic = topic[3:]
	}
	if len(topic) == 0 {
		panic(root)
	}
	if topic[0] == '/' {
		panic(root)
	}
	if len(root) != 0 && root[len(root)-1] != '/' {
		root += "/"
	}
	return root + topic
}
