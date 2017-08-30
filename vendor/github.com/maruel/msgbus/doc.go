// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Package msgbus implements a generic PubSub message bus that follows MQTT
// guidelines.
//
// The main difference with MQTT topic is the support for relative message on
// rebased bus. See RebasePub() for more details.
//
// Spec
//
// The MQTT specification lives at
// http://docs.oasis-open.org/mqtt/mqtt/v3.1.1/mqtt-v3.1.1.html
package msgbus
