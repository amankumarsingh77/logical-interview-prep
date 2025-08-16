
# Scenario: Concurrent Data Aggregator

Imagine you're building a backend service. One of its features is to display aggregated data from several third-party sources. To ensure a snappy user experience, you need to fetch this data concurrently.

## Your Task

Your task is to implement a Go function that fetches "reports" for a list of IDs. However, these external API calls can be slow or fail. Your function must be resilient and efficient.

### Requirements

#### Function Signature
Create a function with the following signature:

```go
func GetAggregatedReports(reportIDs []string) map[string]Report
```

#### Mock API Call
You don't need to call a real API. Create a mock function to simulate the network call:

```go
func fetchReport(reportID string) (Report, error)
```

- This function should have a random delay between 50ms and 500ms to simulate network latency.
- It should also randomly fail (return an error) for about 20% of the calls.
- The Report can be a simple struct, e.g.:
  ```go
  type Report struct { 
      ID string
      Data string 
  }
  ```

#### Concurrency
Your `GetAggregatedReports` function must call `fetchReport` for all `reportIDs` concurrently.

#### Error Handling
If fetching a report for one ID fails, it must not stop the other requests. The final map should only contain reports that were successfully fetched.

#### Timeout
The entire `GetAggregatedReports` operation must have a global timeout of 300ms. If the timeout is hit, the function should return whatever data it has managed to collect up to that point.

---

I encourage you to think out loud as you approach this. How would you structure your Go code to handle the concurrency, error handling, and the timeout? You can start by describing your high-level plan, or you can jump straight into the code.
