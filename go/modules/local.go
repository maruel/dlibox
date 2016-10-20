// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package modules

import (
	"errors"
	"strings"
	"sync"
	"unicode/utf8"
)

// LocalBus is a Bus implementation that runs locally in the process.
// http://www.hivemq.com/blog/mqtt-essentials-part-8-retained-messages
type LocalBus struct {
	mu               sync.Mutex
	persistentTopics map[string][]byte
	subscribers      []subscription
}

func (l *LocalBus) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	for i := range l.subscribers {
		close(l.subscribers[i].channel)
	}
	return nil
}

func (l *LocalBus) Publish(msg Message, qos QOS, retained bool) error {
	p := parseTopic(msg.Topic)
	if p == nil || p.isQuery() {
		return errors.New("invalid topic")
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.persistentTopics == nil {
		l.persistentTopics = map[string][]byte{}
	}
	if len(msg.Payload) == 0 {
		// delete
		delete(l.persistentTopics, msg.Topic)
	} else {
		// Save it first.
		if retained {
			l.persistentTopics[msg.Topic] = msg.Payload
		}
		var chans []chan<- Message
		for i := range l.subscribers {
			if l.subscribers[i].topic.match(msg.Topic) {
				chans = append(chans, l.subscribers[i].channel)
			}
		}
		if len(chans) != 0 {
			go func() {
				for _, c := range chans {
					c <- msg
				}
			}()
		}
	}
	return nil
}

func (l *LocalBus) Subscribe(topic string, qos QOS) (<-chan Message, error) {
	p := parseTopic(topic)
	if p == nil {
		return nil, errors.New("invalid topic")
	}
	c := make(chan Message)
	l.mu.Lock()
	defer l.mu.Unlock()
	l.subscribers = append(l.subscribers, subscription{p, c})
	return c, nil
}

func (l *LocalBus) Unsubscribe(topic string) error {
	p := parseTopic(topic)
	l.mu.Lock()
	defer l.mu.Unlock()
	for i := range l.subscribers {
		if l.subscribers[i].topic.isEqual(p) {
			// Found!
			close(l.subscribers[i].channel)
			copy(l.subscribers[i:], l.subscribers[i+1:])
			l.subscribers = l.subscribers[:len(l.subscribers)-1]
			return nil
		}
	}
	return errors.New("subscription not found")
}

func (l *LocalBus) Get(topic string, qos QOS) ([]Message, error) {
	p := parseTopic(topic)
	if p == nil {
		return nil, errors.New("invalid topic")
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	var out []Message
	for k, v := range l.persistentTopics {
		if p.match(k) {
			out = append(out, Message{k, v})
		}
	}
	return out, nil
}

//

type subscription struct {
	topic   parsedTopic
	channel chan<- Message
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
	t := strings.Split(topic, "/")
	if len(t) == 0 {
		return false
	}
	for i, e := range p {
		if e == "#" {
			return true
		}
		if e == "+" {
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

var _ Bus = &LocalBus{}
