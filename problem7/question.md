## Scenario: Suspicious Activity Detection

You are building a security monitoring tool for a large cloud provider. The tool ingests access logs from various services. Your task is to identify users who might be performing suspicious reconnaissance activities.

### Suspicious User Criteria

A user is considered **"suspicious"** if they access **more than 5 different services** within **any 60-minute window**.

### Log Format

Each log is a string in the format:

```
<timestamp>,<userID>,<serviceID>
```

* **timestamp**: Unix timestamp (integer)
* **userID**: String identifying the user
* **serviceID**: String identifying the service

The logs are **not guaranteed to be sorted by time**.

### Function Signature (in Go)

```go
func findSuspiciousUsers(logs []string) []string
```

### Example Input

```go
[]string{
  "1672531200,user-a,service-auth",    // user-a, 10:00
  "1672531260,user-b,service-auth",    // user-b, 10:01
  "1672531320,user-a,service-storage", // user-a, 10:02
  "1672531380,user-a,service-compute", // user-a, 10:03
  "1672531440,user-a,service-db",      // user-a, 10:04
  "1672534740,user-a,service-network", // user-a, 10:59 (5th service)
  "1672534800,user-a,service-logging", // user-a, 11:00 (6th service in 60 mins) -> user-a is suspicious
  "1672534860,user-b,service-cache",   // user-b, 11:01
}
```

### Expected Output

```go
[]string{"user-a"} // Sorted list of unique suspicious users
```

### Task

* Implement the `findSuspiciousUsers` function.
* The function should return a slice of unique user IDs that have been flagged.
* The result should be **sorted alphabetically**.

### Guidelines

* Think through the data structures required:

    * Map of userID to a list of their access events
    * Use sliding window per user to track distinct services accessed in 60-minute intervals
* Edge cases to consider:

    * Same service accessed multiple times should not be double-counted
    * Sparse access logs (e.g., large gaps between entries)

Let me know when you're ready to start coding or would like help writing the Go implementation.
