package main

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

func processImage(url string) (string, error) {
	fmt.Printf("Processing image: %s\n", url)
	time.Sleep(100 * time.Millisecond)
	if strings.Contains(url, "fail") {
		return "", fmt.Errorf("failed to process %s", url)
	}

	thumbnailPath := "thumbnails/" + strings.Replace(url, "images/", "", 1)
	return thumbnailPath, nil
}

type Result struct {
	URL           string
	ThumbnailPath string
	Error         error
}

type ThumbnailResult struct {
	Successes map[string]string
	Failures  map[string]error
}

func generateThumbnails(imageUrls []string, maxWorker int) *ThumbnailResult {
	jobs := make(chan string, len(imageUrls))
	results := make(chan Result, len(imageUrls))

	var wg sync.WaitGroup
	for i := 0; i < maxWorker; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for url := range jobs {
				thumbnailPath, err := processImage(url)
				results <- Result{
					URL:           url,
					ThumbnailPath: thumbnailPath,
					Error:         err,
				}
			}
		}()
	}

	for _, url := range imageUrls {
		jobs <- url
	}

	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	finalRes := &ThumbnailResult{
		Successes: make(map[string]string),
		Failures:  make(map[string]error),
	}

	for res := range results {
		if res.Error != nil {
			finalRes.Failures[res.URL] = res.Error
		} else {
			finalRes.Successes[res.URL] = res.ThumbnailPath
		}
	}
	return finalRes
}
