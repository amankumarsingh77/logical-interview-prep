package main

import (
	"fmt"
	"sync"
	"time"
)

type Item struct {
	value        interface{}
	expiresAt    time.Time
	expired      bool
	neverExpires bool
}

type cache struct {
	items map[string]*Item
	sync.RWMutex
}

type Cache interface {
	Set(key string, value interface{}, ttl time.Duration)
	Get(key string) (interface{}, bool)
	Delete(key string)
}

func NewMemoryCache(cleanUpDuration time.Duration) Cache {
	memCache := &cache{
		items: make(map[string]*Item),
	}
	go memCache.cleanInactiveItems(cleanUpDuration)
	return memCache
}

func (c *cache) Set(key string, value interface{}, ttl time.Duration) {
	c.Lock()
	defer c.Unlock()
	item := &Item{
		value:     value,
		expiresAt: time.Now().Add(ttl),
	}
	if ttl == 0 {
		item.neverExpires = true
	}
	c.items[key] = item
}

func (c *cache) Get(key string) (interface{}, bool) {
	c.RLock()
	defer c.RUnlock()
	item, ok := c.items[key]
	if !ok {
		return nil, false
	}
	if time.Now().After(item.expiresAt) && !item.neverExpires {
		delete(c.items, key)
		return nil, false
	}
	return item.value, true
}

func (c *cache) Delete(key string) {
	c.Lock()
	defer c.Unlock()
	delete(c.items, key)
}

func (c *cache) cleanInactiveItems(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		c.Lock()
		now := time.Now()
		for key, val := range c.items {
			if now.After(val.expiresAt) && !val.neverExpires {
				delete(c.items, key)
			}
		}
		c.Unlock()
	}
}

func main() {
	memCache := NewMemoryCache(time.Second * 10)

	fmt.Println("## Test 1: Basic Set and Get")
	memCache.Set("name", "gemini", 10*time.Second)
	if val, found := memCache.Get("name"); found && val == "gemini" {
		fmt.Println("PASS: Got correct value.")
	} else {
		fmt.Printf("FAIL: Got val=%v, found=%v\n", val, found)
	}
	fmt.Println("---")

	fmt.Println("## Test 2: Get Non-Existent Key")
	if val, found := memCache.Get("nonexistent"); !found {
		fmt.Println("PASS: Key was not found as expected.")
	} else {
		fmt.Printf("FAIL: Got val=%v, found=%v for non-existent key\n", val, found)
	}
	fmt.Println("---")

	fmt.Println("## Test 3: Key Expiration")
	memCache.Set("short", "lived", 1*time.Second)
	time.Sleep(2 * time.Second)
	if val, found := memCache.Get("short"); !found {
		fmt.Println("PASS: Key expired as expected.")
	} else {
		fmt.Printf("FAIL: Expired key still exists with val=%v, found=%v\n", val, found)
	}
	fmt.Println("---")

	fmt.Println("## Test 4: Key without Expiration (TTL=0)")
	memCache.Set("permanent", 100, 0)
	time.Sleep(1 * time.Second)
	if val, found := memCache.Get("permanent"); found && val == 100 {
		fmt.Println("PASS: Non-expiring key persists.")
	} else {
		fmt.Printf("FAIL: Non-expiring key failed with val=%v, found=%v\n", val, found)
	}
	fmt.Println("---")

	fmt.Println("## Test 5: Update and Delete")
	memCache.Set("todelete", "initial", 10*time.Second)
	memCache.Set("todelete", "updated", 10*time.Second)
	if val, found := memCache.Get("todelete"); found && val == "updated" {
		fmt.Println("PASS: Key value was updated correctly.")
	} else {
		fmt.Printf("FAIL: Key update failed with val=%v, found=%v\n", val, found)
	}
	memCache.Delete("todelete")
	if _, found := memCache.Get("todelete"); !found {
		fmt.Println("PASS: Key was deleted successfully.")
	} else {
		fmt.Println("FAIL: Key still exists after deletion.")
	}
	fmt.Println("---")

	fmt.Println("## Test 6: Concurrency Safety")
	var wg sync.WaitGroup
	numGoroutines := 1000
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(n int) {
			defer wg.Done()
			key := fmt.Sprintf("concurrent-%d", n)
			memCache.Set(key, n, 10*time.Second)
			val, found := memCache.Get(key)
			if !found || val.(int) != n {
				fmt.Printf("FAIL: Concurrency error on key %s\n", key)
			}
		}(i)
	}
	wg.Wait()

	finalVal, finalFound := memCache.Get(fmt.Sprintf("concurrent-%d", numGoroutines-1))
	if finalFound && finalVal.(int) == numGoroutines-1 {
		fmt.Println("PASS: Concurrency test finished successfully.")
	} else {
		fmt.Println("FAIL: A value was missing or incorrect after concurrency test.")
	}
	fmt.Println("---")
}
