package service

import (
	"context"
	"fmt"

	"orders-service/internal/app/model"
	"orders-service/internal/cache"
	"orders-service/internal/db"
)

type OrderService struct {
	db    *db.DB
	cache *cache.Cache
}

func NewOrderService(db *db.DB, c *cache.Cache) *OrderService {
	return &OrderService{
		db:    db,
		cache: c,
	}
}

func (s *OrderService) SaveOrder(ctx context.Context, order *model.Order) error {
	if err := s.validateOrder(order); err != nil {
		return err
	}

	if err := s.db.SaveOrder(ctx, order); err != nil {
		return err
	}

	s.cache.AddOrder(order)

	return nil
}

func (s *OrderService) GetOrder(ctx context.Context, orderUID string) (*model.Order, error) {
	if order, ok := s.cache.GetOrder(orderUID); ok {
		return order, nil
	}

	order, err := s.db.GetOrder(ctx, orderUID)
	if err != nil {
		return nil, fmt.Errorf("getting order frob db: %w", err)
	}

	s.cache.AddOrder(order)

	return order, nil
}

// Возвращает ошибку, если в заказе нет какого-либо существенного поля
func (s *OrderService) validateOrder(order *model.Order) error {
	if order.OrderUID == "" {
		return fmt.Errorf("order_uid cannot be empty")
	}
	if order.TrackNumber == "" {
		return fmt.Errorf("track_number cannot be empty")
	}
	if order.Delivery.Name == "" {
		return fmt.Errorf("delivery name cannot be empty")
	}
	if order.Payment.Transaction == "" {
		return fmt.Errorf("payment transaction cannot be empty")
	}
	if len(order.Items) == 0 {
		return fmt.Errorf("items list cannot be empty")
	}

	if order.Payment.Amount <= 0 {
		return fmt.Errorf("payment amount must be greater than zero")
	}
	if order.Payment.GoodsTotal <= 0 {
		return fmt.Errorf("payment goods_total must be greater than zero")
	}

	for _, item := range order.Items {
		if item.ChrtID == 0 {
			return fmt.Errorf("item chrt_id cannot be zero")
		}
		if item.Price <= 0 {
			return fmt.Errorf("item price must be greater than zero")
		}
	}

	return nil
}
