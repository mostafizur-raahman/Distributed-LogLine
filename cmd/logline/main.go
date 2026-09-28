package main

import (
	"fmt"
	"log"
	"os"

	"logline/internal/config"
	"logline/internal/server"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "configuration error: %v\n", err)
		os.Exit(1)
	}
	srv := server.New(cfg)

	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("logline starting on %s (env=%s, log-level=%s)", addr, cfg.Env, cfg.LogLevel)

	if err := srv.Start(addr); err != nil {
		log.Fatal(err)
	}
}
