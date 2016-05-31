// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"mime"
	"net/http"
	"path"

	"github.com/maruel/dotstar/anim1d"
)

type webServer struct {
	painter *anim1d.Painter
	cache   anim1d.ThumbnailsCache
}

func startWebServer(port int, painter *anim1d.Painter) *webServer {
	ws := &webServer{
		painter: painter,
		cache: anim1d.ThumbnailsCache{
			NumberLEDs:       100,
			ThumbnailHz:      10,
			ThumbnailSeconds: 10,
		},
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", ws.rootHandler)
	mux.HandleFunc("/favicon.ico", ws.faviconHandler)
	mux.HandleFunc("/static/", ws.staticHandler)
	mux.HandleFunc("/switch", ws.switchHandler)
	mux.HandleFunc("/thumbnail/", ws.thumbnailHandler)
	go http.ListenAndServe(fmt.Sprintf(":%d", port), loggingHandler{mux})
	return ws
}

func (s *webServer) rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Ugh", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	if _, err := w.Write(mustRead("root.html")); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *webServer) switchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Ugh", http.StatusMethodNotAllowed)
		return
	}
	b := r.PostFormValue("pattern")
	if len(b) == 0 {
		http.Error(w, "pattern is required", http.StatusBadRequest)
		return
	}
	p, err := base64.URLEncoding.DecodeString(b)
	if err != nil {
		http.Error(w, "pattern is not base64", http.StatusBadRequest)
		return
	}
	p2 := string(p)
	log.Printf("pattern = %q", p2)
	if err := s.painter.SetPattern(p2); err != nil {
		http.Error(w, "invalid JSON pattern", http.StatusBadRequest)
	}
}

func (s *webServer) faviconHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Ugh", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "image/png")
	//w.Header().Set("Cache-Control", "Cache-Control:public, max-age=2592000") // 30d
	w.Write(mustRead("favicon.ico"))
}

func (s *webServer) staticHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Ugh", http.StatusMethodNotAllowed)
		return
	}
	p := r.URL.Path[len("/static/"):]
	w.Header().Set("Content-Type", mime.TypeByExtension(path.Ext(p)))
	//w.Header().Set("Cache-Control", "Cache-Control:public, max-age=2592000") // 30d
	if content := read(p); content != nil {
		w.Write(content)
		return
	}
	http.Error(w, "Not Found", http.StatusNotFound)
}

func (s *webServer) thumbnailHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Ugh", http.StatusMethodNotAllowed)
		return
	}
	b := r.URL.Path[len("/thumbnail/"):]
	p, err := base64.URLEncoding.DecodeString(b)
	if err != nil {
		http.Error(w, "pattern is not base64", http.StatusBadRequest)
		return
	}
	data, err := s.cache.GIF(p)
	if err != nil {
		http.Error(w, "invalid JSON pattern", http.StatusBadRequest)
		return
	}
	_, _ = w.Write(data)
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
