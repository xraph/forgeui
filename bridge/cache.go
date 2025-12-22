package bridge

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// Cache provides result caching
type Cache interface {
	// Get retrieves a cached value
	Get(key string) (any, bool)

	// Set stores a value with TTL
	Set(key string, value any, ttl time.Duration)

	// Delete removes a cached value
	Delete(key string)

	// Clear removes all cached values
	Clear()
}

// MemoryCache is an in-memory cache implementation
type MemoryCache struct {
	mu    sync.RWMutex
	items map[string]*cacheItem
}

// cacheItem holds a cached value with expiration
type cacheItem struct {
	value     any
	expiresAt time.Time
}

// NewMemoryCache creates a new memory cache
func NewMemoryCache() *MemoryCache {
	cache := &MemoryCache{
		items: make(map[string]*cacheItem),
	}

	// Start cleanup goroutine
	go cache.cleanup()

	return cache
}

// Get retrieves a cached value
func (c *MemoryCache) Get(key string) (any, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.items[key]
	if !exists {
		return nil, false
	}

	// Check if expired
	if time.Now().After(item.expiresAt) {
		return nil, false
	}

	return item.value, true
}

// Set stores a value with TTL
func (c *MemoryCache) Set(key string, value any, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = &cacheItem{
		value:     value,
		expiresAt: time.Now().Add(ttl),
	}
}

// Delete removes a cached value
func (c *MemoryCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
}

// Clear removes all cached values
func (c *MemoryCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]*cacheItem)
}

// cleanup removes expired items periodically
func (c *MemoryCache) cleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		now := time.Now()

		for key, item := range c.items {
			if now.After(item.expiresAt) {
				delete(c.items, key)
			}
		}

		c.mu.Unlock()
	}
}

// Count returns the number of cached items
func (c *MemoryCache) Count() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.items)
}

// generateCacheKey generates a cache key from function name and params
func generateCacheKey(funcName string, params json.RawMessage) string {
	h := sha256.New()
	h.Write([]byte(funcName))
	h.Write(params)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// WithBridgeCache adds caching to a bridge
func WithBridgeCache(b *Bridge, cache Cache) *Bridge {
	// Store cache in bridge (would need to add cache field to Bridge struct)
	// This is a helper function for future extension
	return b
}

