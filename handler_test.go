package main

import (
	"encoding/json"
	"io"
	"log/slog"
	"maps"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func request(t *testing.T, h http.Handler, method, path, user string) (int, string) {
	t.Helper()
	req := httptest.NewRequest(method, path, nil)
	if user != "" {
		req.Header.Set("Remote-User", user)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec.Code, rec.Body.String()
}

func get(t *testing.T, h http.Handler, path, user string) (int, string) {
	return request(t, h, http.MethodGet, path, user)
}

func TestHealthz(t *testing.T) {
	h := newHandler(t.TempDir(), slog.New(slog.NewTextHandler(io.Discard, nil)))
	if code, _ := get(t, h, "/healthz", ""); code != http.StatusOK {
		t.Fatalf("GET /healthz = %d, want 200", code)
	}
}

func TestHelloGreetsTheUserAndCounts(t *testing.T) {
	h := newHandler(t.TempDir(), slog.New(slog.NewTextHandler(io.Discard, nil)))
	_, first := get(t, h, "/", "alice")
	code, second := get(t, h, "/", "alice")
	if code != http.StatusOK {
		t.Fatalf("GET / = %d, want 200", code)
	}
	if !strings.Contains(first, "Bonjour alice") || !strings.Contains(first, "Visites : 1") {
		t.Fatalf("first answer = %q", first)
	}
	if !strings.Contains(second, "Visites : 2") {
		t.Fatalf("second answer = %q", second)
	}
}

func TestHelloWithoutUser(t *testing.T) {
	h := newHandler(t.TempDir(), slog.New(slog.NewTextHandler(io.Discard, nil)))
	if _, body := get(t, h, "/", ""); !strings.Contains(body, "Bonjour inconnu") {
		t.Fatalf("answer = %q", body)
	}
}

func TestHelloWithoutWritableData(t *testing.T) {
	h := newHandler("/nonexistent/data", slog.New(slog.NewTextHandler(io.Discard, nil)))
	code, body := get(t, h, "/", "bob")
	if code != http.StatusOK || !strings.Contains(body, "Bonjour bob") {
		t.Fatalf("GET / = %d %q, want 200 and a greeting", code, body)
	}
}

func TestUnknownPath(t *testing.T) {
	h := newHandler(t.TempDir(), slog.New(slog.NewTextHandler(io.Discard, nil)))
	if code, _ := get(t, h, "/nope", ""); code != http.StatusNotFound {
		t.Fatalf("GET /nope = %d, want 404", code)
	}
}

func TestHelloShowsTheVersion(t *testing.T) {
	h := newHandler(t.TempDir(), slog.New(slog.NewTextHandler(io.Discard, nil)))
	if _, body := get(t, h, "/", "alice"); !strings.Contains(body, "Version : "+version) {
		t.Fatalf("answer = %q, want the version %q", body, version)
	}
}

func TestVersion(t *testing.T) {
	h := newHandler(t.TempDir(), slog.New(slog.NewTextHandler(io.Discard, nil)))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/version", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /version = %d, want 200", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("Content-Type = %q, want application/json", ct)
	}
	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("body %q is not a JSON object of strings: %v", rec.Body.String(), err)
	}
	if want := map[string]string{"version": version}; !maps.Equal(body, want) {
		t.Fatalf("body = %v, want exactly %v", body, want)
	}
}

func TestVersionMethods(t *testing.T) {
	h := newHandler(t.TempDir(), slog.New(slog.NewTextHandler(io.Discard, nil)))
	for method, want := range map[string]int{
		http.MethodHead:    http.StatusOK, // answered like GET; a real server drops the body
		http.MethodOptions: http.StatusMethodNotAllowed,
		http.MethodPost:    http.StatusMethodNotAllowed,
		http.MethodPut:     http.StatusMethodNotAllowed,
		http.MethodDelete:  http.StatusMethodNotAllowed,
		http.MethodPatch:   http.StatusMethodNotAllowed,
	} {
		if code, _ := request(t, h, method, "/version", ""); code != want {
			t.Errorf("%s /version = %d, want %d", method, code, want)
		}
	}
}
