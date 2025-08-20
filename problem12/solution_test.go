package main

import (
	"reflect"
	"testing"
)

func TestParseVersion(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		want       *Version
		wantErrMsg string
	}{
		{
			name:  "with only major, minor and patch",
			input: "1.2.3",
			want: &Version{
				Major: 1, Minor: 2, Patch: 3,
			},
		},
		{
			name:  "with main and pre-release",
			input: "1.2.3-alpha",
			want: &Version{
				Major: 1, Minor: 2, Patch: 3, PreRelease: "alpha",
			},
		},
		{
			name:  "with main and metadata",
			input: "1.2.3+build.2",
			want: &Version{
				Major: 1, Minor: 2, Patch: 3, Metadata: "build.2",
			},
		},
		{
			name:  "full version with all parts",
			input: "2.0.1-alpha.1+build.987",
			want: &Version{
				Major: 2, Minor: 0, Patch: 1, PreRelease: "alpha.1", Metadata: "build.987",
			},
		},
		{
			name:  "multi-digit core numbers",
			input: "1.10.20",
			want: &Version{
				Major: 1, Minor: 10, Patch: 20,
			},
		},
		{
			name:  "zero major version",
			input: "0.4.1",
			want: &Version{
				Major: 0, Minor: 4, Patch: 1,
			},
		},
		{
			name:       "error on missing patch version",
			input:      "1.2",
			want:       nil,
			wantErrMsg: "invalid core version format: 1.2",
		},
		{
			name:       "error on empty pre-release identifier",
			input:      "1.2.3-",
			want:       nil,
			wantErrMsg: "pre-release identifier cannot be empty",
		},
		{
			name:       "error on empty build metadata",
			input:      "1.2.3+",
			want:       nil,
			wantErrMsg: "build metadata cannot be empty",
		},
		{
			name:       "error on non-numeric minor version",
			input:      "1.b.3",
			want:       nil,
			wantErrMsg: "minor version is not a number: b",
		},
		{
			name:       "error on leading invalid characters",
			input:      "v1.0.0",
			want:       nil,
			wantErrMsg: "major version is not a number: v1",
		},
		{
			name:       "error on both pre-release and metadata but malformed patch",
			input:      "1.2.alpha-3+build.1",
			want:       nil,
			wantErrMsg: "patch version is not a number: alpha",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseVersion(tt.input)
			if tt.wantErrMsg != "" {
				if err == nil {
					t.Errorf("ParseVersion() error = nil, wantErr %v", tt.wantErrMsg)
					return
				}
				if err.Error() != tt.wantErrMsg {
					t.Errorf("ParseVersion() error = %v, wantErr %v", err, tt.wantErrMsg)
				}
				return
			}
			if err != nil {
				t.Errorf("ParseVersion() unexpected error = %v", err)
				return
			}
			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("ParseVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}
