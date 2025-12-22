package bridge

import (
	"testing"
	"time"
)

func TestMemoryCache_GetSet(t *testing.T) {
	cache := NewMemoryCache()

	// Set a value
	cache.Set("key1", "value1", 1*time.Minute)

	// Get the value
	val, ok := cache.Get("key1")
	if !ok {
		t.Error("Get() should return true for existing key")
	}

	if val != "value1" {
		t.Errorf("Get() = %v, want value1", val)
	}
}

func TestMemoryCache_GetNonExistent(t *testing.T) {
	cache := NewMemoryCache()

	_, ok := cache.Get("nonexistent")
	if ok {
		t.Error("Get() should return false for non-existent key")
	}
}

func TestMemoryCache_Expiration(t *testing.T) {
	cache := NewMemoryCache()

	// Set with very short TTL
	cache.Set("key1", "value1", 100*time.Millisecond)

	// Should be available immediately
	_, ok := cache.Get("key1")
	if !ok {
		t.Error("Get() should return true immediately after Set")
	}

	// Wait for expiration
	time.Sleep(150 * time.Millisecond)

	// Should be expired
	_, ok = cache.Get("key1")
	if ok {
		t.Error("Get() should return false after expiration")
	}
}

func TestMemoryCache_Delete(t *testing.T) {
	cache := NewMemoryCache()

	cache.Set("key1", "value1", 1*time.Minute)
	cache.Delete("key1")

	_, ok := cache.Get("key1")
	if ok {
		t.Error("Get() should return false after Delete")
	}
}

func TestMemoryCache_Clear(t *testing.T) {
	cache := NewMemoryCache()

	cache.Set("key1", "value1", 1*time.Minute)
	cache.Set("key2", "value2", 1*time.Minute)

	cache.Clear()

	if cache.Count() != 0 {
		t.Errorf("Count() after Clear = %d, want 0", cache.Count())
	}
}

func TestGenerateCacheKey(t *testing.T) {
	key1 := generateCacheKey("func1", []byte(`{"param":"value"}`))
	key2 := generateCacheKey("func1", []byte(`{"param":"value"}`))
	key3 := generateCacheKey("func1", []byte(`{"param":"different"}`))

	// Same inputs should generate same key
	if key1 != key2 {
		t.Error("Same inputs should generate same cache key")
	}

	// Different inputs should generate different keys
	if key1 == key3 {
		t.Error("Different inputs should generate different cache keys")
	}
}

