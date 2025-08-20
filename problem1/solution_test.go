package main

import (
	"testing"
)

func TestResolveOrder(t *testing.T) {
	tests := []struct {
		name    string
		tasks   map[string][]string
		wantErr bool
		errMsg  string
	}{
		{
			name: "basic dependency resolution",
			tasks: map[string][]string{
				"compile":       {"download_deps"},
				"run_tests":     {"compile"},
				"create_image":  {"compile"},
				"deploy":        {"run_tests", "create_image"},
				"download_deps": {},
			},
			wantErr: false,
		},
		{
			name: "no dependencies",
			tasks: map[string][]string{
				"task1": {},
				"task2": {},
				"task3": {},
			},
			wantErr: false,
		},
		{
			name: "single task with dependency",
			tasks: map[string][]string{
				"task1": {"task2"},
				"task2": {},
			},
			wantErr: false,
		},
		{
			name: "cycle detection - direct cycle",
			tasks: map[string][]string{
				"task1": {"task2"},
				"task2": {"task1"},
			},
			wantErr: true,
			errMsg:  "cycle detected",
		},
		{
			name: "cycle detection - indirect cycle",
			tasks: map[string][]string{
				"task1": {"task2"},
				"task2": {"task3"},
				"task3": {"task1"},
			},
			wantErr: true,
			errMsg:  "cycle detected",
		},
		{
			name: "missing dependency",
			tasks: map[string][]string{
				"task1": {"task2"},
				"task3": {},
			},
			wantErr: true,
			errMsg:  "could not find task",
		},
		{
			name: "complex dependency graph",
			tasks: map[string][]string{
				"deploy":     {"test", "build"},
				"test":       {"compile"},
				"build":      {"compile"},
				"compile":    {"download"},
				"download":   {},
				"package":    {"build"},
				"distribute": {"package", "deploy"},
			},
			wantErr: false,
		},
		{
			name:    "empty tasks",
			tasks:   map[string][]string{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ResolveOrder(tt.tasks)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ResolveOrder() expected error but got none")
					return
				}
				if tt.errMsg != "" && !contains(err.Error(), tt.errMsg) {
					t.Errorf("ResolveOrder() error = %v, want error containing %v", err, tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("ResolveOrder() unexpected error = %v", err)
				return
			}

			if !isValidTopologicalOrder(tt.tasks, result) {
				t.Errorf("ResolveOrder() result is not a valid topological order: %v", result)
			}

			if len(result) != len(tt.tasks) {
				t.Errorf("ResolveOrder() result length = %v, want %v", len(result), len(tt.tasks))
			}

			taskSet := make(map[string]bool)
			for _, task := range result {
				if taskSet[task] {
					t.Errorf("ResolveOrder() duplicate task in result: %v", task)
				}
				taskSet[task] = true
			}

			for task := range tt.tasks {
				if !taskSet[task] {
					t.Errorf("ResolveOrder() missing task in result: %v", task)
				}
			}
		})
	}
}

func isValidTopologicalOrder(tasks map[string][]string, order []string) bool {
	position := make(map[string]int)
	for i, task := range order {
		position[task] = i
	}

	for task, deps := range tasks {
		taskPos := position[task]
		for _, dep := range deps {
			depPos := position[dep]
			if depPos >= taskPos {
				return false
			}
		}
	}
	return true
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr ||
		(len(s) > len(substr) && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 1; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestResolveOrderDeterministic(t *testing.T) {
	tasks := map[string][]string{
		"compile":       {"download_deps"},
		"run_tests":     {"compile"},
		"create_image":  {"compile"},
		"deploy":        {"run_tests", "create_image"},
		"download_deps": {},
	}

	results := make([][]string, 10)
	for i := 0; i < 10; i++ {
		result, err := ResolveOrder(tasks)
		if err != nil {
			t.Errorf("ResolveOrder() unexpected error = %v", err)
			return
		}
		results[i] = result
	}

	for i := 1; i < len(results); i++ {
		if !isValidTopologicalOrder(tasks, results[i]) {
			t.Errorf("ResolveOrder() run %d produced invalid order: %v", i, results[i])
		}
	}
}
