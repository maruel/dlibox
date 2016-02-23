// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/maruel/interrupt"
)

type WebServer struct {
	cond  sync.Cond
	state string
}

func StartWebServer(port int) *WebServer {
	s := &WebServer{
		cond: *sync.NewCond(&sync.Mutex{}),
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.root)
	mux.HandleFunc("/favicon.ico", s.favicon)
	fmt.Printf("Listening on %d\n", port)
	go http.ListenAndServe(fmt.Sprintf(":%d", port), loggingHandler{mux})
	go func() {
		<-interrupt.Channel
		s.cond.Broadcast()
	}()
	return s
}

func (s *WebServer) Close() error {
	s.cond.Broadcast()
	return nil
}

func (s *WebServer) root(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	s.cond.L.Lock()
	defer s.cond.L.Unlock()
	if _, err := w.Write(read("root.html")); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *WebServer) favicon(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "Cache-Control:public, max-age=2592000") // 30d
	w.Write(read("dotstar.png"))
}

// Private details.

type loggingHandler struct {
	handler http.Handler
}

type loggingResponseWriter struct {
	http.ResponseWriter
	length int
	status int
}

func (l *loggingResponseWriter) Write(data []byte) (size int, err error) {
	size, err = l.ResponseWriter.Write(data)
	l.length += size
	return
}

func (l *loggingResponseWriter) WriteHeader(status int) {
	l.ResponseWriter.WriteHeader(status)
	l.status = status
}

// ServeHTTP logs each HTTP request if -verbose is passed.
func (l loggingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	lrw := &loggingResponseWriter{ResponseWriter: w}
	l.handler.ServeHTTP(lrw, r)
	log.Printf("%s - %3d %6db %4s %s\n", r.RemoteAddr, lrw.status, lrw.length, r.Method, r.RequestURI)
}
