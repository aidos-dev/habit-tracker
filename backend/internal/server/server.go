package server

import (
	"context"
	"log"
	"net/http"

	"github.com/aidos-dev/habit-tracker/backend/internal/config"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(cfg *config.Config, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           ":" + cfg.HTTPServer.Port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20, // 1 MB
		ReadTimeout:    cfg.HTTPServer.Timeout,
		WriteTimeout:   cfg.HTTPServer.Timeout,
		IdleTimeout:    cfg.HTTPServer.IdleTimeout,
	}

	log.Printf("backend server started and listening on port: %v", cfg.HTTPServer.Port)

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
