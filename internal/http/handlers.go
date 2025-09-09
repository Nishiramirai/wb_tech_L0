package http

import (
	"encoding/json"
	"net/http"

	"orders-service/internal/app/service"

	"go.uber.org/zap"
)

type Handlers struct {
	svc    *service.OrderService
	logger *zap.Logger
}

func NewHandlers(svc *service.OrderService, logger *zap.Logger) *Handlers {
	return &Handlers{
		svc:    svc,
		logger: logger,
	}
}

func (h *Handlers) orderHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	orderUID := r.URL.Path[len("/orders/"):]
	if orderUID == "" {
		http.Error(w, "Order UID is missing", http.StatusBadRequest)
		return
	}

	order, err := h.svc.GetOrder(r.Context(), orderUID)
	if err != nil {
		h.logger.Error("Order not found or database error", zap.Error(err), zap.String("order_uid", orderUID))
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(order); err != nil {
		h.logger.Error("Failed to encode response", zap.Error(err), zap.String("order_uid", orderUID))
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
