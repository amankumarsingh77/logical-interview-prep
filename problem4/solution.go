package main

import (
	"context"
	"fmt"
	"log"
	"math/rand/v2"
	"sync"
	"time"
)

type Report struct {
	ID   string
	Data string
}

func main() {
	reportIDs := []string{
		"alpha", "beta", "gamma", "delta",
		"epsilon", "zeta", "eta", "theta",
	}

	log.Println("Starting report aggregation...")
	startTime := time.Now()

	reports := GetAggregatedReports(reportIDs)

	duration := time.Since(startTime)
	log.Printf("Aggregation finished in %v.", duration)
	log.Printf("Global timeout was 300ms.")

	fmt.Println("\n--- Fetched Reports ---")
	if len(reports) == 0 {
		fmt.Println("No reports were successfully fetched.")
	} else {
		for id, report := range reports {
			fmt.Printf("  - ID: %s, Data: '%s'\n", id, report.Data)
		}
	}
	fmt.Println("-----------------------")
}

func GetAggregatedReports(reportIDs []string) map[string]Report {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(300)*time.Millisecond)
	defer cancel()
	var wg sync.WaitGroup
	data := make(chan Report, len(reportIDs))
	resultReports := make(map[string]Report)
	for _, id := range reportIDs {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			if report, err := fetchReport(id); err == nil {
				select {
				case data <- report:
				case <-ctx.Done():
				}
			}
		}(id)
	}
	go func() {
		wg.Wait()
		close(data)
	}()
	for {
		select {
		case <-ctx.Done():
			return resultReports
		case report, ok := <-data:
			if !ok {
				return resultReports
			}
			resultReports[report.ID] = report
		}
	}
}

func fetchReport(reportID string) (Report, error) {
	delay := time.Duration(rand.IntN(500-50+1) + 50)
	time.Sleep(delay * time.Millisecond)
	var report Report
	if rand.Float64() < 0.2 {
		return report, fmt.Errorf("failed to get the report")
	}
	return Report{
		ID:   reportID,
		Data: fmt.Sprintf("Report data for ID: %s", reportID),
	}, nil
}
