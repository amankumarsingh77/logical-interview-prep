package main

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestCache_BasicOperation(t *testing.T) {
	memCache := NewMemoryCache(time.Second * 200)
	t.Run("Set and Get", func(t *testing.T) {
		key := "name"
		value := "aman"
		memCache.Set(key, value, 100*time.Millisecond)
		got, found := memCache.Get(key)
		if !found {
			t.Errorf("expected to find %s, but did not", value)
		}
		if got != value {
			t.Errorf("expected %s for %s", value, got)
		}
	})

	t.Run("Update Value", func(t *testing.T) {
		key := "name"
		value := "gemini"
		memCache.Set(key, value, 0)
		got, found := memCache.Get(key)
		if !found {
			t.Errorf("expected to find %s, but did not", value)
		}
		if got != value {
			t.Errorf("expected %s for %s", value, got)
		}
	})

	t.Run("Delete an element", func(t *testing.T) {
		key := "name"
		memCache.Delete(key)
		if _, found := memCache.Get(key); found {
			t.Errorf("expected to not find any but found")
		}
	})
}

func TestCache_Expiration(t *testing.T) {
	memCache := NewMemoryCache(time.Millisecond * 50)
	t.Run("Item expires after TTL", func(t *testing.T) {
		key := "short"
		value := "i will disappear"
		ttl := 100 * time.Millisecond
		memCache.Set(key, value, ttl)
		if _, found := memCache.Get(key); !found {
			t.Errorf("expected to find %s, but did not", value)
		}
		time.Sleep(ttl + 50*time.Millisecond)
		if _, found := memCache.Get(key); found {
			t.Errorf("expected item to be expired and not found")
		}
	})

	t.Run("Permanent item does not expire", func(t *testing.T) {
		key := "long"
		value := "i wont disappear"
		ttl := 0 * time.Second
		memCache.Set(key, value, ttl)
		got, found := memCache.Get(key)
		if !found {
			t.Errorf("expected permanent item to be found, but it was not")
		}
		if got != value {
			t.Errorf("expected value %v, but got %v", value, got)
		}
	})
}

func TestCache_Concurrency(t *testing.T) {
	memCache := NewMemoryCache(1 * time.Second)
	numGoroutines := 100
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(n int) {
			defer wg.Done()
			key := fmt.Sprintf("key-%d", n)
			value := n
			memCache.Set(key, value, 0)
		}(i)
	}

	wg.Wait()

	for i := 0; i < numGoroutines; i++ {
		key := fmt.Sprintf("key-%d", i)
		expectedValue := i
		got, found := memCache.Get(key)
		if !found {
			t.Errorf("expected to find key %s after concurrent sets, but did not", key)
			continue
		}
		if got != expectedValue {
			t.Errorf("for key %s, expected value %v, but got %v", key, expectedValue, got)
		}
	}
}
