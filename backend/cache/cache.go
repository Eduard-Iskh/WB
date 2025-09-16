package cache

import (
	"sync"
	"wildberies/L0/backend/domain"
)

// getBiIDCache(id){
// 	cache :=map([])
// 	if id  in cache:
// 	retiurmn cahe=
// 	data, err := gGEtById(id)
// 	cahce[id] = data
// 	if er
// 	return data, err
// }

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

// можно убрать
func (c *Cache) GetAll() map[string]domain.Order {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Создаем копию для безопасного доступа
	result := make(map[string]domain.Order)
	for k, v := range c.orders {
		result[k] = v
	}
	return result
}
