package server_test

import (
	"database/sql"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"logline/internal/config"
	"logline/internal/server"
	"logline/internal/store"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func testConfig() config.Config {
	return config.Config{
		Env:        "test",
		LogLevel:   "debug",
		Port:       0,
		DBMaxConns: 5,
		DBMaxIdle:  2,
	}
}

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL not set")
	}

	db, err := sql.Open("pgx", url)
	if err != nil {
		t.Fatalf("opening test db: %v", err)
	}
	if err := db.Ping(); err != nil {
		t.Fatalf("pinging test db: %v", err)
	}

	t.Cleanup(func() { db.Close() })
	return db
}

func TestHandleHealth_DBDown(t *testing.T) {
	db := setupTestDB(t)
	s := store.New(db)
	srv := server.New(testConfig(), testLogger(), s)

	db.Close() // simulate unreachable DB

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503, got %d", rec.Code)
	}

	var resp map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decoding response: %v", err)
	}

	if resp["error"] != "database unreachable" {
		t.Errorf("expected error %q, got %q", "database unreachable", resp["error"])
	}
}
