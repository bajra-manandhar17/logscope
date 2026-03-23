package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"
)

// stubFS provides a minimal in-memory FS with dist/index.html so the static
// handler doesn't panic during tests.
var stubFS = fstest.MapFS{
	"dist/index.html": &fstest.MapFile{Data: []byte("<html></html>")},
}

func TestHealthEndpoint(t *testing.T) {
	srv := New(stubFS, false)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	srv.Handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var body map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if body["status"] != "ok" {
		t.Fatalf("expected status=ok, got %q", body["status"])
	}
}

func TestCORSHeadersPresentInDevMode(t *testing.T) {
	srv := New(stubFS, true)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	srv.Handler.ServeHTTP(rec, req)

	got := rec.Header().Get("Access-Control-Allow-Origin")
	if got != "http://localhost:5173" {
		t.Fatalf("expected CORS origin header, got %q", got)
	}
}

func TestCORSHeadersAbsentInProdMode(t *testing.T) {
	srv := New(stubFS, false)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	srv.Handler.ServeHTTP(rec, req)

	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Fatalf("expected no CORS header in prod mode, got %q", got)
	}
}

func TestCORSPreflightReturns204(t *testing.T) {
	srv := New(stubFS, true)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodOptions, "/api/health", nil)
	srv.Handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rec.Code)
	}
}
