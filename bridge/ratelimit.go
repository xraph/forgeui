package bridge

import (
	"sync"
	"time"
)

// RateLimiter implements token bucket rate limiting
type RateLimiter struct {
	mu      sync.RWMutex
	buckets map[string]*bucket
	rate    int           // tokens per minute
	burst   int           // max tokens
	ttl     time.Duration // bucket cleanup interval
}

// bucket holds tokens for a specific key
type bucket struct {
	tokens    float64
	lastCheck time.Time
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(rate, burst int) *RateLimiter {
	rl := &RateLimiter{
		buckets: make(map[string]*bucket),
		rate:    rate,
		burst:   burst,
		ttl:     5 * time.Minute,
	}

	// Start cleanup goroutine
	go rl.cleanup()

	return rl
}

// Allow checks if a request is allowed for the given key
func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	b, exists := rl.buckets[key]

	if !exists {
		// Create new bucket
		b = &bucket{
			tokens:    float64(rl.burst - 1),
			lastCheck: now,
		}
		rl.buckets[key] = b

		return true
	}

	// Calculate tokens to add based on time elapsed
	elapsed := now.Sub(b.lastCheck)
	tokensToAdd := elapsed.Seconds() * float64(rl.rate) / 60.0

	b.tokens += tokensToAdd
	if b.tokens > float64(rl.burst) {
		b.tokens = float64(rl.burst)
	}

	b.lastCheck = now

	// Check if we have tokens available
	if b.tokens >= 1 {
		b.tokens--
		return true
	}

	return false
}

// Reset resets the rate limit for a key
func (rl *RateLimiter) Reset(key string) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	delete(rl.buckets, key)
}

// cleanup removes stale buckets periodically
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()

		now := time.Now()

		for key, b := range rl.buckets {
			if now.Sub(b.lastCheck) > rl.ttl {
				delete(rl.buckets, key)
			}
		}

		rl.mu.Unlock()
	}
}

// Count returns the number of active buckets
func (rl *RateLimiter) Count() int {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	return len(rl.buckets)
}

// Remaining returns the number of tokens remaining for a key
func (rl *RateLimiter) Remaining(key string) int {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	b, exists := rl.buckets[key]
	if !exists {
		return rl.burst
	}

	return int(b.tokens)
}
