# msgbus

A simplified interface to MQTT.

More than a simple MQTT client implementation, it implements rooting a topic
tree, logging and retained topic retrieval.

Uses https://github.com/maruel/paho.mqtt.golang, a fork of
https://github.com/eclipse/paho.mqtt.golang that drops websocket support to
reduce the dependency set.

[![GoDoc](https://godoc.org/github.com/maruel/msgbus?status.svg)](https://godoc.org/github.com/maruel/msgbus) [![Go Report Card](https://goreportcard.com/badge/github.com/maruel/msgbus)](https://goreportcard.com/report/github.com/maruel/msgbus)
