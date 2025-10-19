package cache

import (
	"sync"
	"time"
)

type unixNano int64

type cacheItem[T any] struct {
	value      T
	expiration unixNano
}

type MemoryCache[T any] struct {
	data            map[string]cacheItem[T]
	cleanupInterval time.Duration
	mu              sync.RWMutex
}

func NewMemoryCache[T any](cleanupInterval time.Duration) *MemoryCache[T] {
	if cleanupInterval == 0 {
		cleanupInterval = 1 * time.Minute
	}

	m := &MemoryCache[T]{
		data:            make(map[string]cacheItem[T], 128),
		cleanupInterval: cleanupInterval,
		mu:              sync.RWMutex{},
	}

	go m.cleanup()

	return m
}

// Delete implements Cache.
func (m *MemoryCache[T]) Delete(key string) {
	m.mu.Lock()
	delete(m.data, key)
	m.mu.Unlock()
}

// Get implements Cache.
func (m *MemoryCache[T]) Get(key string) (T, bool) {
	m.mu.RLock()
	item, ok := m.data[key]
	if !ok {
		return *new(T), false
	}
	m.mu.RUnlock()

	if item.expiration >= unixNano(time.Now().UnixNano()) {
		m.mu.Lock()
		delete(m.data, key)
		m.mu.Unlock()

		return *new(T), false
	}

	m.mu.RLock()
	defer m.mu.RUnlock()
	return item.value, ok
}

// Set implements Cache.
func (m *MemoryCache[T]) Set(key string, value T, ttl time.Duration) {
	var exp unixNano
	if ttl > 0 {
		exp = unixNano(time.Now().Add(ttl).UnixNano())
	}

	m.mu.Lock()
	m.data[key] = cacheItem[T]{value: value, expiration: exp}
	m.mu.Unlock()
}

func (c *MemoryCache[T]) cleanup() {
	ticker := time.NewTicker(c.cleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		now := unixNano(time.Now().UnixNano())
		c.mu.Lock()
		for k, v := range c.data {
			if v.expiration > 0 && now > v.expiration {
				delete(c.data, k)
			}
		}
		c.mu.Unlock()
	}
}

var _ Cache[any] = &MemoryCache[any]{}
