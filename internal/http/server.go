package http

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"
)

type Server struct {
	logger *slog.Logger
	server *http.Server
}

type Config struct {
	Address     string
	ReadTimeout time.Duration
}

func (c Config) Validate() error {
	if c.Address == "" {
		return errors.New("empty server address")
	}

	if c.ReadTimeout == 0 {
		return errors.New("empty read timeout")
	}

	return nil
}

func NewServer(config Config, logger *slog.Logger) *Server {
	mux := newRouter(logger)

	return &Server{
		logger: logger,
		server: &http.Server{ReadTimeout: config.ReadTimeout, Addr: config.Address, Handler: mux},
	}
}

func (s *Server) Start() {
	s.logger.Info("starting server", "address", s.server.Addr)
	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.logger.Error("cannot start server", slog.Any("error", err))
	}
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
