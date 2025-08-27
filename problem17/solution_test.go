package main

import (
	"errors"
	"reflect"
	"testing"
)

type input struct {
	imageUrls  []string
	maxWorkers int
}

func Test_GenerateThumbnails(t *testing.T) {
	tests := []struct {
		name  string
		input input
		want  *ThumbnailResult
	}{
		{
			name: "generate with one worker",
			input: input{
				imageUrls: []string{
					"images/cat.jpg",
					"images/dog.jpg",
					"images/bird-fail.png",
					"images/fish.gif",
					"images/lion.jpeg",
					"images/tiger-fail.bmp",
					"images/bear.svg",
				},
				maxWorkers: 3,
			},
			want: &ThumbnailResult{
				Successes: map[string]string{
					"images/fish.gif":  "thumbnails/fish.gif",
					"images/lion.jpeg": "thumbnails/lion.jpeg",
					"images/bear.svg":  "thumbnails/bear.svg",
					"images/cat.jpg":   "thumbnails/cat.jpg",
					"images/dog.jpg":   "thumbnails/dog.jpg",
				},
				Failures: map[string]error{
					"images/bird-fail.png":  errors.New("failed to process images/bird-fail.png"),
					"images/tiger-fail.bmp": errors.New("failed to process images/tiger-fail.bmp"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := generateThumbnails(tt.input.imageUrls, tt.input.maxWorkers)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("expected %v but got %v", tt.want, got)
			}
		})
	}
}
