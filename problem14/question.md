## Scenario: Concurrent In-Memory Cache with Expiration (TTL) ⏳

### Background

In modern backend systems, performance is critical. A common technique to speed things up is to use a cache to store the results of expensive operations (like database queries or API calls) in memory. To ensure this cached data doesn't become stale, items are often stored with a "Time To Live" (TTL), after which they automatically expire.

### The Task

Your goal is to design and implement a generic, in-memory key-value cache in Golang. This cache must be safe for concurrent use and must support a TTL for each item.

### Core Components

You will need to create a Cache struct and implement the following methods:

- **Set(key string, value interface{}, ttl time.Duration)**
  - This method stores a value associated with a key in the cache.
  - The `ttl` (Time To Live) parameter specifies how long the item should be considered valid.
  - If `ttl` is 0 or negative, the item should never expire.

- **Get(key string) (interface{}, bool)**
  - This method retrieves an item from the cache.
  - It should return the value and `true` if the key exists and has not expired.
  - It should return `nil` and `false` if the key does not exist or if the key has expired.

- **Delete(key string)**
  - This method permanently removes a key and its value from the cache.

### Key Requirements & Constraints

- **Thread Safety:** This is the most critical requirement. Your cache implementation must be safe to use from multiple goroutines simultaneously without causing race conditions.
- **Expiration Logic:** The core of the problem is handling the TTL. An item is considered "expired" if the current time is past its creation time plus its TTL. An attempt to `Get` an expired item should behave as if the item doesn't exist.
- **Memory Management:** Think about what happens to the memory used by expired items. If an expired item is never requested again, will it stay in memory forever? An ideal solution considers this.
