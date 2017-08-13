// Copyright 2016 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"mime"
	"net"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/maruel/anim1d"
	"github.com/maruel/circular"
	"github.com/maruel/dlibox/go/modules"
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
	b      modules.Bus
	l      io.WriterTo
	cache  anim1d.ThumbnailsCache
	config *Config
	ln     net.Listener
	server http.Server
}

func initWeb(b modules.Bus, port int, config *Config, l io.WriterTo) (*webServer, error) {
	s := &webServer{
		b: b,
		l: l,
		cache: anim1d.ThumbnailsCache{
			NumberLEDs:       100,
			ThumbnailHz:      10,
			ThumbnailSeconds: 10,
		},
		config: config,
	}
	// Static replies.
	http.HandleFunc("/", s.rootHandler)
	http.HandleFunc("/favicon.ico", s.faviconHandler)
	http.HandleFunc("/static/", s.staticHandler)
	// Dynamic replies.
	http.HandleFunc("/api/pattern", s.patternHandler)
	http.HandleFunc("/api/patterns", s.patternsHandler)
	http.HandleFunc("/api/publish", s.publishHandler)
	http.HandleFunc("/api/settings", s.settingsHandler)
	http.HandleFunc("/thumbnail/", s.thumbnailHandler)
	http.HandleFunc("/logs", s.logsHandler)

	var err error
	s.ln, err = net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	s.server = http.Server{
		Addr:           s.ln.Addr().String(),
		Handler:        loggingHandler{http.DefaultServeMux},
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
	w.Header().Set("Cache-Control", cacheControl5m)
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
	w.Header().Set("Cache-Control", cacheControl5m)
	w.Write(mustRead("favicon.ico"))
}

func (s *webServer) staticHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Ugh", http.StatusMethodNotAllowed)
		return
	}
	p := r.URL.Path[len("/static/"):]
	w.Header().Set("Content-Type", mime.TypeByExtension(path.Ext(p)))
	w.Header().Set("Cache-Control", cacheControl5m)
	if content := read(p); content != nil {
		_, _ = w.Write(content)
		return
	}
	http.Error(w, "Not Found", http.StatusNotFound)
}

func (s *webServer) patternsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Ugh", http.StatusMethodNotAllowed)
		return
	}
	s.config.LRU.Lock()
	defer s.config.LRU.Unlock()
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "Cache-Control:no-cache, no-store")
	json.NewEncoder(w).Encode(s.config.LRU.Patterns)
}

func (s *webServer) publishHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Ugh", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "Cache-Control:no-cache, no-store")
	state := r.PostFormValue("state")
	if len(state) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "state is required"})
		return
	}
	if !State(state).Valid() {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "state is invalid"})
		return
	}
	if err := s.b.Publish(modules.Message{"//dlibox/halloween/state", []byte(state)}, modules.BestEffort, false); err != nil {
		log.Printf("web: failed to publish: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("failed to publish: %v", err)})
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"ok": "1"})
}

func (s *webServer) settingsHandler(w http.ResponseWriter, r *http.Request) {
	s.config.Settings.Lock()
	defer s.config.Settings.Unlock()
	switch r.Method {
	case "GET":
		data, _ := json.Marshal(s.config.Settings)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "Cache-Control:no-cache, no-store")
		w.Write(data)
	case "POST":
		// TODO(maruel): Accept JSON.
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "Cache-Control:no-cache, no-store")
		settings := Settings{}
		rawEncoded := r.PostFormValue("settings")
		if len(rawEncoded) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "settings is required"})
			return
		}
		raw, err := base64.URLEncoding.DecodeString(rawEncoded)
		if len(raw) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "settings content is required"})
			return
		}
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "settings is not base64"})
			return
		}
		if err := json.Unmarshal(raw, &settings); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("failed to decode: %v", err)})
			return
		}
		// The lock is a problem here so we can't just copy in there. Instead,
		// unmarshal a second time, so the lock is unaffected.
		if err := json.Unmarshal(raw, &s.config.Settings); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("failed to decode: %v", err)})
			return
		}
		// Serialize it again to return the canonical form.
		json.NewEncoder(w).Encode(settings)
	default:
		http.Error(w, "Ugh", http.StatusMethodNotAllowed)
	}
}

func (s *webServer) patternHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "Cache-Control:no-cache, no-store")
		s.config.Settings.Painter.Lock()
		defer s.config.Settings.Painter.Unlock()
		l := s.config.Settings.Painter.Last
		if l == "" {
			l = s.config.Settings.Painter.Startup
		}
		w.Write([]byte(l))
		return
	}
	if r.Method != "POST" {
		http.Error(w, "Ugh", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "Cache-Control:no-cache, no-store")
	// TODO(maruel): Accept JSON.
	rawEncoded := r.PostFormValue("pattern")
	if len(rawEncoded) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "pattern is required"})
		return
	}
	raw, err := base64.URLEncoding.DecodeString(rawEncoded)
	if len(raw) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "pattern content is required"})
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "pattern is not base64"})
		return
	}

	// Reencode in canonical format to send it back to the user.
	var obj anim1d.SPattern
	if err := json.Unmarshal(raw, &obj); err != nil {
		log.Printf("web: invalid JSON pattern: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	raw, err = obj.MarshalJSON()
	if err != nil {
		log.Printf("web: internal error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	if err := s.b.Publish(modules.Message{"painter/setuser", raw}, modules.BestEffort, false); err != nil {
		log.Printf("web: failed to publish: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("failed to publish: %v", err)})
		return
	}
	w.Write(raw)
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
	w.Header().Set("Cache-Control", cacheControl5m)
	_, _ = w.Write(data)
}

func (s *webServer) logsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	// Streams the log buffer over HTTP until Close() is called.
	// AutoFlush ensures the log is not buffered locally indefinitely.
	s.l.WriteTo(circular.AutoFlush(w, time.Second))
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
