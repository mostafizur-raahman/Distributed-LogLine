package server

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"logline/internal/domain"
)

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, domain.StatusResponse{Status: "ok"})
}

func generateTraceID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "unknown"
	}
	return hex.EncodeToString(b)
}
func (s *Server) handleIngest(w http.ResponseWriter, r *http.Request) {
	traceID := generateTraceID()
	logger := s.logger.With(slog.String("trace_id", traceID))

	if r.Header.Get("Content-Type") != "application/json" {
		writeJSON(w, http.StatusUnsupportedMediaType, domain.ErrorResponse{
			Error: "content-type must be application/json",
		})
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)

	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var entry domain.LogEntry
	if err := decoder.Decode(&entry); err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			logger.Warn("unsupported content type",
				slog.String("content_type", r.Header.Get("Content-Type")),
			)
			writeJSON(w, http.StatusRequestEntityTooLarge, domain.ErrorResponse{
				Error: "request body too large",
			})
			return
		}

		logger.Warn("invalid JSON in request",
			slog.String("error", err.Error()),
		)

		writeJSON(w, http.StatusBadRequest, domain.ErrorResponse{
			Error: "invalid JSON: " + err.Error(),
		})
		return
	}

	if msg := domain.ValidateLogEntry(entry); msg != "" {
		logger.Warn("validation failed",
			slog.String("reason", msg),
		)
		writeJSON(w, http.StatusBadRequest, domain.ErrorResponse{
			Error: msg,
		})
		return
	}
	logger.Info("log entry accepted",
		slog.String("service", entry.Service),
		slog.String("entry_level", entry.Level),
		slog.Int64("content_length", r.ContentLength),
	)

	writeJSON(w, http.StatusCreated, domain.StatusResponse{Status: "accepted"})
}

const maxBodySize = 1_048_576 // 1 MB

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
