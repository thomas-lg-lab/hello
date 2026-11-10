package main

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func get(t *testing.T, h http.Handler, path, user string) (int, string) {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, path, nil)
	if user != "" {
		req.Header.Set("Remote-User", user)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec.Code, rec.Body.String()
}

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
	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /version = %d, want 200", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("Content-Type = %q, want %q", ct, "application/json")
	}
	body := rec.Body.String()
	if !strings.Contains(body, `"version":"`+version+`"`) {
		t.Fatalf("body = %q, want {\"version\":\"%s\"}", body, version)
	}
}

func TestVersionMethodNotAllowed(t *testing.T) {
	h := newHandler(t.TempDir(), slog.New(slog.NewTextHandler(io.Discard, nil)))
	for _, method := range []string{"POST", "PUT", "DELETE", "PATCH"} {
		code, _ := request(t, h, method, "/version", "")
		if code != http.StatusMethodNotAllowed {
			t.Fatalf("%s /version = %d, want 405", method, code)
		}
	}
}
