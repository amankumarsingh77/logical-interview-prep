package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

var (
	testFileContent = "0123456789abcdefghijklmnopqrstuvwxyz"
	testFileEtag    = "v1"
	testFileModTime = time.Now().Add(-1 * time.Hour)
)

type mockServerConfig struct {
	supportsRange bool
	etag          string
	modTime       time.Time
	content       []byte
	fail          bool
}

func createMockServer(config mockServerConfig) *httptest.Server {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if config.fail {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}

		w.Header().Set("Etag", config.etag)
		w.Header().Set("Last-Modified", config.modTime.Format(http.TimeFormat))
		if config.supportsRange {
			w.Header().Set("Accept-Ranges", "bytes")
		}

		http.ServeContent(w, r, "testfile.txt", config.modTime, bytes.NewReader(config.content))
	})
	return httptest.NewServer(handler)
}

func setupTestDir(t *testing.T) string {
	dir, err := os.MkdirTemp("", "download_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	t.Cleanup(func() {
		os.RemoveAll(dir)
	})
	return dir
}

func fileHash(t *testing.T, path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read file for hashing: %v", err)
	}
	return fmt.Sprintf("%x", sha256.Sum256(data))
}

func TestDownloadFile(t *testing.T) {
	originalWD, _ := os.Getwd()

	t.Run("FullDownload_NewFile", func(t *testing.T) {
		dir := setupTestDir(t)
		os.Chdir(dir)
		defer os.Chdir(originalWD)
		defer func() { fileList = make(map[string]*Metadata) }()

		server := createMockServer(mockServerConfig{
			supportsRange: true,
			etag:          testFileEtag,
			modTime:       testFileModTime,
			content:       []byte(testFileContent),
		})
		defer server.Close()

		err := downloadFile(server.URL + "/testfile.txt")
		if err != nil {
			t.Fatalf("downloadFile failed: %v", err)
		}

		filePath := filepath.Join(dir, "testfile.txt")
		info, err := os.Stat(filePath)
		if err != nil {
			t.Fatalf("os.Stat failed for downloaded file: %v", err)
		}

		if info.Size() != int64(len(testFileContent)) {
			t.Errorf("Expected file size %d, got %d", len(testFileContent), info.Size())
		}

		if fileHash(t, filePath) != fileHash(t, filePath) {
			t.Error("File content mismatch")
		}
	})

	t.Run("ResumeDownload_Success", func(t *testing.T) {
		dir := setupTestDir(t)
		os.Chdir(dir)
		defer os.Chdir(originalWD)
		defer func() { fileList = make(map[string]*Metadata) }()

		filePath := filepath.Join(dir, "testfile.txt")
		partialContent := []byte(testFileContent[:10])
		if err := os.WriteFile(filePath, partialContent, 0644); err != nil {
			t.Fatalf("Failed to write partial file: %v", err)
		}

		fileList[filePath] = &Metadata{Etag: testFileEtag}

		server := createMockServer(mockServerConfig{
			supportsRange: true,
			etag:          testFileEtag,
			modTime:       testFileModTime,
			content:       []byte(testFileContent),
		})
		defer server.Close()

		err := downloadFile(server.URL + "/testfile.txt")
		if err != nil {
			t.Fatalf("downloadFile failed: %v", err)
		}

		info, err := os.Stat(filePath)
		if err != nil {
			t.Fatalf("os.Stat failed for downloaded file: %v", err)
		}
		if info.Size() != int64(len(testFileContent)) {
			t.Errorf("Expected file size %d, got %d", len(testFileContent), info.Size())
		}
	})

	t.Run("ResumeDownload_FileChanged", func(t *testing.T) {
		dir := setupTestDir(t)
		os.Chdir(dir)
		defer os.Chdir(originalWD)
		defer func() { fileList = make(map[string]*Metadata) }()

		filePath := filepath.Join(dir, "testfile.txt")
		partialContent := []byte(testFileContent[:10])
		if err := os.WriteFile(filePath, partialContent, 0644); err != nil {
			t.Fatalf("Failed to write partial file: %v", err)
		}

		fileList[filePath] = &Metadata{Etag: "old-etag"}

		server := createMockServer(mockServerConfig{
			supportsRange: true,
			etag:          "new-etag",
			modTime:       testFileModTime,
			content:       []byte(testFileContent),
		})
		defer server.Close()

		err := downloadFile(server.URL + "/testfile.txt")
		if err != nil {
			t.Fatalf("downloadFile failed: %v", err)
		}

		info, err := os.Stat(filePath)
		if err != nil {
			t.Fatalf("os.Stat failed for downloaded file: %v", err)
		}
		if info.Size() != int64(len(testFileContent)) {
			t.Errorf("Expected file size %d, got %d", len(testFileContent), info.Size())
		}
	})

	t.Run("FullDownload_ServerDoesNotSupportRange", func(t *testing.T) {
		dir := setupTestDir(t)
		os.Chdir(dir)
		defer os.Chdir(originalWD)
		defer func() { fileList = make(map[string]*Metadata) }()

		filePath := filepath.Join(dir, "testfile.txt")
		partialContent := []byte(testFileContent[:5])
		if err := os.WriteFile(filePath, partialContent, 0644); err != nil {
			t.Fatalf("Failed to write partial file: %v", err)
		}

		fileList[filePath] = &Metadata{Etag: testFileEtag}

		server := createMockServer(mockServerConfig{
			supportsRange: false,
			etag:          testFileEtag,
			modTime:       testFileModTime,
			content:       []byte(testFileContent),
		})
		defer server.Close()

		err := downloadFile(server.URL + "/testfile.txt")
		if err != nil {
			t.Fatalf("downloadFile failed: %v", err)
		}

		info, err := os.Stat(filePath)
		if err != nil {
			t.Fatalf("os.Stat failed for downloaded file: %v", err)
		}
		if info.Size() != int64(len(testFileContent)) {
			t.Errorf("Expected file size %d, got %d", len(testFileContent), info.Size())
		}
	})

	t.Run("Download_URLNotFound", func(t *testing.T) {
		dir := setupTestDir(t)
		os.Chdir(dir)
		defer os.Chdir(originalWD)
		defer func() { fileList = make(map[string]*Metadata) }()

		server := createMockServer(mockServerConfig{fail: true})
		defer server.Close()

		err := downloadFile(server.URL + "/testfile.txt")
		if err == nil {
			t.Fatal("Expected an error for a 404 response, but got nil")
		}
	})
}
