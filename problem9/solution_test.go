package main

import (
	"errors"
	"reflect"
	"testing"
)

func TestParsedCommand(t *testing.T) {
	var tests = []struct {
		name    string
		input   string
		want    *ParsedCommand
		wantErr bool
		error   string
	}{
		{
			name:  "Happy Path - Full Example with Quoted Values",
			input: "publish --path \"/reports/q1 report.pdf\" --user \"alex doe\" -v",
			want: &ParsedCommand{
				Command: "publish",
				Flags: map[string]interface{}{
					"path": "/reports/q1 report.pdf",
					"user": "alex doe",
					"v":    true,
				},
			},
			wantErr: false,
		},
		{
			name:  "Multi-word Command",
			input: "remote add --name \"origin\" -f",
			want: &ParsedCommand{
				Command: "remote add",
				Flags: map[string]interface{}{
					"name": "origin",
					"f":    true,
				},
			},
			wantErr: false,
		},
		{
			name:  "Boolean flag followed by another flag",
			input: "deploy -v --force --no-cache",
			want: &ParsedCommand{
				Command: "deploy",
				Flags: map[string]interface{}{
					"v":        true,
					"force":    true,
					"no-cache": true,
				},
			},
			wantErr: false,
		},
		{
			name:  "No flags, just a command",
			input: "git status",
			want: &ParsedCommand{
				Command: "git status",
				Flags:   map[string]interface{}{},
			},
			wantErr: false,
		},
		{
			name:  "No command, just flags",
			input: "--all -v",
			want: &ParsedCommand{
				Command: "",
				Flags: map[string]interface{}{
					"all": true,
					"v":   true,
				},
			},
			wantErr: false,
		},
		{
			name:  "Handles extra whitespace between tokens",
			input: "  publish     --path \"/reports/q1 report.pdf\"   -v  ",
			want: &ParsedCommand{
				Command: "publish",
				Flags: map[string]interface{}{
					"path": "/reports/q1 report.pdf",
					"v":    true,
				},
			},
			wantErr: false,
		},
		{
			name:  "Edge Case: Empty input string",
			input: "",
			want: &ParsedCommand{
				Command: "",
				Flags:   map[string]interface{}{},
			},
			wantErr: true,
			error:   "empty input",
		},
		// I was not able to think of these edge cases :(
		//{
		//	name:  "Edge Case: Empty quoted value",
		//	input: "update --tag \"\" -f",
		//	want: &ParsedCommand{
		//		Command: "update",
		//		Flags: map[string]interface{}{
		//			"tag": "",
		//			"f":   true,
		//		},
		//	},
		//	wantErr: false,
		//},
		//{
		//	name:  "Invalid Input: Ignores tokens with invalid flag prefixes",
		//	input: "run ---verbose --- -x",
		//	want: &ParsedCommand{
		//		Command: "run ---verbose ---", // '---verbose' and '---' are not flags
		//		Flags: map[string]interface{}{
		//			"x": true,
		//		},
		//	},
		//	wantErr: false,
		//},
		{
			name:    "Error Case: Unclosed quote",
			input:   "commit -m \"feat: add new feature",
			want:    nil,
			wantErr: true,
			error:   "unexpected token \"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseCommand(tt.input)
			if err != nil && errors.Is(err, errors.New(tt.error)) {
				t.Fatalf("error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if got == nil {
					t.Fatalf("got = nil, want = %#v", tt.want)
				}
				if got.Command != tt.want.Command {
					t.Errorf("Command = %q, want %q", got.Command, tt.want.Command)
				}
				if !reflect.DeepEqual(got.Flags, tt.want.Flags) {
					t.Errorf("Flags = %#v, want %#v", got.Flags, tt.want.Flags)
				}
			}
		})
	}
}
