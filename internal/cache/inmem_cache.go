package cache

import (
	"context"
	"sync"
	"time"
)

type InMemoryCache struct {
	data map[string]cacheItem
	mu   sync.RWMutex
}

type cacheItem struct {
	value interface{}
	ttl   time.Time
}

func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{
		data: make(map[string]cacheItem),
	}
}

// Set добавляет элемент в кэш с возможным временем истечения.
func (c *InMemoryCache) Set(_ context.Context, key string, value interface{}, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Если ttl больше 0, устанавливаем время жизни для элемента
	var expiryTime time.Time
	if ttl > 0 {
		expiryTime = time.Now().Add(ttl)
	}

	c.data[key] = cacheItem{
		value: value,
		ttl:   expiryTime,
	}
	return nil
}

// Get возвращает элемент из кэша, если он существует и не истек.
func (c *InMemoryCache) Get(_ context.Context, key string, dest interface{}) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.data[key]
	if !exists {
		return nil // Возвращаем nil, если нет такого ключа
	}

	// Проверяем TTL
	if !item.ttl.IsZero() && time.Now().After(item.ttl) {
		delete(c.data, key) // Удаляем элемент, если он истек
		return nil
	}

	// Преобразуем значение в нужный тип
	*dest.(*interface{}) = item.value
	return nil
}

// Delete удаляет элемент из кэша.
func (c *InMemoryCache) Delete(_ context.Context, key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.data, key)
	return nil
}
