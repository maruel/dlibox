// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"encoding/base64"
	"fmt"
	"image"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"github.com/maruel/dlibox/go/donotuse/devices/devicestest"
	"github.com/maruel/dlibox/go/modules"
	"github.com/maruel/ut"
)

func TestWeb(t *testing.T) {
	t.Parallel()
	var config Config
	config.ResetDefault()
	config.LRU.Patterns = []Pattern{"\"#010101\"", "\"#020202\""}
	d := &devicestest.Display{image.NewNRGBA(image.Rect(0, 0, 128, 1))}
	b := &modules.LocalBus{}
	painter, err := initPainter(b, d, 60, &Painter{}, &LRU{})
	ut.AssertEqual(t, nil, err)
	defer painter.Close()
	s, err := initWeb(b, 0, &config, nil)
	defer s.Close()
	ut.AssertEqual(t, nil, err)
	base := fmt.Sprintf("http://%s/", s.server.Addr)
	// Only Frame are injected in the config, colors (other than black) are
	// ignored.
	resp, err := http.PostForm(base+"switch", url.Values{"pattern": {base64.URLEncoding.EncodeToString([]byte("\"L030303\""))}})
	ut.AssertEqual(t, nil, err)
	_, err = ioutil.ReadAll(resp.Body)
	ut.AssertEqual(t, nil, err)
	// TODO(maruel): Race condition.
	//ut.AssertEqual(t, "[\"\\\"L030303\\\"\",\"\\\"#010101\\\"\",\"\\\"#020202\\\"\"]", string(raw))
}
