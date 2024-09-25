package main

import ("time"; "sync")

type Cache struct {
	cache map[string]cacheEntry
	mu    *sync.Mutex
	interval time.Duration
}
type cacheEntry struct {
	createdAt time.Time
	val interface{}
}

func newCache(interval time.Duration) *Cache {
	cache:= &Cache{
		cache : make(map[string]cacheEntry),
        mu: &sync.Mutex{},
		interval: interval,
	}
	
	go cache.reapLoop()
	return cache
}

func (c *Cache) Add(key string, val interface{}) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.cache[key] = cacheEntry{
        createdAt: time.Now(),
        val:       val,
    }
}

func (c *Cache) Get(key string) (interface{}, bool) {
    c.mu.Lock()
    defer c.mu.Unlock()
    cache, found := c.cache[key]
    if !found {
        return nil, false
    }
    return cache.val, true
}

func (c *Cache) reapLoop() {
    ticker := time.NewTicker(c.interval)
    defer ticker.Stop()

    for {
        <-ticker.C
        c.mu.Lock()
        for key, cacheEntry := range c.cache {
            if time.Since(cacheEntry.createdAt) > c.interval {
                delete(c.cache, key)
            }
        }
        c.mu.Unlock()
    }
}