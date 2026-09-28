package domain

import (
	"testing"
	"time"
)

func TestValidateLogEntry_MissingFields(t *testing.T) {
	tests := []struct {
		name  string
		entry LogEntry
		want  string
	}{
		{"missing level", LogEntry{Message: "ok", Service: "api", Timestamp: time.Now()}, "level is required"},
		{"missing message", LogEntry{Level: "info", Service: "api", Timestamp: time.Now()}, "message is required"},
		{"missing service", LogEntry{Level: "info", Message: "ok", Timestamp: time.Now()}, "service is required"},
		{"missing timestamp", LogEntry{Level: "info", Message: "ok", Service: "api"}, "timestamp is required"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateLogEntry(tt.entry); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestValidateLogEntry_InvalidLevel(t *testing.T) {
	entry := LogEntry{
		Level:     "critical",
		Message:   "test",
		Service:   "api",
		Timestamp: time.Now(),
	}

	want := "level must be one of: debug, info, warn, error, fatal"
	if got := ValidateLogEntry(entry); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestValidateLogEntry_Valid(t *testing.T) {
	entry := LogEntry{
		Level:     "error",
		Message:   "connection refused",
		Service:   "auth",
		Timestamp: time.Now(),
	}

	if got := ValidateLogEntry(entry); got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}
