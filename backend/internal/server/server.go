package server

import (
	"context"
	"net/http"

	"github.com/aidos-dev/habit-tracker/backend/internal/config"
	"golang.org/x/exp/slog"
)

type Server struct {
	httpServer *http.Server
	log        *slog.Logger
}

func (s *Server) Run(cfg *config.Config, log *slog.Logger, handler http.Handler) error {
	s.log = log

	s.httpServer = &http.Server{
		Addr:           ":" + cfg.HTTPServer.Port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20, // 1 MB
		ReadTimeout:    cfg.HTTPServer.Timeout,
		WriteTimeout:   cfg.HTTPServer.Timeout,
		IdleTimeout:    cfg.HTTPServer.IdleTimeout,
	}

	s.log.Info("backend server started and listening", slog.String("port", cfg.HTTPServer.Port))

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
