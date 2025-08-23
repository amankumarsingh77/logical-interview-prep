package main

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func setupTestDir(t *testing.T) string {
	rootDir := t.TempDir()
	content1 := []byte("hello world")
	content2 := []byte("unique content")
	_ = os.WriteFile(filepath.Join(rootDir, "fileA.txt"), content1, 0666)
	_ = os.Mkdir(filepath.Join(rootDir, "subdir"), 0777)
	_ = os.WriteFile(filepath.Join(rootDir, "subdir", "fileC.txt"), content1, 0666)

	_ = os.WriteFile(filepath.Join(rootDir, "empty1.txt"), []byte{}, 0666)
	_ = os.WriteFile(filepath.Join(rootDir, "subdir", "empty2.txt"), []byte{}, 0666)

	_ = os.WriteFile(filepath.Join(rootDir, "fileB.txt"), content2, 0666)

	return rootDir
}

func TestFindDuplicates(t *testing.T) {
	rootDir := setupTestDir(t)

	duplicates, err := FindDuplicates(rootDir)
	if err != nil {
		t.Fatalf("FindDuplicates failed: %v", err)
	}

	expectedGroups := 2
	if len(duplicates) != expectedGroups {
		t.Fatalf("expected %d duplicate groups, but got %d", expectedGroups, len(duplicates))
	}

	h1 := sha256.Sum256([]byte("hello world"))
	hash1 := fmt.Sprintf("%x", h1)

	group1, ok := duplicates[hash1]
	if !ok {
		t.Fatalf("expected to find hash for 'hello world', but did not")
	}

	expectedPaths1 := []string{
		filepath.Join(rootDir, "fileA.txt"),
		filepath.Join(rootDir, "subdir", "fileC.txt"),
	}

	sort.Strings(group1)
	sort.Strings(expectedPaths1)

	if len(group1) != len(expectedPaths1) {
		t.Errorf("expected %d duplicates for hash %s, got %d", len(expectedPaths1), hash1, len(group1))
	}
	for i := range group1 {
		if group1[i] != expectedPaths1[i] {
			t.Errorf("path mismatch: expected %s, got %s", expectedPaths1[i], group1[i])
		}
	}
	h2 := sha256.Sum256([]byte{})
	hash2 := fmt.Sprintf("%x", h2)

	group2, ok := duplicates[hash2]
	if !ok {
		t.Fatalf("expected to find hash for empty files, but did not")
	}

	if len(group2) != 2 {
		t.Errorf("expected 2 duplicates for empty files, got %d", len(group2))
	}
}
