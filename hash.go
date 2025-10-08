package main

import (
	"crypto/md5"
	"strings"
	"path/filepath"
	"io/fs"
	"io"
	"os"
	"fmt"
)

func ListFiles(rootDir string) []string {
	const minSize = 1024 * 1024 * 10 // 10 MiB
	listedFiles := make([]string, 0)
	
	err := filepath.WalkDir(rootDir,
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.IsDir() && strings.HasPrefix(d.Name(), ".") {
				return filepath.SkipDir
			}
			
			if !d.IsDir() {
				fileInfo, err := d.Info()
				if err != nil {
					return err
				}
				if fileInfo.Size() >= minSize {
					listedFiles = append(listedFiles, path)
				}
			}
			
			return nil
		},
	)
	if err != nil {
		fmt.Printf("[ERROR] Could not list directory contents: %v\n", err)
		os.Exit(69)
	}
	
	return listedFiles
}

func HashFile(filePath string) string {
	f, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("[ERROR] Could not open file %s: %v\n", filePath, err)
		os.Exit(69)
	}
	defer f.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, f); err != nil {
		fmt.Printf("[ERROR] Could not copy contents of file %s to hash: %v\n", filePath, err)
		os.Exit(69)
	}

	return fmt.Sprintf("%x", hash.Sum(nil))
}
