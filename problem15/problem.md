# Scenario: The Concurrent Thumbnail Generator 🖼️

Imagine you are building a service for a photo gallery website. When a user uploads a batch of photos, your backend needs to generate a thumbnail for each one. Processing all of them at once could overwhelm the server's CPU and memory.

---

## Your Task

Create a function that takes a list of image URLs and generates thumbnails for them concurrently, but using a limited, fixed number of worker goroutines.

---

## Goal

Implement the `generateThumbnails` function. It should distribute the "work" (processing each image URL) to a pool of workers.

---

## Key Constraints

- **Controlled Concurrency**: The function must not start a new goroutine for every single image. Instead, it should use a specific number of worker goroutines, defined by the `maxWorkers` parameter.
- **Work Distribution**: The image URLs should be distributed among the available workers as they become free.
- **Error Handling**: If processing a specific image fails, it should not stop the other images from being processed. The function should report which images succeeded and which failed.

---

## Setup

Here is a mock function that simulates the work of downloading and processing a single image. You will call this function from your workers.

```go
// --- MOCKED PROCESSING FUNCTION (You don't need to change this) ---

// processImage simulates downloading and creating a thumbnail for an image URL.
// It takes about 100ms and can randomly fail.
func processImage(url string) (string, error) {
    fmt.Printf("Processing image: %s\n", url)
    // Simulate work
    time.Sleep(100 * time.Millisecond)

    // Simulate random failures
    if strings.Contains(url, "fail") {
        return "", fmt.Errorf("failed to process %s", url)
    }

    thumbnailPath := "thumbnails/" + strings.Replace(url, "images/", "", 1)
    return thumbnailPath, nil
}
```

---

## Your Task is to Implement This

```go
type ThumbnailResult struct {
    Successes map[string]string // Key: original URL, Value: thumbnail path
    Failures  map[string]error  // Key: original URL, Value: error
}

// generateThumbnails processes a list of imageURLs using a worker pool.
// It should use exactly `maxWorkers` goroutines to do the processing.
func generateThumbnails(imageURLs []string, maxWorkers int) ThumbnailResult {
    // Your implementation goes here.
}
```
