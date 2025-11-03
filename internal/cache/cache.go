package cache

import (
	"sync"
	"time"
)

type item struct {
	value      interface{}
	expiration int64
}

type Cache struct {
	items    map[string]*item
	mu       sync.RWMutex
	cleanup  *time.Ticker
	defaultTTL time.Duration
}

func New(defaultTTL, cleanupInterval time.Duration) *Cache {
	c := &Cache{
		items:      make(map[string]*item),
		defaultTTL: defaultTTL,
		cleanup:     time.NewTicker(cleanupInterval),
	}

	go c.startCleanup()
	return c
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.items[key]
	if !found {
		return nil, false
	}

	if time.Now().UnixNano() > item.expiration {
		delete(c.items, key)
		return nil, false
	}

	return item.value, true
}

func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ttl == 0 {
		ttl = c.defaultTTL
	}

	c.items[key] = &item{
		value:      value,
		expiration: time.Now().Add(ttl).UnixNano(),
	}
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]*item)
}

func (c *Cache) Stats() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return map[string]interface{}{
		"items":      len(c.items),
		"default_ttl": c.defaultTTL.String(),
	}
}

func (c *Cache) startCleanup() {
	for range c.cleanup.C {
		c.mu.Lock()
		now := time.Now().UnixNano()
		for key, item := range c.items {
			if now > item.expiration {
				delete(c.items, key)
			}
		}
		c.mu.Unlock()
	}
}
