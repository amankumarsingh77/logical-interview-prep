package main

import (
	"testing"
)

func TestResolveTemplate(t *testing.T) {
	data := map[string]interface{}{
		"user": map[string]interface{}{
			"name": "Alex",
			"id":   123,
		},
		"appointment": map[string]interface{}{
			"service": "Dental Check-up",
			"time":    "2025-08-16T14:00:00Z",
		},
		"status": "active",
	}

	tests := []struct {
		name        string
		template    string
		data        map[string]interface{}
		expectedStr string
		expectErr   bool
		errContains string
	}{
		{
			name:        "basic substitution",
			template:    "Hello, ${user.name}!",
			data:        data,
			expectedStr: "Hello, Alex!",
			expectErr:   false,
		},
		{
			name:        "multiple placeholders",
			template:    "Service for ${user.name}: ${appointment.service}.",
			data:        data,
			expectedStr: "Service for Alex: Dental Check-up.",
			expectErr:   false,
		},
		{
			name:        "non-string value",
			template:    "User ID: ${user.id}",
			data:        data,
			expectedStr: "User ID: 123",
			expectErr:   false,
		},
		{
			name:        "missing top-level key",
			template:    "System info: ${system.version}",
			data:        data,
			expectedStr: "System info: ",
			expectErr:   false,
		},
		{
			name:        "missing nested key",
			template:    "User address: ${user.address.city}",
			data:        data,
			expectedStr: "User address: ",
			expectErr:   false,
		},
		{
			name:        "path traverses past final value",
			template:    "Status detail: ${status.detail}",
			data:        data,
			expectedStr: "Status detail: ",
			expectErr:   false,
		},
		{
			name:        "malformed placeholder unclosed",
			template:    "Hello, ${user.name",
			data:        data,
			expectedStr: "",
			expectErr:   true,
			errContains: "missing closing '}'",
		},
		{
			name:        "no placeholders",
			template:    "This is a static string.",
			data:        data,
			expectedStr: "This is a static string.",
			expectErr:   false,
		},
		{
			name:        "empty template",
			template:    "",
			data:        data,
			expectedStr: "",
			expectErr:   false,
		},
		{
			name:        "adjacent placeholders",
			template:    "${user.name}${status}",
			data:        data,
			expectedStr: "Alexactive",
			expectErr:   false,
		},
		{
			name:        "placeholder at start and end",
			template:    "${status} user ${user.name}",
			data:        data,
			expectedStr: "active user Alex",
			expectErr:   false,
		},
		{
			name:        "empty placeholder",
			template:    "Hello ${}!",
			data:        data,
			expectedStr: "Hello !",
			expectErr:   false,
		},
		{
			name:        "nested placeholder syntax error",
			template:    "Hello ${user.${name}}!",
			data:        data,
			expectedStr: "",
			expectErr:   true,
			errContains: "malformed template",
		},
		{
			name:        "only placeholder opening",
			template:    "${",
			data:        data,
			expectedStr: "",
			expectErr:   true,
			errContains: "missing closing '}'",
		},
		{
			name:        "only placeholder closing",
			template:    "}",
			data:        data,
			expectedStr: "}",
			expectErr:   false,
		},
		{
			name:        "multiple same placeholders",
			template:    "${user.name} and ${user.name} again",
			data:        data,
			expectedStr: "Alex and Alex again",
			expectErr:   false,
		},
		{
			name:        "placeholder with spaces around",
			template:    "Hello ${ user.name }!",
			data:        data,
			expectedStr: "Hello !",
			expectErr:   false,
		},
		{
			name:     "deep nesting",
			template: "Value: ${level1.level2.level3}",
			data: map[string]interface{}{
				"level1": map[string]interface{}{
					"level2": map[string]interface{}{
						"level3": "deep value",
					},
				},
			},
			expectedStr: "Value: deep value",
			expectErr:   false,
		},
		{
			name:     "boolean value",
			template: "Active: ${active}",
			data: map[string]interface{}{
				"active": true,
			},
			expectedStr: "Active: true",
			expectErr:   false,
		},
		{
			name:     "nil value",
			template: "Null: ${null}",
			data: map[string]interface{}{
				"null": nil,
			},
			expectedStr: "Null: <nil>",
			expectErr:   false,
		},
		{
			name:     "zero value",
			template: "Zero: ${zero}",
			data: map[string]interface{}{
				"zero": 0,
			},
			expectedStr: "Zero: 0",
			expectErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ResolveTemplate(tt.template, tt.data)

			if tt.expectErr {
				if err == nil {
					t.Errorf("ResolveTemplate() expected error but got none")
					return
				}
				if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("ResolveTemplate() error = %v, want error containing %v", err.Error(), tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("ResolveTemplate() unexpected error = %v", err)
				return
			}

			if result != tt.expectedStr {
				t.Errorf("ResolveTemplate() = %q, want %q", result, tt.expectedStr)
			}
		})
	}
}

func TestGetFormatedTemplate(t *testing.T) {
	data := map[string]interface{}{
		"user": map[string]interface{}{
			"name": "Alex",
			"id":   123,
		},
		"status": "active",
	}

	tests := []struct {
		name     string
		path     string
		data     map[string]interface{}
		expected interface{}
	}{
		{
			name:     "simple key",
			path:     "status",
			data:     data,
			expected: "active",
		},
		{
			name:     "nested key",
			path:     "user.name",
			data:     data,
			expected: "Alex",
		},
		{
			name:     "nested number",
			path:     "user.id",
			data:     data,
			expected: 123,
		},
		{
			name:     "missing key",
			path:     "missing",
			data:     data,
			expected: "",
		},
		{
			name:     "missing nested key",
			path:     "user.missing",
			data:     data,
			expected: "",
		},
		{
			name:     "path past non-map",
			path:     "status.detail",
			data:     data,
			expected: "",
		},
		{
			name:     "empty path",
			path:     "",
			data:     data,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getFormatedTemplate(tt.path, tt.data)
			if result != tt.expected {
				t.Errorf("getFormatedTemplate() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestResolveTemplateEmptyData(t *testing.T) {
	emptyData := map[string]interface{}{}

	result, err := ResolveTemplate("Hello ${name}!", emptyData)
	if err != nil {
		t.Errorf("ResolveTemplate() unexpected error = %v", err)
	}
	if result != "Hello !" {
		t.Errorf("ResolveTemplate() = %q, want %q", result, "Hello !")
	}
}

func TestResolveTemplateComplexScenario(t *testing.T) {
	data := map[string]interface{}{
		"config": map[string]interface{}{
			"db": map[string]interface{}{
				"host": "localhost",
				"port": 5432,
			},
			"app": map[string]interface{}{
				"name":    "MyApp",
				"version": "1.0.0",
			},
		},
		"env": "production",
	}

	template := "Connecting to ${config.db.host}:${config.db.port} for ${config.app.name} v${config.app.version} in ${env} environment"
	expected := "Connecting to localhost:5432 for MyApp v1.0.0 in production environment"

	result, err := ResolveTemplate(template, data)
	if err != nil {
		t.Errorf("ResolveTemplate() unexpected error = %v", err)
	}
	if result != expected {
		t.Errorf("ResolveTemplate() = %q, want %q", result, expected)
	}
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
