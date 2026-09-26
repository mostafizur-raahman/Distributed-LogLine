package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handleHealth)
	mux.HandleFunc("POST /ingest", handleIngest)
	return mux
}

func TestHealthEndpoint(t *testing.T) {
	mux := newMux()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	body := strings.TrimSpace(rec.Body.String())

	if body != "healthy" {
		t.Errorf("expected body %q, got %q", "ok", body)
	}

}

func TestIngestEndpoint(t *testing.T) {
	mux := newMux()

	payload := `{"service":"auth","level":"error","message":"connection refused"}`
	req := httptest.NewRequest(http.MethodPost, "/ingest", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", rec.Code)
	}
}

func TestUnknownRouteReturns404(t *testing.T) {
	mux := newMux()

	req := httptest.NewRequest(http.MethodGet, "/nonexistent", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", rec.Code)
	}
}

func TestWrongMethodReturns405(t *testing.T) {
	mux := newMux()

	req := httptest.NewRequest(http.MethodGet, "/ingest", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", rec.Code)
	}
}
