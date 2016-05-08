// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"image/color"
	"log"
	"mime"
	"net/http"
	"path"
	"sort"

	"github.com/maruel/dotstar/anim1d"
)

type webServer struct {
	painter  *anim1d.Painter
	registry *anim1d.PatternRegistry
}

func startWebServer(port int, painter *anim1d.Painter, registry *anim1d.PatternRegistry) *webServer {
	ws := &webServer{painter: painter, registry: registry}
	mux := http.NewServeMux()
	mux.HandleFunc("/", ws.rootHandler)
	mux.HandleFunc("/favicon.ico", ws.faviconHandler)
	mux.HandleFunc("/patterns", ws.patternsHandler)
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
	if n := r.PostFormValue("mode"); len(n) != 0 {
		log.Printf("mode = %s", n)
		if p := s.registry.Patterns[n]; p != nil {
			s.painter.SetPattern(p)
			return
		}
		// TODO(maruel): return an error.
		return
	}

	if n := r.PostFormValue("color"); len(n) != 0 {
		log.Printf("color = %s", n)
		if len(n) != 7 || n[0] != '#' {
			// TODO(maruel): return an error.
			return
		}
		b, err := hex.DecodeString(n[1:])
		if err != nil {
			// TODO(maruel): return an error.
			return
		}
		s.painter.SetPattern(&anim1d.StaticColor{color.NRGBA{b[0], b[1], b[2], 255}})
		return
	}

	// TODO(maruel): return an error.
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

func (s *webServer) patternsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Ugh", http.StatusMethodNotAllowed)
		return
	}
	out := make([]string, 0, len(s.registry.Patterns))
	for name := range s.registry.Patterns {
		out = append(out, name)
	}
	sort.Strings(out)
	bytes, err := json.Marshal(out)
	if err != nil {
		http.Error(w, "Ugh", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
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
	if data := s.registry.Thumbnail(r.URL.Path[len("/thumbnail/"):]); len(data) != 0 {
		_, _ = w.Write(data)
	} else {
		http.Error(w, "Not Found", http.StatusNotFound)
	}
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
