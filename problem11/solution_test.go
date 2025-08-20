package main

import "testing"

func TestDecompress(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
		err     error
	}{
		{
			name:    "with only string",
			input:   "asjnajksdnasdadsad",
			want:    "asjnajksdnasdadsad",
			wantErr: false,
		},
		{
			name:    "with basic compression",
			input:   "2[ab]2[cd]ef",
			want:    "ababcdcdef",
			wantErr: false,
		},
		{
			name:    "with nested string",
			input:   "2[a3[b]]",
			want:    "abbbabbb",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Decompress(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("expect %s but got an err %v", tt.want, err)
			}
			if tt.want != got {
				t.Fatalf("expect %s but got  %v", tt.want, got)
			}
		})
	}
}
