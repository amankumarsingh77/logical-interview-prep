package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

type Metadata struct {
	Name          string    `json:"name"`
	Path          string    `json:"path"`
	Etag          string    `json:"etag"`
	LastModified  time.Time `json:"lastModified"`
	Size          int64     `json:"size"`
	StartByte     int64     `json:"startByte"`
	Exists        bool      `json:"exists"`
	IsModified    bool      `json:"isModified"`
	SupportsRange bool      `json:"supportsRange"`
	DownloadUrl   string    `json:"downloadUrl"`
}

var fileList map[string]*Metadata

func init() {
	fileList = make(map[string]*Metadata)
	loadMetadata()
}

func downloadFile(urlStr string) error {
	metadata, err := getFileMetadata(urlStr)
	if err != nil {
		return err
	}

	log.Println("Metadata:", metadata)

	if metadata.Exists {
		if metadata.IsModified {
			res, err := http.Get(urlStr)
			if err != nil {
				return fmt.Errorf("failed to download the file: %w", err)
			}
			if err := saveFile(res.Body, metadata.Path, true); err != nil {
				return err
			}
		} else {
			req, err := http.NewRequest("GET", urlStr, nil)
			if err != nil {
				return fmt.Errorf("failed to create request: %w", err)
			}
			req.Header.Set("Range", fmt.Sprintf("bytes=%d-", metadata.StartByte))
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				return fmt.Errorf("failed to download the file: %w", err)
			}
			if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
				return errors.New("could not find the file")
			}

			if err := saveFile(resp.Body, metadata.Path, false); err != nil {
				return err
			}
		}
	} else {
		res, err := http.Get(urlStr)
		if err != nil {
			return fmt.Errorf("failed to download the file: %w", err)
		}
		if err := saveFile(res.Body, metadata.Path, true); err != nil {
			return err
		}
	}

	fileList[metadata.Path] = metadata

	return nil
}

func loadMetadata() {

	existing := make(map[string]*Metadata)
	if data, err := os.ReadFile("metadata.json"); err == nil {
		_ = json.Unmarshal(data, &existing)
	}
	for k, v := range existing {
		fileList[k] = v
	}
	log.Println(fileList)
}

func saveFile(body io.ReadCloser, path string, new bool) error {
	defer body.Close()

	var file *os.File
	var err error

	if new {
		file, err = os.Create(path)
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
	} else {
		file, err = os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			return fmt.Errorf("failed to open file for append: %w", err)
		}
	}
	defer file.Close()

	if _, err := io.Copy(file, body); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func getFileMetadata(urlLink string) (*Metadata, error) {
	metadata := &Metadata{}
	link, err := url.Parse(urlLink)
	if err != nil {
		return nil, fmt.Errorf("invalid url: %w", err)
	}

	dir, _ := os.Getwd()
	metadata.Name = filepath.Base(link.Path)
	metadata.Path = filepath.Join(dir, metadata.Name)
	metadata.DownloadUrl = link.String()

	if info, err := os.Stat(metadata.Path); err == nil {
		metadata.Size = info.Size()
		metadata.StartByte = info.Size()
		metadata.Exists = true
	}

	res, err := http.Head(link.String())
	if err != nil {
		return nil, fmt.Errorf("failed to send HEAD request: %w", err)
	}
	defer res.Body.Close()

	if res.Header.Get("Accept-Ranges") != "" {
		metadata.SupportsRange = true
	}

	if etag := res.Header.Get("Etag"); etag != "" {
		if file, ok := fileList[metadata.Path]; ok {
			if file.Etag != etag {
				metadata.IsModified = true
			}
		}
		metadata.Etag = etag
	}
	if val := res.Header.Get("Last-Modified"); val != "" {
		if t, err := time.Parse(time.RFC1123, val); err == nil {
			if file, ok := fileList[metadata.Path]; ok {
				if !file.LastModified.Equal(t) {
					metadata.IsModified = true
				}
			}
			metadata.LastModified = t
		}
	}
	fileList[metadata.Path] = metadata

	return metadata, nil
}

func saveMetaData() error {
	metadataFile := "metadata.json"
	existing := make(map[string]*Metadata)
	if data, err := os.ReadFile(metadataFile); err == nil {
		_ = json.Unmarshal(data, &existing)
	}

	for k, v := range fileList {
		existing[k] = v
	}

	file, err := os.Create(metadataFile)
	if err != nil {
		return fmt.Errorf("failed to create metadata file: %w", err)
	}
	defer file.Close()

	data, err := json.MarshalIndent(existing, "", "  ")
	if err != nil {
		return fmt.Errorf("failed marshalling the data: %w", err)
	}

	if _, err := file.Write(data); err != nil {
		return fmt.Errorf("failed to write the data: %w", err)
	}

	return nil
}

func main() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	done := make(chan error)
	go func() {
		done <- downloadFile("")
	}()
	select {
	case err := <-done:
		if err != nil {
			log.Println("Download failed:", err)
		} else {
			fmt.Println("Download completed successfully")
		}
	case sig := <-sigChan:
		fmt.Println("\nDownload interrupted:", sig)
	}
	if err := saveMetaData(); err != nil {
		log.Fatal("Failed to save metadata:", err)
	}
	fmt.Println("Metadata saved.")
}
