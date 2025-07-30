
# Scenario: In-Memory LRU Cache

Many high-performance applications, from databases to web servers, use an in-memory cache to avoid re-computing or re-fetching frequently accessed data. To prevent the cache from growing indefinitely, they use an "eviction policy." One of the most common policies is LRU (Least Recently Used).

## Your Task

Your task is to implement an LRU cache in Go from scratch. It will store string key-value pairs and must operate with a fixed capacity.

### Requirements

#### Structure
You'll need to define a `LRUCache` struct.

#### Constructor
Create a constructor function:

```go
NewLRUCache(capacity int) *LRUCache
```

#### Core Methods

Implement the following two methods on your struct:

```go
Get(key string) (string, bool)
```
- If the key exists in the cache, it should return the corresponding value and `true`.
- A `Get` operation counts as using an item, so it should now be considered the most recently used.
- If the key does not exist, it should return an empty string and `false`.

```go
Put(key, value string)
```
- If the key already exists, update its value and mark it as the most recently used.
- If the key is new:
    - If the cache is already at full capacity, evict the least recently used item before inserting the new one.
    - Add the new key-value pair to the cache. It is now the most recently used item.

---

## Hints

How would you approach this? Think about the data structures you would need to achieve both fast lookups (`Get`) and efficient removal of the least-used item.

As before, feel free to outline your plan first or dive right into the code.
