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
    mutex     sync.Mutex
    entries   map[string]cacheEntry
    interval  time.Duration
}

func NewCache(interval time.Duration) *Cache {
    c := &Cache{
        entries:  make(map[string]cacheEntry),
        interval: interval,
    }

    go c.reapLoop()
    return c
}

func (c *Cache) Add(key string, val []byte) {
    c.mutex.Lock()
    defer c.mutex.Unlock()

    c.entries[key] = cacheEntry{
        createdAt: time.Now(),
        val:       val,
    }
}

func (c *Cache) Get(key string) ([]byte, bool) {
    c.mutex.Lock()
    defer c.mutex.Unlock()
    entry, exists := c.entries[key]
    if exists {
        return entry.val, true
    }
    return nil, false
}

func (c *Cache) reapLoop() {
    ticker := time.NewTicker(c.interval)
    defer ticker.Stop()
    for {
        <-ticker.C
        c.mutex.Lock()
        for key, entry := range c.entries {
            if time.Now().After(entry.createdAt.Add(c.interval)) {
                delete(c.entries, key)
            }
        }
        c.mutex.Unlock()
    }
}
