package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", handleHealth)
	mux.HandleFunc("POST /ingest", handleIngest)
	addr := ":4000"
	log.Printf("logline starting on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

func writeJSON(w http.ResponseWriter, status int, body string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	fmt.Fprint(w, body)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, `{"status": "ok"}`)
}

func handleIngest(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		writeJSON(w, http.StatusUnsupportedMediaType,
			`{"error": "content-type must be application/json"}`)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeJSON(w, http.StatusBadRequest,
			`{"error": "failed to read request body"}`)
		return
	}

	if len(body) == 0 {
		writeJSON(w, http.StatusBadRequest,
			`{"error": "request body must not be empty"}`)
		return
	}

	fmt.Println(string(body))
	writeJSON(w, http.StatusCreated, `{"status": "accepted"}`)
}
