package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

const maxBodySize = 1_048_576 // 1 MB

type LogEntry struct {
	Level     string         `json:"level"`
	Message   string         `json:"message"`
	Service   string         `json:"service"`
	Timestamp time.Time      `json:"timestamp"`
	Data      map[string]any `json:"data,omitempty"`
}

type errorResponse struct {
	Error string `json:"error"`
}

type statusResponse struct {
	Status string `json:"status"`
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", handleHealth)
	mux.HandleFunc("POST /ingest", handleIngest)

	addr := ":4000"

	log.Printf("logline starting on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, statusResponse{Status: "ok"})
}

func handleIngest(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		writeJSON(w, http.StatusUnsupportedMediaType, errorResponse{
			Error: "content-type must be application/json",
		})
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var entry LogEntry
	if err := decoder.Decode(&entry); err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			writeJSON(w, http.StatusRequestEntityTooLarge, errorResponse{
				Error: "request body too large",
			})
			return
		}

		writeJSON(w, http.StatusBadRequest, errorResponse{
			Error: "invalid JSON: " + err.Error(),
		})
		return
	}

	if msg := validateLogEntry(entry); msg != "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{
			Error: msg,
		})
		return
	}

	fmt.Printf("received: level=%s service=%s message=%s\n",
		entry.Level, entry.Service, entry.Message)

	writeJSON(w, http.StatusCreated, statusResponse{Status: "accepted"})
}

func validateLogEntry(entry LogEntry) string {
	if entry.Level == "" {
		return "level is required"
	}
	if entry.Message == "" {
		return "message is required"
	}
	if entry.Service == "" {
		return "service is required"
	}
	if entry.Timestamp.IsZero() {
		return "timestamp is required"
	}

	switch entry.Level {
	case "debug", "info", "warn", "error", "fatal":
		// valid
	default:
		return "level must be one of: debug, info, warn, error, fatal"
	}

	return ""
}
