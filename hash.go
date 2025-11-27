package main

import (
	"strings"
	"path/filepath"
	"io/fs"
	"io"
	"os"
	"fmt"

	"lukechampine.com/blake3"
)

type FileInfo struct {
	Path string
	Size int64
}


func ListFiles(rootDir string) []FileInfo {
	const minSize = 1024 * 1024 * 10 // 10 MiB
	var listedFiles []FileInfo
	
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
					listedFiles = append(listedFiles, FileInfo{
						Path: path,
						Size: fileInfo.Size(),
					})
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

func readFileRange(filePath string, offset, length int64) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	_, err = file.Seek(offset, io.SeekStart)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, length)
	bytesRead, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return buf[:bytesRead], nil
}

func compareFileRanges(path1, path2 string, rangeSize int64) (bool, error) {
	first1, err := readFileRange(path1, 0, rangeSize)
	if err != nil {
		return false, err
	}

	first2, err := readFileRange(path2, 0, rangeSize)
	if err != nil {
		return false, err
	}

	if string(first1) != string(first2) {
		return false, nil
	}

	info1, err := os.Stat(path1)
	if err != nil {
		return false, err
	}

	info2, err := os.Stat(path2)
	if err != nil {
		return false, err
	}

	last1, err := readFileRange(path1, info1.Size()-rangeSize, rangeSize)
	if err != nil {
		return false, err
	}

	last2, err := readFileRange(path2, info2.Size()-rangeSize, rangeSize)
	if err != nil {
		return false, err
	}

	return string(last1) == string(last2), nil
}

/* This function needs some more mental effort... */
func getFileType(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))
	if ext == "" {
		return "no_extension"
	}
	return ext
}

func HashFile(filePath string) (string, error) {
	f, err := os.Open(filePath)
 	if err != nil {
 		fmt.Printf("[ERROR] Could not open file %s: %v\n", filePath, err)
		return "", err
 	}
 	defer f.Close()	

	hash := blake3.New(32, nil)
	
	if _, err := io.Copy(hash, f); err != nil {
		fmt.Printf("[ERROR] Could not copy contents of file %s to hash: %v\n", filePath, err)
	 	return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
