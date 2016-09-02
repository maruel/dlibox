// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"mime"
	"net"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/maruel/dlibox/go/anim1d"
)

var (
	hostName string
	rootTmpl *template.Template
)

func init() {
	var err error
	if hostName, err = os.Hostname(); err != nil {
		panic(err)
	}
	rootTmpl = template.Must(template.New("name").Parse(string(mustRead("root.html"))))
}

type webServer struct {
	painter *anim1d.Painter
	cache   anim1d.ThumbnailsCache
	config  *Config
	ln      net.Listener
	server  http.Server
}

func startWebServer(port int, painter *anim1d.Painter, config *Config) (*webServer, error) {
	s := &webServer{
		painter: painter,
		cache: anim1d.ThumbnailsCache{
			NumberLEDs:       100,
			ThumbnailHz:      10,
			ThumbnailSeconds: 10,
		},
		config: config,
	}
	mux := http.NewServeMux()
	// Static replies.
	mux.HandleFunc("/", s.rootHandler)
	mux.HandleFunc("/favicon.ico", s.faviconHandler)
	mux.HandleFunc("/static/", s.staticHandler)
	// Dynamic replies.
	mux.HandleFunc("/config", s.configHandler)
	mux.HandleFunc("/switch", s.switchHandler)
	mux.HandleFunc("/thumbnail/", s.thumbnailHandler)

	var err error
	s.ln, err = net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	s.server = http.Server{
		Addr:           s.ln.Addr().String(),
		Handler:        loggingHandler{mux},
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 16,
	}
	go s.server.Serve(s.ln)
	return s, nil
}

func (s *webServer) Close() error {
	err := s.ln.Close()
	s.ln = nil
	return err
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
	//w.Header().Set("Cache-Control", "Cache-Control:public, max-age=2592000") // 30d
	keys := struct {
		Host string
	}{
		hostName,
	}
	if err := rootTmpl.Execute(w, keys); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *webServer) faviconHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Ugh", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "Cache-Control:public, max-age=2592000") // 30d
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

func (s *webServer) configHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Ugh", http.StatusMethodNotAllowed)
		return
	}

	data, _ := json.Marshal(s.config)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (s *webServer) switchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Ugh", http.StatusMethodNotAllowed)
		return
	}
	rawEncoded := r.PostFormValue("pattern")
	if len(rawEncoded) == 0 {
		http.Error(w, "pattern is required", http.StatusBadRequest)
		return
	}
	raw, err := base64.URLEncoding.DecodeString(rawEncoded)
	if len(raw) == 0 {
		http.Error(w, "pattern content is required", http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, "pattern is not base64", http.StatusBadRequest)
		return
	}
	// Reformat the pattern in canonical format.
	// TODO(maruel): Change SetPattern() to accept a anim1d.Pattern.
	var p anim1d.SPattern
	if err := json.Unmarshal(raw, &p); err != nil {
		http.Error(w, "invalid JSON pattern", http.StatusBadRequest)
		return
	}
	b, err := p.MarshalJSON()
	if err != nil {
		http.Error(w, "internal error", http.StatusBadRequest)
		return
	}
	pattern := string(b)
	log.Printf("pattern = %q", pattern)
	if err := s.painter.SetPattern(pattern); err != nil {
		http.Error(w, "invalid JSON pattern", http.StatusBadRequest)
		return
	}

	s.config.Inject(pattern)
	_, _ = w.Write([]byte("success"))
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
	w.Header().Set("Content-Type", "image/gif")
	//w.Header().Set("Cache-Control", "Cache-Control:public, max-age=2592000") // 30d
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
