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

	"github.com/maruel/dlibox/go/anim1d"
	"github.com/maruel/dlibox/go/pio/fakes"
	"github.com/maruel/ut"
)

func TestWeb(t *testing.T) {
	var config Config
	config.ResetDefault()
	d := &fakes.Display{image.NewNRGBA(image.Rect(0, 0, 128, 1))}
	painter := anim1d.MakePainter(d, 60)
	defer painter.Close()
	s, err := startWebServer(0, painter, &config)
	defer s.Close()
	ut.AssertEqual(t, nil, err)
	base := fmt.Sprintf("http://%s/", s.server.Addr)
	resp, err := http.PostForm(base+"switch", url.Values{"pattern": {base64.URLEncoding.EncodeToString([]byte("\"#010203\""))}})
	ut.AssertEqual(t, nil, err)
	raw, err := ioutil.ReadAll(resp.Body)
	ut.AssertEqual(t, nil, err)
	ut.AssertEqual(t, "success", string(raw))
}
