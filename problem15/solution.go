package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// this is just a simple solution will work on it tomorrow ;)

func FindDuplicates(rootDir string) (map[string][]string, error) {
	duplicates := make(map[string][]string)
	sizeOfFiles := make(map[int64][]string)
	if err := filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			info, err := d.Info()
			if err == nil {
				sizeOfFiles[info.Size()] = append(sizeOfFiles[info.Size()], path)
			}
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("error occured : %v", err)
	}
	for _, files := range sizeOfFiles {
		if len(files) < 2 {
			continue
		}
		hashes := make(map[string][]string)
		for _, f := range files {
			file, err := os.Open(f)
			if err != nil {
				continue
			}
			hasher := sha256.New()
			if _, err = io.Copy(hasher, file); err != nil {
				file.Close()
				continue
			}
			file.Close()
			hashStr := hex.EncodeToString(hasher.Sum(nil))
			hashes[hashStr] = append(hashes[hashStr], f)
		}
		for hash, paths := range hashes {
			if len(paths) > 1 {
				duplicates[hash] = paths
			}
		}
	}
	return duplicates, nil
}
