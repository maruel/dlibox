// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package msgbus

import (
	"fmt"
	"io"
	"log"
	"strings"
	"unicode/utf8"
)

// QOS defines the quality of service to use when publishing and subscribing to
// messages.
type QOS int8

const (
	// BestEffort means the broker/client will deliver the message once, with no
	// confirmation.
	//
	// This enables asynchronous operation.
	BestEffort QOS = 0
	// MinOnce means the broker/client will deliver the message at least once,
	// with confirmation required.
	//
	// Do not use if message duplication is problematic.
	MinOnce QOS = 1
	// ExactlyOnce means the broker/client will deliver the message exactly once
	// by using a four step handshake.
	//
	// This enforces synchronous operation.
	ExactlyOnce QOS = 2
)

const qosName = "BestEffortMinOnceExactlyOnce"

var qosIndex = [...]uint8{0, 10, 17, 28}

func (i QOS) String() string {
	if i < 0 || i >= QOS(len(qosIndex)-1) {
		return fmt.Sprintf("QOS(%d)", i)
	}
	return qosName[qosIndex[i]:qosIndex[i+1]]
}

// Message represents a single message to a single topic.
type Message struct {
	// Topic is the MQTT topic. It may have a prefix stripped by RebaseSub() or
	// inserted by RebasePub().
	Topic string
	// Payload is the application specific data.
	//
	// Publishing a message with no Payload deleted a retained Topic, and has no
	// effect on non-retained topic.
	Payload []byte
}

// Bus is a publisher-subscriber bus.
//
// The topics are expected to use the MQTT definition. "Mosquitto" has good
// documentation about this: https://mosquitto.org/man/mqtt-7.html
//
// For more information about retained message behavior, see
// http://www.hivemq.com/blog/mqtt-essentials-part-8-retained-messages
type Bus interface {
	io.Closer
	fmt.Stringer

	// Publish publishes a message to a topic.
	//
	// If msg.Payload is empty, the topic is deleted if it was retained.
	//
	// It is not guaranteed that messages are propagated in order, unless
	// qos ExactlyOnce is used.
	Publish(msg Message, qos QOS, retained bool) error

	// Subscribe sends updates to this topic query through the returned channel.
	Subscribe(topicQuery string, qos QOS) (<-chan Message, error)

	// Unsubscribe removes a previous subscription.
	//
	// Trying to unsubscribe from an invalid topic or a topic not currently
	// subscribed is ignored.
	//
	// BUG: while Subscribe() can be called multiple times with a topic query, a
	// single Unsubscribe() call will unregister all subscriptions.
	Unsubscribe(topicQuery string)

	// Retained retrieves a copy of all matching messages for a retained topic
	// query.
	Retained(topicQuery string) ([]Message, error)
}

// Log returns a Bus that logs all operations done on it, via log standard
// package.
func Log(b Bus) Bus {
	return &logging{b}
}

// RebasePub rebases a Bus when publishing messages.
//
// All Message published have their Topic prefixed with root.
//
// Messages retrieved are unaffected.
//
// It is possible to publish a message topic outside of root with:
//  - "../" to backtrack closer to root
//  - "//" to ignore the root
func RebasePub(b Bus, root string) Bus {
	if len(root) != 0 && root[len(root)-1] != '/' {
		root += "/"
	}
	return &rebasePublisher{b, root}
}

// RebaseSub rebases a Bus when subscribing or getting topics.
//
// All Message retrieved have their Topic prefix root stripped.
//
// Messages published are unaffected.
//
// It is possible to subscribe to a message topic outside of root with:
//  - "../" to backtrack closer to root
//  - "//" to ignore the root
func RebaseSub(b Bus, root string) Bus {
	if len(root) != 0 && root[len(root)-1] != '/' {
		root += "/"
	}
	return &rebaseSubscriber{b, root}
}

// Private code.

type logging struct {
	bus Bus
}

func (l *logging) String() string {
	return l.bus.String()
}

func (l *logging) Close() error {
	log.Printf("%s.Close()", l)
	return l.bus.Close()
}

func (l *logging) Publish(msg Message, qos QOS, retained bool) error {
	log.Printf("%s.Publish({%s, %q}, %s, %t)", l, msg.Topic, string(msg.Payload), qos, retained)
	return l.bus.Publish(msg, qos, retained)
}

func (l *logging) Subscribe(topicQuery string, qos QOS) (<-chan Message, error) {
	log.Printf("%s.Subscribe(%s, %s)", l, topicQuery, qos)
	c, err := l.bus.Subscribe(topicQuery, qos)
	if err != nil {
		return c, err
	}
	c2 := make(chan Message)
	go func() {
		defer close(c2)
		for msg := range c {
			log.Printf("%s <- Message{%s, %q}", l, msg.Topic, string(msg.Payload))
			c2 <- msg
		}
	}()
	return c2, nil
}

