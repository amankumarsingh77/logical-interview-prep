package main

import (
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
