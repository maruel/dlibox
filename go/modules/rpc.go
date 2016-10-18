// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package modules

import "net/rpc"

// Use JSONRPC for notifications only message. The "rpc" keyword is misleading,
// this code never waits for a reply.

type Notifications struct {
}

func ToRPC(m Bus) *rpc.Server {
}

func FromSubscription(c <-chan *Message) {
}

type Command string

const (
	SetPattern Command = "setpattern"
)

type CommandMsg struct {
	Base Command
	Call string
}

func (c *CommandMsg) ToMsg(root string) Message {
	return Message{Topic: root + "/" + c.Base, Payload: nil}
}
