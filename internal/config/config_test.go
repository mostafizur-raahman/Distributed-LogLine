package config

import (
	"testing"
)

func TestDefaults(t *testing.T) {
	cfg := Config{
		Port:     4000,
		Env:      "development",
		LogLevel: "info",
	}

	if err := cfg.validate(); err != nil {
		t.Fatalf("defaults should be valid: %v", err)
	}
}

func TestValidatePort(t *testing.T) {
	tests := []struct {
		name    string
		port    int
		wantErr bool
	}{
		{"valid port", 4000, false},
		{"port 1", 1, false},
		{"port 65535", 65535, false},
		{"port 0", 0, true},
		{"negative port", -1, true},
		{"port too high", 65536, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Config{Port: tt.port, Env: "development", LogLevel: "info"}
			err := cfg.validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("port=%d: wantErr=%v, got %v", tt.port, tt.wantErr, err)
			}
		})
	}
}

func TestValidateEnv(t *testing.T) {
	tests := []struct {
		name    string
		env     string
		wantErr bool
	}{
		{"development", "development", false},
		{"production", "production", false},
		{"staging", "staging", true},
		{"empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Config{Port: 4000, Env: tt.env, LogLevel: "info"}
			err := cfg.validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("env=%q: wantErr=%v, got %v", tt.env, tt.wantErr, err)
			}
		})
	}
}

func TestValidateLogLevel(t *testing.T) {
	tests := []struct {
		name    string
		level   string
		wantErr bool
	}{
		{"debug", "debug", false},
		{"info", "info", false},
		{"warn", "warn", false},
		{"error", "error", false},
		{"trace", "trace", true},
		{"empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Config{Port: 4000, Env: "development", LogLevel: tt.level}
			err := cfg.validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("level=%q: wantErr=%v, got %v", tt.level, tt.wantErr, err)
			}
		})
	}
}

func TestLoadEnvOverridesDefaults(t *testing.T) {
	t.Setenv("PORT", "9090")
	t.Setenv("APP_ENV", "production")
	t.Setenv("LOG_LEVEL", "debug")

	cfg := Config{
		Port:     4000,
		Env:      "development",
		LogLevel: "info",
	}

	cfg.loadEnv()

	if cfg.Port != 9090 {
		t.Errorf("expected port 9090, got %d", cfg.Port)
	}
	if cfg.Env != "production" {
		t.Errorf("expected env production, got %q", cfg.Env)
	}
	if cfg.LogLevel != "debug" {
		t.Errorf("expected log level debug, got %q", cfg.LogLevel)
	}
}

func TestLoadEnvIgnoresInvalidPort(t *testing.T) {
	t.Setenv("PORT", "not-a-number")

	cfg := Config{Port: 4000, Env: "development", LogLevel: "info"}
	cfg.loadEnv()

	if cfg.Port != 4000 {
		t.Errorf("expected default port 4000, got %d", cfg.Port)
	}
}
