package main

import (
	"errors"
	"reflect"
	"testing"
)

func TestParseQuery(t *testing.T) {
	ErrMalformedInput := errors.New("malformed query")
	tests := []struct {
		name    string
		input   string
		want    map[string]interface{}
		wantErr bool
		err     error
	}{
		{
			name:  "basic input with strings and interface",
			input: "user.name=Alex%20Doe&user.id=123&roles[0]=admin&roles[1]=editor&active=true",
			want: map[string]interface{}{
				"user": map[string]interface{}{
					"id":   "123",
					"name": "Alex Doe",
				},
				"roles": []interface{}{
					"admin",
					"editor",
				},
				"active": "true",
			},
			wantErr: false,
		},
		{
			name:  "nested object",
			input: "address.city=New%20York&address.zip=10001&address.country=USA",
			want: map[string]interface{}{
				"address": map[string]interface{}{
					"city":    "New York",
					"zip":     "10001",
					"country": "USA",
				},
			},
			wantErr: false,
		},
		{
			name:  "array of numbers",
			input: "scores[0]=10&scores[1]=20&scores[2]=30",
			want: map[string]interface{}{
				"scores": []interface{}{
					"10",
					"20",
					"30",
				},
			},
			wantErr: false,
		},
		{
			name:  "multiple nested arrays and objects",
			input: "team.members[0].name=Alice&team.members[0].role=dev&team.members[1].name=Bob&team.members[1].role=qa",
			want: map[string]interface{}{
				"team": map[string]interface{}{
					"members": []interface{}{
						map[string]interface{}{
							"name": "Alice",
							"role": "dev",
						},
						map[string]interface{}{
							"name": "Bob",
							"role": "qa",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:  "empty values",
			input: "foo=&bar=",
			want: map[string]interface{}{
				"foo": "",
				"bar": "",
			},
			wantErr: false,
		},
		{
			name:    "malformed query missing key",
			input:   "=value",
			want:    nil,
			wantErr: true,
			err:     ErrMalformedInput,
		},
		{
			name:    "malformed query missing value",
			input:   "onlykey",
			want:    nil,
			wantErr: true,
			err:     ErrMalformedInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseQuery(tt.input)
			if err != nil && errors.Is(err, errors.New("malformed input")) {
				t.Fatalf("error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if got == nil {
					t.Fatalf("got = nil, want = %#v", tt.want)
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Flags = %#v, want %#v", got, tt.want)
				}
			}
		})
	}
}
