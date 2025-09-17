package cache

import (
	"sync"
	domain "wildberies/L0/backend/internal/entify"
)

type Cache struct {
	mu     sync.RWMutex
	orders map[string]domain.Order
}

func NewCache() *Cache {
	return &Cache{
		orders: make(map[string]domain.Order),
	}
}

func (c *Cache) Set(id string, order domain.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.orders[id] = order
}

func (c *Cache) Get(id string) (domain.Order, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	order, exists := c.orders[id]
	return order, exists
}
