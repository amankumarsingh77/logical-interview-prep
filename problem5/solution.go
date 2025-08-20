package main

import (
	"container/list"
	"fmt"
)

type Node struct {
	key   string
	value string
}

type LRUCache struct {
	capacity int
	list     *list.List
	cache    map[string]*list.Element
}

func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		list:     list.New(),
		cache:    make(map[string]*list.Element),
	}
}

func (l *LRUCache) Get(key string) (string, bool) {
	ele, ok := l.cache[key]
	if !ok {
		return "", false
	}
	l.list.MoveToFront(ele)
	return ele.Value.(Node).value, true
}

func (l *LRUCache) Put(key string, value string) {
	if l.capacity == 0 {
		return
	}
	ele, ok := l.cache[key]
	if ok {
		l.list.MoveToFront(ele)
		ele.Value = Node{key, value}
	} else {
		if l.capacity == l.list.Len() {
			back := l.list.Back()
			if back != nil {
				l.list.Remove(back)
				delete(l.cache, back.Value.(Node).key)
			}
		}
		pushedElement := l.list.PushFront(Node{key, value})
		l.cache[key] = pushedElement
	}
}

func (l *LRUCache) printState(action string) {
	fmt.Printf("--- After %s ---\n", action)
	fmt.Print("Order (MRU -> LRU): ")
	for e := l.list.Front(); e != nil; e = e.Next() {
		fmt.Printf("[%s: %s] ", e.Value.(Node).key, e.Value.(Node).value)
	}
	fmt.Println("\n--------------------")
}

func main() {
	lru := NewLRUCache(3)
	fmt.Println("Created Cache with Capacity 3")
	lru.printState("Init")

	lru.Put("A", "1")
	lru.printState(`Put("A", "1")`)
	lru.Put("B", "2")
	lru.printState(`Put("B", "2")`)
	lru.Put("C", "3")
	lru.printState(`Put("C", "3")`)

	val, found := lru.Get("A")
	fmt.Printf("Get(\"A\") -> Value: %s, Found: %v\n", val, found)
	lru.printState(`Get("A")`)

	lru.Put("D", "4")
	lru.printState(`Put("D", "4")`)

	val, found = lru.Get("B")
	fmt.Printf("Get(\"B\") -> Value: %s, Found: %v\n", val, found)
	lru.printState(`Get("B")`)
}
