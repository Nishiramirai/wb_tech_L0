package http

import (
	"fmt"
	"net/http"

	"orders-service/internal/app/service"

	"go.uber.org/zap"
)

type Server struct {
	handlers *Handlers
	logger   *zap.Logger
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

	s.logger.Info("Starting HTTP server", zap.Int("port", port))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), mux); err != nil {
		s.logger.Fatal("Failed to start HTTP server", zap.Error(err))
	}
}
