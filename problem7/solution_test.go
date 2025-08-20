package main

import (
	"reflect"
	"sort"
	"testing"
)

func TestFindSuspiciousUsers(t *testing.T) {
	tests := []struct {
		name     string
		logs     []string
		expected []string
	}{
		{
			name: "basic suspicious user",
			logs: []string{
				"1672531200,user-a,service-auth",
				"1672531260,user-a,service-storage",
				"1672531320,user-a,service-compute",
				"1672531380,user-a,service-db",
				"1672531440,user-a,service-network",
				"1672531500,user-a,service-logging",
			},
			expected: []string{"user-a"},
		},
		{
			name: "exactly 5 services in 60 minutes - not suspicious",
			logs: []string{
				"1672531200,user-a,service-auth",
				"1672531260,user-a,service-storage",
				"1672531320,user-a,service-compute",
				"1672531380,user-a,service-db",
				"1672531440,user-a,service-network",
			},
			expected: []string{},
		},
		{
			name: "6 services but spans more than 60 minutes",
			logs: []string{
				"1672531200,user-a,service-auth",
				"1672531260,user-a,service-storage",
				"1672531320,user-a,service-compute",
				"1672531380,user-a,service-db",
				"1672531440,user-a,service-network",
				"1672534801,user-a,service-logging",
			},
			expected: []string{},
		},
		{
			name: "multiple users, one suspicious",
			logs: []string{
				"1672531200,user-a,service-auth",
				"1672531260,user-b,service-auth",
				"1672531320,user-a,service-storage",
				"1672531380,user-a,service-compute",
				"1672531440,user-a,service-db",
				"1672531500,user-a,service-network",
				"1672531560,user-a,service-logging",
				"1672531620,user-b,service-storage",
			},
			expected: []string{"user-a"},
		},
		{
			name: "multiple suspicious users",
			logs: []string{
				"1672531200,user-a,service-1",
				"1672531200,user-b,service-1",
				"1672531260,user-a,service-2",
				"1672531260,user-b,service-2",
				"1672531320,user-a,service-3",
				"1672531320,user-b,service-3",
				"1672531380,user-a,service-4",
				"1672531380,user-b,service-4",
				"1672531440,user-a,service-5",
				"1672531440,user-b,service-5",
				"1672531500,user-a,service-6",
				"1672531500,user-b,service-6",
			},
			expected: []string{"user-a", "user-b"},
		},
		{
			name: "same service accessed multiple times",
			logs: []string{
				"1672531200,user-a,service-auth",
				"1672531260,user-a,service-auth",
				"1672531320,user-a,service-auth",
				"1672531380,user-a,service-auth",
				"1672531440,user-a,service-auth",
				"1672531500,user-a,service-auth",
			},
			expected: []string{},
		},
		{
			name: "sliding window test",
			logs: []string{
				"1672531200,user-a,service-1",
				"1672531260,user-a,service-2",
				"1672531320,user-a,service-3",
				"1672531380,user-a,service-4",
				"1672531440,user-a,service-5",
				"1672534900,user-a,service-6",
				"1672534960,user-a,service-7",
			},
			expected: []string{},
		},
		{
			name:     "empty logs",
			logs:     []string{},
			expected: []string{},
		},
		{
			name: "malformed logs ignored",
			logs: []string{
				"1672531200,user-a,service-1",
				"invalid-log",
				"1672531260,user-a,service-2",
				"incomplete,log",
				"1672531320,user-a,service-3",
				"1672531380,user-a,service-4",
				"1672531440,user-a,service-5",
				"1672531500,user-a,service-6",
			},
			expected: []string{"user-a"},
		},
		{
			name: "unsorted timestamps",
			logs: []string{
				"1672531500,user-a,service-6",
				"1672531200,user-a,service-1",
				"1672531440,user-a,service-5",
				"1672531260,user-a,service-2",
				"1672531380,user-a,service-4",
				"1672531320,user-a,service-3",
			},
			expected: []string{"user-a"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := findSuspiciousUsers(tt.logs)
			sort.Strings(result)
			sort.Strings(tt.expected)

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("findSuspiciousUsers() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetLogsByUser(t *testing.T) {
	logs := []string{
		"1672531200,user-a,service-auth",
		"1672531100,user-a,service-storage",
		"1672531300,user-b,service-compute",
	}

	result := getLogsByUser(logs)

	if len(result) != 2 {
		t.Errorf("getLogsByUser() returned %d users, want 2", len(result))
	}

	userALogs := result["user-a"]
	if len(userALogs) != 2 {
		t.Errorf("getLogsByUser() user-a has %d logs, want 2", len(userALogs))
	}

	if userALogs[0].Timestamp != 1672531100 {
		t.Errorf("getLogsByUser() first log timestamp = %d, want 1672531100 (should be sorted)", userALogs[0].Timestamp)
	}
	if userALogs[1].Timestamp != 1672531200 {
		t.Errorf("getLogsByUser() second log timestamp = %d, want 1672531200", userALogs[1].Timestamp)
	}

	userBLogs := result["user-b"]
	if len(userBLogs) != 1 {
		t.Errorf("getLogsByUser() user-b has %d logs, want 1", len(userBLogs))
	}
}

func TestGetLogsByUserMalformedInput(t *testing.T) {
	logs := []string{
		"1672531200,user-a,service-auth",
		"invalid",
		"incomplete,log",
		"1672531300,user-b,service-compute",
	}

	result := getLogsByUser(logs)

	if len(result) != 2 {
		t.Errorf("getLogsByUser() returned %d users, want 2 (malformed logs should be skipped)", len(result))
	}

	if _, exists := result["user-a"]; !exists {
		t.Errorf("getLogsByUser() user-a not found")
	}
	if _, exists := result["user-b"]; !exists {
		t.Errorf("getLogsByUser() user-b not found")
	}
}

func TestSlidingWindowAlgorithm(t *testing.T) {
	logs := []string{
		"1672531200,user-a,service-1",
		"1672531800,user-a,service-2",
		"1672532400,user-a,service-3",
		"1672533000,user-a,service-4",
		"1672533600,user-a,service-5",
		"1672534200,user-a,service-6",
		"1672534801,user-a,service-7",
	}

	result := findSuspiciousUsers(logs)

	if len(result) != 1 || result[0] != "user-a" {
		t.Errorf("findSuspiciousUsers() = %v, want [user-a] (6 services within first 60-minute window)", result)
	}
}

func TestExactlyAtBoundary(t *testing.T) {
	logs := []string{
		"1672531200,user-a,service-1",
		"1672531260,user-a,service-2",
		"1672531320,user-a,service-3",
		"1672531380,user-a,service-4",
		"1672531440,user-a,service-5",
		"1672534800,user-a,service-6",
	}

	result := findSuspiciousUsers(logs)

	if len(result) != 1 || result[0] != "user-a" {
		t.Errorf("findSuspiciousUsers() = %v, want [user-a] (exactly 3600 seconds = 60 minutes)", result)
	}
}

func TestDuplicateServices(t *testing.T) {
	logs := []string{
		"1672531200,user-a,service-auth",
		"1672531260,user-a,service-storage",
		"1672531320,user-a,service-compute",
		"1672531380,user-a,service-auth",
		"1672531440,user-a,service-db",
		"1672531500,user-a,service-network",
	}

	result := findSuspiciousUsers(logs)

	if len(result) != 0 {
		t.Errorf("findSuspiciousUsers() = %v, want [] (duplicate services shouldn't count)", result)
	}
}
