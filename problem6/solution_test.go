package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestUniqNoOptions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "basic deduplication",
			input:    "apple\napple\nbanana\ncherry\ncherry\napple",
			expected: "apple\nbanana\ncherry\napple\n",
		},
		{
			name:     "no consecutive duplicates",
			input:    "apple\nbanana\ncherry",
			expected: "apple\nbanana\ncherry\n",
		},
		{
			name:     "all same lines",
			input:    "apple\napple\napple",
			expected: "apple\n",
		},
		{
			name:     "single line",
			input:    "apple",
			expected: "apple\n",
		},
		{
			name:     "empty input",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			var writer bytes.Buffer
			opts := Options{}

			err := Uniq(reader, &writer, opts)
			if err != nil {
				t.Errorf("Uniq() unexpected error = %v", err)
				return
			}

			if writer.String() != tt.expected {
				t.Errorf("Uniq() output = %q, want %q", writer.String(), tt.expected)
			}
		})
	}
}

func TestUniqCountOption(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "count basic",
			input:    "apple\napple\nbanana\ncherry\ncherry\ncherry",
			expected: "2 apple\n1 banana\n3 cherry\n",
		},
		{
			name:     "count with non-consecutive duplicates",
			input:    "apple\nbanana\napple\nbanana",
			expected: "2 apple\n2 banana\n",
		},
		{
			name:     "count single occurrences",
			input:    "apple\nbanana\ncherry",
			expected: "1 apple\n1 banana\n1 cherry\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			var writer bytes.Buffer
			opts := Options{Count: true}

			err := Uniq(reader, &writer, opts)
			if err != nil {
				t.Errorf("Uniq() unexpected error = %v", err)
				return
			}

			if writer.String() != tt.expected {
				t.Errorf("Uniq() output = %q, want %q", writer.String(), tt.expected)
			}
		})
	}
}

func TestUniqDuplicatesOption(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "duplicates basic",
			input:    "apple\napple\nbanana\ncherry\ncherry\nmango",
			expected: "apple\ncherry\n",
		},
		{
			name:     "no duplicates",
			input:    "apple\nbanana\ncherry",
			expected: "",
		},
		{
			name:     "all duplicates",
			input:    "apple\napple\nbanana\nbanana",
			expected: "apple\nbanana\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			var writer bytes.Buffer
			opts := Options{Duplicates: true}

			err := Uniq(reader, &writer, opts)
			if err != nil {
				t.Errorf("Uniq() unexpected error = %v", err)
				return
			}

			if writer.String() != tt.expected {
				t.Errorf("Uniq() output = %q, want %q", writer.String(), tt.expected)
			}
		})
	}
}

func TestUniqUniqueOption(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "unique basic",
			input:    "apple\napple\nbanana\ncherry\ncherry\nmango",
			expected: "banana\nmango\n",
		},
		{
			name:     "all unique",
			input:    "apple\nbanana\ncherry",
			expected: "apple\nbanana\ncherry\n",
		},
		{
			name:     "no unique lines",
			input:    "apple\napple\nbanana\nbanana",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			var writer bytes.Buffer
			opts := Options{Unique: true}

			err := Uniq(reader, &writer, opts)
			if err != nil {
				t.Errorf("Uniq() unexpected error = %v", err)
				return
			}

			if writer.String() != tt.expected {
				t.Errorf("Uniq() output = %q, want %q", writer.String(), tt.expected)
			}
		})
	}
}

func TestUniqCountAndDuplicates(t *testing.T) {
	input := "apple\napple\nbanana\ncherry\ncherry\ncherry\nmango"
	expected := "2 apple\n3 cherry\n"

	reader := strings.NewReader(input)
	var writer bytes.Buffer
	opts := Options{Count: true, Duplicates: true}

	err := Uniq(reader, &writer, opts)
	if err != nil {
		t.Errorf("Uniq() unexpected error = %v", err)
		return
	}

	if writer.String() != expected {
		t.Errorf("Uniq() output = %q, want %q", writer.String(), expected)
	}
}

func TestUniqCountAndUnique(t *testing.T) {
	input := "apple\napple\nbanana\ncherry\ncherry\nmango"
	expected := "1 banana\n1 mango\n"

	reader := strings.NewReader(input)
	var writer bytes.Buffer
	opts := Options{Count: true, Unique: true}

	err := Uniq(reader, &writer, opts)
	if err != nil {
		t.Errorf("Uniq() unexpected error = %v", err)
		return
	}

	if writer.String() != expected {
		t.Errorf("Uniq() output = %q, want %q", writer.String(), expected)
	}
}

func TestUniqEmptyLines(t *testing.T) {
	input := "apple\n\n\nbanana\n\ncherry"
	expected := "apple\n\nbanana\n\ncherry\n"

	reader := strings.NewReader(input)
	var writer bytes.Buffer
	opts := Options{}

	err := Uniq(reader, &writer, opts)
	if err != nil {
		t.Errorf("Uniq() unexpected error = %v", err)
		return
	}

	if writer.String() != expected {
		t.Errorf("Uniq() output = %q, want %q", writer.String(), expected)
	}
}

func TestUniqWhitespaceLines(t *testing.T) {
	input := "apple\n  \n  \nbanana"
	expected := "apple\n  \nbanana\n"

	reader := strings.NewReader(input)
	var writer bytes.Buffer
	opts := Options{}

	err := Uniq(reader, &writer, opts)
	if err != nil {
		t.Errorf("Uniq() unexpected error = %v", err)
		return
	}

	if writer.String() != expected {
		t.Errorf("Uniq() output = %q, want %q", writer.String(), expected)
	}
}

func TestUniqLongInput(t *testing.T) {
	var inputBuilder strings.Builder
	for i := 0; i < 1000; i++ {
		inputBuilder.WriteString("line\n")
	}
	input := inputBuilder.String()
	expected := "line\n"

	reader := strings.NewReader(input)
	var writer bytes.Buffer
	opts := Options{}

	err := Uniq(reader, &writer, opts)
	if err != nil {
		t.Errorf("Uniq() unexpected error = %v", err)
		return
	}

	if writer.String() != expected {
		t.Errorf("Uniq() output length = %d, want %d", len(writer.String()), len(expected))
	}
}
