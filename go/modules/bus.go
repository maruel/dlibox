// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package modules

import "io"

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
	Publish(msg *Message, qos QOS, retained bool) error

	// Subscribe sends back updates to this topic query.
	Subscribe(topic string, qos QOS) (<-chan *Message, error)
	// Unsubscribe removes a previous subscription.
	Unsubscribe(topic string) error
	// Get retrieves matching messages for a retained topic query.
	Get(topic string, qos QOS) ([]*Message, error)
}