func (l *logging) Unsubscribe(topicQuery string) {
	log.Printf("%s.Unsubscribe(%s)", l, topicQuery)
	l.bus.Unsubscribe(topicQuery)
}

func (l *logging) Retained(topicQuery string) ([]Message, error) {
	log.Printf("%s.Retained(%s)", l, topicQuery)
	return l.bus.Retained(topicQuery)
}

// Rebase support.

type rebasePublisher struct {
	Bus
	root string
}

func (r *rebasePublisher) String() string {
	return r.Bus.String() + "/" + r.root
}

func (r *rebasePublisher) Publish(msg Message, qos QOS, retained bool) error {
	msg.Topic = mergeTopic(r.root, msg.Topic)
	return r.Bus.Publish(msg, qos, retained)
}

type rebaseSubscriber struct {
	Bus
	root string
}

func (r *rebaseSubscriber) String() string {
	return r.Bus.String() + "/" + r.root
}

func (r *rebaseSubscriber) Subscribe(topicQuery string, qos QOS) (<-chan Message, error) {
	if strings.HasPrefix(topicQuery, "//") {
		return r.Bus.Subscribe(topicQuery[2:], qos)
	}
	// BUG: Support mergeTopic().
	actual := r.root + topicQuery
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
				panic(fmt.Errorf("bus: unexpected topic prefix %q, expected %q", msg.Topic, actual))
			}
			c2 <- Message{msg.Topic[offset:], msg.Payload}
		}
	}()
	return c2, nil
}

func (r *rebaseSubscriber) Unsubscribe(topicQuery string) {
	if strings.HasPrefix(topicQuery, "//") {
		r.Bus.Unsubscribe(topicQuery[2:])
		return
	}
	// BUG: Support mergeTopic().
	r.Bus.Unsubscribe(r.root + topicQuery)
}

func (r *rebaseSubscriber) Retained(topicQuery string) ([]Message, error) {
	// BUG: Support mergeTopic().
	msgs, err := r.Bus.Retained(r.root + topicQuery)
	if err != nil {
		return msgs, err
	}
	offset := len(r.root)
	for i := range msgs {
		msgs[i].Topic = msgs[i].Topic[offset:]
	}
	return msgs, err
}

// Topic parsing.

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

// parsedTopic is either a query or a static topic.
type parsedTopic []string

func parseTopic(topic string) parsedTopic {
	if len(topic) == 0 || len(topic) > 65535 || strings.ContainsRune(topic, rune(0)) || !utf8.ValidString(topic) {
		return nil
	}
	p := parsedTopic(strings.Split(topic, "/"))
	if !p.isValid() {
		return nil
	}
	return p
}

func (p parsedTopic) isValid() bool {
	// http://docs.oasis-open.org/mqtt/mqtt/v3.1.1/errata01/os/mqtt-v3.1.1-errata01-os-complete.html#_Toc442180921
	// section 4.7.2 about '$' prefix and section 4.7.3
	if len(p[0]) != 0 && p[0][0] == '$' {
		return false
	}
	for i, e := range p {
		// As per the spec, empty sections are valid.
		if i != len(p)-1 && e == "#" {
			// # can only appear at the end.
			return false
		} else if e != "+" && e != "#" {
			if strings.HasSuffix(e, "#") || strings.HasSuffix(e, "+") {
				return false
			}
		}
	}
	return true
}

func (p parsedTopic) isQuery() bool {
	for _, e := range p {
		if e == "#" || e == "+" {
			return true
		}
	}
	return false
}

func (p parsedTopic) isEqual(other parsedTopic) bool {
	if len(other) != len(p) {
		return false
	}
	for i := range p {
		if p[i] != other[i] {
			return false
		}
	}
	return true
}

// match follows rules as defined at section 4.7:
// http://docs.oasis-open.org/mqtt/mqtt/v3.1.1/errata01/os/mqtt-v3.1.1-errata01-os-complete.html#_Toc442180919
func (p parsedTopic) match(topic string) bool {
	if len(topic) == 0 {
		return false
	}
	t := strings.Split(topic, "/")
	// 4.7.2
	isPrivate := strings.HasPrefix(t[len(t)-1], "$")
	for i, e := range p {
		if e == "#" {
			return !isPrivate
		}
		if e == "+" {
			if isPrivate {
				return false
			}
			if i == len(p)-1 && len(t) == len(p) || len(t) == len(p)-1 {
				return true
			}
			continue
		}
		if len(t) <= i || t[i] != e {
			return false
		}
	}
	return len(t) == len(p)
}

var _ Bus = &local{}
