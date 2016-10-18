// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package jsonnotif implements a JSON-RPC ClientCodec and ServerCodec
// for the rpc package that only supports notifications.
package jsonnotif

import (
	"encoding/json"
	"io"
	"net/rpc"
)

type clientCodec struct {
	dec *json.Decoder // for reading JSON values
	enc *json.Encoder // for writing JSON values
	c   io.Closer

	// temporary work space
	req clientRequest
}

// NewClientCodec returns a new rpc.ClientCodec using JSON-RPC on conn.
func NewClientCodec(conn io.ReadWriteCloser) rpc.ClientCodec {
	return &clientCodec{
		dec: json.NewDecoder(conn),
		enc: json.NewEncoder(conn),
		c:   conn,
	}
}

type clientRequest struct {
	Method string         `json:"method"`
	Params [1]interface{} `json:"params"`
}

func (c *clientCodec) WriteRequest(r *rpc.Request, param interface{}) error {
	c.req.Method = r.ServiceMethod
	c.req.Params[0] = param
	return c.enc.Encode(&c.req)
}

func (c *clientCodec) ReadResponseHeader(r *rpc.Response) error {
	return nil
}

func (c *clientCodec) ReadResponseBody(x interface{}) error {
	return nil
}

func (c *clientCodec) Close() error {
	return c.c.Close()
}

func NewClient(conn io.ReadWriteCloser) *rpc.Client {
	return rpc.NewClientWithCodec(NewClientCodec(conn))
}
