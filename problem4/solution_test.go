package main

import (
	"fmt"
	"testing"
	"time"
)

func TestGetAggregatedReportsTimeout(t *testing.T) {
	reportIDs := []string{"test1", "test2", "test3"}

	start := time.Now()
	reports := GetAggregatedReports(reportIDs)
	duration := time.Since(start)

	if duration > 350*time.Millisecond {
		t.Errorf("GetAggregatedReports() took %v, should timeout around 300ms", duration)
	}

	if len(reports) > len(reportIDs) {
		t.Errorf("GetAggregatedReports() returned %d reports, max should be %d", len(reports), len(reportIDs))
	}

	for id, report := range reports {
		if report.ID != id {
			t.Errorf("GetAggregatedReports() report ID mismatch: got %s, want %s", report.ID, id)
		}
		expectedData := fmt.Sprintf("Report data for ID: %s", id)
		if report.Data != expectedData {
			t.Errorf("GetAggregatedReports() report data = %s, want %s", report.Data, expectedData)
		}
	}
}

func TestGetAggregatedReportsEmptyInput(t *testing.T) {
	reportIDs := []string{}

	start := time.Now()
	reports := GetAggregatedReports(reportIDs)
	duration := time.Since(start)

	if duration > 50*time.Millisecond {
		t.Errorf("GetAggregatedReports() with empty input took %v, should be very fast", duration)
	}

	if len(reports) != 0 {
		t.Errorf("GetAggregatedReports() returned %d reports, want 0", len(reports))
	}
}

func TestGetAggregatedReportsSingleID(t *testing.T) {
	reportIDs := []string{"single"}

	reports := GetAggregatedReports(reportIDs)

	if len(reports) > 1 {
		t.Errorf("GetAggregatedReports() returned %d reports, max should be 1", len(reports))
	}

	if report, exists := reports["single"]; exists {
		if report.ID != "single" {
			t.Errorf("GetAggregatedReports() report ID = %s, want single", report.ID)
		}
		expectedData := "Report data for ID: single"
		if report.Data != expectedData {
			t.Errorf("GetAggregatedReports() report data = %s, want %s", report.Data, expectedData)
		}
	}
}

func TestGetAggregatedReportsConcurrency(t *testing.T) {
	reportIDs := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}

	start := time.Now()
	reports := GetAggregatedReports(reportIDs)
	duration := time.Since(start)

	if duration > 350*time.Millisecond {
		t.Errorf("GetAggregatedReports() took %v, concurrent execution should complete within timeout", duration)
	}

	for id, report := range reports {
		if report.ID != id {
			t.Errorf("GetAggregatedReports() report ID mismatch: got %s, want %s", report.ID, id)
		}
	}
}

func mockFetchReportSuccess(reportID string) (Report, error) {
	time.Sleep(50 * time.Millisecond)
	return Report{
		ID:   reportID,
		Data: fmt.Sprintf("Report data for ID: %s", reportID),
	}, nil
}

func mockFetchReportFailure(reportID string) (Report, error) {
	time.Sleep(50 * time.Millisecond)
	return Report{}, fmt.Errorf("failed to fetch report %s", reportID)
}

func mockFetchReportSlow(reportID string) (Report, error) {
	time.Sleep(400 * time.Millisecond)
	return Report{
		ID:   reportID,
		Data: fmt.Sprintf("Report data for ID: %s", reportID),
	}, nil
}

func TestGetAggregatedReportsWithMockSuccess(t *testing.T) {
	reportIDs := []string{"test1", "test2", "test3"}
	reports := GetAggregatedReports(reportIDs)

	if len(reports) > 3 {
		t.Errorf("GetAggregatedReports() returned %d reports, max should be 3", len(reports))
	}

	for id, report := range reports {
		if report.ID != id {
			t.Errorf("GetAggregatedReports() report ID = %s, want %s", report.ID, id)
		}
	}
}

func TestGetAggregatedReportsBasicFunctionality(t *testing.T) {
	reportIDs := []string{"alpha", "beta", "gamma"}
	reports := GetAggregatedReports(reportIDs)

	for id, report := range reports {
		if report.ID != id {
			t.Errorf("GetAggregatedReports() report ID = %s, want %s", report.ID, id)
		}
		expectedData := fmt.Sprintf("Report data for ID: %s", id)
		if report.Data != expectedData {
			t.Errorf("GetAggregatedReports() report data = %s, want %s", report.Data, expectedData)
		}
	}
}

func TestGetAggregatedReportsTimeoutBehavior(t *testing.T) {
	reportIDs := []string{"test1", "test2"}

	start := time.Now()
	reports := GetAggregatedReports(reportIDs)
	duration := time.Since(start)

	if duration > 350*time.Millisecond {
		t.Errorf("GetAggregatedReports() took %v, should timeout around 300ms", duration)
	}

	if len(reports) > len(reportIDs) {
		t.Errorf("GetAggregatedReports() returned %d reports, max should be %d", len(reports), len(reportIDs))
	}
}

func TestGetAggregatedReportsLargeInput(t *testing.T) {
	reportIDs := make([]string, 100)
	for i := 0; i < 100; i++ {
		reportIDs[i] = fmt.Sprintf("report_%d", i)
	}

	start := time.Now()
	reports := GetAggregatedReports(reportIDs)
	duration := time.Since(start)

	if duration > 350*time.Millisecond {
		t.Errorf("GetAggregatedReports() took %v, should complete within timeout", duration)
	}

	if len(reports) > 100 {
		t.Errorf("GetAggregatedReports() returned %d reports, max should be 100", len(reports))
	}

	for id, report := range reports {
		if report.ID != id {
			t.Errorf("GetAggregatedReports() report ID mismatch: got %s, want %s", report.ID, id)
		}
	}
}
