// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package msgbus

import (
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
	"time"
	"unicode/utf8"
)

// QOS defines the quality of service to use when publishing and subscribing to
// messages.
//
// The normative definition is
// http://docs.oasis-open.org/mqtt/mqtt/v3.1.1/errata01/os/mqtt-v3.1.1-errata01-os-complete.html#_Toc442180912
type QOS int8

const (
	// BestEffort means the broker/client will deliver the message at most once,
	// with no confirmation.
	BestEffort QOS = 0
	// MinOnce means the broker/client will deliver the message at least once,
	// potentially duplicate.
	//
	// Do not use if message duplication is problematic.
	MinOnce QOS = 1
	// ExactlyOnce means the broker/client will deliver the message exactly once
	// by using a four step handshake.
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
	// Retained signifies that the message is permanent until explicitly changed.
	// Otherwise it is ephemeral.
	Retained bool
}

// Bus is a publisher-subscriber bus.
//
// The topics are expected to use the MQTT definition. "Mosquitto" has good
// documentation about this: https://mosquitto.org/man/mqtt-7.html
//
// For more information about retained message behavior, see
// http://www.hivemq.com/blog/mqtt-essentials-part-8-retained-messages
//
// Implementation of Bus are expected to implement fmt.Stringer.
type Bus interface {
	io.Closer

	// Publish publishes a message to a topic.
	//
	// If msg.Payload is empty, the topic is deleted if it was retained.
	//
	// It is not guaranteed that messages are propagated in order, unless
	// qos ExactlyOnce is used.
	Publish(msg Message, qos QOS) error

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
// Returns nil if root is an invalid topic or if it is a topic query.
//
// It is possible to publish a message topic outside of root with:
//  - "../" to backtrack closer to root
//  - "//" to ignore the root
func RebasePub(b Bus, root string) Bus {
	if len(root) != 0 && root[len(root)-1] != '/' {
		root += "/"
	}
	t := root[:len(root)-1]
	p, err := parseTopic(t)
	if err != nil {
		log.Printf("RebasePub(%s, %q): %v", b, t, err)
		return nil
	}
	if p.isQuery() {
		log.Printf("RebasePub(%s, %q): cannot use topic query", b, t)
		return nil
	}
	return &rebasePublisher{b, root}
}

// RebaseSub rebases a Bus when subscribing or getting topics.
//
// All Message retrieved have their Topic prefix root stripped.
//
// Messages published are unaffected.
//
// Returns nil if root is an invalid topic or if it is a topic query.
//
// It is possible to subscribe to a message topic outside of root with:
//  - "../" to backtrack closer to root
//  - "//" to ignore the root
func RebaseSub(b Bus, root string) Bus {
	if len(root) != 0 && root[len(root)-1] != '/' {
		root += "/"
	}
	t := root[:len(root)-1]
	p, err := parseTopic(t)
	if err != nil {
		log.Printf("RebaseSub(%s, %q): %v", b, t, err)
		return nil
	}
	if p.isQuery() {
		log.Printf("RebaseSub(%s, %q): cannot use topic query", b, t)
		return nil
	}
	return &rebaseSubscriber{b, root}
}

// Retained retrieves all matching messages for one or multiple topics.
//
// Topic queries cannot be used.
//
// If a topic is missing, will wait for up to d for it to become available. If
// all topics are available, returns as soon as they are all retrieved.
func Retained(b Bus, d time.Duration, topic ...string) (map[string][]byte, error) {
	// Quick local check.
	var ps []parsedTopic
	for i, t := range topic {
		p, err := parseTopic(t)
		if err != nil {
			return nil, fmt.Errorf("invalid topic %q: %v", t, err)
		}
		if p.isQuery() {
			return nil, fmt.Errorf("cannot use topic query %q", t)
		}
		for j := 0; j < i; j++ {
			if topic[j] == topic[i] {
				return nil, fmt.Errorf("cannot specify topic %q twice", t)
			}
		}
		ps = append(ps, p)
	}

	// Subscribes to all topics concurrently. This reduces the effect of round
	// trip latency.
	type result struct {
		c   <-chan Message
		err error
	}
	channels := make(chan result, len(topic))
	for _, t := range topic {
		go func(t string) {
			c, err := b.Subscribe(t, MinOnce)
			channels <- result{c, err}
		}(t)
	}

	// Ensure all topics are unsubscribed even in case of error. This also
	// ensures the channels are closed, so the goroutine started below do not
	// leak.
	defer func() {
		for _, t := range topic {
			b.Unsubscribe(t)
		}
	}()

	// Look at all topic subscription to ensures they all succeeded.
	master := make(chan Message)
	for range topic {
		r := <-channels
		if r.err != nil {
			return nil, r.err
		}
		go func(c <-chan Message) {
			v, ok := <-c
			if !ok || !v.Retained {
				return
			}
			master <- v
		}(r.c)
	}

	// Retrieve results.
	out := map[string][]byte{}
	for loop := true; loop && len(out) < len(topic); {
		// Reset the timer after every message retrieved.
		a := time.After(d)
		select {
		case v := <-master:
			out[v.Topic] = v.Payload
		case <-a:
			loop = false
		}
	}
	return out, nil
}

// Private code.

type logging struct {
	Bus
}

func (l *logging) Close() error {
	log.Printf("%s.Close()", l)
	return l.Bus.Close()
}

func (l *logging) Publish(msg Message, qos QOS) error {
	log.Printf("%s.Publish({%s, %q, %t}, %s)", l, msg.Topic, string(msg.Payload), msg.Retained, qos)
	return l.Bus.Publish(msg, qos)
}

func (l *logging) Subscribe(topicQuery string, qos QOS) (<-chan Message, error) {
	log.Printf("%s.Subscribe(%s, %s)", l, topicQuery, qos)
	c, err := l.Bus.Subscribe(topicQuery, qos)
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
	l.Bus.Unsubscribe(topicQuery)
}

// Rebase support.

type rebasePublisher struct {
	Bus
	root string
}

func (r *rebasePublisher) String() string {
	return fmt.Sprintf("%s/%s", r.Bus, r.root)
}

func (r *rebasePublisher) Publish(msg Message, qos QOS) error {
	msg.Topic = mergeTopic(r.root, msg.Topic)
	return r.Bus.Publish(msg, qos)
}

type rebaseSubscriber struct {
	Bus
	root string
}

func (r *rebaseSubscriber) String() string {
	return fmt.Sprintf("%s/%s", r.Bus, r.root)
}

func (r *rebaseSubscriber) Subscribe(topicQuery string, qos QOS) (<-chan Message, error) {
	if strings.HasPrefix(topicQuery, "//") {
		return r.Bus.Subscribe(topicQuery[2:], qos)
	}
	// BUG: Support mergeTopic().
	actual := r.root + topicQuery
	p, err := parseTopic(actual)
	if err != nil {
		return nil, err
	}
	c, err := r.Bus.Subscribe(actual, qos)
	if err != nil {
		return nil, err
	}
	c2 := make(chan Message)
	offset := len(r.root)
	go func() {
		defer close(c2)
		// Translate the topics.
		for msg := range c {
			if !p.isQuery() && !p.match(msg.Topic) {
				panic(fmt.Errorf("bus: unexpected topic prefix %q, expected %q", msg.Topic, actual))
			}
			c2 <- Message{Topic: msg.Topic[offset:], Payload: msg.Payload, Retained: msg.Retained}
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

func parseTopic(topic string) (parsedTopic, error) {
	if len(topic) == 0 {
		return nil, errors.New("empty topic")
	}
	if len(topic) > 65535 {
		return nil, fmt.Errorf("topic length %d over 65535 characters", len(topic))
	}
	if strings.ContainsRune(topic, rune(0)) || !utf8.ValidString(topic) {
		return nil, errors.New("topic must be valid UTF-8")
	}
	p := parsedTopic(strings.Split(topic, "/"))
	if err := p.Validate(); err != nil {
		return nil, err
	}
	return p, nil
}

func (p parsedTopic) String() string {
	return strings.Join(p, "/")
}

func (p parsedTopic) Validate() error {
	for i, e := range p {
		if i != len(p)-1 && e == "#" {
			return errors.New("wildcard # can only appear at the end of a topic query")
		} else if e != "+" && e != "#" {
			if strings.HasSuffix(e, "#") {
				return errors.New("wildcard # can not appear inside a topic section")
			}
			if strings.HasSuffix(e, "+") {
				return errors.New("wildcard + can not appear inside a topic section")
			}
		}
	}
	return nil
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
