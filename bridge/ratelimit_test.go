package bridge

import (
	"testing"
	"time"
)

func TestRateLimiter_Allow(t *testing.T) {
	rl := NewRateLimiter(10, 10) // 10 per minute, burst 10

	// First 10 requests should be allowed
	for i := 0; i < 10; i++ {
		if !rl.Allow("test-key") {
			t.Errorf("request %d should be allowed", i+1)
		}
	}

	// 11th request should be denied
	if rl.Allow("test-key") {
		t.Error("11th request should be denied")
	}
}

func TestRateLimiter_Reset(t *testing.T) {
	rl := NewRateLimiter(10, 10)

	// Use up the bucket
	for i := 0; i < 10; i++ {
		rl.Allow("test-key")
	}

	// Reset
	rl.Reset("test-key")

	// Should be allowed again
	if !rl.Allow("test-key") {
		t.Error("request after reset should be allowed")
	}
}

func TestRateLimiter_Remaining(t *testing.T) {
	rl := NewRateLimiter(10, 10)

	initial := rl.Remaining("test-key")
	if initial != 10 {
		t.Errorf("initial remaining = %d, want 10", initial)
	}

	rl.Allow("test-key")

	remaining := rl.Remaining("test-key")
	if remaining != 9 {
		t.Errorf("remaining after 1 request = %d, want 9", remaining)
	}
}

func TestRateLimiter_Refill(t *testing.T) {
	rl := NewRateLimiter(60, 10) // 60 per minute = 1 per second

	// Use up the bucket
	for i := 0; i < 10; i++ {
		rl.Allow("test-key")
	}

	// Should be denied
	if rl.Allow("test-key") {
		t.Error("request should be denied when bucket empty")
	}

	// Wait for refill (1 second = 1 token)
	time.Sleep(1100 * time.Millisecond)

	// Should be allowed after refill
	if !rl.Allow("test-key") {
		t.Error("request should be allowed after refill")
	}
}

func TestRateLimiter_Count(t *testing.T) {
	rl := NewRateLimiter(10, 10)

	if rl.Count() != 0 {
		t.Errorf("initial count = %d, want 0", rl.Count())
	}

	rl.Allow("key1")
	rl.Allow("key2")

	if rl.Count() != 2 {
		t.Errorf("count after 2 keys = %d, want 2", rl.Count())
	}
}

