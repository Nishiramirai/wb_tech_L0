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

// Создает кеш с указанным лимитом и заполняет его данными из БД
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

// Заполняет кэш данными из БД
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

// Добавляет заказ в кэш
func (c *Cache) AddOrder(order *model.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.orders) >= c.maxSize {
		c.EvictElement()
	}
	c.orders[order.OrderUID] = order
}

// Возвращает структуру с заказом и bool, было ли что-то в кэше по переданному uid
func (c *Cache) GetOrder(orderUID string) (*model.Order, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	order, ok := c.orders[orderUID]
	return order, ok
}

// Удаляет элемент из кэша
func (c *Cache) EvictElement() {
	// Лучше конечно сделать LRU кэш, но пока как есть
	for uid := range c.orders {
		delete(c.orders, uid)
		break
	}
}
