package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	entries map[string]cacheEntry
	mu      sync.RWMutex
	ticker  *time.Ticker
	done    chan struct{}
}

func NewCache(interval time.Duration) *Cache {
	// Initialize map
	initEntries := make(map[string]cacheEntry)

	res := Cache{
		entries: initEntries,
		done:    make(chan struct{}),
	}

	// Start auto cleanup
	res.ticker = time.NewTicker(interval)

	go func() {
		for {
			select {
			case <-res.ticker.C:
				res.reapLoop(interval)
			case <-res.done:
				return
			}
		}
	}()

	return &res
}

// Stop terminates the background reap goroutine and releases resources.
func (c *Cache) Stop() {
	c.ticker.Stop()
	close(c.done)
}

// Create or update a cache entry
func (c *Cache) Add(key string, myVal []byte) {

	// take lock
	c.mu.Lock()
	defer c.mu.Unlock()

	ce := cacheEntry{
		createdAt: time.Now(),
		val:       myVal,
	}

	c.entries[key] = ce
}

// Retrieve a cache entry
func (c *Cache) Get(key string) ([]byte, bool) {

	// Take read lock
	c.mu.RLock()
	defer c.mu.RUnlock()

	res, ok := c.entries[key]
	if !ok {
		return nil, false
	}

	return res.val, true

}

// Automatically clean old cache entries
func (c *Cache) reapLoop(interval time.Duration) {

	// take lock
	c.mu.Lock()
	defer c.mu.Unlock()

	expireTime := time.Now().Add(-interval)

	for k, v := range c.entries {
		createTime := v.createdAt

		if createTime.Before(expireTime) {
			// Remove old entry
			delete(c.entries, k)
		}
	}
}
