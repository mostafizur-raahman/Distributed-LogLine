package server

import (
	"logline/internal/config"
	"net/http"
	"time"
)

type Server struct {
	mux *http.ServeMux
	cfg config.Config
}

func New(cfg config.Config) *Server {
	s := &Server{
		mux: http.NewServeMux(),
		cfg: cfg,
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
