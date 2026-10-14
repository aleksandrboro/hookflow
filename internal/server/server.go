package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/aleksandrboro/hookflow/internal/config"
)

type Server struct {
	log  *slog.Logger
	addr string
	srv  *http.Server
}

func NewServer(log *slog.Logger, cfg config.ServerCfg, handler http.Handler) (*Server, error) {
	if handler == nil {
		return nil, errors.New("HTTP handler is required")
	}

	addr := strings.TrimSpace(cfg.Addr)
	if addr == "" {
		return nil, errors.New("server address is required")
	}

	if log == nil {
		log = slog.Default()
	}

	return &Server{
		log:  log,
		addr: cfg.Addr,
		srv: &http.Server{
			Addr:         cfg.Addr,
			WriteTimeout: cfg.WriteTimeout,
			ReadTimeout:  cfg.ReadTimeout,
			Handler:      handler,
		},
	}, nil
}

func (s *Server) Start() error {
	if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("error starting server on %q: %w", s.srv.Addr, err)
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutdown error: %w", err)
	}

	s.log.Info("HTTP server stopped")

	return nil
}

func (s *Server) Close() error {
	return s.srv.Close()
}
