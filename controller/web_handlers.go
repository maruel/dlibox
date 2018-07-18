// Copyright 2018 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

//go:generate go run internal/static_gen.go -o static_prod.go

package controller

import (
	"encoding/base64"
	"mime"
	"net/http"
	"path/filepath"

	"github.com/maruel/anim1d"
)

const cacheControl5m = "Cache-Control:public,max-age=300"

func (s *webServer) addOtherHandlers() {
	s.cache = anim1d.ThumbnailsCache{
		NumberLEDs:       100,
		ThumbnailHz:      10,
		ThumbnailSeconds: 10,
	}
	http.HandleFunc("/raw/dlibox/v1/thumbnail/", noContent(s.apithumbnailtokenhandler))
	http.HandleFunc("/raw/dlibox/v1/xsrf_token", noContent(s.apiXSRFTokenHandler))
	http.HandleFunc("/raw/dlibox/v1/log", s.logHandler)
	http.HandleFunc("/favicon.ico", getOnly(s.getFavicon))
	http.HandleFunc("/static/", getOnly(s.getStatic))
	// Do not use getOnly here as it is the 'catch all, one and we want to check
	// that before the method.
	http.HandleFunc("/", noContent(s.getRoot))

}

// Static handlers

// /
func (s *webServer) getRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" && r.Method != "HEAD" {
		http.Error(w, "Only GET is allowed", http.StatusMethodNotAllowed)
		return
	}
	s.setXSRFCookie(r.RemoteAddr, w)
	content := getContent("static/index.html")
	if content == nil {
		http.Error(w, "Content missing", 500)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Cache-Control", cacheControlContent)
	w.Write(content)
}

// /favicon.ico
func (s *webServer) getFavicon(w http.ResponseWriter, r *http.Request) {
	content := getContent("static/favicon.ico")
	if content == nil {
		http.Error(w, "Content missing", 500)
		return
	}
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", cacheControlContent)
	w.Write(content)
}

// /static/
func (s *webServer) getStatic(w http.ResponseWriter, r *http.Request) {
	content := getContent(r.URL.Path[1:])
	if content == nil {
		http.Error(w, "Content missing", 404)
		return
	}
	m := "application/octect-stream"
	if ext := filepath.Ext(r.URL.Path); len(ext) > 1 {
		if m2 := mime.TypeByExtension(ext); m2 != "" {
			m = m2
		}
	}
	w.Header().Set("Content-Type", m)
	w.Header().Set("Cache-Control", cacheControlContent)
	w.Write(content)
}

// Other handlers

// /raw/dlibox/v1/xsrf_token
func (s *webServer) apiXSRFTokenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST is allowed", http.StatusMethodNotAllowed)
		return
	}
	t := s.setXSRFCookie(r.RemoteAddr, w)
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Cache-Control", cacheControlNone)
	w.WriteHeader(200)
	w.Write([]byte(t))
}

// /raw/dlibox/v1/log
func (s *webServer) logHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	// Streams the log buffer over HTTP until Close() is called.
	// AutoFlush ensures the log is not buffered locally indefinitely.
	//j.l.WriteTo(circular.AutoFlush(w, time.Second))
	w.Write([]byte("TODO"))
}

// /raw/dlibox/v1/thumbnail
func (s *webServer) apithumbnailtokenhandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Ugh", http.StatusMethodNotAllowed)
		return
	}
	b := r.URL.Path[len("//raw/dlibox/v1/thumbnail/"):]
	if len(b) == 0 {
		http.Error(w, "Ugh", 404)
		return
	}
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
	w.Header().Set("Content-Type", "image/gif")
	w.Header().Set("Cache-Control", cacheControl5m)
	_, _ = w.Write(data)
}
