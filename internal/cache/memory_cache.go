package cache

import (
	"sync"
	"time"
)

type CacheItem struct {
	StatusCode int
	Header     map[string][]string
	Body       []byte
	Expiration time.Time
}

type MemoryCache struct {
	items map[string]CacheItem
	mu    sync.RWMutex
	ttl   time.Duration
}

func NewMemoryCache(ttl time.Duration) *MemoryCache {
	c := &MemoryCache{
		items: make(map[string]CacheItem),
		ttl:   ttl,
	}
	go c.cleanup()
	return c
}

func (c *MemoryCache) Get(key string) (CacheItem, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, found := c.items[key]
	if !found || time.Now().After(item.Expiration) {
		return CacheItem{}, false
	}
	return item, true
}

func (c *MemoryCache) Set(key string, status int, header map[string][]string, body []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = CacheItem{
		StatusCode: status,
		Header:     header,
		Body:       body,
		Expiration: time.Now().Add(c.ttl),
	}
}

func (c *MemoryCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]CacheItem)
}

func (c *MemoryCache) cleanup() {
	for {
		time.Sleep(c.ttl)
		c.mu.Lock()
		for k, v := range c.items {
			if time.Now().After(v.Expiration) {
				delete(c.items, k)
			}
		}
		c.mu.Unlock()
	}
}
