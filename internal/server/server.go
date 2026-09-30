package server

import (
	"log/slog"
	"logline/internal/config"
	"net/http"
	"os"
	"time"
)

type Server struct {
	cfg    config.Config
	logger *slog.Logger
	mux    *http.ServeMux
}

func New(cfg config.Config) *Server {
	logger := newLogger(cfg.Env)

	s := &Server{
		mux:    http.NewServeMux(),
		cfg:    cfg,
		logger: logger,
	}

	s.registerRoutes()

	return s
}

func (s *Server) Start(addr string) error {
	httpSrv := &http.Server{
		Addr:         addr,
		Handler:      s,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	return httpSrv.ListenAndServe()
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) registerRoutes() {
	s.mux.HandleFunc("GET /health", s.handleHealth)
	s.mux.HandleFunc("POST /ingest", s.handleIngest)
}

func newLogger(env string) *slog.Logger {
	var handler slog.Handler

	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}

	switch env {
	case "production":
		handler = slog.NewJSONHandler(os.Stdout, opts)
	default:
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)

	return logger
}
