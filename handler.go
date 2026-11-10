package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

type app struct {
	dataDir string
	log     *slog.Logger
	mu      sync.Mutex
}

func newHandler(dataDir string, log *slog.Logger) http.Handler {
	a := &app{dataDir: dataDir, log: log}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", a.healthz)
	mux.HandleFunc("GET /{$}", a.hello)
	return mux
}

func (a *app) healthz(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintln(w, "ok")
}

func (a *app) hello(w http.ResponseWriter, r *http.Request) {
	user := strings.TrimSpace(r.Header.Get("Remote-User"))
	if user == "" {
		user = "inconnu"
	}
	visits, err := a.countVisit()
	if err != nil {
		a.log.Error("visit counter", "err", err)
	}
	a.log.Info("hello", "user", user, "visits", visits)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "Bonjour %s\n", user)
	if err == nil {
		fmt.Fprintf(w, "Visites : %d\n", visits)
	}
}

// countVisit increments the counter kept in <dataDir>/visits; a missing file is a first start.
func (a *app) countVisit() (int, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	path := filepath.Join(a.dataDir, "visits")
	n := 0
	if b, err := os.ReadFile(path); err == nil {
		n, _ = strconv.Atoi(strings.TrimSpace(string(b)))
	} else if !os.IsNotExist(err) {
		return 0, err
	}
	n++
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, []byte(strconv.Itoa(n)+"\n"), 0o600); err != nil {
		return 0, err
	}
	return n, os.Rename(tmp, path)
}
