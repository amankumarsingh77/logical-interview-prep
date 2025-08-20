# Scenario: Resilient API Client with Exponential Backoff

Imagine you are writing a client that needs to fetch data from an external API. This API is known to be unreliable; sometimes requests fail due to temporary issues like network timeouts or server overload. Instead of failing immediately, you want your client to automatically retry the request a few times.

## Your Task

Write a generic Retry function in Go that takes another function as an argument and executes it. If the function fails with a transient (retryable) error, your Retry function should wait for a period and try again, implementing an exponential backoff strategy for the delay.

### Requirements

#### Function Signature

```go
func Retry(fn func() (string, error), retries int, delay time.Duration) (string, error)
```

#### Retry Logic

- If `fn` executes successfully, `Retry` should return its result immediately.
- If `fn` returns a retryable error, the `Retry` function should wait and then call `fn` again.
- If `fn` returns a non-retryable error (e.g., "400 Bad Request"), `Retry` should stop immediately and return that error.

#### Exponential Backoff

The delay between retries should double after each failed attempt. For example, if the initial delay is 100ms, the sequence of delays should be 100ms, 200ms, 400ms, and so on.

#### Max Retries

If the function still fails after the specified number of retries, `Retry` should give up and return the last error it received.

---

## Provided for You

```go
import (
    "errors"
    "fmt"
    "math/rand"
    "time"
)

// Define specific error types to distinguish between failures.
var ErrTransient = errors.New("a transient error occurred")
var ErrPermanent = errors.New("a permanent error occurred")

// UnreliableAPICall simulates a function that might fail.
func UnreliableAPICall() (string, error) {
    // Simulate random failures
    r := rand.Intn(10)

    if r < 3 { // 30% chance of permanent error
        fmt.Println("API call failed with a permanent error.")
        return "", ErrPermanent
    } else if r < 7 { // 40% chance of transient error
        fmt.Println("API call failed with a transient error.")
        return "", ErrTransient
    }

    fmt.Println("API call succeeded!")
    return "Success!", nil
}
```

---

## Thought Process

- How would you structure the retry loop?
- How do you decide whether to retry or return an error?
- How do you double the delay between retries?
- How do you ensure it doesn’t retry more than `retries` times?