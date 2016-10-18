// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package jsonnotif

import (
	"encoding/json"
	"errors"
	"io"
	"net/rpc"
)

var errMissingParams = errors.New("jsonrpc: request body missing params")

type serverCodec struct {
	dec *json.Decoder // for reading JSON values
	enc *json.Encoder // for writing JSON values
	c   io.Closer

	// temporary work space
	req serverRequest
}

// NewServerCodec returns a new rpc.ServerCodec using JSON-RPC on conn.
func NewServerCodec(conn io.ReadWriteCloser) rpc.ServerCodec {
	return &serverCodec{
		dec: json.NewDecoder(conn),
		enc: json.NewEncoder(conn),
		c:   conn,
	}
}

type serverRequest struct {
	Method string           `json:"method"`
	Params *json.RawMessage `json:"params"`
}

func (r *serverRequest) reset() {
	r.Method = ""
	r.Params = nil
}

func (c *serverCodec) ReadRequestHeader(r *rpc.Request) error {
	c.req.reset()
	if err := c.dec.Decode(&c.req); err != nil {
		return err
	}
	r.ServiceMethod = c.req.Method
	return nil
}

func (c *serverCodec) ReadRequestBody(x interface{}) error {
	if x == nil {
		return nil
	}
	if c.req.Params == nil {
		return errMissingParams
	}
	// JSON params is array value.
	// RPC params is struct.
	// Unmarshal into array containing struct for now.
	// Should think about making RPC more general.
	var params [1]interface{}
	params[0] = x
	return json.Unmarshal(*c.req.Params, &params)
}

var null = json.RawMessage([]byte("null"))

func (c *serverCodec) WriteResponse(r *rpc.Response, x interface{}) error {
	return nil
}

func (c *serverCodec) Close() error {
	return c.c.Close()
}

// ServeConn runs the JSON-RPC server on a single connection.
// ServeConn blocks, serving the connection until the client hangs up.
// The caller typically invokes ServeConn in a go statement.
func ServeConn(conn io.ReadWriteCloser) {
	rpc.ServeCodec(NewServerCodec(conn))
}
