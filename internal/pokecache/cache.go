package pokecache

import (
	"time"
)

type Cache struct {
	cache map[string]cacheEntry
	//mu       sync.Mutex
	interval time.Duration
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(seconds time.Duration) Cache {
	nC := Cache{
		cache:    make(map[string]cacheEntry),
		interval: time.Second * seconds,
	}
	go nC.reapLoop()
	return nC
}

func (c *Cache) Add(key string, val []byte) {
	c.cache[key] = cacheEntry{
		val:       val,
		createdAt: time.Now().UTC(),
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	cacheEntry, ok := c.cache[key]
	return cacheEntry.val, ok
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(time.Millisecond * 500)
	for range ticker.C {
		c.reap()
	}
}

func (c *Cache) reap() {
	cutOffTime := time.Now().UTC().Add(-c.interval)
	for key, v := range c.cache {
		if v.createdAt.Before(cutOffTime) {
			delete(c.cache, key)
		}
	}
}
