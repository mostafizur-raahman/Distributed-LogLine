package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port     int
	Env      string
	LogLevel string
}

func LoadConfig() (Config, error) {
	cfg := Config{
		Port:     4000,
		Env:      "development",
		LogLevel: "info",
	}

	cfg.loadEnv()

	if err := cfg.loadFlags(); err != nil {
		return Config{}, err
	}

	if err := cfg.validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (c *Config) loadEnv() {
	if v := os.Getenv("PORT"); v != "" {
		if p, err := strconv.Atoi(v); err == nil {
			c.Port = p
		}
	}

	if v := os.Getenv("APP_ENV"); v != "" {
		c.Env = v
	}

	if v := os.Getenv("LOG_LEVEL"); v != "" {
		c.LogLevel = v
	}
}

func (c *Config) loadFlags() error {
	fs := flag.NewFlagSet("logline", flag.ContinueOnError)

	fs.IntVar(&c.Port, "port", c.Port, "HTTP server port")
	fs.StringVar(&c.Env, "env", c.Env, "Environment (development, production)")
	fs.StringVar(&c.LogLevel, "log-level", c.LogLevel, "Log level (debug, info, warn, error)")

	return fs.Parse(os.Args[1:])
}

func (c *Config) validate() error {
	if c.Port < 1 || c.Port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535, got %d", c.Port)
	}

	switch c.Env {
	case "development", "production":
		// valid
	default:
		return fmt.Errorf("env must be 'development' or 'production', got %q", c.Env)
	}

	switch c.LogLevel {
	case "debug", "info", "warn", "error":
		// valid
	default:
		return fmt.Errorf("log-level must be one of: debug, info, warn, error — got %q", c.LogLevel)
	}

	return nil
}
