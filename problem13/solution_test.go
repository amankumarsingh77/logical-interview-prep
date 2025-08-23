package main

import (
	"testing"
	"time"
)

func TestRateLimiter_Functionality(t *testing.T) {
	limit := 3
	window := 2 * time.Second
	limiter := NewRateLimiter(window, limit, 10*time.Second)
	ip := "192.168.1.100"
	t.Run("Allows requests under limit", func(t *testing.T) {
		for i := 0; i < limit; i++ {
			if !limiter.IsAllowed(ip) {
				t.Errorf("request %d should have been allowed, but was blocked", i+1)
			}
		}
	})

	t.Run("Blocks request over limit", func(t *testing.T) {
		if limiter.IsAllowed(ip) {
			t.Error("request over the limit should have been blocked, but was allowed")
		}
	})

	t.Run("Allows request after window expires", func(t *testing.T) {
		// Wait for a period longer than the window duration
		time.Sleep(window + 100*time.Millisecond)

		if !limiter.IsAllowed(ip) {
			t.Error("request after window expiration should be allowed, but was blocked")
		}
	})
}

func TestRateLimiter_MultipleIPs(t *testing.T) {
	limit := 5
	window := 5 * time.Second
	limiter := NewRateLimiter(window, limit, 10*time.Second)

	ip1 := "10.0.0.1"
	ip2 := "10.0.0.2"

	for i := 0; i < limit; i++ {
		limiter.IsAllowed(ip1)
	}
	if limiter.IsAllowed(ip1) {
		t.Fatalf("ip1 should be blocked after reaching its limit")
	}

	t.Run("Does not block a different IP", func(t *testing.T) {
		if !limiter.IsAllowed(ip2) {
			t.Error("a different IP was blocked, but should have been allowed")
		}
	})
}
