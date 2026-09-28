package domain

import "time"

type LogEntry struct {
	Level     string         `json:"level"`
	Message   string         `json:"message"`
	Service   string         `json:"service"`
	Timestamp time.Time      `json:"timestamp"`
	Data      map[string]any `json:"data,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type StatusResponse struct {
	Status string `json:"status"`
}

func ValidateLogEntry(entry LogEntry) string {
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
