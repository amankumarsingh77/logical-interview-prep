package main

import (
	"errors"
	"testing"
	"time"
)

func TestRetrySuccess(t *testing.T) {
	callCount := 0
	successFn := func() (string, error) {
		callCount++
		return "success", nil
	}

	result, err := Retry(successFn, 3, 10*time.Millisecond)

	if err != nil {
		t.Errorf("Retry() unexpected error = %v", err)
	}
	if result != "success" {
		t.Errorf("Retry() result = %v, want success", result)
	}
	if callCount != 1 {
		t.Errorf("Retry() call count = %v, want 1", callCount)
	}
}

func TestRetryPermanentError(t *testing.T) {
	callCount := 0
	permanentErrorFn := func() (string, error) {
		callCount++
		return "", ErrPermanent
	}

	result, err := Retry(permanentErrorFn, 3, 10*time.Millisecond)

	if !errors.Is(err, ErrPermanent) {
		t.Errorf("Retry() error = %v, want ErrPermanent", err)
	}
	if result != "" {
		t.Errorf("Retry() result = %v, want empty string", result)
	}
	if callCount != 1 {
		t.Errorf("Retry() call count = %v, want 1 (no retries for permanent error)", callCount)
	}
}

func TestRetryTransientErrorEventualSuccess(t *testing.T) {
	callCount := 0
	eventualSuccessFn := func() (string, error) {
		callCount++
		if callCount < 3 {
			return "", ErrTransient
		}
		return "success after retries", nil
	}

	start := time.Now()
	result, err := Retry(eventualSuccessFn, 5, 10*time.Millisecond)
	duration := time.Since(start)

	if err != nil {
		t.Errorf("Retry() unexpected error = %v", err)
	}
	if result != "success after retries" {
		t.Errorf("Retry() result = %v, want 'success after retries'", result)
	}
	if callCount != 3 {
		t.Errorf("Retry() call count = %v, want 3", callCount)
	}

	expectedMinDuration := 10*time.Millisecond + 20*time.Millisecond
	if duration < expectedMinDuration {
		t.Errorf("Retry() duration = %v, want at least %v (exponential backoff)", duration, expectedMinDuration)
	}
}

func TestRetryExhaustsAllRetries(t *testing.T) {
	callCount := 0
	alwaysTransientFn := func() (string, error) {
		callCount++
		return "", ErrTransient
	}

	start := time.Now()
	result, err := Retry(alwaysTransientFn, 3, 5*time.Millisecond)
	duration := time.Since(start)

	if !errors.Is(err, ErrTransient) {
		t.Errorf("Retry() error = %v, want ErrTransient", err)
	}
	if result != "" {
		t.Errorf("Retry() result = %v, want empty string", result)
	}
	if callCount != 4 {
		t.Errorf("Retry() call count = %v, want 4 (initial + 3 retries)", callCount)
	}

	expectedMinDuration := 5*time.Millisecond + 10*time.Millisecond + 20*time.Millisecond
	if duration < expectedMinDuration {
		t.Errorf("Retry() duration = %v, want at least %v (exponential backoff)", duration, expectedMinDuration)
	}
}

func TestRetryZeroRetries(t *testing.T) {
	callCount := 0
	transientErrorFn := func() (string, error) {
		callCount++
		return "", ErrTransient
	}

	result, err := Retry(transientErrorFn, 0, 10*time.Millisecond)

	if !errors.Is(err, ErrTransient) {
		t.Errorf("Retry() error = %v, want ErrTransient", err)
	}
	if result != "" {
		t.Errorf("Retry() result = %v, want empty string", result)
	}
	if callCount != 1 {
		t.Errorf("Retry() call count = %v, want 1 (no retries allowed)", callCount)
	}
}

func TestRetryExponentialBackoff(t *testing.T) {
	callCount := 0

	transientErrorFn := func() (string, error) {
		callCount++
		return "", ErrTransient
	}

	start := time.Now()
	Retry(transientErrorFn, 3, 10*time.Millisecond)
	duration := time.Since(start)

	expectedMinDuration := 70 * time.Millisecond
	if duration < expectedMinDuration {
		t.Errorf("Retry() duration = %v, want at least %v (exponential backoff)", duration, expectedMinDuration)
	}

	expectedMaxDuration := 200 * time.Millisecond
	if duration > expectedMaxDuration {
		t.Errorf("Retry() duration = %v, should not exceed %v", duration, expectedMaxDuration)
	}

	if callCount != 4 {
		t.Errorf("Retry() call count = %v, want 4 (initial + 3 retries)", callCount)
	}
}

func TestRetryCustomError(t *testing.T) {
	customErr := errors.New("custom error")
	callCount := 0
	customErrorFn := func() (string, error) {
		callCount++
		return "", customErr
	}

	result, err := Retry(customErrorFn, 3, 10*time.Millisecond)

	if !errors.Is(err, customErr) {
		t.Errorf("Retry() error = %v, want custom error", err)
	}
	if result != "" {
		t.Errorf("Retry() result = %v, want empty string", result)
	}
	if callCount != 1 {
		t.Errorf("Retry() call count = %v, want 1 (no retries for non-transient error)", callCount)
	}
}

func TestRetryNegativeRetries(t *testing.T) {
	callCount := 0
	transientErrorFn := func() (string, error) {
		callCount++
		return "", ErrTransient
	}

	result, err := Retry(transientErrorFn, -1, 10*time.Millisecond)

	if !errors.Is(err, ErrTransient) {
		t.Errorf("Retry() error = %v, want ErrTransient", err)
	}
	if result != "" {
		t.Errorf("Retry() result = %v, want empty string", result)
	}
	if callCount != 1 {
		t.Errorf("Retry() call count = %v, want 1 (negative retries treated as zero)", callCount)
	}
}
