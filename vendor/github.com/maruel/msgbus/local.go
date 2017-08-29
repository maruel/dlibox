// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package msgbus

import (
	"errors"
	"log"
	"sync"
)

// New returns a local thread safe memory backed Bus.
//
// This Bus is thread safe. It is useful for unit tests or as a local broker.
func New() Bus {
	return &local{
		persistentTopics: map[string][]byte{},
	}
}

type local struct {
	mu               sync.Mutex
	persistentTopics map[string][]byte
	subscribers      []*subscription
}

func (l *local) String() string {
	return "LocalBus"
}

func (l *local) Close() error {
	l.mu.Lock()
	subs := l.subscribers
	l.persistentTopics = map[string][]byte{}
	l.subscribers = nil
	l.mu.Unlock()

	for _, s := range subs {
		s.closeSub()
	}
	return nil
}

func (l *local) Publish(msg Message, qos QOS, retained bool) error {
	p := parseTopic(msg.Topic)
	if p == nil || p.isQuery() {
		return errors.New("invalid topic")
	}
	subscribers := func() []*subscription {
		l.mu.Lock()
		defer l.mu.Unlock()
		if len(msg.Payload) == 0 {
			delete(l.persistentTopics, msg.Topic)
			return nil
		}
		if retained {
			b := make([]byte, len(msg.Payload))
			copy(b, msg.Payload)
			l.persistentTopics[msg.Topic] = b
		}
		return l.getSubscribers(msg.Topic)
	}()

	// Do the rest unlocked.
	if qos > BestEffort {
		// Synchronous.
		var wg sync.WaitGroup
		for i := range subscribers {
			wg.Add(1)
			go func(s *subscription) {
				defer wg.Done()
				s.publish(msg)
			}(subscribers[i])
		}
		wg.Wait()
	}

	// Asynchronous.
	for i := range subscribers {
		go subscribers[i].publish(msg)
	}
	return nil
}

func (l *local) Subscribe(topicQuery string, qos QOS) (<-chan Message, error) {
	// QOS is ignored. Eventually it could be used to make the channel buffered.
	p := parseTopic(topicQuery)
	if p == nil {
		return nil, errors.New("invalid topic")
	}
	s := &subscription{topicQuery: p, channel: make(chan Message)}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.subscribers = append(l.subscribers, s)
	return s.channel, nil
}

func (l *local) Unsubscribe(topicQuery string) {
	p := parseTopic(topicQuery)
	if p == nil {
		log.Printf("%s.Unsubscribe(%s): invalid topic", l, topicQuery)
		return
	}
	subscribers := func() []*subscription {
		l.mu.Lock()
		defer l.mu.Unlock()
		var out []*subscription
		for i := 0; i < len(l.subscribers); {
			if l.subscribers[i].topicQuery.isEqual(p) {
				out = append(out, l.subscribers[i])
				copy(l.subscribers[i:], l.subscribers[i+1:])
				l.subscribers = l.subscribers[:len(l.subscribers)-1]
			} else {
				i++
			}
		}
		// Compact array if necessary.
		if cap(l.subscribers) > 16 && cap(l.subscribers) >= 2*len(l.subscribers) {
			s := l.subscribers
			l.subscribers = make([]*subscription, len(s))
			copy(l.subscribers, s)
		}
		return out
	}()
	for _, s := range subscribers {
		s.closeSub()
	}
}

func (l *local) Retained(topicQuery string) ([]Message, error) {
	ps := parseTopic(topicQuery)
	if ps == nil {
		return nil, errors.New("invalid topic")
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	var out []Message
	for t, payload := range l.persistentTopics {
		if ps.match(t) {
			p := make([]byte, len(payload))
			copy(p, payload)
			out = append(out, Message{t, p})
		}
	}
	return out, nil
}

func (l *local) getSubscribers(t string) []*subscription {
	// Must be called with lock held.
	var out []*subscription
	for i := range l.subscribers {
		if l.subscribers[i].topicQuery.match(t) {
			out = append(out, l.subscribers[i])
		}
	}
	return out
}

//

type subscription struct {
	topicQuery parsedTopic

	mu      sync.RWMutex
	channel chan Message
}

// publish synchronously sends the message.
func (s *subscription) publish(msg Message) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.channel == nil {
		return
	}
	s.channel <- msg
}

func (s *subscription) closeSub() {
	s.mu.RLock()
	c := s.channel
	s.mu.RUnlock()
	if c == nil {
		return
	}

	done := make(chan struct{})
	var wg sync.WaitGroup
	defer wg.Wait()
	wg.Add(1)
	go func() {
		defer wg.Done()
		// Empty the channel to unblock publish() call(s) if any.
		for {
			select {
			case <-c:
			case <-done:
				return
			}
		}
	}()

	s.mu.Lock()
	defer func() {
		s.mu.Unlock()
		done <- struct{}{}
	}()
	if s.channel == nil {
		return
	}
	// It's now guaranteed there's no pending publish() call.
	close(s.channel)
	s.channel = nil
}
