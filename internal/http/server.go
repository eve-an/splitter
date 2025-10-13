package http

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
)

type Server struct {
	logger *slog.Logger
	server *http.Server
}

func NewServer(address string, logger *slog.Logger) *Server {
	mux := newRouter(logger)

	return &Server{
		logger: logger,
		server: &http.Server{Addr: address, Handler: mux},
	}
}

func (s *Server) Start() {
	go func() {
		s.logger.Info("starting server", "address", s.server.Addr)
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error("cannot start server", slog.Any("error", err))
		}
	}()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
