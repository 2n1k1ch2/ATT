package http

import (
	"AvitoTestTask/internal/config"
	"context"
	"errors"
	"net/http"
	"time"
)

var (
	ErrServerClosed = errors.New("server closed")
)

type Server struct {
	httpServer *http.Server
	cfg        *config.Config
}

func NewServer(cfg *config.Config, handler http.Handler) *Server {
	return &Server{
		cfg: cfg,
		httpServer: &http.Server{
			Addr:         ":" + cfg.AppPort,
			Handler:      handler,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  120 * time.Second,
		},
	}
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
