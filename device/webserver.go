// Copyright 2017 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package device

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/maruel/dlibox/shared"
	"github.com/maruel/serve-dir/loghttp"
)

// webServer is the device's web server. It is quite simple.
func webServer(server string, port int) error {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	if server != "" {
		http.DefaultServeMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "http://"+server, 302)
		})
	}
	s := http.Server{
		Addr:           ln.Addr().String(),
		Handler:        &loghttp.Handler{Handler: http.DefaultServeMux},
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 16,
	}
	go s.Serve(ln)
	log.Printf("Visit: http://%s:%d/debug/pprof for debugging", shared.Hostname(), port)
	return nil
}
