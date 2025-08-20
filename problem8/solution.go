package main

import (
	"errors"
	"fmt"
	"strings"
)

func main() {
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

	testCases := []struct {
		name        string
		template    string
		data        map[string]interface{}
		expectedStr string
		expectErr   bool
	}{
		{
			name:        "Happy Path - Basic Substitution",
			template:    "Hello, ${user.name}!",
			data:        data,
			expectedStr: "Hello, Alex!",
			expectErr:   false,
		},
		{
			name:        "Multiple and Nested Placeholders",
			template:    "Service for ${user.name}: ${appointment.service}.",
			data:        data,
			expectedStr: "Service for Alex: Dental Check-up.",
			expectErr:   false,
		},
		{
			name:        "Resolving to a Non-String Value (Integer)",
			template:    "User ID: ${user.id}",
			data:        data,
			expectedStr: "User ID: 123",
			expectErr:   false,
		},
		{
			name:        "Missing Path - Top-Level Key",
			template:    "System info: ${system.version}",
			data:        data,
			expectedStr: "System info: ",
			expectErr:   false,
		},
		{
			name:        "Missing Path - Nested Key",
			template:    "User address: ${user.address.city}",
			data:        data,
			expectedStr: "User address: ",
			expectErr:   false,
		},
		{
			name:        "Path Traverses Past a Final Value",
			template:    "Status detail: ${status.detail}",
			data:        data,
			expectedStr: "Status detail: ",
			expectErr:   false,
		},
		{
			name:        "Error Case - Malformed Placeholder (Unclosed)",
			template:    "Hello, ${user.name",
			data:        data,
			expectedStr: "",
			expectErr:   true,
		},
		{
			name:        "Edge Case - Template with No Placeholders",
			template:    "This is a static string.",
			data:        data,
			expectedStr: "This is a static string.",
			expectErr:   false,
		},
		{
			name:        "Edge Case - Empty Template String",
			template:    "",
			data:        data,
			expectedStr: "",
			expectErr:   false,
		},
		{
			name:        "Edge Case - Adjacent Placeholders",
			template:    "${user.name}${status}",
			data:        data,
			expectedStr: "Alexactive",
			expectErr:   false,
		},
		{
			name:        "Edge Case - Placeholder at Start and End",
			template:    "${status} user ${user.name}",
			data:        data,
			expectedStr: "active user Alex",
			expectErr:   false,
		},
	}

	fmt.Println("Running test cases...")
	for _, tc := range testCases {
		resultStr, err := ResolveTemplate(tc.template, tc.data)

		if !tc.expectErr && err != nil {
			fmt.Printf("❌ FAILED: %s\n", tc.name)
			fmt.Printf("   Template: %#v\n", tc.template)
			fmt.Printf("   Expected no error, but got: %v\n\n", err)
			continue
		}

		if tc.expectErr {
			if err == nil {
				fmt.Printf("❌ FAILED: %s\n", tc.name)
				fmt.Printf("   Template: %#v\n", tc.template)
				fmt.Printf("   Expected an error, but got nil\n\n")
			} else {
				fmt.Printf("✅ PASSED: %s (Correctly returned an error)\n\n", tc.name)
			}
			continue
		}

		if resultStr != tc.expectedStr {
			fmt.Printf("❌ FAILED: %s\n", tc.name)
			fmt.Printf("   Template: %#v\n", tc.template)
			fmt.Printf("   Expected: %#v\n", tc.expectedStr)
			fmt.Printf("   Got:      %#v\n\n", resultStr)
		} else {
			fmt.Printf("✅ PASSED: %s\n\n", tc.name)
		}
	}
}

func ResolveTemplate(text string, data map[string]interface{}) (string, error) {
	var formatedString strings.Builder
	formatedString.Grow(len(text))
	start := 0

	for {
		open := strings.Index(text[start:], "${")
		if open == -1 {
			formatedString.WriteString(text[start:])
			break
		}
		open += start
		formatedString.WriteString(text[start:open])

		end := strings.Index(text[open:], "}")
		if end == -1 {
			return "", errors.New("malformed placeholder: missing closing '}'")
		}
		end += open

		nested := strings.Index(text[open+2:end], "${")
		if nested != -1 {
			return "", errors.New("malformed template")
		}
		ans := getFormatedTemplate(text[open+2:end], data)
		formatedString.WriteString(fmt.Sprintf("%v", ans))
		start = end + 1
	}

	return formatedString.String(), nil
}

func getFormatedTemplate(unformattedStr string, data map[string]interface{}) interface{} {
	var current interface{} = data
	parts := strings.Split(unformattedStr, ".")
	for _, ele := range parts {
		m, ok := current.(map[string]interface{})
		if !ok {
			return ""
		}
		val, ok := m[ele]
		if !ok {
			return ""
		}
		current = val
	}
	return current
}
