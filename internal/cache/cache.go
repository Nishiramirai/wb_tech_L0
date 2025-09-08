package cache

import (
	"context"
	"sync"

	"orders-service/internal/app/model"
	"orders-service/internal/db"
)

type Cache struct {
	mu      sync.RWMutex
	orders  map[string]*model.Order
	db      *db.DB
	maxSize int
}

func NewCache(ctx context.Context, db *db.DB, maxSize int) (*Cache, error) {
	cache := &Cache{
		orders:  make(map[string]*model.Order, maxSize),
		db:      db,
		maxSize: maxSize,
	}

	if err := cache.PopulateFromDB(ctx); err != nil {
		return nil, err
	}

	return cache, nil
}

func (c *Cache) PopulateFromDB(ctx context.Context) error {
	orders, err := c.db.GetRecentOrders(ctx, c.maxSize)
	if err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	for _, order := range orders {
		c.orders[order.OrderUID] = order
	}

	return nil
}

func (c *Cache) AddOrder(order *model.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.orders) >= c.maxSize {
		for uid := range c.orders {
			delete(c.orders, uid)
			break
		}
	}
	c.orders[order.OrderUID] = order
}

func (c *Cache) GetOrder(orderUID string) (*model.Order, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	order, ok := c.orders[orderUID]
	return order, ok
}
