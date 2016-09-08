// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package host

// Key represents one of the supported key press.
type Key string

const (
	Key100Plus     Key = "KEY_100PLUS"
	Key200Plus     Key = "KEY_200PLUS"
	KeyChannel     Key = "KEY_CHANNEL"
	KeyChannelDown Key = "KEY_CHANNELDOWN"
	KeyChannelUp   Key = "KEY_CHANNELUP"
	KeyEQ          Key = "KEY_EQ"
	KeyNext        Key = "KEY_NEXT"
	KeyNumeric0    Key = "KEY_NUMERIC_0"
	KeyNumeric1    Key = "KEY_NUMERIC_1"
	KeyNumeric2    Key = "KEY_NUMERIC_2"
	KeyNumeric3    Key = "KEY_NUMERIC_3"
	KeyNumeric4    Key = "KEY_NUMERIC_4"
	KeyNumeric5    Key = "KEY_NUMERIC_5"
	KeyNumeric6    Key = "KEY_NUMERIC_6"
	KeyNumeric7    Key = "KEY_NUMERIC_7"
	KeyNumeric8    Key = "KEY_NUMERIC_8"
	KeyNumeric9    Key = "KEY_NUMERIC_9"
	KeyPlayPause   Key = "KEY_PLAYPAUSE"
	KeyPrevious    Key = "KEY_PREVIOUS"
	KeyVolumeDown  Key = "KEY_VOLUMEDOWN"
	KeyVolumeUp    Key = "KEY_VOLUMEUP"
)

type Message struct {
	Key        Key
	RemoteType string // Remote type name
	Repeat     bool   // true if the button press is a repeated key press; i.e. the user holds the button
}

// IR defines an infrared receiver and emitter.
type IR interface {
	// Channel returns a channel that is used to listen to new messages capted by
	// the IR receiver. It will be closed when the device is closed.
	Channel() <-chan Message
	// Emit emits a key press.
	Emit(remote string, key Key) error
}
