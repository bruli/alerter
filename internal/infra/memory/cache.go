package memory

import (
	"context"
	"sync"
	"time"
)

type Cache struct {
	m map[string]time.Time
	sync.RWMutex
	ttl time.Duration
}

func (c *Cache) Remove(v string) {
	delete(c.m, v)
}

func (c *Cache) RunTTL(ctx context.Context, refresh time.Duration) {
	tc := time.NewTicker(refresh)
	for {
		select {
		case <-ctx.Done():
			return
		case <-tc.C:
			c.Lock()
			for k, v := range c.m {
				if time.Now().After(v) {
					c.Remove(k)
				}
			}
			c.Unlock()
		}
	}
}

func (c *Cache) Set(v string) {
	c.Lock()
	defer c.Unlock()
	c.m[v] = time.Now().UTC()
}

func (c *Cache) Exists(v string) bool {
	c.RLock()
	defer c.RUnlock()
	_, ok := c.m[v]
	return ok
}

func NewCache(ttl time.Duration) *Cache {
	m := make(map[string]time.Time)
	return &Cache{ttl: ttl, m: m}
}
