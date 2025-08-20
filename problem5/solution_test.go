package main

import (
	"testing"
)

func TestNewLRUCache(t *testing.T) {
	cache := NewLRUCache(3)

	if cache.capacity != 3 {
		t.Errorf("NewLRUCache() capacity = %v, want 3", cache.capacity)
	}
	if cache.list.Len() != 0 {
		t.Errorf("NewLRUCache() initial list length = %v, want 0", cache.list.Len())
	}
	if len(cache.cache) != 0 {
		t.Errorf("NewLRUCache() initial cache map length = %v, want 0", len(cache.cache))
	}
}

func TestLRUCacheGetMiss(t *testing.T) {
	cache := NewLRUCache(3)

	value, found := cache.Get("nonexistent")

	if found {
		t.Errorf("Get() found = true, want false for nonexistent key")
	}
	if value != "" {
		t.Errorf("Get() value = %v, want empty string for nonexistent key", value)
	}
}

func TestLRUCachePutAndGet(t *testing.T) {
	cache := NewLRUCache(3)

	cache.Put("key1", "value1")
	value, found := cache.Get("key1")

	if !found {
		t.Errorf("Get() found = false, want true")
	}
	if value != "value1" {
		t.Errorf("Get() value = %v, want value1", value)
	}
	if cache.list.Len() != 1 {
		t.Errorf("Cache list length = %v, want 1", cache.list.Len())
	}
}

func TestLRUCacheUpdateExistingKey(t *testing.T) {
	cache := NewLRUCache(3)

	cache.Put("key1", "value1")
	cache.Put("key1", "value2")

	value, found := cache.Get("key1")

	if !found {
		t.Errorf("Get() found = false, want true")
	}
	if value != "value2" {
		t.Errorf("Get() value = %v, want value2", value)
	}
	if cache.list.Len() != 1 {
		t.Errorf("Cache list length = %v, want 1 (no duplicate entries)", cache.list.Len())
	}
}

func TestLRUCacheEviction(t *testing.T) {
	cache := NewLRUCache(2)

	cache.Put("key1", "value1")
	cache.Put("key2", "value2")
	cache.Put("key3", "value3")

	_, found1 := cache.Get("key1")
	value2, found2 := cache.Get("key2")
	value3, found3 := cache.Get("key3")

	if found1 {
		t.Errorf("Get(key1) found = true, want false (should be evicted)")
	}
	if !found2 {
		t.Errorf("Get(key2) found = false, want true")
	}
	if value2 != "value2" {
		t.Errorf("Get(key2) value = %v, want value2", value2)
	}
	if !found3 {
		t.Errorf("Get(key3) found = false, want true")
	}
	if value3 != "value3" {
		t.Errorf("Get(key3) value = %v, want value3", value3)
	}
	if cache.list.Len() != 2 {
		t.Errorf("Cache list length = %v, want 2", cache.list.Len())
	}
}

func TestLRUCacheGetPromotesToFront(t *testing.T) {
	cache := NewLRUCache(3)

	cache.Put("A", "1")
	cache.Put("B", "2")
	cache.Put("C", "3")

	cache.Get("A")

	cache.Put("D", "4")

	_, foundA := cache.Get("A")
	_, foundB := cache.Get("B")
	_, foundC := cache.Get("C")
	_, foundD := cache.Get("D")

	if !foundA {
		t.Errorf("Get(A) found = false, want true (A should be promoted)")
	}
	if foundB {
		t.Errorf("Get(B) found = true, want false (B should be evicted)")
	}
	if !foundC {
		t.Errorf("Get(C) found = false, want true")
	}
	if !foundD {
		t.Errorf("Get(D) found = false, want true")
	}
}

func TestLRUCachePutPromotesToFront(t *testing.T) {
	cache := NewLRUCache(3)

	cache.Put("A", "1")
	cache.Put("B", "2")
	cache.Put("C", "3")

	cache.Put("A", "1_updated")

	cache.Put("D", "4")

	valueA, foundA := cache.Get("A")
	_, foundB := cache.Get("B")
	_, foundC := cache.Get("C")
	_, foundD := cache.Get("D")

	if !foundA {
		t.Errorf("Get(A) found = false, want true (A should be promoted)")
	}
	if valueA != "1_updated" {
		t.Errorf("Get(A) value = %v, want 1_updated", valueA)
	}
	if foundB {
		t.Errorf("Get(B) found = true, want false (B should be evicted)")
	}
	if !foundC {
		t.Errorf("Get(C) found = false, want true")
	}
	if !foundD {
		t.Errorf("Get(D) found = false, want true")
	}
}

func TestLRUCacheCapacityZero(t *testing.T) {
	cache := NewLRUCache(0)

	cache.Put("key1", "value1")
	value, found := cache.Get("key1")

	if found {
		t.Errorf("Get() found = true, want false (capacity is 0)")
	}
	if value != "" {
		t.Errorf("Get() value = %v, want empty string", value)
	}
}

func TestLRUCacheCapacityOne(t *testing.T) {
	cache := NewLRUCache(1)

	cache.Put("key1", "value1")
	cache.Put("key2", "value2")

	value1, found1 := cache.Get("key1")
	value2, found2 := cache.Get("key2")

	if found1 {
		t.Errorf("Get(key1) found = true, want false (should be evicted)")
	}
	if value1 != "" {
		t.Errorf("Get(key1) value = %v, want empty string", value1)
	}
	if !found2 {
		t.Errorf("Get(key2) found = false, want true")
	}
	if value2 != "value2" {
		t.Errorf("Get(key2) value = %v, want value2", value2)
	}
}

func TestLRUCacheComplexSequence(t *testing.T) {
	cache := NewLRUCache(3)

	cache.Put("1", "one")
	cache.Put("2", "two")
	cache.Put("3", "three")

	cache.Get("1")
	cache.Put("4", "four")

	_, found1 := cache.Get("1")
	_, found2 := cache.Get("2")
	_, found3 := cache.Get("3")
	_, found4 := cache.Get("4")

	if !found1 {
		t.Errorf("Get(1) found = false, want true")
	}
	if found2 {
		t.Errorf("Get(2) found = true, want false (should be evicted)")
	}
	if !found3 {
		t.Errorf("Get(3) found = false, want true")
	}
	if !found4 {
		t.Errorf("Get(4) found = false, want true")
	}

	cache.Get("3")
	cache.Put("5", "five")

	_, found1After := cache.Get("1")
	_, found3After := cache.Get("3")
	_, found4After := cache.Get("4")
	_, found5 := cache.Get("5")

	if found1After {
		t.Errorf("Get(1) found = true, want false (should be evicted after second round)")
	}
	if !found3After {
		t.Errorf("Get(3) found = false, want true")
	}
	if !found4After {
		t.Errorf("Get(4) found = false, want true")
	}
	if !found5 {
		t.Errorf("Get(5) found = false, want true")
	}
}
