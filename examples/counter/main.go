package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"sync"
	"time"
)

type state struct {
	mu        sync.Mutex
	remaining int
	limit     int
	racey     bool
	delay     time.Duration
}

func main() {
	addr := flag.String("addr", "127.0.0.1:8090", "listen address (loopback)")
	limit := flag.Int("limit", 10, "quota limit")
	racey := flag.Bool("racey", false, "naive check-then-act (TOCTOU)")
	delayMS := flag.Int("delay-ms", 8, "sleep after check in racey mode")
	flag.Parse()

	s := &state{remaining: *limit, limit: *limit, racey: *racey, delay: time.Duration(*delayMS) * time.Millisecond}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /take", s.take)
	mux.HandleFunc("GET /status", s.status)
	mux.HandleFunc("POST /reset", s.reset)
	log.Printf("counter listening on http://%s  limit=%d racey=%v", *addr, *limit, *racey)
	log.Fatal(http.ListenAndServe(*addr, mux))
}

func (s *state) take(w http.ResponseWriter, _ *http.Request) {
	if s.racey {
		if s.remaining <= 0 {
			http.Error(w, `{"ok":false,"reason":"exhausted"}`, http.StatusConflict)
			return
		}
		time.Sleep(s.delay)
		s.remaining--
		writeJSON(w, http.StatusOK, map[string]any{"ok": true, "remaining": s.remaining})
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if s.remaining <= 0 {
		http.Error(w, `{"ok":false,"reason":"exhausted"}`, http.StatusConflict)
		return
	}
	s.remaining--
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "remaining": s.remaining})
}

func (s *state) status(w http.ResponseWriter, _ *http.Request) {
	s.mu.Lock()
	rem := s.remaining
	lim := s.limit
	s.mu.Unlock()
	if s.racey {
		rem = s.remaining
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"remaining": rem,
		"used":      lim - rem,
		"limit":     lim,
		"racey":     s.racey,
	})
}

func (s *state) reset(w http.ResponseWriter, _ *http.Request) {
	s.mu.Lock()
	s.remaining = s.limit
	s.mu.Unlock()
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "remaining": s.limit})
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}
