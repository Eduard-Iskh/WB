package cache

import (
	"sync"
	domain "wildberies/L0/backend/internal/domain"
)

type Cache struct {
	mu        sync.RWMutex
	orders    map[string]domain.Order
	ordersKey []string
	maxItems  int
}

func NewCache(maxItems int) *Cache {

	if maxItems <= 0 {
		maxItems = 100
	}
	return &Cache{
		orders:    make(map[string]domain.Order),
		ordersKey: make([]string, 0),
		maxItems:  maxItems,
	}
}

func (c *Cache) Set(id string, order domain.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.orders[id]; exists {
		c.orders[id] = order
		return
	}

	if len(c.ordersKey) >= c.maxItems {
		orderId := c.ordersKey[0]
		c.ordersKey = c.ordersKey[1:]
		delete(c.orders, orderId)
	}
	c.orders[id] = order
	c.ordersKey = append(c.ordersKey, id)
}

func (c *Cache) Get(id string) (domain.Order, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	order, exists := c.orders[id]
	return order, exists
}
