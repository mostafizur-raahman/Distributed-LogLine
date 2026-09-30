package main

import (
	"fmt"
	"log/slog"
	"logline/internal/config"
	"logline/internal/server"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}
	srv := server.New(cfg)

	addr := fmt.Sprintf(":%d", cfg.Port)

	slog.Info("logline starting",
		slog.String("addr", addr),
		slog.String("env", cfg.Env),
		slog.String("log-level", cfg.LogLevel),
	)
	return srv.Start(addr)
}
