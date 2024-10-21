package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"forum/pkg/logger"
)

type Server struct {
	httpServer *http.Server
	logger     logger.Logger
}

func NewServer(l *logger.Logger) *Server {
	return &Server{
		logger: *l,
	}
}

func (s *Server) RunServer(port string, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	s.logger.Info("starting server at http://localhost%s\n", s.httpServer.Addr)
	if err := s.httpServer.ListenAndServe(); err != nil {
		return fmt.Errorf("http server listen and serve: %w", err)
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	fmt.Println("gracefully shutting down server...")
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("http server shutdown: %w", err)
	}
	return nil
}
