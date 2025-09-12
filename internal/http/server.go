package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"orders-service/internal/app/service"

	"go.uber.org/zap"
)

type Server struct {
	handlers   *Handlers
	logger     *zap.Logger
	httpServer *http.Server
}

func NewServer(svc *service.OrderService, logger *zap.Logger) (*Server, error) {
	return &Server{
		handlers: NewHandlers(svc, logger),
		logger:   logger,
	}, nil
}

func (s *Server) Start(port int) {
	mux := http.NewServeMux()

	mux.HandleFunc("/orders/", s.handlers.orderHandler)

	mux.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./web"))))

	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	s.logger.Info("Starting HTTP server", zap.Int("port", port))
	if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.logger.Fatal("Failed to start HTTP server", zap.Error(err))
	}
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Shutting down HTTP server...")
	return s.httpServer.Shutdown(ctx)
}
