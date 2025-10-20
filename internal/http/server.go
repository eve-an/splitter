package http

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/eve-an/splitter/internal/config"
)

type Server struct {
	logger *slog.Logger
	server *http.Server
}

func NewServer(
	config config.Server,
	logger *slog.Logger,
	handler http.Handler,
) *Server {
	return &Server{
		logger: logger,
		server: &http.Server{ReadTimeout: config.ReadTimeout, Addr: config.Address, Handler: handler},
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
