## Scenario: API Rate Limiter 🚦

### Imagine
You are building a backend service that exposes a public API. To prevent abuse and ensure fair usage for all clients, you need to implement a rate limiter.

### The Goal
Create a rate limiter that restricts each unique client (identified by an IP address) to a maximum of **10 requests per minute**.

### The Function
Implement a function or a method:

```go
IsAllowed(ipAddress string) bool
```

### Behavior
- When a request comes in from a given IP address, this function is called.
- If the IP address has made **fewer than 10 requests in the last 60 seconds**, the function should return `true` (allowing the request) **and record the request**.
- If the IP address has **already made 10 or more requests in the last 60 seconds**, the function should return `false` (blocking the request).

### Assumptions
- This will run on a **single server instance**, so you don't need to worry about distributed systems complexity for now.
- The system should be **memory-efficient**, as it might handle **thousands of unique IP addresses**.
